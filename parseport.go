package ppp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PortRange 80 -> 80-80,  80-81
type PortRange struct {
	begin uint16
	end   uint16
}

type intRange struct {
	begin int
	end   int
}

func ParsePort(s string) ([]uint16, error) {
	prs, err := portArrayList(s)
	if err != nil {
		return nil, err
	}
	ports := make([]uint16, 0, len(prs)*16)
	for _, pr := range prs {
		ports = append(ports, pr.Ports()...)
	}
	return uniquePortRange(ports), nil
}

func portArrayList(s string) ([]*PortRange, error) {
	results := make([]*PortRange, 0, len(s)/8) // 这里len(s)/8忘了
	s = strings.ReplaceAll(s, " ", "")
	sections := strings.Split(s, ",")

	for _, section := range sections {
		if section == "" {
			continue
		}
		pr, err := portRangeSwitch(section)
		if err != nil {
			return nil, err
		}
		results = append(results, pr)
	}

	if len(results) < 1 {
		return nil, errors.New("empty port range")
	}
	return results, nil
}

func portRangeSwitch(s string) (*PortRange, error) {
	var pr *PortRange

	switch {
	case strings.Contains(s, "-"):
		// 必须把pr, err = ir.ToUint16()放到err!=nil后面，
		//因为如果有err代表ir为nil,ir为nil就会导致.ToUint16报错空指针引用
		ir, err := parsePortRange(s)
		if err != nil {
			return nil, err
		}
		pr, err = ir.ToUint16()
		if err != nil {
			return nil, err
		}
		return pr, nil
	default:
		ir, err := SinglePort(s)

		if err != nil {
			return nil, err
		}
		pr, err = ir.ToUint16()

		if err != nil {
			return nil, err
		}
		return pr, nil
	}

}

func parsePortRange(s string) (*intRange, error) {
	index := strings.Index(s, "-")
	count := strings.Count(s, "-")
	if index == 0 || count != 1 {
		return nil, fmt.Errorf("invalid port parse: %s", s) // errors.new 带数据返回？
	}

	begin, err := strconv.Atoi(s[:index])
	if err != nil {
		return nil, fmt.Errorf("invaild port range parse begin or end：%s", s)
	}
	end, err := strconv.Atoi(s[index+1:])
	if err != nil {
		return nil, fmt.Errorf("invaild port range parse begin or end：%s", s)
	}

	return &intRange{
		begin: begin,
		end:   end,
	}, nil
}

func SinglePort(s string) (*intRange, error) {
	begin, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("invaild single port parse：%s", s)
	}
	end, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("invaild single port parse：%s", s)
	}

	return &intRange{
		begin: begin,
		end:   end,
	}, nil
}

func (ir *intRange) ToUint16() (*PortRange, error) {
	pr := new(PortRange)
	switch {
	case ir.begin > 65535 || ir.end > 65535:
		return nil, fmt.Errorf("port out of range:%v-%v", ir.begin, ir.end)
	case ir.begin > ir.end:
		return nil, fmt.Errorf("the begin port is larger than the end port: %v-%v", ir.begin, ir.end)
	}
	pr.begin = uint16(ir.begin)
	pr.end = uint16(ir.end)
	return pr, nil
}

func (pr *PortRange) Ports() []uint16 {
	result := make([]uint16, 0, pr.end-pr.begin+1)
	for i := pr.begin; i <= pr.end; i++ {
		result = append(result, i)
	}
	return result
}

func uniquePortRange(ports []uint16) []uint16 {
	uniquePort := make(map[uint16]bool)
	var result []uint16
	for _, v := range ports {
		if _, ok := uniquePort[v]; ok {
			continue
		}
		uniquePort[v] = true
		result = append(result, v)
	}
	return result
}
