package redisCluster

import (
	"fmt"

	"strconv"

	"gopkg.in/redis.v5"
	"log"
	"strings"
)

func InsertData(address string, count int) error {

	//连接到redis分片
	addresses := strings.Split(address, ",")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addresses,
	})
	if client == nil {
		err := fmt.Errorf("InsertData: Create Client: Failed to create client")
		return err
	}
	defer client.Close()

	pong, err := client.Ping().Result()
	if err != nil {
		log.Println("InsertData: Cannot ping this client")
		return err
	}
	log.Printf("InsertData: Ping this client... [ %s ]\n", pong)

	err = client.Set("key"+strconv.Itoa(count), "value"+strconv.Itoa(count), 0).Err()
	if err != nil {
		log.Printf("InsertData: [ %d ] times \n", count)
		return err
	}

	log.Println("InsertData: Insert data key" + strconv.Itoa(count) + " completely")
	return nil
}

func FindKey(address string) {

	addresses := strings.Split(address, ",")

	for _, ad := range addresses {
		//创建客户端
		client := redis.NewClient(&redis.Options{
			Addr: ad,
			DB:   0,
		})
		if client == nil {
			log.Printf("FindKey [ %s ]: Create Client: Failed to create client\n", ad)
			continue
		}
		defer client.Close()

		pong, err := client.Ping().Result()
		if err != nil {
			log.Printf("FindKey [ %s ]: Cannot ping this client:%s\n", ad, err)
			continue
		}
		log.Printf("FindKey [ %s ]: Ping this client... [ %s ]\n", ad, pong)

		//查找keys
		results, err := client.Keys("*").Result()
		log.Println("FindKey: results: ", results)
		if err != nil {
			log.Printf("FindKey [ %s ]: Failed to run command keys * : %s\n", ad, err)
			continue
		}
		log.Printf("FindKey [ %s ]: find [ %d ] data successfully\n", ad, len(results))
	}
}
