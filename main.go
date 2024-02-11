package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	queueCon := setupRabbitMQConnection()
	initQueue(queueCon)

	router := chi.NewRouter()

	router.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		if err := produceUpdate(queueCon, body); err != nil {
			println(err.Error())
		}

		render.Status(r, http.StatusOK)
	})

	log.Print("listen")
	if err := http.ListenAndServe("0.0.0.0:8443", router); err != nil {
		log.Fatal(err)
	}
}

func produceUpdate(queueCon *amqp.Connection, update json.RawMessage) error {
	ch, err := queueCon.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) { _ = ch.Close() }(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"",
		"updates",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return err
	}
	return nil
}

func initQueue(con *amqp.Connection) {
	ch, err := con.Channel()
	if err != nil {
		panic(err)
	}
	defer func(ch *amqp.Channel) { _ = ch.Close() }(ch)

	_, err = ch.QueueDeclare(
		"updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
}

func setupRabbitMQConnection() *amqp.Connection {
	con, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		println(err.Error())
		panic(err)
	}
	return con
}
