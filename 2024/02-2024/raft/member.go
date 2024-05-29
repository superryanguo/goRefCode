package main

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	log.Printf("Member begin to join the cluster")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	for i := 0; i < 19; i++ {
		log.Printf("Member doing sth......\n")
		time.Sleep(5 * time.Second)
	}
	log.Printf("Member left the cluster\n")
}
