package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"log"
	"os"
	"time"
)

// channel to publish kafka messages
var kchan = make(chan Kmsg, 10)

func initKafkaz() {
	//setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// consuner
	cg := initConzumer()
	go conzume(cg)

	// producer
	pr := initProduzer()
	go produze(pr)
}

func initConzumer() *consumergroup.ConsumerGroup {
	// consumer config
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	// join to consumer group
	zookeeperConn := kafkaConfig.zhost + ":" + kafkaConfig.zport
	cg, err := consumergroup.JoinConsumerGroup(kafkaConfig.cgroup,
		[]string{kafkaConfig.topic},
		[]string{zookeeperConn},
		config)
	if err != nil {
		fmt.Println("Error consumer group: ", err.Error())
		os.Exit(1)
	}

	return cg
}

func initProduzer() sarama.SyncProducer {
	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	kafkaConn := kafkaConfig.khost + ":" + kafkaConfig.kport
	pr, err := sarama.NewSyncProducer([]string{kafkaConn}, config)
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	return pr
}

func conzume(cg *consumergroup.ConsumerGroup) {
	for {
		select {
		case msg := <-cg.Messages():
			// messages coming through chanel
			// only take messages from subscribed topic
			if msg.Topic != kafkaConfig.topic {
				continue
			}

			z := string(msg.Value)
			fmt.Println("Received topic, msg: ", msg.Topic, z)

			// commit to zookeeper that message is read
			// this prevent read message multiple times after restart
			err := cg.CommitUpto(msg)
			if err != nil {
				fmt.Println("Error commit zookeeper: ", err.Error())
			}

			// send message to corresponding channel in rchan
			rchans[z.Attr["uid"]] <- z
		}
	}
}

func produze(pr sarama.SyncProducer) {
	for {
		select {
		case kmsg := <-kchan:
			// received kafka message to send
			// publish sync
			msg := &sarama.ProducerMessage{
				Topic: kmsg.Topic,
				Value: sarama.StringEncoder(kmsg.Msg),
			}
			p, o, err := pr.SendMessage(msg)
			if err != nil {
				fmt.Println("Error publish: ", err.Error())
			}
			fmt.Println("Published msg, partition, offset: ", kmsg.Msg, p, o)
		}
	}
}
