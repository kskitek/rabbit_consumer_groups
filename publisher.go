package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(string) error
	PublishHTTP(http.ResponseWriter, *http.Request)
}

type channelPublisher struct {
	name string
	ch   *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) (Publisher, error) {
	topic := os.Getenv("RABBIT_TOPIC_NAME")
	if topic == "" {
		return nil, fmt.Errorf("`RABBIT_TOPIC_NAME` cannot be empty")
	}

	err := declareExchange(topic, ch)
	if err != nil {
		return nil, err
	}

	return &channelPublisher{name: topic, ch: ch}, nil
}

func declareExchange(topic string, ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		topic,   // name
		"topic", // type
		false,   // durable
		true,    // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
}

func (p *channelPublisher) Publish(payload string) error {
	logrus.WithField("payload", payload).Debug("publishing")
	pub := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(payload),
	}
	return p.ch.Publish(
		p.name, // exchange
		p.name, // routing key
		false,  // mandatory
		false,  // immediate
		pub)
}

func (p *channelPublisher) PublishHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
