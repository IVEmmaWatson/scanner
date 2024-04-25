package asd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type create interface {
	ipv4ARP() error
	ipv6NDP() error

	ipv4GetAddr() error
	ipv6GetAddr() error

	ipv4SynRequest() error
	ipv6SynRequest() error

	tcpRead() error

	route() error
	interfaceName() error
	ipCheck() bool
}

type Scanner struct {
	srcMAC  net.HardwareAddr
	dstMAC  net.HardwareAddr
	srcIP   net.IP
	dstIP   net.IP
	gwIP    net.IP
	handle  *pcap.Handle
	buf     gopacket.SerializeBuffer
	opt     gopacket.SerializeOptions
	CIDR    string
	iFName  string
	tcpPort []uint16
}

func SynScan(dstIP net.IP, gateway, srcCIDR string, port []uint16) error {

	gw := net.ParseIP(gateway)
	cidr := srcCIDR
	
	s, err := newSyn(cidr, dstIP, gw, port)
	if err != nil {
		return err
	}

	device := getDevice(s.srcIP) // 获取发送和接收数据包流量接口的名称
	interfaceHandle, err := pcap.OpenLive(device, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	s.handle = interfaceHandle
	defer interfaceHandle.Close()

	switch s.ipCheck() {
	// true 为ipv4
	case true:
		if err := s.ipv4ARP(); err != nil {
			return errors.New("ipv4 arp request error")
		}
		time.Sleep(time.Second * 3)
		if err = s.ipv4GetAddr(); err != nil {
			return errors.New("ipv4 get addr error")
		}
		if err = s.ipv4SynRequest(); err != nil {
			return errors.New("ipv4 syn request error")
		}
	// false 为ipv6
	case false:
		if err = s.ipv6NDP(); err != nil {
			return errors.New("ipv6 ndp request error")
		}

		time.Sleep(time.Second * 3)
		if err = s.ipv6GetAddr(); err != nil {
			return errors.New("ipv6 get addr error")
		}

		if err = s.ipv6SynRequest(); err != nil {
			return errors.New("ipv6 syn request error")
		}
	}

	time.Sleep(time.Second * 5)
	if err = s.tcpRead(); err != nil {
		fmt.Println("syn read error:", err)
	}
	return nil
}

func newSyn(cidr string, dst, gw net.IP, port []uint16) (*Scanner, error) {
	s := &Scanner{
		dstIP: dst,
		CIDR:  cidr,
	}

	tcpPort := port

	s.tcpPort = tcpPort
	err := s.interfaceName()
	if err != nil {
		return nil, err
	}

	s.gwIP = gw
	err = s.route()
	if err != nil {
		return nil, err
	}

	option := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	s.opt = option

	return s, nil
}

func getDevice(src net.IP) string {
	var device string
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, iFace := range interfaces {
		for _, address := range iFace.Addresses {
			if address.IP.Equal(src) {
				device = iFace.Name
			}
		}
	}

	return device
}

func (s *Scanner) ipv4ARP() error {
	if s.gwIP == nil {
		s.gwIP = s.dstIP
	}

	eth := layers.Ethernet{
		SrcMAC:       s.srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet, // 表示连接上层链路层的类型为以太网
		Protocol:          layers.EthernetTypeIPv4, // arp工作在网络层，所以协议是ip层的ipv4
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest, // 1为request，2为reply
		SourceHwAddress:   []byte(s.srcMAC),
		SourceProtAddress: []byte(s.srcIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(s.gwIP.To4()),
	}

	buf := gopacket.NewSerializeBuffer()

	if err := gopacket.SerializeLayers(buf, s.opt, &eth, &arp); err != nil {
		return err
	}

	if err := s.handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}

	fmt.Println("ARP request sent successfully.")
	return nil
}

func (s *Scanner) ipv6NDP() error {
	if s.gwIP == nil {
		s.gwIP = s.dstIP
	}

	eth := layers.Ethernet{
		SrcMAC:       []byte(s.srcMAC),
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv6,
	}

	ipv6 := layers.IPv6{
		Version:    6,
		NextHeader: layers.IPProtocolICMPv6,
		HopLimit:   255,
		SrcIP:      []byte(s.srcIP.To16()),
		DstIP:      []byte(s.dstIP.To16()),
	}

	icmp := layers.ICMPv6{
		TypeCode: layers.CreateICMPv6TypeCode(layers.ICMPv6TypeNeighborSolicitation, 0),
		Checksum: 0,
	}

	opt := layers.ICMPv6Option{
		Type: layers.ICMPv6OptSourceAddress, // 用于指定发送者的地址类型
		Data: []byte(s.srcMAC),              // data是一个字节数组，表示发送者的 MAC 地址
	}

	ndp := layers.ICMPv6NeighborSolicitation{
		TargetAddress: []byte(s.dstIP.To16()),
	}
	ndp.Options = append(ndp.Options, opt) // 无法将option赋值给options

	err := icmp.SetNetworkLayerForChecksum(&ipv6)
	if err != nil {
		return err
	}

	buf := gopacket.NewSerializeBuffer()

	if err := gopacket.SerializeLayers(buf, s.opt, &eth, &ipv6, &icmp, &ndp); err != nil {
		return err
	}

	if err := s.handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}

	fmt.Println("ARP request sent successfully.")
	return nil
}

