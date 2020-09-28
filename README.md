# locust-http
     base on locust and boomer, use etcd and gRPC to push http request and make locust test tasks。
     基于locust和boomer核心，使用基于etcd做为压测机服务发现。
     使用gRPC推送http请求事务描述信息，让压测机自己构造http接口测试任务。
     使用postman类似的界面管理http请求事务描述信息。
     目前是核心且基本框架及功能的完成，如有更多的要求比如任务隔离，执行历史管理......这就是你的事了。

## 说明
     拷贝了原版locust的main.py及webUI和前端部分代码进行修改。
     在boomer之上增加了gRPC服务，能解析master发来的http接口测试任务描述信息并生成boomer的任务
     boomer及grequests源码部分地方做了小改动--主要避免异常退出。
     如果work端编译失败，请将错的地方(主要是配置连接数上限配置)删除即可
     【1】worker/boomerHazardServer.go 的 baseReqOptions := &grequests.RequestOptions{ // 基本的http请求配置
     下面要删除
		MaxIdleConnsPerHost: 2000, // 限制连接数
		DisableKeepAlives:   false,
     【2】boomer源码的runner中close channel某些情况化因为重复关闭会造成异常关闭。这里最好加上recover处理


## 启动参考：这里的ip、port都是例子，请根据实际情况设置
  ### 1-先下载etcd并启动etcd：
  etcd.exe --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379
      etcd下载地址： https://github.com/etcd-io/etcd/releases/download/v3.3.25/etcd-v3.3.25-windows-amd64.zip  
  ### 2-在服务上启动matser（假定服务器ip：192.168.23.222）：
       python3 main.py --master-host=192.168.23.222
  ### 3-可以直接从web页面的压测机管理下载编译好的exe（windows专用）
       根据操作系统(windows,linux)编译对应压测机（worker）端应用程序：go build boomerHazardServer.go
  ### 4-压测机端执行此程序（假定etcd的ip：192.168.23.222）
      boomerHazardServer.exe -EtcdAddr 192.168.23.222:2379 [-Host 压测机ip] [-Port 3000]
    
# based open source project
  ## golang
    "github.com/levigross/grequests"
    "github.com/myzhan/boomer"
    "go.etcd.io/etcd/v3/clientv3"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
  ## python
    pip安装：
    pip install grpcio
    pip install grpcio-tools
    pip install etcd3
    pip install locust==1.2.3
  ## 前端
    layui
# 部分代码摘抄和参考了网络上的文章，致谢
   go-etcd-grpc :"https://www.cnblogs.com/wujuntian/p/12838041.html"
