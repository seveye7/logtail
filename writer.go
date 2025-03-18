package main

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type SinkWriter interface {
	Write(string) (int64, int64, error)
}

type ConsoleWriter struct {
	index int
}

func (c *ConsoleWriter) Write(s string) (int64, int64, error) {
	if s == "" {
		return 0, 0, nil
	}
	textLogger.Info("---------->>>", "text", s)
	return 1, int64(len(s)), nil
}

type KafkaWriter struct {
	Config *Config
	w      *kafka.Writer
	msgs   []kafka.Message
}

func (k *KafkaWriter) Write(s string) (int64, int64, error) {
	batchMax := k.Config.Out.Kafka.BatchMax
	if batchMax == 0 {
		batchMax = 100
	}

	if s != "" {
		k.msgs = append(k.msgs, kafka.Message{Value: []byte(s)})

		if int64(len(k.msgs)) < batchMax {
			return 0, 0, nil
		}
	}

	if len(k.msgs) == 0 {
		return 0, 0, nil
	}

	var msgs []kafka.Message
	if int64(len(k.msgs)) > batchMax {
		msgs = k.msgs[:batchMax]
		k.msgs = k.msgs[batchMax:]
	} else {
		msgs = k.msgs
		k.msgs = []kafka.Message{}
	}
	err := k.w.WriteMessages(context.Background(), msgs...)
	if err != nil && strings.Contains(err.Error(), "Unknown Topic Or Partition") {
		time.Sleep(5 * time.Second)
		err = k.w.WriteMessages(context.Background(), msgs...)
	}
	if err != nil {
		return 0, 0, err
	}
	n := len(msgs)
	var l int64
	for _, v := range msgs {
		l += int64(len(v.Value))
	}
	return int64(n), l, err
}

func NewSinkWriter(config *Config, index int) SinkWriter {
	if config.Files[index].Topic != "" {
		w := &kafka.Writer{
			Addr:                   kafka.TCP(config.Out.Kafka.Hosts...),
			Topic:                  config.Files[index].Topic,
			AllowAutoTopicCreation: true,
		}
		if config.Out.Kafka.Username != "" {
			if config.Out.Kafka.Sasl == "scram" {
				mechanism, _ := scram.Mechanism(scram.SHA256, config.Out.Kafka.Username, config.Out.Kafka.Password)
				w.Transport = &kafka.Transport{
					SASL: mechanism,
					TLS:  &tls.Config{},
				}
			} else {
				mechanism := plain.Mechanism{Username: config.Out.Kafka.Username, Password: config.Out.Kafka.Password}
				w.Transport = &kafka.Transport{
					SASL: mechanism,
				}
			}
		}
		return &KafkaWriter{
			Config: config,
			w:      w,
		}
	}

	return &ConsoleWriter{index: index}
}
