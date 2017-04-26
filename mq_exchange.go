package appcomm

import (
	"github.com/yylq/log"

	"errors"

	"github.com/streadway/amqp"
)

type MQExchange struct {
	Host     string
	Ch       *amqp.Channel
	Name     string
	Key      string
	IsClosed bool
}

var (
	err_closed = errors.New("channel is closed")
)

func MakeExchange(host, name, key string) (*MQExchange, error) {

	return GetMQExchange(host, name, key)
}

func GetMQExchange(host, name, key string) (*MQExchange, error) {
	log.Debugf("host:%s name:%s key:%s \n", host, name, key)
	ex := &MQExchange{Host: host, Name: name, Key: key}
	err := ex.Open()
	if err != nil {
		return nil, err
	}
	return ex, nil
}
func (this *MQExchange) Open() error {
	conn, err := amqp.Dial(this.Host)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()

	if err != nil {

		return err
	}
	this.Ch = ch
	return nil
}
func (this *MQExchange) TryConnect() bool {
	conn, err := amqp.Dial(this.Host)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
func (this *MQExchange) ReOpen() error {
	if !this.IsClosed {
		this.Close()
	}
	return this.Open()
}
func (mch *MQExchange) MQSend(buf []byte) error {
	if mch.IsClosed {
		return err_closed
	}
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
	if mch.IsClosed {
		return err_closed
	}
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
	if mch.IsClosed {
		return err_closed
	}
	err := mch.Ch.Close()
	mch.IsClosed = true
	return err
}
