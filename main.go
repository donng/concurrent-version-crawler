package main

import (
	"crawler/concurrent-version-crawler/engine"
	"crawler/concurrent-version-crawler/zhenai/parser"
	"crawler/concurrent-version-crawler/scheduler"
)

const url = "http://www.zhenai.com/zhenghun"

func main() {
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		WorkerCount: 100,
	}

	e.Run(engine.Request{
		Url:        url,
		ParserFunc: parser.ParseCityList,
	})
}
