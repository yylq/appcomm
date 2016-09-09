package common

import (
	"fmt"

	"github.com/yylq/log"

	"github.com/streadway/amqp"
)

type MQChannel struct {
	Ch   *amqp.Channel
	Name string
}

func GetIDChannel(id string) string {
	return fmt.Sprintf("mon_%s", id)
}
func MakeIDChannel(host string, id string) (*MQChannel, error) {
	if id == "" {
		return GetMQChannel(host, "")
	}
	return GetMQChannel(host, GetIDChannel(id))
}
func GetMQChannel(host string, name string) (*MQChannel, error) {
	log.Debugf("host:%s name:%s\n", host, name)
	conn, err := amqp.Dial(host)
	if err != nil {

		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {

		return nil, err
	}
	if name != "" {
		_, err = ch.QueueDeclare(
			name,  // name
			false, // durable
			false, // delete when usused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}
	}
	return &MQChannel{Ch: ch, Name: name}, nil
}
func (mch *MQChannel) MQSend(buf []byte) error {
	err := mch.Ch.Publish(
		"",       // exchange
		mch.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(buf),
		})
	return err
}
func (mch *MQChannel) MQSendBeforeDeclare(name string, buf []byte) error {
	_, err := mch.Ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = mch.Ch.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(buf),
		})
	return err
}
func (mch *MQChannel) MQClear() error {
	_, err := mch.Ch.QueuePurge(mch.Name, true)
	return err
}
func (mch *MQChannel) Close() error {
	err := mch.Ch.Close()
	return err
}
