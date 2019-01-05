package ipquery

import (
	"log"
	"os"
)

var InitIpList *IpList

func init() {
	InitIpList = NewIpList()
	url := "./testdata/ip.txt"
	err := IpLoad(url)
	if err == nil {
		println("init success")
	} else {
		log.Fatal(err)
	}
}

func IpLoad(url string) error{
	reader, err:= os.Open(url)
	if err != nil {
		return err
	}
	return InitIpList.Load(reader)
}

func IpFind(ip string) ([]byte, error) {
	ipMiddle, err := InitIpList.Find(ip)
	if err != nil {
		return nil, err
	}
	return ipMiddle.Data, nil
}

func Length() int {
	return InitIpList.Length()
}