package contracts

import (
	"database/sql"
	"github.com/Shopify/sarama"
)

type PgDriver interface {
	Query(qString string, fields ...interface{}) (*sql.Rows, error)
}

type TransportClient interface {
	InitConsumer(topic string) (sarama.PartitionConsumer, error)
	InitProducer() (sarama.AsyncProducer, error)
}
