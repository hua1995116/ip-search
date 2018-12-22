package main

import(
	"fmt"
	"html/template"
    "log"
    "io"
	"net/http"
    "regexp"
    "strings"
    "strconv"
    "crypto/md5"
    "time"
    // "encoding/base64"
    "encoding/hex"
)

const (
    MAX_COUNT = 2000
)

var ipMap = make(map[string]int)

func getIp(w http.ResponseWriter, r *http.Request) {
    ip := r.Header.Get("X-Real-IP")
    if ip == ""{
	    // 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
    }

    r.ParseForm(); 
    key := r.Form.Get("key")
    
    if key == "" {
        // 没有key的情况，公用池
        fmt.Println("公用情况")
        if ipMap["public"] > 2000 {
            fmt.Fprintf(w, "接口调用频繁")
            return;
        }
        if ipMap["public"] == 0 {
            ipMap["public"] = 1
        } else {
            ipMap["public"] = ipMap["public"] + 1
        }
        fmt.Println(ipMap)
        fmt.Fprintf(w, "公用通道")
        return;
    } else {
        fmt.Println(key)
        if ipMap[ip] > 2000 {
            fmt.Fprintf(w, "接口调用频繁")
        }
        if ipMap[ip] == 0 {
            ipMap[ip] = 1
        } else {
            ipMap[ip] = ipMap[ip] + 1
        }
        fmt.Println(ipMap)
        fmt.Fprintf(w, "单用通道")
    }
}

func resetPublicMap() {
    const MINUTE = time.Second * 60
    t1 := time.NewTimer(MINUTE * 5)
    for {
        select {
        case <-t1.C:
            ipMap["public"] = 0
            t1.Reset(MINUTE * 5)
        }
    }
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        t, _ := template.ParseFiles("./search.html")
        log.Println(t.Execute(w, nil))
    } else {
        Squery := r.FormValue("search")
        fmt.Println(Squery)

        if len(Squery) > 50 && len(Squery) < 1 {
            fmt.Fprintf(w, "请填写正确的长度!")
            return
        }

        if m, _ := regexp.MatchString(`^\d{1,4}\.\d{1,4}\.\d{1,4}.\d{1,4}$`, Squery); !m {
            fmt.Fprintf(w, "格式错误!")
            return
        }
        fmt.Fprintf(w, "ip填写正确")
    }
}

func handleApply (w http.ResponseWriter, r *http.Request) {
    crutime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(crutime, 10))
    sum := hex.EncodeToString(h.Sum(nil))
    fmt.Fprintf(w, sum)
}

func main() {
	http.HandleFunc("/ip", getIp)       //设置访问的路由
    http.HandleFunc("/search", handleSearch)         //设置访问的路由
    http.HandleFunc("/apply", handleApply)
    err := http.ListenAndServe(":9090", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    resetPublicMap()
}