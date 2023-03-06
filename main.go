package main

import (
    "fmt"
	"flag"
	"os"
	"regexp"
    "net/http"
	"strings"
	"path/filepath"
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
	var tls_crt  string
	var tls_key  string
	var payload string
	var close_tls bool


	flag.StringVar(&server, "server", "", "https 服务")
	flag.BoolVar(&close_tls, "close_tls", false, "关闭 tls")
	flag.StringVar(&payload, "payload", "", "简单payload")
	flag.StringVar(&tls_crt, "crt", "server-crt.pem", "TLS crt")
	flag.StringVar(&tls_key, "key", "server-key.pem", "TLS key")
	// 解析命令行参数写入注册的flag里
	flag.Parse()

	
	// 判断域名是否合规
	if server!= "" {
		var host string
		var port string

		server_split := strings.Split(server, ":")
		
		host = server_split[0]
		if len(server_split)>1{
			port = server_split[1]
		}
		
		if is_host, _ := regexp.MatchString(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`, host); !is_host {
			return
		}
		if is_port, _ := regexp.MatchString(`[0-9]+`, port); !is_port {
			return
		}

	}else{
		if close_tls{
			server = "0.0.0.0:80"
		}else{
			server = "0.0.0.0:443"
		}
	}


	// tls 证书路径
	current_path,_ := get_current_path()
	tls_crt = current_path +"/"+ tls_crt
	tls_key = current_path +"/"+ tls_key

	// 开始启动服务
	fmt.Println("[*] Starting server ", server, "...")


	if payload!="" {

		payload_server(server, tls_crt, tls_key, payload, close_tls)
	}else{

		http_server(server, tls_crt, tls_key, close_tls)
	}

}
	

// 开启payload 模式服务
func payload_server(server string,tls_crt string,tls_key string, payload string, close_tls bool){
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){	
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Write([]byte(payload));
	})

	if close_tls{
		// 使用http
		http.ListenAndServe(server, nil)
	}else{
		// 使用https
		http.ListenAndServeTLS(server, tls_crt, tls_key, nil)
	}

}

// 开启文件类型模式
func http_server(server string,tls_crt string,tls_key string, close_tls bool){

	mux := http.NewServeMux()
  	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	if close_tls{
		// 使用http
		http.ListenAndServe(server, mux)
	}else{
		// 使用https
		http.ListenAndServeTLS(server, tls_crt, tls_key, mux)
	}
}

func get_current_path() (string, error) {
    path, err := filepath.Abs(filepath.Dir(os.Args[0]))
    return string(path), err
}


// 设置hosts域名绑定
func set_hosts(host string){

	// etc/hosts

}

// 取消hosts域名绑定
func unload_hosts(host string){

}

