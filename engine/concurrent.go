package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
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
		e.Scheduler.Submit(r)
	}

	itemCount := 0
	for {
		result := <-out
		// 打印 Items 结果
		for _, item := range result.Items {
			log.Printf("Got item #%d: %v",itemCount, item)
			itemCount++
		}

		// 新的 request 提交到 scheduler
		for _, request := range result.Requests {
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
