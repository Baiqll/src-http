package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/arews-cn/src-http/src/cert"
	"github.com/thinkeridea/go-extend/exnet"
)

func main() {

	var banner = `

    _____ ____  ________    __  __      
   / ___// __ \/ ____/ /_  / /_/ /_____ 
   \__ \/ /_/ / /   / __ \/ __/ __/ __ \
  ___/ / _, _/ /___/ / / / /_/ /_/ /_/ /
 /____/_/ |_|\____/_/ /_/\__/\__/ .___/ 
                                /_/       v1.0
   
	Enabling https service dedicated to SRC testing
    `
	fmt.Println(string(banner))

	// now:=time.Now().Format("2006-01-02 15:04:05")

	var server string
	var payload string
	var close_tls bool
	var tls_path = "/tmp/.src-http/"

	flag.StringVar(&server, "server", "", "https 服务")
	flag.BoolVar(&close_tls, "close_tls", false, "关闭 tls")
	flag.StringVar(&payload, "payload", "", "payload")

	// 解析命令行参数写入注册的flag里
	flag.Parse()

	// 判断域名是否合规
	if server != "" {
		var host string
		var port string

		server_split := strings.Split(server, ":")

		host = server_split[0]
		if len(server_split) > 1 {
			port = server_split[1]
		}

		if is_host, _ := regexp.MatchString(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`, host); !is_host {
			return
		}
		if is_port, _ := regexp.MatchString(`[0-9]+`, port); !is_port {
			return
		}

	} else {
		if close_tls {
			server = "0.0.0.0:80"
		} else {
			server = "0.0.0.0:443"
		}
	}

	// 开始启动服务
	fmt.Println("[*] Starting server ", server, "...")

	err := create_tls_cert(tls_path)
	if err != nil {
		fmt.Println("TLS Cert Error")
	}

	http_server(server, tls_path+"server.pem", tls_path+"server.key", payload, close_tls)

}

// 开启文件类型模式
func http_server(server string, tls_crt string, tls_key string, payload string, close_tls bool) {

	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	// payload server
	mux.HandleFunc("/payload/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("\nDate: ", time.Now())
		fmt.Println("From: ", get_remote_ip(r))
		fmt.Println("Method:", r.Method)
		fmt.Println("URL: ", r.URL)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Write([]byte(payload))
	})

	// message server
	mux.HandleFunc("/message/", func(w http.ResponseWriter, r *http.Request) {

		content := make([]byte, r.ContentLength)
		r.Body.Read(content)

		fmt.Println("\nDate: ", time.Now())
		fmt.Println("From: ", get_remote_ip(r))
		fmt.Println("Method:", r.Method)
		fmt.Println("URL: ", r.URL)
		fmt.Println("Param: ", r.URL.RawQuery)
		fmt.Println("Body: ", string(content))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Write([]byte(`{"message": "OK"}`))
	})

	if close_tls {
		// 使用http
		http.ListenAndServe(server, mux)
	} else {
		// 使用https
		err := http.ListenAndServeTLS(server, tls_crt, tls_key, mux)
		if err != nil {
			fmt.Println("TLS Cert Error:", err.Error())
		}
	}
}

func create_tls_cert(tls_path string) (err error) {

	// tls 证书路径
	if !exists(tls_path+"server.pem") || !exists(tls_path+"server.key") {

		// 创建证书
		err = cert.CreateCert(tls_path)

	}

	return
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func get_remote_ip(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := exnet.ClientPublicIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := exnet.ClientIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// 设置hosts域名绑定
func set_hosts(host string) {

	// etc/hosts

}

// 取消hosts域名绑定
func unload_hosts(host string) {

}
