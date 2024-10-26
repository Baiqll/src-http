package httpx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/baiqll/src-http/src/lib"
)

type Uri struct{
	Host string
	Port string
	ServerBanner string
	IsNewDomain bool
	InternetIp string
	EnableTls bool
	ShowInternetServer bool
}

func Parse(str string, enable_tls bool)Uri{

	httpx := Uri{
		Host: "127.0.0.1",
		InternetIp: lib.GetInternetIP(),
		ShowInternetServer: true,
		EnableTls:  enable_tls,
	}

	httpx.init(str)

	return httpx

}

func(u *Uri) init_port(){

	if u.Port == ""{
		if u.EnableTls{
			u.Port = "443"
		}else{
			u.Port = "80"
		}
	}
}

func(u *Uri) init_server_banner(){
	method := "http"
	banner := ""
	Listening_port := ":"+ u.Port

	if u.EnableTls{
		method = "https"
	}
	if u.Port == "443" || u.Port == "80"{
		Listening_port = ""
		
	}

	banner += fmt.Sprintf("[*] Starting server %s://%s%s ...\n", method, u.Host, Listening_port)
	if u.ShowInternetServer{
		banner += fmt.Sprintf("[*] Internet server %s://%s%s ...\n", method, u.InternetIp, Listening_port)
	}

	banner += fmt.Sprintf("[*] Listening 0.0.0.0:%s",Listening_port)

	u.ServerBanner = banner
	
}

func (u *Uri)init(str string){
	
	// 判断域名是否合规
	if str != "" {
		
		server_split := strings.Split(str, ":")
		u.Host = server_split[0]
		if len(server_split) > 1 {
			u.Port = server_split[1]
		}

		if is_host, _ := regexp.MatchString(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`, u.Host); !is_host {

			return
		}

		if u.Host!= "0.0.0.0"{
			/*
				设置本地域名解析
			*/
			u.IsNewDomain = lib.NewDNS(u.Host)
			// show_internet_server = false
			u.ShowInternetServer = false

		}
		
	}

	lib.NewDNS(u.InternetIp)

	u.init_port()

	u.init_server_banner()

}
