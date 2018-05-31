package persist

import (
	"log"
	"github.com/olivere/elastic"
	"context"
	"crawler/concurrent-version-crawler/engine"
	"github.com/kataras/iris/core/errors"
)

// 1. 连接 elasticsearch
// 2. goroutine 开启死循环， 接收 channel 传送过来的 Item,并存储
// 3. 返回接收 Item 的 channel
func ItemSaver(index string) (chan engine.Item, error) {

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver : got item # %d: %v", itemCount, item)
			itemCount++

			err := save(client, item, index)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()

	return out, nil
}

// 保存用户信息，返回存储成功的 ID
func save(client *elastic.Client, item engine.Item, index string) (err error) {

	if item.Type == "" {
		return errors.New("Must supply Type")
	}

	// 第一个 Index 指创建

	indexService := client.Index().
		Index(index).
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
