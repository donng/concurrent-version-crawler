package main

import (
	"crawler/concurrent-version-crawler/engine"
	"crawler/concurrent-version-crawler/zhenai/parser"
	"crawler/concurrent-version-crawler/scheduler"
	"crawler/concurrent-version-crawler/persist"
)

const url = "http://www.zhenai.com/zhenghun"

func main() {
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan: persist.ItemSaver(),
	}

	e.Run(engine.Request{
		Url:        url,
		ParserFunc: parser.ParseCityList,
	})

	//e.Run(engine.Request{
	//	Url:        "http://www.zhenai.com/zhenghun/shanghai",
	//	ParserFunc: parser.ParseCity,
	//})
}
