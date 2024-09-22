# src-http

Goalng 编写的简易 web 服务，一键部署 exp

#### 安装

```shell

go install github.com/baiqll/src-http@latest

```

#### 使用

```shell

# 默认启动
src-http

# 加载 payload
src-http -payload "<img/src/onerror=alert(1)>"

# 指定域名(自动本地dns域名绑定)、端口
src-http -server example.com:8080

# 使用https
src-http -tls

# 使用自定义证书
src-http -crt /home/crt.pem -key /home/key.pem

# 指定默认模版
src-http -default default.html

```

默认开启 https 服务

- https://0.0.0.0/redirect?url=* 重定向
- https://0.0.0.0/message 接收信息接口
- https://0.0.0.0/payload 返回指定 payload 信息接口
- https://0.0.0.0/ 当前路径下的 ftp 服务

信息接收服务

<img width="791" alt="image" src="https://user-images.githubusercontent.com/77313240/226531697-b5cf2d15-ed04-4006-ac91-1f552536d124.png">

ftp 服务

<img width="791" alt="image" src="https://user-images.githubusercontent.com/77313240/226533305-6e2a9c8c-a5d3-4309-9c17-a5e66c7f1baa.png">

payload

<img width="791" alt="image" src="https://user-images.githubusercontent.com/77313240/226548810-e64cd8a2-879f-4259-9505-ca5d479ced3b.png">
