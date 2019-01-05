package ipquery

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"sync"
)

type Tools struct {
}

var (
	Tool = New()
	once sync.Once
)

/**
 * 返回单例实例
 * @method New
 */
func New() (t *Tools) {
	once.Do(func() { //只执行一次
		t = &Tools{}
	})
	return t
}

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

func (t *Tools) Ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}