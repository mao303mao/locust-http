# locust-hazard
base on locust and boomer, use etcd and gRPC to push http request and make locust test tasks。
基于locust和boomer核心，使用基于etcd做为压测机服务发现，使用gRPC推送http请求事务描述信息，让压测机自己构造http接口测试任务。
使用postman类似的节目管理http请求事务描述信息。
目前是核心且基本框架及功能的完成，如有更多的要求比如任务隔离，执行历史管理......这就是你的事了。

## 说明
拷贝了原版locust的main.py及webUI和前端部分代码进行修改。


## 启动参考：这里的ip、port都是例子，请根据实际情况设置
  1-先启动etcd（假定服务器ip：192.168.23.222）：etcd.exe --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379
  2-在启动matser（假定服务器ip：192.168.23.222）：python3 main.py --master-host=192.168.23.222
  3-编译压测机端应用程序：go build boomerHazardServer
  4-压测机端执行此程序：boomerHazardServer.exe -EtcdAddr 192.168.23.222:2379 [-Host 压测机ip] [-Port 3000]
    
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
    pip install locust==1.1.1
  ## 前端
    layui
# some code consulted article  from website
   go-etcd-grpc :"https://www.cnblogs.com/wujuntian/p/12838041.html"
