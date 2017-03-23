package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"strings"
	"testDBs/mongo"
	"time"
	"testDBs/rabbitmq"
	"testDBs/redis"
	"testDBs/mysql"
	"testDBs/redisCluster"
	"log"
	"testDBs/mongoCluster"
)

func main() {
	app := cli.NewApp()
	app.Name = "testDBs"
	app.Usage = "testDBs"
	app.Version = "1.0"
	app.Author = "Kaesa Li"
	app.Email = "kaesa.li@daocloud.io"

	app.Commands = []cli.Command{
		{
			Name: "mongo",
			Usage: "Input the mongo args",
			Action: Mongo,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the mongo service address",
					EnvVar: "MONGO_ADDRESS",
				},
				cli.StringFlag{
					Name: "u",
					Usage: "the mysql service username",
					EnvVar: "MYSQL_USERNAME",
				},
				cli.StringFlag{
					Name: "p",
					Usage: "the mysql service password",
					EnvVar: "MYSQL_PASSWORD",
				},
				cli.StringFlag{
					Name: "db",
					Usage: "the mysql service database",
					EnvVar: "MYSQL_DATABASE",
				},
			},
		},

		{
			Name: "mysql",
			Usage: "Input the mysql args",
			Action: Mysql,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "master",
					Usage: "the mysql service master address",
					EnvVar: "MYSQL_MASTER_ADDRESS",
				},
				cli.StringFlag{
					Name: "slave",
					Usage: "the mysql service slave address",
					EnvVar: "MYSQL_SLAVE_ADDRESS",
				},
				cli.StringFlag{
					Name: "u",
					Usage: "the mysql service username",
					EnvVar: "MYSQL_USERNAME",
				},
				cli.StringFlag{
					Name: "p",
					Usage: "the mysql service password",
					EnvVar: "MYSQL_PASSWORD",
				},
				cli.StringFlag{
					Name: "db",
					Usage: "the mysql service database",
					EnvVar: "MYSQL_DATABASE",
				},
				cli.StringFlag{
					Name: "tb",
					Usage: "the mysql service table",
					EnvVar: "MYSQL_TABLE",
				},
			},
		},

		{
			Name: "rabbitmq",
			Usage: "Input the rabbitmq args",
			Action: Rabbitmq,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the rabbitmq service address",
					EnvVar: "RABBITMQ_ADDRESS",
				},
				cli.StringFlag{
					Name: "p",
					Usage: "the rabbitmq service password",
					EnvVar: "MYSQL_PASSWORD",
				},
			},
		},

		{
			Name: "redis",
			Usage: "Input the redis-ha args",
			Action: Redis,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the redis-ha service address",
					EnvVar: "REDIS_ADDRESS",
				},
				cli.StringFlag{
					Name: "p",
					Usage: "the redis password",
					EnvVar: "REDIS_PASSWORD",
				},
			},
		},

		{
			Name: "redisCluster",
			Usage: "Input the redis-cluster args",
			Action: RedisCluster,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad1",
					Usage: "the redis-cluster service master1 & slave1 address",
					EnvVar: "REDIS_CLUSTER_ADDRESS_1",
				},
				cli.StringFlag{
					Name: "ad2",
					Usage: "the redis-cluster service master2 & slave2 address",
					EnvVar: "REDIS_CLUSTER_ADDRESS_2",
				},
				cli.StringFlag{
					Name: "ad3",
					Usage: "the redis-cluster service master3 & slave3 address",
					EnvVar: "REDIS_CLUSTER_ADDRESS_3",
				},
			},
		},

		{
			Name: "mongoCluster",
			Usage: "Input the mongo-cluster args",
			Action: MongoCluster,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the mongo-cluster service mongos address",
					EnvVar: "MONGO_CLUSTER_ADDRESS",
				},
				//cli.StringFlag{
				//	Name: "db",
				//	Usage: "the mongo-cluster service database",
				//	EnvVar: "MONGO_CLUSTER_DATABASE",
				//},
				//cli.StringFlag{
				//	Name: "co",
				//	Usage: "the mongo-cluster service collection",
				//	EnvVar: "MONGO_CLUSTER_COLLECTION",
				//},
			},
		},
	}

	app.Run(os.Args)
}

