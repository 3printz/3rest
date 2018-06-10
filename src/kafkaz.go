package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"log"
	"os"
	"time"
)

func initKafkaz() {
	//setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// consumer config
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	// join to consumer group
	zookeeperConn := kafkaConfig.zhost + ":" + kafkaConfig.zport
	cg, err := consumergroup.JoinConsumerGroup(kafkaConfig.cgroup, []string{kafkaConfig.topic}, []string{zookeeperConn}, config)
	if err != nil {
		fmt.Println("Error consumer group: ", err.Error())
		os.Exit(1)
	}

	// run consumer
	consume(cg)
}

func consume(cg *consumergroup.ConsumerGroup) {
	for {
		select {
		case msg := <-cg.Messages():
			// messages coming through chanel
			// only take messages from subscribed topic
			if msg.Topic != kafkaConfig.topic {
				continue
			}

			fmt.Println("Topic: ", msg.Topic)
			fmt.Println("Value: ", string(msg.Value))

			// commit to zookeeper that message is read
			// this prevent read message multiple times after restart
			err := cg.CommitUpto(msg)
			if err != nil {
				fmt.Println("Error commit zookeeper: ", err.Error())
			}

			// TODO start goroutene to handle the senz message
		}
	}
}
