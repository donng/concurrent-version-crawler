package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request) // 配置 worker 的 channel
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)
	out := make(chan ParseResult)
	e.Scheduler.ConfigureMasterWorkerChan(in)
	// 创建 worker
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, out)
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
func createWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			request := <-in
			// fetch and parse request url
			result, e := worker(request)
			if e != nil {
				continue
			}

			// 通过 channel 输出结果
			out <- result
		}
	}()
}
