package rabbitmq

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	shared "micron/shared"
	sql2 "micron/sql"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consumer(db *sql.DB) {
	conn, err := amqp.Dial(shared.Config.AMQPURL)
	handleError(err, "Can not connect to AMQP")
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("game", false, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer read, PID: %d", os.Getegid())
		for d := range messageChannel {
			fmt.Println(strings.Repeat("-", 100))
			log.Printf("Received a message: %s", d.Body)

			addGame := &shared.AddGame{}

			err := json.Unmarshal(d.Body, addGame)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			fmt.Println(strings.Repeat("-", 100))
			fmt.Printf("Game :%s - Publisher: %s\n", addGame.Name, addGame.Publisher)

			res, err2 := sql2.IntersSqlContent(db, addGame)
			handleError(err2, "Could not Insert Game to Sql")
			log.Printf("Inserted Game ID : %d", res)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

			sql2.GetGames(db)
		}
	}()

	// Stop for program termination
	<-stopChan
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
