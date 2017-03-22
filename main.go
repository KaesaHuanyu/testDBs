package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"strings"
	"testDBs/mongo"
	"time"
	"testDBs/rabbitmq"
	"log"
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
					Name: "ad",
					Usage: "the mysql service address",
					EnvVar: "MYSQL_ADDRESS",
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
			Name: "redis-ha",
			Usage: "Input the redis-ha args",
			Action: Redis_ha,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the redis-ha service address",
					EnvVar: "REDIS_HA_ADDRESS",
				},
			},
		},

		{
			Name: "redis-cluster",
			Usage: "Input the redis-cluster args",
			Action: Redis_cluster,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the redis-cluster service address",
					EnvVar: "REDIS_CLUSTER_ADDRESS",
				},
			},
		},

		{
			Name: "mongo-cluster",
			Usage: "Input the mongo-cluster args",
			Action: Mongo_cluster,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "ad",
					Usage: "the mongo-cluster service address",
					EnvVar: "MONGO_CLUSTER_ADDRESS",
				},
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

	//add data
	err := mongo.AddData(address, username, password, database)
	if err != nil {
		log.Printf("Fail to add data: %s\n", err)
	}

	for i := 0; i < 30; i++{
		err := mongo.FindPrimary(address, username, password, database)
		if err != nil {
			log.Printf("[ %d ] times : %s", i, err)
			continue
		}
		err = mongo.FindData(address, username, password, database)
		if err != nil {
			log.Printf("[ %d ] times : %s", i, err)
		}
		fmt.Println("Sleeping 10 seconds...you can down / up some point")
		time.Sleep(10 * time.Second)
	}
	fmt.Printf("\n\ntime out!...if you want to test more, please restart this program\n")

}

//test mysql
func Mysql(c *cli.Context) {
	//address := mustGetStringVar(c, "ad")
	//username := mustGetStringVar(c, "u")
	//password := mustGetStringVar(c, "p")
	//database := mustGetStringVar(c, "db")
	//err := mysql.initMySQL()
}

//test rabbitmq
func Rabbitmq(c *cli.Context) {
	address := mustGetStringVar(c, "ad")
	password := mustGetStringVar(c, "p")

	for i := 0; i < 30; i++ {
		stop := make(chan bool)

		go rabbitmq.ProducerGo(address, password, stop)
		rabbitmq.ConsumerGo(address, password, stop)

		log.Printf("Sleeping 10 seconds...you can down / up some point\n")
		time.Sleep(10 * time.Second)
	}
	fmt.Printf("\n\ntime out!...if you want to test more, please restart this program\n")
}

//test redis_ha
func Redis_ha(c *cli.Context) {
	//address := mustGetStringVar(c, "ad")
	//addresses := strings.Split(address, ",")
	//for _ = range addresses {
	//
	//}
}

//test redis_cluster
func Redis_cluster(c *cli.Context) {

}

//test mongo_cluster
func Mongo_cluster(c *cli.Context) {

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