# locust-http
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
     【2】boomer源码的github.com\myzhan\boomer\runner.go(行221)
     中close channel某些情况化因为重复关闭会造成异常关闭, 这里最好加上recover处理：
     defer func(){
      	r:=recover();if r!=nil{
	   fmt.Println("处理Boomer关闭遇到异常:",r)
      }}()

         
## 启动参考：这里的ip、port都是例子，请根据实际情况设置
  ### 1-先下载etcd并在服务器上启动etcd：
  etcd下载地址： https://github.com/etcd-io/etcd/releases/download/v3.3.25/etcd-v3.3.25-windows-amd64.zip  
  	etcd.exe --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379
       
  ### 2-在服务器上启动matser（假定服务器ip：192.168.23.222）：
  #### python3-pip安装：
    pip install grpcio
    pip install grpcio-tools
    pip install etcd3
    pip install locust==1.2.3
  #### 启动命令：
    python3 main.py --master-host=192.168.23.222 [--step-load]
  ### 3-压力器上，可以直接从web页面的压测机管理下载编译好的执行程序（windows64，及linux64）
  压力器上，执行此程序（假定etcd的ip：192.168.23.222）
      boomerHazardServer -EtcdAddr 192.168.23.222:2379 [-Host 压力器自己的ip] [-Port 3000]
    
## 如果想要自己编译woker端，需要以下依赖库
  ### golang：go get命令，一些库可能会出现与本地冲突（etcd的），需要自己删除冲突的
    "github.com/levigross/grequests"
    "github.com/antchfx/htmlquery"
    "github.com/myzhan/boomer"
    "go.etcd.io/etcd/v3/clientv3"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/protobuf"
  ### 编译
  #### 根据操作系统(windows,linux)编译对应压测机（worker）端应用程序：
  	go build boomerHazardServer.go
  注：Windows编译linux，先在cmd执行set GOOS=linux 以及 set GOARCH=amd64
  
 
## 部分代码摘抄和参考了网络上的文章，致谢
   go-etcd-grpc :"https://www.cnblogs.com/wujuntian/p/12838041.html"
   
