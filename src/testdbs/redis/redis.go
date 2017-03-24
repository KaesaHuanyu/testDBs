package redis

import (
	"fmt"

	"strconv"

	"gopkg.in/redis.v5"
	"log"
)

func InsertData(address, password string) error {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
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

	for i := 1; i <= 1000; i++ {
		err := client.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 0).Err()
		if err != nil {
			log.Printf("InsertData: [ %d ] times : %s\n", i, err)
			continue
		}
	}

	log.Println("InsertData: Insert datas completely")
	return nil
}

func FindKey(address, password string) error {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	if client == nil {
		err := fmt.Errorf("FindKey: Create Client: Failed to create client")
		return err
	}
	defer client.Close()

	pong, err := client.Ping().Result()
	if err != nil {
		log.Println("FindKey: Cannot ping this client")
		return err
	}
	log.Printf("FindKey: Ping this client... [ %s ]\n", pong)

	results, err := client.Keys("*").Result()
	log.Println("FindKey: results: ", results)
	if err != nil {
		log.Println("FindKey: Failed to run command keys *")
		return err
	}
	log.Printf("FindKey: find [ %d ] data successfully\n", len(results))
	return nil
}
