package messages

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	callbacks     ProducerCallbacks
	topic         string
}

func NewProducer(topic string, brokerlist []string) *Producer {
	callbacks := ProducerCallbacks{
		OnError: onProducerError,
	}
	producer := Producer{callbacks: callbacks, topic: topic}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	saramaProducer, err := sarama.NewAsyncProducer(brokerlist, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama Producer:", err)
		panic(err)
	}
	go func() {
		for err := range saramaProducer.Errors() {
			if producer.callbacks.OnError != nil {
				producer.callbacks.OnError(err)
			}
		}
	}()
	producer.asyncProducer = saramaProducer
	return &producer
}

func (p *Producer) Produce(payload []byte) {
	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(payload),
	}
}

func onProducerError(err error) {
	log.Println("onProducerError: ", err)
}

func (p *Producer) Close() error {
	log.Println("Producer.Close()")
	if err := p.asyncProducer.Close(); err != nil {
		return err
	}
	return nil
}

type ProducerCallbacks struct {
	OnError func(error)
}
