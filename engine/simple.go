package engine

import (
	"crawler/concurrent-version-crawler/fetcher"
	"log"
)

type SimpleEngine struct {

}

func (e SimpleEngine) Run(seeds ...Request)  {
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

// 将 fetcher 和 parser 合并为 worker
func worker(r Request)  (ParseResult, error) {
	log.Printf("Fetching %s", r.Url)
	// 获得 url 内容
	content, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetchting url %s, %v", r.Url, err)
		return ParseResult{}, err
	}

	// 解析 url 内容
	return r.ParserFunc(content), nil
}