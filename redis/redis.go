package main

import (
	"fmt"

	"strconv"

	"time"

	"gopkg.in/redis.v5"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.100.103:30014",
		Password: "dangerous",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	i := 1
	for {
		err = client.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 0).Err()
		if err != nil {
			fmt.Println("insert failed")
		} else {
			fmt.Println("insert success")
		}
		i++
		time.Sleep(1 * time.Second)
	}
}