func (s *Scanner) ipv4GetAddr() error {
	var Mac net.HardwareAddr

	start := time.Now()
	for {
		data, _, err := s.handle.ReadPacketData()
		if err != nil {
			if err == pcap.NextErrorTimeoutExpired {
				fmt.Println("ARP 响应超时")
				break
			}
			return err
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp := arpLayer.(*layers.ARP)
			if net.IP(arp.SourceProtAddress).Equal(s.gwIP) {
				fmt.Println("ARP 响应接收成功")
				// fmt.Println("ARP 获取的目标 MAC 地址:", arp.SourceHwAddress)
				Mac = arp.SourceHwAddress
				break
			}
		}

		// 检查超时
		if time.Since(start) >= time.Second*20 {
			fmt.Println("ARP 响应超时")
			break
		}
	}
	s.dstMAC = Mac
	return nil
}

func (s *Scanner) ipv6GetAddr() error {
	var Mac net.HardwareAddr

	start := time.Now()
	for {
		data, _, err := s.handle.ReadPacketData()
		if err != nil {
			if err == pcap.NextErrorTimeoutExpired {
				fmt.Println("NDP 响应超时")
				break
			}
			return err
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		if icmpV6 := packet.Layer(layers.LayerTypeICMPv6); icmpV6 != nil {
			icmp := icmpV6.(*layers.ICMPv6)
			if icmp.TypeCode.Type() == layers.ICMPv6TypeNeighborAdvertisement {
				// NDP协议数据都是28字节，最后6字节为mac地址，中间16字节为ipv6地址
				ip6 := net.IP(icmp.Payload[4:20])
				if ip6.Equal(s.dstIP) {
					Mac = net.HardwareAddr(icmp.Payload[len(icmp.Payload)-6:])
					break
				}

			}
		}

		// 检查超时
		if time.Since(start) >= time.Second*20 {
			fmt.Println("NDP 响应超时")
			break
		}
	}
	s.dstMAC = Mac
	return nil
}

func (s *Scanner) ipv4SynRequest() error {
	srcMAC := s.srcMAC
	gatewayMAC := s.dstMAC

	eth := layers.Ethernet{
		DstMAC:       gatewayMAC, // 网关mac
		SrcMAC:       srcMAC,     // 本机mac
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := layers.IPv4{
		Version:  4,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    []byte(s.srcIP.To4()), // 本机ip
		DstIP:    []byte(s.dstIP.To4()), // 目标ip
	}

	tcp := layers.TCP{
		SrcPort: layers.TCPPort(12345),
		SYN:     true,
		// Window:  8192,
	}
	err := tcp.SetNetworkLayerForChecksum(&ip4)
	if err != nil {
		return err
	}

	for _, port := range s.tcpPort {
		tcp.DstPort = layers.TCPPort(port)
		buffer := gopacket.NewSerializeBuffer()
		if err := gopacket.SerializeLayers(buffer, s.opt, &eth, &ip4, &tcp); err != nil {
			return err
		}

		packetData := buffer.Bytes()
		if err := s.handle.WritePacketData(packetData); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) ipv6SynRequest() error {
	srcMAC := s.srcMAC
	gatewayMAC := s.dstMAC

	eth := layers.Ethernet{
		DstMAC:       gatewayMAC, // 网关mac
		SrcMAC:       srcMAC,     // 本机mac
		EthernetType: layers.EthernetTypeIPv6,
	}

	ip6 := layers.IPv6{
		Version:    6,
		NextHeader: layers.IPProtocolTCP,
		HopLimit:   255,
		SrcIP:      []byte(s.srcIP.To16()), // 本机ip
		DstIP:      []byte(s.dstIP.To16()), // 目标ip
	}

	tcp := layers.TCP{
		SrcPort: layers.TCPPort(12345),
		SYN:     true,
		// Window:  8192,
	}
	err := tcp.SetNetworkLayerForChecksum(&ip6)
	if err != nil {
		return err
	}

	for _, port := range s.tcpPort {
		tcp.DstPort = layers.TCPPort(port)
		buffer := gopacket.NewSerializeBuffer()
		if err := gopacket.SerializeLayers(buffer, s.opt, &eth, &ip6, &tcp); err != nil {
			return err
		}

		packetData := buffer.Bytes()
		if err := s.handle.WritePacketData(packetData); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) tcpRead() error {

	start := time.Now()
	for {
		data, _, err := s.handle.ReadPacketData()
		if err != nil {
			if err == pcap.NextErrorTimeoutExpired {
				return errors.New("TCP 读取错误错误")
			}
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			if tcp.DstPort == 12345 && tcp.SYN && tcp.ACK {
				fmt.Printf("%v port to %v\n", tcp.SrcPort, tcp.DstPort)
				fmt.Printf("%v port is open\n", tcp.SrcPort)
				continue
			}
		}

		if time.Since(start) > time.Second*5 {
			return errors.New("tcp port check is end")

		}

	}

}

func (s *Scanner) route() error {
	ip, ipNet, err := net.ParseCIDR(s.CIDR)
	if err != nil {
		return err
	}
	mask := ipNet.Mask
	src := ip.Mask(mask)
	dst := s.dstIP.Mask(mask)
	if src.Equal(dst) {
		s.gwIP = nil
	}
	s.srcIP = ip
	return nil
}

func (s *Scanner) interfaceName() error {
	iFace, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, face := range iFace {
		Addr, err := face.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range Addr {
			if addr.String() == s.CIDR {
				s.iFName = face.Name
				s.srcMAC = face.HardwareAddr
			}
		}
	}
	return nil
}

func (s *Scanner) ipCheck() bool {
	if s.dstIP.To4() != nil {
		return true
	} else {
		return false
	}
}
