package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
)

// if you want to add a member you need to add it with below calling to inform existing cluster
// etcdctl member add infra3 --peer-urls=http://10.0.1.13:2380
// etcdctl --endpoints=http://192.21.38.16:2380 member add node162 --peer-urls=http://192.21.38.16:2384
//etcdctl --endpoints=http://192.21.38.16:2380 member add node162 --peer-urls=http://192.21.38.16:2384                                                                                               ght@ght
//Member 7e7af65b071921ff added to cluster 95a4eb0e242da7f8

//ETCD_NAME="node162"
//ETCD_INITIAL_CLUSTER="node7=http://192.21.38.7:2380,node71=http://192.21.38.7:2382,node162=http://192.21.38.16:2384,node16=http://192.21.38.16:2380"
//ETCD_INITIAL_ADVERTISE_PEER_URLS="http://192.21.38.16:2384"
//ETCD_INITIAL_CLUSTER_STATE="existing"

// etcdctl --endpoints=http://192.21.38.7:2380 endpoint status
//etcdctl --endpoints="http://192.21.38.16:2380" member list

func main() {
	var (
		dataDir  = flag.String("data-dir", "data", "data directory")
		name     = flag.String("name", "node162", "node name")
		joinAddr = flag.String("join", "http://192.21.38.16:2384", "join address (e.g., http://192.21.38.16:2384)")
	)
	flag.Parse()

	if *joinAddr == "" {
		log.Fatalf("join address is required")
	}

	log.Println("joinaddr=", *joinAddr)
	// 配置 etcd 节点
	cfg := embed.NewConfig()
	cfg.Dir = *dataDir
	cfg.Name = *name
	cfg.ClusterState = embed.ClusterStateFlagExisting
	cl := fmt.Sprintf("%s=%s", *name, *joinAddr)

	//cfg.InitialCluster = "node71=http://192.21.38.7:2382,node7=http://192.21.38.7:2380,node16=http://192.21.38.16:2380"
	cfg.InitialCluster = "node71=http://192.21.38.7:2382,node7=http://192.21.38.7:2380,node16=http://192.21.38.16:2380" + "," + cl

	cfg.ListenPeerUrls = []url.URL{{Scheme: "http", Host: "192.21.38.16:2384"}}
	cfg.ListenClientUrls = []url.URL{{Scheme: "http", Host: "192.21.38.16:2383"}}
	cfg.AdvertiseClientUrls = []url.URL{{Scheme: "http", Host: "192.21.38.16:2383"}}
	cfg.AdvertisePeerUrls = []url.URL{{Scheme: "http", Host: "192.21.38.16:2384"}}

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

	log.Fatal(<-e.Err())

}
