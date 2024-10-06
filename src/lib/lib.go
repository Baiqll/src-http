package lib

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"

	"github.com/thinkeridea/go-extend/exnet"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func GetRemoteIp(req *http.Request) string {
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

func HomeDir() (dir string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not get user home directory: ", err)
	}
	dir = usr.HomeDir
	return
}

func GetInternetIP() (ip string) {
	/*
		查看主机对应的外网IP
	*/
	resp, _ := http.Get("http://ifconfig.me")

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	ip = string(body)

	return
}

// 设置hosts域名绑定
func SetHosts(host string) {

	// etc/hosts

}

// 取消hosts域名绑定
func UnloadHosts(host string) {

}

func NewDNS(domain string)(is_new bool){

	is_new = false

	cmd := exec.Command("sh", "-c",fmt.Sprintf("cat /etc/hosts | grep %s",domain))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	if string(out) =="" {

		is_new = true

		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '127.0.0.1	%s	# src-http' | sudo tee -a /etc/hosts > /dev/null",domain))
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}


