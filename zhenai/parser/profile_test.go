package parser

import (
	"testing"
	"io/ioutil"
	"crawler/concurrent-version-crawler/engine"
	"crawler/concurrent-version-crawler/model"
)

func TestParseProfile(t *testing.T) {
	// 读取文件内容
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	// 解析文件内容
	result := ParseProfile(contents, "http://album.zhenai.com/u/1446528158", "阿兰")

	// 对比结果
	if len(result.Items) != 1 {
		t.Errorf("Items should contain 1 element, But was %v", result.Items)
	}

	actual := result.Items[0]

	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/1446528158",
		Type: "zhenai",
		Id:   "1446528158",
		Payload: model.Profile{
			Name:       "阿兰",
			Gender:     "女",
			Age:        27,
			Height:     158,
			Weight:     0,
			Income:     "3001-5000元",
			Marriage:   "未婚",
			Education:  "中专",
			Occupation: "--",
			Hokou:      "四川阿坝",
			Xingzuo:    "双子座",
			House:      "租房",
			Car:        "未购车",
		},
	}

	if actual != expected {
		t.Errorf("expected %v; but was %v", expected, actual)
	}
}
