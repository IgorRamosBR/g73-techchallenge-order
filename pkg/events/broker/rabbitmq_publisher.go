package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQPublisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQPublisher(channel *amqp.Channel, exchange string) Publisher {
	return &rabbitMQPublisher{channel: channel, exchange: exchange}
}

func (c *rabbitMQPublisher) Publish(ctx context.Context, destination string, message []byte) error {
	err := c.channel.Publish(
		c.exchange,  //exchange,
		destination, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *rabbitMQPublisher) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}

	return nil
}
