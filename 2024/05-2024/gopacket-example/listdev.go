package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// 打印设备信息
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}

//// Interface describes a single network interface on a machine.
//type Interface struct {
//Name        string //设备名称
//Description string //设备描述信息
//Flags       uint32
//Addresses   []InterfaceAddress //网口的地址信息列表
//}
//// InterfaceAddress describes an address associated with an Interface.
//// Currently, it's IPv4/6 specific.
//type InterfaceAddress struct {
//IP        net.IP
//Netmask   net.IPMask // Netmask may be nil if we were unable to retrieve it.
//Broadaddr net.IP     // Broadcast address for this IP may be nil
//P2P       net.IP     // P2P destination address for this IP may be nil
//}
