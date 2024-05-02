package main

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"strings"
)

type ipv4 struct {
	begin uint32
	end   uint32
}

type ipv6 struct {
	begin *big.Int
	end   *big.Int
}

type IpRange struct {
	begin net.IP
	end   net.IP
}

type ipManager interface {
	IsIPv4() bool
	IsIPv6() bool
	Begin() net.IP
	End() net.IP
	Num() *big.Int
	List() []net.IP
	String() string
}

func (ipv4 *ipv4) IsIPv4() bool {
	return true
}

func (ipv4 *ipv4) IsIPv6() bool {
	return false
}

func (ipv4 *ipv4) Begin() net.IP {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, ipv4.begin)
	return buf
}

func (ipv4 *ipv4) End() net.IP {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, ipv4.end)
	return buf
}

func (ipv4 *ipv4) Num() *big.Int {
	num := ipv4.end - ipv4.begin
	return big.NewInt(int64(num) + 1)
}

func (ipv4 *ipv4) List() []net.IP {
	unitList := make([]uint32, 0)
	ipv4List := make([]net.IP, 0)
	for ip := ipv4.begin; ip <= ipv4.end; ip++ {
		unitList = append(unitList, ip)
	}

	buf := make([]byte, 4)
	for _, v := range unitList {
		binary.BigEndian.PutUint32(buf, v)
		ipv4List = append(ipv4List, buf)
	}
	return ipv4List
}

func (ipv4 *ipv4) String() string {
	return fmt.Sprintf("%s-%s", ipv4.Begin(), ipv4.End())
}

func (ipv6 *ipv6) IsIPv4() bool {
	return false
}

func (ipv6 *ipv6) IsIPv6() bool {
	return true
}

func (ipv6 *ipv6) Begin() net.IP {
	ipAddr := make(net.IP, net.IPv6len)
	return ipv6.begin.FillBytes(ipAddr)
}

func (ipv6 *ipv6) End() net.IP {
	ipAddr := make(net.IP, net.IPv6len)
	return ipv6.end.FillBytes(ipAddr)
}

func (ipv6 *ipv6) Num() *big.Int {
	num := new(big.Int).Sub(ipv6.end, ipv6.begin)
	return num.Add(num, big.NewInt(1))
}

func (ipv6 *ipv6) List() []net.IP {
	bigIntList := make([]*big.Int, 0)
	ipv6List := make([]net.IP, 0)
	for ipv6.begin.Cmp(ipv6.end) <= 0 {
		buf := big.NewInt(0)
		buf.Set(ipv6.begin)
		bigIntList = append(bigIntList, buf)
		ipv6.begin = ipv6.begin.Add(ipv6.begin, big.NewInt(1))
	}

	for _, re := range bigIntList {
		ipAddr := make(net.IP, net.IPv6len)
		res := re.FillBytes(ipAddr)
		ipv6List = append(ipv6List, res)
	}
	return ipv6List
}

func (ipv6 *ipv6) String() string {
	return fmt.Sprintf("%s-%s", ipv6.Begin(), ipv6.End())
}

func newIPv4(nip *IpRange) (*ipv4, error) {
	begin := nip.begin.To4()
	end := nip.end.To4()
	if begin == nil || end == nil {
		return nil, fmt.Errorf("ipv4 address %s or %s parse wrong", begin, end)
	}
	b := binary.BigEndian.Uint32(begin)
	e := binary.BigEndian.Uint32(end)
	if b > e {
		return nil, fmt.Errorf("ipv4 address %v lager than %v", nip.begin, nip.end)
	}
	return &ipv4{
		begin: b,
		end:   e,
	}, nil
}

