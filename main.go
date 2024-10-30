package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/baiqll/src-http/src/cert"
	"github.com/baiqll/src-http/src/httpx"
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
	var config_path = filepath.Join(lib.HomeDir(), ".config/src-http")
	var enable_file_server bool


	flag.StringVar(&server, "server", "", "https 服务")
	flag.BoolVar(&enable_tls, "tls", false, "是否开启tls，默认关闭")
	flag.StringVar(&payload, "payload", "", "payload")
	flag.StringVar(&default_file, "df", "", "default_file")
	flag.BoolVar(&enable_file_server,"fm", false, "是否开启文件服务，默认关闭")

	// 解析命令行参数写入注册的flag里
	flag.Parse()

	uri := httpx.Parse(server, enable_tls)

	fmt.Println(uri.ServerBanner)

	err := cert.CreateTlsCert(config_path, []string{uri.Host}, uri.InternetIp, uri.IsNewDomain)
	if err != nil {
		fmt.Println("TLS Cert Error")
	}

	HttpServer(httpx.Config{
		Server: fmt.Sprintf("0.0.0.0:%s",uri.Port),
		Host: uri.Host,
		ConfigPath: config_path,
		Payload: payload,
		DefaultFile: default_file,
		EnableTLS: enable_tls,
		EnableFileServer: enable_file_server,
	})

}

func SetResponseHeader(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func HttpWrite(w http.ResponseWriter, res_data []byte){
	SetResponseHeader(w)
	w.Write(res_data)
}

// 定义一个可以接受额外参数的HTTP处理函数
func SRCHandler(config httpx.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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

			data, err := ioutil.ReadFile(config.DefaultFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			HttpWrite(w,data)
			
		}else if(strings.HasPrefix(r.URL.String(), "/payload")){
			// 自定义返固定内容

			cmdtype := r.URL.Query().Get("type");

			new_payload := r.URL.Query().Get("payload");

			switch(cmdtype){
			case "jscmd":
				
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if _, err := fmt.Fprintln(w, `<!DOCTYPE html>`+
					`<html>`+
					`<head><title>xss</title></head>`+
					`<body>`+
					`<script>` + new_payload + `</script>`+
					`</body>`+
					`</html>`); err != nil {
						fmt.Println("jscmd Error:", err.Error())
				}
			case "html":
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if _, err := fmt.Fprintln(w, new_payload); err != nil {
					fmt.Println("html Error:", err.Error())
				}
			case "jscode":
				w.Header().Set("Content-Type", "application/javascript;")
				if _, err := fmt.Fprintln(w, new_payload); err != nil {
					fmt.Println("html Error:", err.Error())
				}

			default:
				HttpWrite(w,[]byte(config.Payload))
			}


		}else if(strings.HasPrefix(r.URL.String(), "/message")){
			// 返回全内容（接收消息）

			HttpWrite(w,[]byte(`{"message": "OK"}`))
			
		}else if(config.EnableFileServer){
			// 文件系统
			SetResponseHeader(w)
			http.FileServer(http.Dir("./")).ServeHTTP(w, r)

		}
    }
}

// 开启文件类型模式
func HttpServer(config httpx.Config) {

	mux := http.NewServeMux()

	 // 创建一个中间件链
	handlerWithMiddleware := httpx.HostMiddleware(config.Host, http.HandlerFunc(SRCHandler(config)))

	 // 设置路由并启动服务器
	mux.Handle("/", handlerWithMiddleware)

	if config.EnableTLS {
		// 使用https
		tls_crt := filepath.Join(config.ConfigPath, "server.pem")
		tls_key := filepath.Join(config.ConfigPath, "server.key")

		err := http.ListenAndServeTLS(config.Server, tls_crt, tls_key, mux)
		if err != nil {
			fmt.Println("TLS Cert Error:", err.Error())
		}
		
	} else {
		// 使用http
		http.ListenAndServe(config.Server, mux)
	}
}

