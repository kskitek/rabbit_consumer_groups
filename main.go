package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/avast/retry-go"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var conn *amqp.Connection
	err := retry.Do(connFunc(conn), retry.Attempts(5), retry.Delay(2), retry.Units(time.Second))
	defer closeConn(conn)
	if err != nil {
		logrus.WithError(err).Error("when connecting to queue")
		return
	}
	logrus.Info("Started")

	ch, err := declareChannel(conn)
	defer ch.Close()

	p, err := NewPublisher(ch)
	if err != nil {
		logrus.WithError(err).Error("when creating publisher")
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/publish/{payload}", p.PublishHTTP) // not clean but simple :)
	// r.HandleFunc("/addServiceQueue/{name}", nil)

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

func declareChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		logrus.WithError(err).Error("cannot create channel")
		return nil, err
	}
	return ch, nil
}

func closeConn(conn *amqp.Connection) {
	if conn != nil {
		conn.Close()
	}
}

func connFunc(conn *amqp.Connection) func() error {
	return func() error {
		c, err := amqp.Dial(getRabbitURL())
		conn = c
		return err
	}
}

func getRabbitURL() string {
	host := os.Getenv("RABBIT_HOST")
	port := os.Getenv("RABBIT_PORT")
	user := os.Getenv("RABBIT_USER")
	password := os.Getenv("RABBIT_PASSWORD")
	urlPattern := "amqp://%s:%s@%s:%s"

	logrus.WithFields(logrus.Fields{"host": host, "port": port, "user": user}).Info("Connecting to RabbitMQ using")

	return fmt.Sprintf(urlPattern, user, password, host, port)
}
