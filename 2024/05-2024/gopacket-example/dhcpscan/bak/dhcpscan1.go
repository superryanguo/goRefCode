package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// 打开网络接口并将其设置为捕获模式
	handle, err := pcap.OpenLive("ens8", 1500, true, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器以只捕获 DHCP 响应数据包
	filter := "udp and (port 67 or port 68)"
	//filter := "udp and dst port 67"
	if err = handle.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}

	// 捕获 DHCP 响应数据包并解析其中的信息
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("message processing...")
		dhcpPacket := packet.Layer(layers.LayerTypeDHCPv4)
		if dhcpPacket == nil {
			continue
		}

		// 解析 DHCP 响应数据包的 IP 地址、MAC 地址、主机名、域名和配置参数等信息
		dhcpLayer := dhcpPacket.(*layers.DHCPv4)
		fmt.Printf("IPc 地址: %s\n", dhcpLayer.ClientIP)
		fmt.Printf("IPYc 地址: %s\n", dhcpLayer.YourClientIP)
		fmt.Printf("IP next server 地址: %s\n", dhcpLayer.NextServerIP)
		fmt.Printf("IP RelayAgentIP 地址: %s\n", dhcpLayer.RelayAgentIP)
		fmt.Printf("client MAC 地址: %s\n", dhcpLayer.ClientHWAddr)
		fmt.Printf("dhcp server主机名: %s\n", dhcpLayer.ServerName)
		fmt.Printf("dhcp server config file: %s\n", dhcpLayer.File)
		fmt.Printf("配置参数: %v\n", dhcpLayer.Options)
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			fmt.Println("Ethernet layer detected.")
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
			fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
			fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
			// Ethernet type is typically IPv4 but could be ARP or other
			fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
			fmt.Println()
		}
	}
}
