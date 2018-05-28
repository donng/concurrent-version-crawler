package parser

import (
	"crawler/concurrent-version-crawler/engine"
	"regexp"
)

const cityRe = `<a href="(http://album.zhenai.com/u/[\d]+)"[^>]*>([^<]+)</a>`

func ParseCity(content []byte) engine.ParseResult {
	re := regexp.MustCompile(cityRe)
	matches := re.FindAllSubmatch(content, -1)

	result := engine.ParseResult{}
	for _, m := range matches {
		name := string(m[2])
		result.Items = append(result.Items, "User "+ name)
		result.Requests = append(
			result.Requests, engine.Request{
				Url: string(m[1]),
				// 注意： 闭包用法
				ParserFunc: func(bytes []byte) engine.ParseResult {
					return ParseProfile(bytes, name)
				},
			},
		)
	}

	return result
}
