package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	// 设置Kafka Broker地址
	brokers := []string{"localhost:9092"}

	// 设置消费者组ID
	groupID := "my-consumer-group"

	// 设置要订阅的主题(topic)
	topic := "my-topic"

	// 创建配置对象
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 创建消费者实例
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	deferfunc() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	// 加入消费者组
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, consumer)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	deferfunc() {
		if err := consumerGroup.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		}
	}()

	// 创建等待组，用于协调关闭信号和消费完成信号
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// 开启消费者协程
	gofunc() {
		defer wg.Done()
		for {
			topics := []string{topic}
			handler := &consumerHandler{}

			// 消费消息
			err := consumerGroup.Consume(topics, handler)
			if err != nil {
				log.Printf("Error consuming message: %v", err)
			}

			// 检查是否有中断信号
			if handler.interrupted() {
				return
			}
		}
	}()

	// 监听关闭信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// 等待关闭信号
	<-signals

	// 关闭消费者组，等待消费完成
	handler := &consumerHandler{}
	handler.shutdown()
	wg.Wait()

	fmt.Println("Consumer stopped")
}

// 自定义消费者处理程序
type consumerHandler struct {
	ready      chanbool
	interruptedFlag bool
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	h.ready = make(chanbool)
	returnnil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	returnnil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// 处理收到的消息
		fmt.Printf("Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
	}

	returnnil
}

func (h *consumerHandler) interrupted() bool {
	return h.interruptedFlag
}

func (h *consumerHandler) shutdown() {
	h.interruptedFlag = true
	close(h.ready)
}
