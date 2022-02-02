package main

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/rs/cors"
	"net/http"
	"github.com/gorilla/mux"

)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err)
	}
}

func main() {
	go listenMessageBroker()
	router := mux.NewRouter()
	router.HandleFunc("/joiner", AddJoiner).Methods("POST")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func listenMessageBroker() {
	conn, err := amqp.Dial("amqps://admin:UliML6QEYZ6Jbz4Ji4f6kbnW4nxy2sgw@g1leyd.stackhero-network.com:5671")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	q, err := ch.QueueDeclare(
		"Profiles", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"", // consumer
		true, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// log.Printf("Received a message: %s", d.Body)
			insertEmployee(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL-C")
	<-forever
}