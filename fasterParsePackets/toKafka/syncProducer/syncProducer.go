package syncProducer

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func NewProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.ClientID = "newsDataSource"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Printf("producer_test create producer error :%s\n", err.Error())
		return nil, err
	}
	return producer, err
}
