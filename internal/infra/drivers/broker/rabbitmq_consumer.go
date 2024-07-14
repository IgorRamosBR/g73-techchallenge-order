package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type rabbitConsumer struct {
	channel    *amqp.Channel
	messagesCh <-chan amqp.Delivery
	queueName  string
}

func NewRabbitMQConsumer(channel *amqp.Channel, queueName string) (Consumer, error) {
	messagesCh, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack, set to false for manual ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer for the queue [%s], error: [%w]", queueName, err)
	}

	return &rabbitConsumer{channel: channel, queueName: queueName, messagesCh: messagesCh}, nil
}

func (c *rabbitConsumer) StartConsumer(processMessage func(message []byte) error) {
	log.Infof("Starting consuming queue [%s]", c.queueName)
	forever := make(chan bool)
	go func() {
		for msg := range c.messagesCh {
			log.Debugf("Received a message: %s", msg.Body)

			err := processMessage(msg.Body)
			if err != nil {
				log.Errorf("failed to process message, error: %s", err.Error())
				msg.Nack(false, true)
				continue
			}

			// Acknowledge the message after successful processing
			msg.Ack(false)
		}
	}()

	<-forever
	log.Errorf("queue [%s] consumer stopped working", c.queueName)
}
