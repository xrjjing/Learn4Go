package todo

import (
	"sync"
	"testing"
)

// TestStore_ConcurrentCreate 测试并发创建 TODO 的 ID 唯一性
func TestStore_ConcurrentCreate(t *testing.T) {
	store := NewStore()
	const numGoroutines = 100
	const itemsPerGoroutine = 10

	// 用于收集所有创建的 TODO ID
	var mu sync.Mutex
	ids := make(map[int]bool)

	// 并发创建 TODO
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				todo, err := store.Create("test item", userID)
				if err != nil {
					t.Errorf("Create failed: %v", err)
					return
				}

				mu.Lock()
				if ids[todo.ID] {
					t.Errorf("Duplicate ID detected: %d", todo.ID)
				}
				ids[todo.ID] = true
				mu.Unlock()
			}
		}(uint(i))
	}

	wg.Wait()

	// 验证所有 ID 都是唯一的
	expectedCount := numGoroutines * itemsPerGoroutine
	if len(ids) != expectedCount {
		t.Errorf("Expected %d unique IDs, got %d", expectedCount, len(ids))
	}

	// 验证数据库中的条目数
	items, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != expectedCount {
		t.Errorf("Expected %d items in store, got %d", expectedCount, len(items))
	}
}

// TestStore_IDMonotonicity 测试 ID 单调递增特性
func TestStore_IDMonotonicity(t *testing.T) {
	store := NewStore()

	var lastID int
	for i := 0; i < 100; i++ {
		todo, err := store.Create("test", 1)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if i > 0 && todo.ID <= lastID {
			t.Errorf("ID not monotonic: previous=%d, current=%d", lastID, todo.ID)
		}
		lastID = todo.ID
	}
}

// TestStore_ListByUser 测试按用户过滤
func TestStore_ListByUser(t *testing.T) {
	store := NewStore()

	// 为不同用户创建 TODO
	for i := 0; i < 10; i++ {
		userID := uint(i % 3) // 3 个不同的用户
		_, err := store.Create("test", userID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// 测试用户 0 的 TODO
	user0Todos, err := store.ListByUser(0)
	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}
	expectedCount := 4 // 0, 3, 6, 9
	if len(user0Todos) != expectedCount {
		t.Errorf("Expected %d todos for user 0, got %d", expectedCount, len(user0Todos))
	}

	// 验证所有 TODO 都属于用户 0
	for _, todo := range user0Todos {
		if todo.UserID != 0 {
			t.Errorf("Expected UserID=0, got %d", todo.UserID)
		}
	}
}

// TestStore_ListPaged 测试分页功能
func TestStore_ListPaged(t *testing.T) {
	store := NewStore()

	// 创建 25 个 TODO
	for i := 0; i < 25; i++ {
		_, err := store.Create("test", 1)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	tests := []struct {
		page      int
		pageSize  int
		wantLen   int
		wantTotal int
	}{
		{page: 1, pageSize: 10, wantLen: 10, wantTotal: 25},
		{page: 2, pageSize: 10, wantLen: 10, wantTotal: 25},
		{page: 3, pageSize: 10, wantLen: 5, wantTotal: 25},
		{page: 4, pageSize: 10, wantLen: 0, wantTotal: 25},
		{page: 1, pageSize: 20, wantLen: 20, wantTotal: 25},
		{page: 0, pageSize: 10, wantLen: 10, wantTotal: 25}, // page < 1 默认为 1
	}

	for _, tt := range tests {
		items, total, err := store.ListPaged(tt.page, tt.pageSize)
		if err != nil {
			t.Errorf("ListPaged(%d, %d) failed: %v", tt.page, tt.pageSize, err)
			continue
		}
		if len(items) != tt.wantLen {
			t.Errorf("ListPaged(%d, %d) returned %d items, want %d",
				tt.page, tt.pageSize, len(items), tt.wantLen)
		}
		if total != tt.wantTotal {
			t.Errorf("ListPaged(%d, %d) returned total=%d, want %d",
				tt.page, tt.pageSize, total, tt.wantTotal)
		}
	}
}

// TestStore_Toggle 测试切换完成状态
func TestStore_Toggle(t *testing.T) {
	store := NewStore()

	todo, err := store.Create("test", 1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 初始状态应该是未完成
	if todo.Done {
		t.Error("New todo should not be done")
	}

	// 标记为已完成
	updated, ok, err := store.Toggle(todo.ID, true)
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}
	if !ok {
		t.Error("Toggle should return ok=true for existing item")
	}
	if !updated.Done {
		t.Error("Toggle should set Done=true")
	}

	// 再次切换为未完成
	updated, ok, err = store.Toggle(todo.ID, false)
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}
	if !ok {
		t.Error("Toggle should return ok=true for existing item")
	}
	if updated.Done {
		t.Error("Toggle should set Done=false")
	}

	// 测试不存在的 ID
	_, ok, err = store.Toggle(999999, true)
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}
	if ok {
		t.Error("Toggle should return ok=false for non-existent item")
	}
}

// TestStore_Delete 测试删除操作
func TestStore_Delete(t *testing.T) {
	store := NewStore()

	todo, err := store.Create("test", 1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 删除存在的 TODO
	ok, err := store.Delete(todo.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !ok {
		t.Error("Delete should return ok=true for existing item")
	}

	// 验证已删除
	_, exists, err := store.Get(todo.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if exists {
		t.Error("Get should return exists=false after deletion")
	}

	// 再次删除应该返回 false
	ok, err = store.Delete(todo.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if ok {
		t.Error("Delete should return ok=false for non-existent item")
	}
}

// BenchmarkStore_Create 性能测试：创建 TODO
func BenchmarkStore_Create(b *testing.B) {
	store := NewStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Create("test", 1)
	}
}

// BenchmarkStore_CreateParallel 性能测试：并发创建 TODO
func BenchmarkStore_CreateParallel(b *testing.B) {
	store := NewStore()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.Create("test", 1)
		}
	})
}

// BenchmarkStore_List 性能测试：列出所有 TODO
func BenchmarkStore_List(b *testing.B) {
	store := NewStore()
	// 预填充 1000 个 TODO
	for i := 0; i < 1000; i++ {
		_, _ = store.Create("test", 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.List()
	}
}
