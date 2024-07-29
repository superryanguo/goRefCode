package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "ens8"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 2 * time.Second
	handle       *pcap.Handle
)

func main() {
	// 打开某一网络设备
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	i := 0
	for packet := range packetSource.Packets() {
		// Process packet here
		i++
		fmt.Println("Packet:-----")
		fmt.Println(packet)
		if i > 5 { //just print 5 packets
			break
		}
	}
}