func newIPv6(nip *IpRange) (*ipv6, error) {
	ip := new(ipv6)
	begin := nip.begin.To16()
	end := nip.end.To16()

	if begin == nil || end == nil {
		return nil, fmt.Errorf("ipv6 address %s or %s parse wrong", begin, end)
	}

	b := new(big.Int).SetBytes(begin)
	e := new(big.Int).SetBytes(end)

	if b.Cmp(e) == 1 {
		return nil, fmt.Errorf("ipv6 address %v lager than %v", nip.begin, nip.end)
	}

	ip.begin = new(big.Int).SetBytes(b.Bytes())
	ip.end = new(big.Int).SetBytes(e.Bytes())

	return ip, nil

}

func ParseCIDRIP(s string) (*IpRange, error) {
	section := new(IpRange)
	count := strings.Count(s, "/")
	if count != 1 {
		return nil, fmt.Errorf("incorrect ipCIDR parse:%v", s)
	}
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("parse CIDR wrong:%v", err)
	}
	section.begin = ip.Mask(ipNet.Mask)
	for i := range section.begin {
		section.end = append(section.end, section.begin[i]|^ipNet.Mask[i])
	}
	return section, nil
}

func ParseRangeIP(s string) (*IpRange, error) {
	section := new(IpRange)
	hyphenIdx := strings.Index(s, "-")
	numHyphen := strings.Count(s, "-")
	section.begin = net.ParseIP(s[:hyphenIdx])
	section.end = net.ParseIP(s[hyphenIdx+1:])
	if section.begin == nil || section.end == nil || numHyphen != 1 {
		return nil, fmt.Errorf("wrong ip range,%v\n", s)
	}
	return section, nil
}

func SingleIP(s string) (*IpRange, error) {
	section := new(IpRange)
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("wrong single ip %v", s)
	}
	section.begin = ip
	section.end = ip

	return section, nil
}

func newIpRange(s string) ([]*IpRange, error) {
	var result []*IpRange
	var ir *IpRange
	var err error
	if s == "" {
		return nil, fmt.Errorf("nothing input:%s", s)
	}
	raw := strings.Split(s, ",")
	for _, v := range raw {
		switch {
		case strings.Contains(v, "-"):
			ir, err = ParseRangeIP(v)
		case strings.Contains(v, "/"):
			ir, err = ParseCIDRIP(v)
		default:
			ir, err = SingleIP(v)

		}
		if err != nil {
			return nil, err
		}
		result = append(result, ir)
	}
	return result, err
}

func (ip *IpRange) checkIPType() (*ipv4, *ipv6, error) {
	var ipv4 *ipv4
	var ipv6 *ipv6

	var err error
	if ip.begin.To4() != nil {
		ipv4, err = newIPv4(ip)
		if err != nil {
			return nil, nil, err
		}
		return ipv4, nil, nil
	}

	ipv6, err = newIPv6(ip)
	if err != nil {
		return nil, nil, err
	}

	return nil, ipv6, nil
}

func ipToString(ip []net.IP) []string {
	ipList := make([]string, 0)
	for _, v := range ip {
		ipList = append(ipList, v.String())
	}
	return ipList
}

func ParseIP(s string) ([]net.IP, error) {
	ips, err := newIpRange(s)
	var IP []net.IP
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		ipv4, ipv6, err := ip.checkIPType()
		if err != nil {
			return nil, err
		}

		switch {
		case ipv4 != nil:
			fmt.Printf("IPv4: %s\n", ipv4.String())
			fmt.Printf("IPv4 address Number : %s\n", ipv4.Num())
			fmt.Printf("IPv4 address List: %v\n", ipToString(ipv4.List()))
			fmt.Println("---------------------------------")
			IP = append(IP, ipv4.List()...)
		case ipv6 != nil:
			fmt.Printf("IPv6: %s\n", ipv6.String())
			fmt.Printf("IPv6 address Number: %s\n", ipv6.Num())
			fmt.Printf("IPv6 address List List: %v\n", ipToString(ipv6.List()))
			fmt.Println("---------------------------------")
			IP = append(IP, ipv6.List()...)
		default:
			panic("unreachable code")
		}

	}
	return IP, nil
}
