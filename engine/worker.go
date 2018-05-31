package engine

import (
	"log"
	"crawler/concurrent-version-crawler/fetcher"
)

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
	return r.ParserFunc(content, r.Url), nil
}