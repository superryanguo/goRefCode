package main

import (
	"log"
	"net/url"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
)

func main() {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	cfg.Name = "node7"
	cfg.InitialCluster = "node7=http://168.21.38.7:2380,node16=http://168.21.38.16:2380"
	cfg.ListenPeerUrls = []url.URL{{Scheme: "http", Host: "168.21.38.7:2380"}}
	cfg.ListenClientUrls = []url.URL{{Scheme: "http", Host: "168.21.38.7:2379"}}
	//ListenPeerUrls, ListenClientUrls, ListenClientHttpUrls []url.URL
	//AdvertisePeerUrls, AdvertiseClientUrls                 []url.URL
	cfg.AdvertisePeerUrls = []url.URL{{Scheme: "http", Host: "168.21.38.7:2380"}}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}

	for {
	}
}
