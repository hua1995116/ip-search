package main

import (
    "crypto/md5"
    "fmt"
    "github.com/hua1995116/ip-search/ipquery"
    "html/template"
    "io"
    "log"
    "net"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
    // "encoding/base64"
    "encoding/hex"
)

const (
    MAX_COUNT = 2000
    PUBLIC = "public"
    MINUTE = time.Second * 60
)

type IpItem struct {
    lock sync.Mutex
    count int
}

var ipMap = make(map[string]*IpItem)

func getIp(w http.ResponseWriter, r *http.Request) {
    ip := r.Header.Get("X-Real-IP")
    if ip == "" {
	    // 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
    }

    r.ParseForm()
    key := r.Form.Get("key")
    fmt.Println("key is ", key)
    if key == "" {
        ip = PUBLIC
        fmt.Println("公用通道")
    } else {
        fmt.Println("单用通道")
    }
    if ipMap[ip] == nil {
        ipMap[ip] = &IpItem{count: 0}
    }
    ipMap[ip].lock.Lock()
    ipMap[ip].count ++
    ipMap[ip].lock.Unlock()
    for key, value := range ipMap {
        fmt.Println(key, "数量为: ",value.count)
    }
    if ipMap[ip].count > MAX_COUNT {
        fmt.Println( "接口调用频繁")
        fmt.Fprintf(w, "接口调用频繁")
        return
    }
    fmt.Println( "接口调用正常")
    //fmt.Fprintf(w, "接口调用正常")
    res, err := ipquery.IpFind(ip)

    if err != nil {
        fmt.Fprintf(w, "抱歉库内暂无收入该ip")
    } else {
        fmt.Fprintf(w, string(res))
    }

}

func resetPublicMap() {
    fmt.Println("重置装置启动")
    t1 := time.NewTimer(MINUTE * 5)
    for {
        select {
            case <-t1.C:
                fmt.Println("成功重置1次")
                ipMap[PUBLIC].count = 0
                t1.Reset(MINUTE * 5)
        }
    }
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        t, _ := template.ParseFiles("./view/search.html")
        log.Println(t.Execute(w, nil))
    } else {
        sQuery := r.FormValue("search")
        fmt.Println(sQuery)

        if len(sQuery) > 50 && len(sQuery) < 1 {
            fmt.Fprintf(w, "请填写正确的长度!")
            return
        }

        addr := net.ParseIP(sQuery)
        if addr == nil {
            fmt.Fprintf(w, "格式错误!")
            return
        }
        res, err := ipquery.IpFind(sQuery)

        if err != nil {
            fmt.Fprintf(w, "抱歉库内暂无收入该ip")
        } else {
            fmt.Fprintf(w, string(res))
        }

        //fmt.Fprintf(w, "ip填写正确")

    }
}

func handleApply (w http.ResponseWriter, r *http.Request) {
    curtime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(curtime, 10))
    sum := hex.EncodeToString(h.Sum(nil))
    fmt.Fprintf(w, sum)
}


func handleRegister(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {

    }
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

}

func init () {
    go resetPublicMap()
}

func main() {
	http.HandleFunc("/ip", getIp)       //设置访问的路由
    http.HandleFunc("/search", handleSearch) //设置访问的路由
    http.HandleFunc("/apply", handleApply) //设置访问的路由
    http.HandleFunc("/register", handleRegister)
    http.HandleFunc("/login", handleLogin)
    err := http.ListenAndServe(":7878", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}