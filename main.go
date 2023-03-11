package main

import (
	rabbitMq "micron/RabbitMQ"
	sql "micron/sql"
)

func main() {

	var db = sql.SqlOpen()
	defer db.Close()

	sql.GetGames(db)

	rabbitMq.Consumer(db)
}
