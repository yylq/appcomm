package appcomm

import (
	"fmt"

	"github.com/yylq/log"

	"github.com/streadway/amqp"
)

type MQChannel struct {
	Ch       *amqp.Channel
	Name     string
	Host     string
	IsClosed bool
}

func GetIDChannel(id string) string {
	return fmt.Sprintf("mon_%s", id)
}
func MakeIDChannel(host string, id string) (*MQChannel, error) {
	if id == "" {
		return GetMQChannel(host, "", false)
	}
	return GetMQChannel(host, GetIDChannel(id), false)
}
func GetMQChannel(host string, name string, declare bool) (*MQChannel, error) {
	log.Debugf("host:%s name:%s\n", host, name)
	Channel := &MQChannel{Host: host, Name: name}
	err := Channel.Open()
	if err != nil {
		return nil, err
	}
	if declare {
		ch := Channel.Ch
		_, err = ch.QueueDeclare(
			name,  // name
			false, // durable
			false, // delete when usused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			Channel.Close()
			return nil, err
		}
	}
	return Channel, nil
}
func (this *MQChannel) Open() error {
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
func (this *MQChannel) TryConnect() bool {
	conn, err := amqp.Dial(this.Host)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
func (this *MQChannel) ReOpen() error {
	if !this.IsClosed {
		this.Close()
	}
	return this.Open()
}
func (this *MQChannel) MQSend(buf []byte) error {
	if this.IsClosed {
		return err_closed
	}
	err := this.Ch.Publish(
		"",        // exchange
		this.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(buf),
		})
	return err
}
func (this *MQChannel) MQSendWith(buf []byte, name string) error {
	if this.IsClosed {
		return err_closed
	}
	err := this.Ch.Publish(
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
func (this *MQChannel) GetMsgHander() (<-chan amqp.Delivery, error) {
	if this.IsClosed {
		return nil, err_closed
	}
	ch := this.Ch
	return ch.Consume(
		this.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

func (this *MQChannel) MQClear() error {
	if this.IsClosed {
		return err_closed
	}
	_, err := this.Ch.QueuePurge(this.Name, true)
	return err
}
func (this *MQChannel) MQClearWith(name string) error {
	if this.IsClosed {
		return err_closed
	}
	_, err := this.Ch.QueuePurge(name, true)
	return err
}
func (this *MQChannel) Close() error {
	if this.IsClosed {
		return err_closed
	}
	err := this.Ch.Close()
	this.IsClosed = true
	return err
}
