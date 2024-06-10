package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/baiqll/src-http/src/cert"
	"github.com/baiqll/src-http/src/lib"
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
	var default_file string
	var tls_path = filepath.Join(lib.HomeDir(), ".config/src-http")

	flag.StringVar(&server, "server", "", "https 服务")
	flag.BoolVar(&close_tls, "distls", false, "关闭 tls")
	flag.StringVar(&payload, "payload", "", "payload")
	flag.StringVar(&default_file, "f", "", "default_file")

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

	err := cert.CreateTlsCert(tls_path, lib.GetInternetIP())
	if err != nil {
		fmt.Println("TLS Cert Error")
	}

	http_server(server, filepath.Join(tls_path, "server.pem"), filepath.Join(tls_path, "server.key"), payload, default_file, close_tls)

}

fun http_write(w http.ResponseWriter, res_data []byte){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write(res_data)
}

// 开启文件类型模式
func http_server(server string, tls_crt string, tls_key string, payload string, default_file string, close_tls bool) {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("%s %s %s\n",r.Method, r.URL,r.Proto)
		fmt.Printf("Host: %s\n",r.Host)
		fmt.Printf("From: %s\n",lib.GetRemoteIp(r))
		

		// 打印请求头
		for key, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

			// 打印请求体，如果请求体是可读的
		if r.Method == "POST" || r.Method == "PUT" {
			content := make([]byte, r.ContentLength)
		r.Body.Read(content)
			fmt.Printf("\n%s\n", content)
		}
		fmt.Print("\n")

		if(strings.HasPrefix(r.URL, "/default")){
			data, err := ioutil.ReadFile(default_file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http_write(w,data)
			
		}else if(strings.HasPrefix(r.URL, "/Payload")){

			http_write(w,[]byte(payload))

		}else if(strings.HasPrefix(r.URL, "/message")){

			http_write(w,[]byte(`{"message": "OK"}`))
			
		}else{

			http.FileServer(http.Dir("./")).ServeHTTP(w, r)

		}

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

