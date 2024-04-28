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

type NetIP struct {
	begin net.IP
	end   net.IP
}

type IPRange interface {
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
	result := make([]uint32, 0)
	results := make([]net.IP, 0)
	for ip := ipv4.begin; ip <= ipv4.end; ip++ {
		result = append(result, ip)
	}

	buf := make([]byte, 4)
	for _, v := range result {
		binary.BigEndian.PutUint32(buf, v)
		results = append(results, buf)
	}
	return results
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
	result := make([]*big.Int, 0)
	results := make([]net.IP, 0)
	for ipv6.begin.Cmp(ipv6.end) <= 0 {
		buf := big.NewInt(0)
		buf.Set(ipv6.begin)
		result = append(result, buf)
		ipv6.begin = ipv6.begin.Add(ipv6.begin, big.NewInt(1))
	}

	for _, re := range result {
		ipAddr := make(net.IP, net.IPv6len)
		res := re.FillBytes(ipAddr)
		results = append(results, res)
	}
	return results
}

func (ipv6 *ipv6) String() string {
	return fmt.Sprintf("%s-%s", ipv6.Begin(), ipv6.End())
}

func newIPv4(nip *NetIP) (*ipv4, error) {
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

func newIPv6(nip *NetIP) (*ipv6, error) {
	ip6 := new(ipv6)
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

	ip6.begin = new(big.Int).SetBytes(b.Bytes())
	ip6.end = new(big.Int).SetBytes(e.Bytes())

	return ip6, nil

}

func ParseCIDRIP(s string) (*NetIP, error) {
	cidr := new(NetIP)
	count := strings.Count(s, "/")
	if count != 1 {
		return nil, fmt.Errorf("incorrect ipCIDR parse:%v", s)
	}
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("parse CIDR wrong:%v", err)
	}
	cidr.begin = ip.Mask(ipNet.Mask)
	for i := range cidr.begin {
		cidr.end = append(cidr.end, cidr.begin[i]|^ipNet.Mask[i])
	}
	return cidr, nil
}

func ParseRangeIP(s string) (*NetIP, error) {
	rp := new(NetIP)
	hyphenIdx := strings.Index(s, "-")
	numHyphen := strings.Count(s, "-")
	rp.begin = net.ParseIP(s[:hyphenIdx])
	rp.end = net.ParseIP(s[hyphenIdx+1:])
	if rp.begin == nil || rp.end == nil || numHyphen != 1 {
		return nil, fmt.Errorf("wrong ip range,%v\n", s)
	}
	return rp, nil
}

func Single(s string) (*NetIP, error) {
	single := new(NetIP)
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("wrong single ip %v", s)
	}
	single.begin = ip
	single.end = ip

	return single, nil
}

func parseIP(s string) ([]*NetIP, error) {
	var netIP []*NetIP
	var ip *NetIP
	var err error
	if s == "" {
		return nil, fmt.Errorf("nothing input:%s", s)
	}
	r := strings.Split(s, ",")
	for _, a := range r {
		switch {
		case strings.Contains(a, "-"):
			ip, err = ParseRangeIP(a)
		case strings.Contains(a, "/"):
			ip, err = ParseCIDRIP(a)
		default:
			ip, err = Single(a)

		}
		if err != nil {
			return nil, err
		}
		netIP = append(netIP, ip)
	}
	return netIP, err
}

func (np *NetIP) checkIPType() (*ipv4, *ipv6, error) {
	var ipv4 *ipv4
	var ipv6 *ipv6

	var err error
	if np.begin.To4() != nil {
		ipv4, err = newIPv4(np)
		if err != nil {
			return nil, nil, err
		}
		return ipv4, nil, nil
	}

	ipv6, err = newIPv6(np)
	if err != nil {
		return nil, nil, err
	}

	return nil, ipv6, nil
}

func ipToString(ip []net.IP) []string {
	ips := make([]string, 0)
	for _, v := range ip {
		ips = append(ips, v.String())
	}
	return ips
}

func ParseIP(s string) ([]net.IP, error) {
	ips, err := parseIP(s)
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
