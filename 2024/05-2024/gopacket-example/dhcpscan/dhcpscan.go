package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	//handle, err := pcap.OpenLive("ens8", 1024, false, pcap.BlockForever)
	handle, err := pcap.OpenLive("ens8", 1900, true, 10*time.Second)
	if err != nil {
		log.Fatalf("无法打开网络接口: %v", err)
	}
	defer handle.Close()

	filter := "udp and (port 67 or port 68)"
	if err = handle.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}

	go sendDHCPDiscover(handle)
	captureDHCPOffer(handle)
}

func captureDHCPOffer(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4)
		if dhcpLayer != nil {
			fmt.Println("+++++++++++++++++++++++Capturing dhcp packet, processing...+++++++++++++++++++++++")
			dhcp, _ := dhcpLayer.(*layers.DHCPv4)
			fmt.Println("捕获到 DHCPv4 数据包:")
			fmt.Printf("  - 事务 ID: %08x\n", dhcp.Xid)
			fmt.Printf("  - YourClientIP 地址: %s\n", dhcp.YourClientIP)
			fmt.Printf("  - ClientIP  地址: %s\n", dhcp.ClientIP)
			fmt.Printf("  - NextServerIP  地址: %s\n", dhcp.NextServerIP)
			fmt.Printf("  - RelayAgentIP  地址: %s\n", dhcp.RelayAgentIP)
			fmt.Printf("  - File: %s\n", string(dhcp.File))
			fmt.Printf("  - ServerName: %s\n", string(dhcp.ServerName))
			fmt.Printf("  - ClientHWAddr: %s\n", dhcp.ClientHWAddr)
			fmt.Printf("  - DHCP配置参数字符: %s\n", dhcp.Options.String())
			if dhcp.Operation == layers.DHCPOpReply {
				fmt.Println("######捕获到 DHCP Reply 数据包")
				for _, v := range dhcp.Options {
					//fmt.Printf("------dhcpOption k=%d, v=%s\n", k, v)
					if v.Type == layers.DHCPOptMessageType && v.Length == 1 && v.Data[0] == 0x02 {
						fmt.Println("########捕获到 DHCP Offer 数据包:")
						fmt.Printf("########OFFER配置参数字符: %s\n", dhcp.Options.String())
					}
				}
			}
		}
	}
}

func sendDHCPDiscover(handle *pcap.Handle) {
	// 创建 DHCP 数据包
	dhcpLayer := &layers.DHCPv4{
		Operation:    layers.DHCPOpRequest,
		HardwareType: layers.LinkTypeEthernet,
		HardwareLen:  6,
		Xid:          uint32(time.Now().UnixNano()),
		Options: []layers.DHCPOption{
			{
				Type:   layers.DHCPOptMessageType,
				Length: 1,
				Data:   []byte{byte(layers.DHCPMsgTypeDiscover)},
			},
			{
				Type:   layers.DHCPOptParamsRequest,
				Length: 7,
				Data:   []byte{byte(layers.DHCPOptSubnetMask), byte(layers.DHCPOptRouter), byte(layers.DHCPOptDNS), byte(layers.DHCPOptLeaseTime), byte(layers.DHCPOptDomainName), byte(layers.DHCPOptVendorOption), byte(layers.DHCPOptRouterDiscovery)},
			},
		},
	}

	// 创建以太网帧
	ethernetLayer := &layers.Ethernet{
		//SrcMAC: net.HardwareAddr{0xcc, 0x5e, 0xf8, 0xb4, 0xb7, 0x85},
		SrcMAC:       net.HardwareAddr{0x52, 0x54, 0x00, 0x82, 0xe3, 0xc9},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}

	// 创建 IP 层
	ipLayer := &layers.IPv4{
		Version:    4,
		IHL:        5,
		TOS:        0,
		Length:     20,
		Id:         1,
		Flags:      layers.IPv4DontFragment,
		FragOffset: 0,
		TTL:        64,
		Protocol:   layers.IPProtocolUDP,
		SrcIP:      net.IP{0, 0, 0, 0},
		DstIP:      net.IP{255, 255, 255, 255},
	}

	// 创建 UDP 层
	udpLayer := &layers.UDP{
		SrcPort: 68,
		DstPort: 67,
	}
	udpLayer.SetNetworkLayerForChecksum(ipLayer)

	// 组装数据包
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, ethernetLayer, ipLayer, udpLayer, dhcpLayer); err != nil {
		log.Fatalf("无法序列化数据包: %v", err)
	}

	// 发送数据包
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		log.Fatalf("无法发送数据包: %v", err)
	}

	fmt.Println("已发送 DHCP Discover 数据包")
}
