package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var Version string
var BuildTime string

type mqResources struct {
	queueName                     string
	messageBody                   string
	expectedMessageCountPublished int
	expectedMessageCountConsumed  int
}

type settings struct {
	Password string
	Username string
	Port     int
	Server   string
}

func defaultConnStr(s *settings) string {
	return "amqp://" + s.Username + ":" + s.Password + "@" + s.Server + ":" + string(s.Port) + "/"
}

func RabbitMQAMQPConnection(s *settings) (*amqp.Connection, error) {
	conn, err := amqp.Dial(defaultConnStr(s))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to rabbitMQ %s", err)
	}
	return conn, nil
}

func main() {
	var s settings
	err := envconfig.Process("secret_demo", &s)
	if err != nil {
		log.Fatal("failed to initialize settings ", err.Error())
	}
	log.Debugf("You are running buildversion %v: compiled at: %v ", Version, BuildTime)
	qn := &mqResources{
		queueName:                     "MQTestQueue",
		messageBody:                   "This is a test message",
		expectedMessageCountPublished: 1,
		expectedMessageCountConsumed:  0,
	}

	conn, err := RabbitMQAMQPConnection(&s)
	if err != nil {
		log.Fatalf("Could not create RabbitMQ connection %v", err)
	}

	if err := RabbitMQCreateQueue(conn, qn); err != nil {
		log.Fatalf("Error connect to RabbitMQ and create a queue %v", err)
	}

	if err := RabbitMQPublishMessage(conn, qn); err != nil {
		log.Fatalf("Error publish message to RabbitMQ %v", err)
	}

	if err := RabbitMQConsumeMessage(conn, qn); err != nil {
		log.Fatalf("Error consuming message from RabbitMQ %v", err)
	}

}

func RabbitMQCreateQueue(conn *amqp.Connection, m *mqResources) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("could not open channel %v", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(m.queueName, false, false, false, false, nil)
	if err != nil {
		log.Errorf("failed to declare RabbitMQ queue %v", err)
	}
	if m.queueName != q.Name {
		return fmt.Errorf("expected queue name (%v) got (%v)", m.queueName, q.Name)
	}
	fmt.Printf("[OK] Create message queue %v\n", q.Name)
	return nil
}

func RabbitMQPublishMessage(conn *amqp.Connection, m *mqResources) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("could not open channel %s", err)
	}
	defer ch.Close()
	if err = ch.Publish(
		"",
		m.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(m.messageBody),
		}); err != nil {
		log.Errorf("failed to publish message %s", err)
	}
	if err := ch.Confirm(false); err != nil {
		log.Errorf("message published could not be confirmed %s", err)
	}
	result, err := ch.QueueInspect(m.queueName)
	if err != nil {
		return err
	}
	if result.Messages != m.expectedMessageCountPublished {
		return fmt.Errorf("expected messagecount %v got messageCount %v", m.expectedMessageCountPublished, result.Messages)
	}
	fmt.Printf("[OK] Published message. Message count: %v\n", result.Messages)
	return nil

}

func RabbitMQConsumeMessage(conn *amqp.Connection, m *mqResources) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Could not open channel %s", err)
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		m.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("failed to register a consumer %s", err)
	}
	for d := range msgs {
		if string(d.Body) != m.messageBody {
			log.Errorf("Expected message body %v got %v", d.Body, m.messageBody)
		}
		if err := d.Ack(false); err != nil {
			log.Errorf("%v", err)
		}
	}
	result, _ := ch.QueueInspect(m.queueName)
	if result.Messages != m.expectedMessageCountConsumed {
		return fmt.Errorf("expected messagecount %v got messageCount %v", m.expectedMessageCountConsumed, result.Messages)
	}
	fmt.Printf("[OK] Consumed message. Message count: %v\n", result.Messages)
	return nil
}
