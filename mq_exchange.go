package appcomm

import (
	"github.com/yylq/log"

	"github.com/streadway/amqp"
)

type MQExchange struct {
	Ch   *amqp.Channel
	Name string
	Key  string
}

func MakeExchange(host, name, key string) (*MQExchange, error) {

	return GetMQExchange(host, name, key)
}

func GetMQExchange(host, name, key string) (*MQExchange, error) {
	log.Debugf("host:%s name:%s key:%s \n", host, name, key)
	conn, err := amqp.Dial(host)
	if err != nil {

		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {

		return nil, err
	}
	return &MQExchange{Ch: ch, Name: name, Key: key}, nil
}
func (mch *MQExchange) MQSend(buf []byte) error {
	err := mch.Ch.Publish(
		mch.Name, // exchange
		mch.Key,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(buf),
		})
	return err
}

func (mch *MQExchange) MQSendWith(buf []byte, name, key string) error {
	err := mch.Ch.Publish(
		name,  // exchange
		key,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(buf),
		})
	return err
}
func (mch *MQExchange) Close() error {
	err := mch.Ch.Close()
	return err
}
