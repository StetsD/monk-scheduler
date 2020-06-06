package infrastructure

import (
	"github.com/Shopify/sarama"
	config "github.com/stetsd/monk-conf"
	"log"
)

type KafkaClient struct {
	config config.Config
}

func NewKafkaClient(conf config.Config) *KafkaClient {
	return &KafkaClient{config: conf}
}

func (kc *KafkaClient) InitConsumer(topic string) (sarama.PartitionConsumer, error) {
	configSarama := sarama.NewConfig()
	configSarama.Producer.Return.Successes = true
	configSarama.Net.SASL.Password = kc.config.Get(config.TransportPass)
	configSarama.Net.SASL.User = kc.config.Get(config.TransportUser)
	consumer, err := sarama.NewConsumer([]string{kc.config.Get(config.TransportHost) + ":" + kc.config.Get(config.TransportPort)}, configSarama)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)

	if err != nil {
		return nil, err
	}

	return partitionConsumer, nil
}

func (kc *KafkaClient) InitProducer() (sarama.AsyncProducer, error) {
	configSarama := sarama.NewConfig()
	configSarama.Producer.Return.Successes = true
	configSarama.Net.SASL.Password = kc.config.Get(config.TransportPass)
	configSarama.Net.SASL.User = kc.config.Get(config.TransportUser)
	producer, err := sarama.NewAsyncProducer([]string{kc.config.Get(config.TransportHost) + ":" + kc.config.Get(config.TransportPort)}, configSarama)

	if err != nil {
		return nil, err
	}

	return producer, nil
}
