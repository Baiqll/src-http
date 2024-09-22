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
	var enable_tls bool
	var default_file string
	var tls_path = filepath.Join(lib.HomeDir(), ".config/src-http")
	var internet_ip = lib.GetInternetIP()
	var domain string
	var port string
	var method string
	var web_server string
	var is_new_domain = false
	var show_internet_server = true

	flag.StringVar(&server, "server", "", "https 服务")
	flag.BoolVar(&enable_tls, "tls", false, "是否开启tls，默认关闭")
	flag.StringVar(&payload, "payload", "", "payload")
	flag.StringVar(&default_file, "f", "", "default_file")

	// 解析命令行参数写入注册的flag里
	flag.Parse()

	// 判断域名是否合规
	if server != "" {
		

		server_split := strings.Split(server, ":")
		domain = server_split[0]
		if len(server_split) > 1 {
			port = server_split[1]
		}

		if is_host, _ := regexp.MatchString(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`, domain); !is_host {

			return
		}

		if domain!= "0.0.0.0"{
			/*
				设置本地域名解析
			*/
			is_new_domain = lib.NewDNS(domain)
			show_internet_server = false

		}
		
	}

	lib.NewDNS(internet_ip)

	if port == ""{
		if enable_tls{
			port = "443"
		}else{
			port = "80"
		}
	}

	server = "0.0.0.0:"+ port


	if enable_tls{
		method = "https"
	}else{
		method = "http"
	}

	if domain !=""{
		web_server = method + "://" + domain + ":" + port
	}else{
		web_server = method + "://127.0.0.1:" + port
	}


	// 开始启动服务
	fmt.Println("[*] Starting server ",web_server, "...")
	if show_internet_server{
		fmt.Println("[*] Internet server ", method + "://" + internet_ip + ":" + port )
	}
	
	fmt.Println("[*] Listening ", server)

	err := cert.CreateTlsCert(tls_path,[]string{domain},internet_ip, is_new_domain)
	if err != nil {
		fmt.Println("TLS Cert Error")
	}

	http_server(server, filepath.Join(tls_path, "server.pem"), filepath.Join(tls_path, "server.key"), payload, default_file, enable_tls)

}

func http_write(w http.ResponseWriter, res_data []byte){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write(res_data)
}

// 开启文件类型模式
func http_server(server string, tls_crt string, tls_key string, payload string, default_file string, enable_tls bool) {

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

		if (strings.HasPrefix(r.URL.String(), "/redirect")){
			// 设置重定向

			location := r.URL.Query().Get("url")

			http.Redirect(w, r, location, http.StatusFound)

		}else if(strings.HasPrefix(r.URL.String(), "/default")){
			// 设置默认信息

			data, err := ioutil.ReadFile(default_file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http_write(w,data)
			
		}else if(strings.HasPrefix(r.URL.String(), "/Payload")){
			// 自定义返回内容

			http_write(w,[]byte(payload))

		}else if(strings.HasPrefix(r.URL.String(), "/message")){
			// 返回全内容（接收消息）

			http_write(w,[]byte(`{"message": "OK"}`))
			
		}else{
			// 文件系统

			http.FileServer(http.Dir("./")).ServeHTTP(w, r)

		}

	})

	if enable_tls {
		// 使用https
		err := http.ListenAndServeTLS(server, tls_crt, tls_key, mux)
		if err != nil {
			fmt.Println("TLS Cert Error:", err.Error())
		}
		
	} else {
		// 使用http
		http.ListenAndServe(server, mux)
	}
}

