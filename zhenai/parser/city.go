package parser

import (
	"crawler/concurrent-version-crawler/engine"
	"regexp"
)

var profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[\d]+)"[^>]*>([^<]+)</a>`)
var cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)

func ParseCity(content []byte) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(content, -1)

	result := engine.ParseResult{}
	for _, m := range matches {
		name := string(m[2])
		url := string(m[1])
		result.Requests = append(
			result.Requests, engine.Request{
				Url: url,
				// 注意： 闭包用法
				ParserFunc: func(bytes []byte) engine.ParseResult {
					return ParseProfile(bytes, url, name)
				},
			},
		)
	}

	matches = cityUrlRe.FindAllSubmatch(content, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseCity,
		})
	}

	return result
}
