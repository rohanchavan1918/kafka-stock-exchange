package conf

import (
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Host  string `viper:"string" validate:"required" mapstructure:"kafka_host"`
	Port  int64  `viper:"string" validate:"required" mapstructure:"kafka_port"`
	Topic string `viper:"string" validate:"required" mapstructure:"topic"`
}

func (c *KafkaConfig) getKafkaHost() string {
	// Get a kafka producer from the config
	if c.Host == "" || c.Port == 0 || c.Topic == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *KafkaConfig) GetProducer() (*kafka.Writer, error) {
	// Get a kafka producer from the config

	kafkaHost := c.getKafkaHost()
	if kafkaHost == "" {
		err := errors.New("Kafka host, port or topic cannot be empty")
		return nil, err
	}
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaHost},
		Topic:   c.Topic,
	})
	return w, nil
}

func (c *KafkaConfig) GetConsumer() (*kafka.Reader, error) {

	kafkaHost := c.getKafkaHost()
	if kafkaHost == "" {
		err := errors.New("Kafka host, port or topic cannot be empty")
		return nil, err
	}

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaHost},
		Topic:   c.Topic,
	})

	if consumer != nil {
		return consumer, nil
	} else {
		err := errors.New("Failed to create kafka consumer")
		return nil, err
	}

}
