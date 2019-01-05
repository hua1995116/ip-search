package ipquery

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const (
	MAX_ITEM = 3
)

var ErrorNotFound = errors.New("can not find this ip! ")

type IpItem struct {
	Begin uint32
	End   uint32
	Data  []byte
}

type IpList []*IpItem

func NewIpList() *IpList {
	return &IpList{}
}

func (item *IpList) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		list := strings.SplitN(line, "\t", MAX_ITEM)
		if len(list) < MAX_ITEM {
			continue
		}

		begin := Tool.ParseInt(list[0], 0)
		end := Tool.ParseInt(list[1], 0)
		ipIt := &IpItem{
			Begin: uint32(begin),
			End:   uint32(end),
			Data:  []byte(list[2]),
		}

		*item = append(*item, ipIt)
	}

	return scanner.Err()
}

func (item *IpList) Length() int {
	return len(*item)
}

func (item *IpList) Find(ip string) (*IpItem, error) {
	ipMiddle, err := item.Search(ip)
	if err != nil {
		return nil, err
	}

	return ipMiddle, nil
}

func (item *IpList) Search(ip string) (*IpItem, error) {
	var low, high int = 0, (item.Length() - 1)

	ipPointer := *item
	ipLong := Tool.Ip2Long(ip)

	if ipLong <= 0 {
		return nil, ErrorNotFound
	}

	for low <= high {
		var middle int = (high-low)/2 + low

		ipMiddle := ipPointer[middle]

		if ipLong >= ipMiddle.Begin && ipLong <= ipMiddle.End {
			return ipMiddle, nil
		} else if ipLong < ipMiddle.Begin {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}

	return nil, ErrorNotFound
}
