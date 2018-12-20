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
    r.ParseForm(); 
    id := r.Form.Get("id")
    fmt.Println(id)
    crutime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(crutime, 10))
    sum := hex.EncodeToString(h.Sum(nil))
    fmt.Println(sum)

    if ip == ""{
	    // 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
    }

    if ipMap[ip] > 2000 {
        fmt.Fprintf(w, "接口调用频繁")
    }

    if ipMap[ip] == 0 {
        ipMap[ip] = 1
    } else {
        ipMap[ip] = ipMap[ip] + 1
    }
    fmt.Println(ipMap)
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

func main() {
	http.HandleFunc("/ip", getIp)       //设置访问的路由
    http.HandleFunc("/search", handleSearch)         //设置访问的路由
    err := http.ListenAndServe(":9090", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}