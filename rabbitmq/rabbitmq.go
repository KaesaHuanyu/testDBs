package rabbitmq

import (
	"strconv"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func ProducerGo(address, password string, stop chan bool) {

	link := "amqp://admin:" + password + "@" + address + "/"
	conn, err := amqp.Dial(link)
	failOnError(err, "Producer: Failed to connect to RabbitMQ")
	defer conn.Close()
	if conn == nil {
		log.Printf("Producer: Dial: conn is nil, please wait next time\n")
		time.Sleep(2 * time.Second)
		stop <- true
		return
	} else if conn.Properties == nil {
		log.Printf("Producer: Dial: conn.Properties is nil, please wait next time\n")
		time.Sleep(2 * time.Second)
		stop <- true
		return
	}


	ch, err := conn.Channel()
	failOnError(err, "Producer: Failed to open a channel")
	defer ch.Close()
	if ch == nil {
		log.Printf("Producer: Channel: ch is nil, please wait next time\n")
		time.Sleep(2 * time.Second)
		stop <- true
		return
	}

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Producer: Failed to declare a queue")

	for i := 1; i <= 10; i++{
		body := "hello" + strconv.Itoa(i)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Printf("Producer: Failed to Publish a message")
		} else {
			log.Printf("Producer: Sent hello" + strconv.Itoa(i) + "\n")
		}
	}

	time.Sleep(2 * time.Second)
	stop <- true
	return
}

func ConsumerGo(address, password string, stop chan bool) {

	link := "amqp://admin:" + password + "@" + address + "/"
	conn, err := amqp.Dial(link)
	failOnError(err, "Consumer: Failed to connect to RabbitMQ")
	defer conn.Close()

	//检查指针是否为空
	if conn == nil {
		log.Printf("Producer: Dial: conn is nil, please wait next time\n")
		return
	} else if conn.Properties == nil {
		log.Printf("Producer: Dial: conn.Properties is nil, please wait next time\n")
		return
	}

	ch, err := conn.Channel()
	failOnError(err, "Consumer: Failed to open a channel")
	defer ch.Close()

	if ch == nil {
		log.Printf("Producer: Channel: ch is nil, please wait next time\n")
		return
	}

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Consumer: Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Consumer: Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Consumer: Received a message: %s\n", d.Body)
		}
	}()

	<-stop
	return
}
