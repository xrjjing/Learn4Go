# Java → Go 范式转换速查表

> 面向 3+ 年 Java 经验开发者的 Go 语言快速对照指南

## 快速对照表

| 主题 | Java | Go |
|------|------|-----|
| 包管理 | Maven/Gradle | Go Modules (`go mod`) |
| 可见性 | public/private/protected | 大写=导出，小写=包内 |
| 类定义 | `class` | `struct` |
| 继承 | `extends` | 嵌入 (embedding) |
| 接口 | `implements` (显式) | 隐式实现 |
| 异常 | try-catch-finally | `if err != nil` |
| 并发 | Thread/ExecutorService | goroutine + channel |
| 同步 | synchronized | sync.Mutex |
| 依赖注入 | Spring @Autowired | Wire / 手动注入 |
| 泛型 | Type Erasure | Go 1.18+ 类型参数 |
| 测试 | JUnit | `go test` |
| Web | Spring MVC | Gin/Echo/net/http |
| ORM | JPA/Hibernate | GORM/Ent |
| 配置 | application.yml | Viper |

---

## 1. 类与结构体

**核心区别**：Go 将 "类" 拆分为数据 (struct) 和行为 (方法)

### Java
```java
public class User {
    private String name;
    private int age;

    public User(String name, int age) {
        this.name = name;
        this.age = age;
    }

    public String getName() {
        return this.name;
    }

    public void setAge(int age) {
        this.age = age;
    }
}
```

### Go
```go
type User struct {
    name string  // 小写=私有
    Age  int     // 大写=公开
}

// 构造函数 (工厂模式)
func NewUser(name string, age int) *User {
    return &User{name: name, Age: age}
}

// 方法定义在结构体外部
// (u *User) 是接收者，类似 this
func (u *User) GetName() string {
    return u.name
}

func (u *User) SetAge(age int) {
    u.Age = age
}
```

**常见错误**：试图在 struct {} 内部定义方法

---

## 2. 继承 vs 组合

**核心区别**：Go 没有继承，使用嵌入 (embedding) 实现组合

### Java (继承)
```java
class Animal {
    public void eat() {
        System.out.println("Eating");
    }
}

class Dog extends Animal {
    public void bark() {
        System.out.println("Woof!");
    }
}

// 使用
Dog dog = new Dog();
dog.eat();  // 继承的方法
dog.bark();
```

### Go (嵌入)
```go
type Animal struct{}

func (a *Animal) Eat() {
    fmt.Println("Eating")
}

type Dog struct {
    Animal  // 匿名嵌入，方法被"提升"
}

func (d *Dog) Bark() {
    fmt.Println("Woof!")
}

// 使用
dog := Dog{}
dog.Eat()   // 从 Animal 提升的方法
dog.Bark()
```

**常见错误**：期望嵌入能实现多态 (Dog 不是 Animal 的子类型)

---

## 3. 接口实现

**核心区别**：Java 显式声明 implements，Go 隐式满足

### Java
```java
interface Writer {
    void write(String data);
}

// 必须显式声明 implements
class ConsoleWriter implements Writer {
    @Override
    public void write(String data) {
        System.out.println(data);
    }
}
```

### Go
```go
type Writer interface {
    Write(data string)
}

type ConsoleWriter struct{}

// 只要方法签名匹配，自动实现接口
// 无需声明 implements
func (c ConsoleWriter) Write(data string) {
    fmt.Println(data)
}

// 使用
var w Writer = ConsoleWriter{}  // 自动满足
w.Write("Hello")
```

**最佳实践**：在使用方定义接口，而非实现方

---

## 4. 错误处理

**核心区别**：Go 没有异常，错误是普通返回值

### Java
```java
public String readFile(String path) throws IOException {
    try {
        return Files.readString(Path.of(path));
    } catch (IOException e) {
        logger.error("读取失败: " + e.getMessage());
        throw e;
    } finally {
        // 清理资源
    }
}
```

### Go
```go
func readFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("读取 %s 失败: %w", path, err)
    }
    return string(data), nil
}

// 调用方
content, err := readFile("config.txt")
if err != nil {
    log.Printf("错误: %v", err)
    return  // 或 return err
}
// 继续处理 content
```

**常见错误**：
- 使用 `_` 忽略错误：`content, _ := readFile(path)` ❌
- 滥用 panic 模拟异常

---

## 5. 并发模型

