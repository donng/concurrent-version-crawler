package engine

import (
	"log"
)

type SimpleEngine struct {
}

func (e SimpleEngine) Run(seeds ...Request) {
	// 生成请求的队列
	var requests []Request
	for _, seed := range seeds {
		requests = append(requests, seed)
	}

	for len(requests) > 0 {
		// 取出第一个 request 请求
		r := requests[0]
		requests = requests[1:]

		log.Printf("fetching url %s", r.Url)

		parseResult, err := worker(r)
		if err != nil {
			continue
		}

		// 将新的 request 添加到队列中
		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			log.Printf("Got items %v", item)
		}
	}
}
