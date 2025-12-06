// Package workerpool 提供一个简单的并发 worker 池实现。
//
// Worker 池是 Go 并发编程中的常见模式，用于限制同时执行的 goroutine 数量，
// 避免资源耗尽（如文件句柄、数据库连接等）。
//
// 基本使用流程：
//  1. 创建池: pool := workerpool.New(workerCount, queueSize)
//  2. 启动 worker: pool.Run()
//  3. 提交任务: pool.Submit(task)
//  4. 关闭投递: pool.Close()
//  5. 等待完成: pool.Wait()
//
// 示例:
//
//	pool := workerpool.New(3, 10) // 3 个 worker，队列缓冲 10
//	pool.Run()
//	for i := 0; i < 100; i++ {
//	    pool.Submit(func() error {
//	        // 执行任务...
//	        return nil
//	    })
//	}
//	pool.Close()
//	pool.Wait()
package workerpool

import "sync"

// Task 定义一个可执行的任务。
// 任务是一个函数，执行后返回 error（如果执行成功返回 nil）。
// 这种设计允许任务携带闭包捕获的数据，同时提供错误反馈机制。
type Task func() error

// Pool 是一个固定大小的 worker 池。
//
// 核心组件说明：
//   - size: worker（goroutine）数量，决定了并发度
//   - tasks: 任务通道，生产者通过它投递任务，worker 从中消费
//   - wg: WaitGroup，用于等待所有 worker 完成工作
//
// 工作原理：
//   - 调用 Run() 后，会启动 size 个 goroutine 作为 worker
//   - 每个 worker 不断从 tasks 通道读取任务并执行
//   - 当 tasks 通道关闭且为空时，worker 自动退出
type Pool struct {
	size  int       // worker 数量（并发 goroutine 数）
	tasks chan Task // 任务队列（带缓冲的通道）
	wg    sync.WaitGroup
}

// New 创建一个新的 worker 池。
//
// 参数说明：
//   - size: worker 数量，即同时执行任务的 goroutine 数量
//   - queueSize: 任务队列缓冲大小，决定了可以暂存多少待处理任务
//
// 缓冲大小的选择：
//   - queueSize = 0: 无缓冲，Submit 会阻塞直到有 worker 空闲
//   - queueSize > 0: 有缓冲，可以快速投递任务而不阻塞（直到缓冲满）
//
// 推荐做法：queueSize 设为 size 的 2-3 倍，平衡内存使用与投递效率。
func New(size int, queueSize int) *Pool {
	return &Pool{
		size:  size,
		tasks: make(chan Task, queueSize),
	}
}

// Run 启动所有 worker goroutine。
//
// 调用此方法后，worker 池开始工作：
//  1. 创建 size 个 goroutine
//  2. 每个 goroutine 循环从 tasks 通道读取任务
//  3. 读取到任务后立即执行
//  4. 当 tasks 通道关闭且为空时，goroutine 退出
//
// 注意：必须在 Submit 之前调用 Run()，否则任务会阻塞在通道中。
func (p *Pool) Run() {
	for i := 0; i < p.size; i++ {
		p.wg.Add(1) // 为每个 worker 增加计数
		go func() {
			defer p.wg.Done() // worker 退出时减少计数
			// for-range 会持续从通道读取，直到通道关闭且为空
			for task := range p.tasks {
				// 执行任务，这里忽略返回值
				// 生产环境中应该记录或处理错误
				_ = task()
			}
		}()
	}
}

// Submit 向池中提交一个任务。
//
// 行为说明：
//   - 如果队列未满，任务立即入队，Submit 返回
//   - 如果队列已满，Submit 阻塞，直到有 worker 取走任务腾出空间
//
// 这种阻塞行为被称为"背压"（backpressure），是一种流量控制机制，
// 防止生产者速度远超消费者导致内存无限增长。
//
// 注意：不要在调用 Close() 之后再调用 Submit()，会引发 panic。
func (p *Pool) Submit(t Task) {
	p.tasks <- t
}

// Close 关闭任务通道，表示不再有新任务。
//
// 调用后：
//   - 不能再调用 Submit()（会 panic）
//   - worker 会处理完队列中剩余的任务
//   - 队列清空后 worker 自动退出
//
// 典型调用顺序：Submit 完所有任务 -> Close() -> Wait()
func (p *Pool) Close() {
	close(p.tasks)
}

// Wait 阻塞等待所有 worker 退出。
//
// 调用此方法前必须先调用 Close()，否则 Wait 会永远阻塞
// （因为 worker 在等待新任务，而没有新任务也没有关闭信号）。
//
// Wait 返回时，保证：
//   - 所有已提交的任务都已执行完毕
//   - 所有 worker goroutine 都已退出
func (p *Pool) Wait() {
	p.wg.Wait()
}
