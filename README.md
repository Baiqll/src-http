# src-http
Goalng 编写的简易web 服务，一键部署exp

#### 安装
```shell

go install github.com/arews-cn/src-http@latest

```

#### 使用
```shell

src-http

# src-http payload <img/src/onerror=alert(1)>
```
默认开启 https 服务
* https://0.0.0.0/message 接收信息接口
* https://0.0.0.0/payload 返回指定payload信息接口
* https://0.0.0.0/        当前路径下的ftp 服务

信息接收服务

<img width="791" alt="image" src="https://user-images.githubusercontent.com/77313240/226531697-b5cf2d15-ed04-4006-ac91-1f552536d124.png">

ftp 服务

<img width="791" alt="image" src="https://user-images.githubusercontent.com/77313240/226533305-6e2a9c8c-a5d3-4309-9c17-a5e66c7f1baa.png">
