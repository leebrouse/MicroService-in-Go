package event

import amqp "github.com/rabbitmq/amqp091-go"

// rabbitMQ Exchange config
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //exchange name
		"topic",      //type
		true,         //durable
		false,        //autoDelete
		false,        //interal
		false,        //nowait?
		nil,          //extra arguments
	)
}

// rabbitMQ queue config
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
