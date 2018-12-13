package nsq

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bitly/go-nsq"
	"github.com/maps90/nucleus/config"
)

var (
	nullLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	publisher  *Produce
)

// Config struct
type Config struct {
	Host    string
	LogPath string
}

// Publisher interface
type Publisher interface {
	send(string, interface{}) error
}

// Produce struct
type Produce struct {
	Instance *nsq.Producer
}

func (n *Produce) send(topic string, payload interface{}) error {
	payloadJSON, _ := json.Marshal(payload)

	err := n.Instance.Publish(topic, []byte(payloadJSON))
	if err != nil {
		log.Println("PublishNSQ: failed to publish", err.Error())
		return err
	}
	return nil
}

// Publish message to NSQ required topic and payload
func Publish(topic string, payload interface{}) error {
	if publisher == nil {
		connect()
	}

	return publisher.send(topic, payload)
}

// connect return a publisher
func connect() (Publisher, error) {
	if publisher != nil {
		return publisher, nil
	}
	configNSQ := nsq.NewConfig()
	client, err := nsq.NewProducer(config.GetString("nsq.publisher"), configNSQ)
	if err != nil {
		return nil, err
	}

	client.SetLogger(nullLogger, nsq.LogLevelInfo)
	publisher = &Produce{
		Instance: client,
	}
	return publisher, nil
}