**核心区别**：goroutine 比 Thread 轻量 1000 倍，通过 channel 通信

### Java
```java
// 线程池
ExecutorService executor = Executors.newFixedThreadPool(10);
Future<Integer> future = executor.submit(() -> {
    return compute();
});
Integer result = future.get();  // 阻塞等待

// CompletableFuture
CompletableFuture.supplyAsync(() -> compute())
    .thenAccept(result -> System.out.println(result));
```

### Go
```go
// goroutine + channel
ch := make(chan int)

go func() {
    result := compute()
    ch <- result  // 发送结果
}()

result := <-ch  // 接收结果

// 多任务等待
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        process(id)
    }(i)
}
wg.Wait()  // 等待所有完成
```

**核心原则**：不要通过共享内存来通信，而要通过通信来共享内存

---

## 6. 同步机制

### Java
```java
public class Counter {
    private int count = 0;

    public synchronized void increment() {
        count++;
    }

    // 或使用显式锁
    private final Lock lock = new ReentrantLock();

    public void incrementWithLock() {
        lock.lock();
        try {
            count++;
        } finally {
            lock.unlock();
        }
    }
}
```

### Go
```go
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()  // 确保解锁
    c.count++
}

// 读写锁 (多读单写)
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Get(key string) int {
    sm.mu.RLock()         // 读锁
    defer sm.mu.RUnlock()
    return sm.m[key]
}

func (sm *SafeMap) Set(key string, val int) {
    sm.mu.Lock()          // 写锁
    defer sm.mu.Unlock()
    sm.m[key] = val
}
```

**常见错误**：复制包含 Mutex 的结构体 (必须传指针)

---

## 7. 依赖注入

### Java (Spring)
```java
@Service
public class UserService {
    @Autowired
    private UserRepository repository;

    @Autowired
    private EmailService emailService;
}

// Spring 自动装配
```

### Go (手动注入)
```go
type UserService struct {
    repo  UserRepository
    email EmailService
}

// 构造函数注入 (推荐)
func NewUserService(repo UserRepository, email EmailService) *UserService {
    return &UserService{
        repo:  repo,
        email: email,
    }
}

// main.go 中显式组装
func main() {
    db := initDB()
    repo := NewUserRepository(db)
    email := NewEmailService(smtpConfig)
    service := NewUserService(repo, email)
    // ...
}
```

**优点**：依赖关系一目了然，无魔法

---

## 8. 常见陷阱

| 陷阱 | 说明 | 解决方案 |
|------|------|----------|
| nil map 写入 | `var m map[string]int; m["a"]=1` panic | 使用 `make(map[string]int)` |
| 切片共享底层数组 | append 可能影响原切片 | 需要时用 `copy` |
| goroutine 泄漏 | 未关闭的 channel 或无限阻塞 | 使用 context 控制生命周期 |
| 循环变量捕获 | for 循环中 goroutine 捕获变量 | 传参或 Go 1.22+ 自动修复 |
| defer 在循环中 | 每次迭代创建 defer 开销大 | 提取为函数 |
| 值接收者修改 | `func (u User) SetName()` 不会修改原对象 | 使用指针接收者 `*User` |

---

## 9. 框架对照

| 场景 | Java | Go |
|------|------|-----|
| Web 框架 | Spring Boot | Gin, Echo, Fiber |
| REST 注解 | @RestController, @GetMapping | `r.GET("/path", handler)` |
| 请求绑定 | @RequestBody | `c.ShouldBindJSON(&req)` |
| 中间件 | Filter, Interceptor | `r.Use(middleware)` |
| ORM | JPA/Hibernate | GORM, Ent |
| 数据库迁移 | Flyway, Liquibase | golang-migrate, goose |
| 配置 | application.yml | Viper |
| 日志 | Logback, Log4j | Zap, slog |
| 依赖注入 | Spring IoC | Wire, Fx |
| 微服务 | Spring Cloud | go-micro, Kratos |
| RPC | Feign | gRPC |

---

## 10. 学习建议

1. **忘掉 OOP 思维**：Go 是面向组合的，不是面向继承的
2. **拥抱显式**：显式错误处理、显式依赖注入、显式并发控制
3. **保持简单**：Go 的哲学是"少即是多"
4. **阅读标准库**：Go 标准库是最好的学习材料
5. **使用 go fmt**：代码风格不要争论，让工具统一
