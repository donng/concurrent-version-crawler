package persist

import (
	"log"
	"github.com/olivere/elastic"
	"context"
	"crawler/concurrent-version-crawler/engine"
	"github.com/kataras/iris/core/errors"
)

func ItemSaver() chan engine.Item {
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver : got item # %d: %v", itemCount, item)
			itemCount++

			err := save(item)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()

	return out
}

// 保存用户信息，返回存储成功的 ID
func save(item engine.Item) (err error) {
	client, e := elastic.NewClient(
		// must turn off sniff in docker
		elastic.SetSniff(false),
	)

	if e != nil {
		panic(e)
	}

	if item.Type == "" {
		return errors.New("Must supply Type")
	}

	// 第一个 Index 指创建

	indexService := client.Index().
		Index("dating_profile").
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err = indexService.Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}