//test mongodb
func Mongo(c *cli.Context) {

	address := mustGetStringVar(c, "ad")
	username := mustGetStringVar(c, "u")
	password := mustGetStringVar(c, "p")
	database := mustGetStringVar(c, "db")

	//insert data
	err := mongo.InsertData(address, username, password, database)
	if err != nil {
		log.Printf("InsertData: Fail to insert data: %s\n", err)
	}

	for i := 1; i <= 30; i++{

		//Find Primary Point
		err := mongo.FindPrimary(address, username, password, database)
		if err != nil {
			log.Printf("[ %d ] times : %s", i, err)
			continue
		}

		//Find Datas' number
		err = mongo.FindData(address, username, password, database)
		if err != nil {
			log.Printf("[ %d ] times : %s", i, err)
		}

		fmt.Println("Sleeping 10 seconds...now you can down / up some points")
		time.Sleep(10 * time.Second)
	}

	log.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test mysql
func Mysql(c *cli.Context) {

	masteraddress := mustGetStringVar(c, "master")
	slaveaddress := mustGetStringVar(c, "slave")
	username := mustGetStringVar(c, "u")
	password := mustGetStringVar(c, "p")
	database := mustGetStringVar(c, "db")
	table := mustGetStringVar(c, "tb")

	var index int
	for i := 1; i <= 30; i++ {

		//insert data
		err := mysql.InsertData(masteraddress, username, password, database, table)
		if err != nil {
			log.Printf("InsertData: Failed to insert data: %s\n", err)
			continue
		}

		index++

		//Find master point datas
		log.Println("Find data in master point...")
		err = mysql.FindData(masteraddress, username, password, database, table, index)
		if err != nil {
			log.Printf("[ %d ] times : %s...please wait for next time", i, err)
		}

		//Find slave point datas
		log.Println("Find data in slave point...")
		err = mysql.FindData(slaveaddress, username, password, database, table, index)
		if err != nil {
			log.Printf("[ %d ] times : %s...please wait for next time", i, err)
		}

		log.Printf("Sleeping 10 seconds...you can down / up some points\n")
		time.Sleep(10 * time.Second)
	}

	log.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test rabbitmq
func Rabbitmq(c *cli.Context) {
	address := mustGetStringVar(c, "ad")
	password := mustGetStringVar(c, "p")

	for i := 1; i <= 30; i++ {
		stop := make(chan bool)

		go rabbitmq.ProducerGo(address, password, stop)
		rabbitmq.ConsumerGo(address, password, stop)

		log.Printf("Sleeping 10 seconds...you can down / up some points\n")
		time.Sleep(10 * time.Second)
	}

	log.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test redis_ha
func Redis(c *cli.Context) {
	address := mustGetStringVar(c, "ad")
	password := mustGetStringVar(c, "p")

	//insert data
	err := redis.InsertData(address, password)
	if err != nil {
		log.Printf("InsertData: Failed to insert data: %s\n", err)
	}

	for i := 1; i <= 30; i++ {

		//Find Keys
		err := redis.FindKey(address, password)
		if err != nil {
			log.Printf("[ %d ] times : %s...please wait for next time", i, err)
		}

		log.Printf("Sleeping 10 seconds...you can down / up some points\n")
		time.Sleep(10 * time.Second)
	}

	log.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test redis_cluster
func RedisCluster(c *cli.Context) {
	address1 := mustGetStringVar(c, "ad1")
	address2 := mustGetStringVar(c, "ad2")
	address3 := mustGetStringVar(c, "ad3")

	for i := 1; i <= 30; i++ {

		//Insert Kay-Value
		err := redisCluster.InsertData(address1 + "," + address2 + "," + address3, i)
		if err != nil {
			log.Printf("InsertData: Failed to insert data: %s\n", err)
		}

		//Find cluster 1 keys
		log.Println("Find keys in cluster 1...")
		redisCluster.FindKey(address1)

		//Find cluster 2 keys
		log.Println("Find keys in cluster 2...")
		redisCluster.FindKey(address2)

		//Find cluster 1 keys
		log.Println("Find keys in cluster 3...")
		redisCluster.FindKey(address3)

		log.Printf("Sleeping 10 seconds...you can down / up some points\n")
		time.Sleep(10 * time.Second)
	}

	log.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test mongo_cluster
func MongoCluster(c *cli.Context) {
	address := mustGetStringVar(c, "ad")

	for count := 1; count <= 30; count ++ {

		//Insert datas
		err := mongoCluster.InsertData(address, count)
		if err != nil {
			log.Printf("InsertData: Failed to insert data: %s\n", err)
		}

		//Find shard1

	}
}

func errExit(code int, format string, val ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", val...)
	os.Exit(code)
}

func mustGetStringVar(c *cli.Context, key string) string {
	v := strings.TrimSpace(c.String(key))
	if v == "" {
		errExit(1, "%s must be provided", key)
	}
	return v
}