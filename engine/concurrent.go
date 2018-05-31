package engine

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan    chan Item
}

type Scheduler interface {
	Submit(Request)
	WorkerChan() chan Request
	ReadyNotifier
	ConfigureMasterWorkerChan(chan Request) // 配置 worker 的 channel
	Run()
}

// 重构： 将 ready 拆分出来，减轻 createWorker ，其不再需要 scheduler ，只需要 ReadyNotifier
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	// out channel 接收 worker 的返回值
	out := make(chan ParseResult)
	// 创建 队列 并持续处理
	e.Scheduler.Run()
	// 创建 worker
	for i := 0; i < e.WorkerCount; i++ {
		// 此处重构， in 的 channel 不再由 engine 创建，因为 engine 不知道有一个还是多个 channel
		// engine 只需要跟 scheduler 要 channel 就行了。
		// out 用于 worker 向外发送结果
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	// 提交 request
	for _, r := range seeds {
		// 判断 url 是否存在
		if isDuplicate(r.Url) {
			// log.Printf("Duplicate request %s", r.Url)
			continue
		}
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		// 向 itemsaver 发送 items
		for _, item := range result.Items {
			go func() { e.ItemChan <- item }()
		}

		// 新的 request 提交到 scheduler
		for _, request := range result.Requests {
			// 判断 url 是否存在
			if isDuplicate(request.Url) {
				// log.Printf("Duplicate request %s", request.Url)
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

// 创建 worker
func createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// 配置 channel
			ready.WorkerReady(in)
			// 接收 scheduler 分配的 request 请求（发送在 scheduler 的 Run 中）
			request := <-in
			// fetch and parse request url
			result, e := worker(request)
			if e != nil {
				continue
			}

			// 通过 channel 输出结果（接收在 Run 中）
			out <- result
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}
