package scheduler

import "crawler/concurrent-version-crawler/engine"

// 为每个 request 开启 goroutine 的方法不错，但是不可控。
// 所以我们创建了 request 队列  和  worker 队列

type QueuedScheduler struct {
	requestChan chan engine.Request      // request 的 channel,接收 Submit 的值
	workerChan  chan chan engine.Request // 每个 worker 有自己的 channel
}

// 获得 worker 的 channel
func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (*QueuedScheduler) ConfigureMasterWorkerChan(chan engine.Request) {
	panic("implement me")
}

// 1. 建立两个 channel
// 2. 开启 go func 和 for 循环，等待任务
// 3. 开启两个队列，当两个队列同时有值时，取 request 发送到 存在worker
// 4. 没有则监听来的是 request 还是 worker， 分别添加到队列中
func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		// 维护两个队列
		var requestQ []engine.Request
		var workerQ [] chan engine.Request

		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			// 用 select 判断来的是哪个 channel的数据, 并加入相应的队列
			select {
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
