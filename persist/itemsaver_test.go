package persist

import (
	"testing"
	"crawler/concurrent-version-crawler/model"
	"github.com/olivere/elastic"
	"context"
	"encoding/json"
	"crawler/concurrent-version-crawler/engine"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/1446528158",
		Type: "zhenai",
		Id:   "AWKtk0Z0DVFuU4T8-6jI",
		Payload:model.Profile{
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
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	const index = "dating_test"

	err = save(client, expected, index)
	if err != nil {
		panic(err)
	}


	result, e := client.Get().Index(index).Type(expected.Type).Id(expected.Id).Do(context.Background())
	if e != nil {
		panic(e)
	}

	t.Logf("%s", result.Source)
	var actual engine.Item
	err = json.Unmarshal([]byte(*result.Source), &actual)
	if err != nil {
		panic(err)
	}
	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile

	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}
}
