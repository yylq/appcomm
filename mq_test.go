package appcomm

import (
	"testing"
)

var (
	rabbit_host string = "amqp://root:root@192.168.176.7:5672/"
)

func TestChannel(t *testing.T) {
	ch, err := GetMQChannel(rabbit_host, "que_hello_1", false)
	if err != nil {
		t.Fatal(err)
	}
	defer ch.Close()
	err = ch.MQClear()
	if err != nil {
		t.Fatal(err)
	}
	err = ch.MQSend([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
}
func TestChannelWith(t *testing.T) {
	ch, err := GetMQChannel(rabbit_host, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer ch.Close()
	names := []string{"que_hello_1", "que_hello_2", "que_hello_3", "que_hello_4"}
	for _, name := range names {
		err = ch.MQClearWith(name)
		if err != nil {
			t.Fatal(err)
		}
		err = ch.MQSendWith([]byte("hello world"), name)
		if err != nil {
			t.Fatal(err)
		}
	}

}
func testExchange(t *testing.T) {
	ch, err := MakeExchange(rabbit_host, "ex_test_hello", "")
	if err != nil {
		t.Fatal(err)
	}
	defer ch.Close()
	err = ch.MQSend([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	err = ch.MQSendWith([]byte("hello world"), "ex_test_hello_2", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ch.MQSendWith([]byte("hello world"), "ex_test_hello", "hell_1")
	if err != nil {
		t.Fatal(err)
	}
}
