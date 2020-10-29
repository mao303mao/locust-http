package main

import (
	"flag"
	"fmt"
	proto "locust_http/proto"
	"locust_http/utils/boomerWrap"
	"locust_http/utils/gRpcEtcd"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/levigross/grequests"
	"github.com/myzhan/boomer"
	"go.etcd.io/etcd/v3/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var host = gRpcEtcd.GetLocalNetIpV4() // 本机ip
var ServiceName = "boomer_service"
var (
	Host        = flag.String("Host", host, "listening host")                           // 服务器Port
	Port        = flag.Int("Port", 3000, "listening port")                           // 服务器Port
	EtcdAddr    = flag.String("EtcdAddr", "127.0.0.1:2379", "register etcd address") // etcd的地址
)

// etcd客户端
var cli *clientv3.Client
// 全局boomer
var globalBoomer *boomer.Boomer
// 接受来自EndBoomer的请求处理
var quitSignal chan int=make(chan int)
// boomer运行状态
var boomerStatus = false
// 是否master关闭
var quitByClient =false
// boomer等待退出
func waitForQuit(gboomer *boomer.Boomer) {
	quitByClient = false
	boomer.Events.SubscribeOnce("boomer:quit", func() {
		defer func(){
			r:=recover();if r!=nil{
			fmt.Println("处理Boomer关闭遇到异常:",r)
		}}()
		// fmt.Println("call the boomer:quit",quitByClient)
		boomerStatus=false // 结束运行
		if !quitByClient{
			quitSignal<-1 // 释放下面EndBommer处理协程
			fmt.Println("事件订阅中获取了非Client关闭Boomer的消息")
		}
		fmt.Println("boomer服务已经关闭")
	})
	go func(){ // 此处添加通过EndBommer获取关闭boomer信号处理的代码
		<-quitSignal
		fmt.Println("从管理机client获取了关闭Boomer的消息")
		if boomerStatus{
			quitByClient=true
			gboomer.Quit()
		}
	}()
}

// rpc服务接口
type BoomerCallService struct{}

func (hes *BoomerCallService) InitBommer(ctx context.Context, req *proto.InitBommerRequest) (*proto.BoomerCallResponse, error) {
	if boomerStatus{ // 如果运行，结束（暂时不加全局锁）
		return &proto.BoomerCallResponse{
			Status:  false,
			Message: fmt.Sprintf("%s--当前boomer正在运行，需要重置后再创建新的事务", time.Now().Format("2006-01-02 15:04:05")),
		}, nil
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano()) // 每次初始化时决定随机
	baseReqOptions := &grequests.RequestOptions{ // 基本的http请求配置
		InsecureSkipVerify: true, // 忽略https证书错误
		MaxIdleConnsPerHost: 2000, // 限制连接数
		DisableKeepAlives:   false,
		RedirectLimit:       -1, // 不进行重定向
	}
	if req.GetHttpProxy()!=""{
		httpProxyUrl, err := url.Parse(fmt.Sprintf("http://%s",req.GetHttpProxy()))
		if err!=nil{
			return &proto.BoomerCallResponse{
				Status:  false,
				Message: fmt.Sprintf("%s--代理设置有误，请检查", time.Now().Format("2006-01-02 15:04:05")),
			}, nil
		}
		httpsProxyUrl, _ := url.Parse(fmt.Sprintf("http://%s",req.GetHttpProxy()))
		baseReqOptions.Proxies=map[string]*url.URL{
			"http":  httpProxyUrl,
			"https": httpsProxyUrl,
		}
	}
	storedParamValues := map[string]string{} // 全局变量存储器
	var reqSession *grequests.Session // 全局http session
	if req.IsSession {
		reqSession = grequests.NewSession(baseReqOptions)
	} else {
		reqSession = nil
	}
	if req.GetPreTask() != nil {  // 执行前置任务，主要用于更新session及存储一些全局量，比如一些登录相关的token
		err:=boomerWrap.DoPreTask(baseReqOptions, reqSession, req.GetPreTask(), storedParamValues)
		if err!=nil{
			return &proto.BoomerCallResponse{
				Status:  false,
				Message: fmt.Sprintf("%s--创建Boomer事务失败，前置步骤执行失败[%s]",
					time.Now().Format("2006-01-02 15:04:05"),err.Error()),
			}, nil
		}
	}

	taskList := []*boomer.Task{}
	// fmt.Println(req.GetLocustMaster())
	globalBoomer = boomer.NewBoomer(req.GetLocustMaster(), 5557)
	for _, testTask := range req.MainTask { // 创建boomer任务
		task := &boomer.Task{
			Name:   testTask.GetTaskName(),
			Weight: int(testTask.GetTaskWeight()),
			Fn:     boomerWrap.MakeTestTask(baseReqOptions, reqSession, testTask, globalBoomer, storedParamValues),
		}
		fmt.Println(*task)
		taskList = append(taskList, task)
	}
	globalBoomer.Run(taskList...)
	waitForQuit(globalBoomer)
	boomerStatus = true // boomer正在运行
	return &proto.BoomerCallResponse{
		Status:  true,
		Message: fmt.Sprintf("%s--已经成功获取创建Boomer事务请求", time.Now().Format("2006-01-02 15:04:05")),
	}, nil
}

func (hes *BoomerCallService) EndBommer(ctx context.Context, req *proto.EndBommerRequest) (*proto.BoomerCallResponse, error) {
	if !boomerStatus{ // 如果非运行，结束
		return &proto.BoomerCallResponse{
			Status:  false,
			Message: fmt.Sprintf("%s--当前boomer没有运行，不需要关闭", time.Now().Format("2006-01-02 15:04:05")),
		}, nil
	}
	quitSignal<-1 // 写入关闭信号
	// 此处添加关闭Boomer的信号处理
	return &proto.BoomerCallResponse{
		Status:  true,
		Message: fmt.Sprintf("%s--已经成功获取关闭Boomer信息", time.Now().Format("2006-01-02 15:04:05")),
	}, nil
}

func main() {
	flag.Parse()

	// 监听网络
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *Port))
	if err != nil {
		fmt.Println("监听网络失败：", err)
		return
	}
	defer listener.Close()

	// 创建grpc句柄
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	// 将Server结构体注册到grpc服务中
	proto.RegisterBoomerCallServiceServer(srv, &BoomerCallService{})

	// 将服务地址注册到etcd中

	cli, err = gRpcEtcd.NewEtcdClient(*EtcdAddr)
	if err != nil {
		fmt.Println("创建ETCD客户端失败：", err)
		return
	}
	// println(fmt.Sprintf("%n"),cli.ActiveConnection())

	serverAddr := fmt.Sprintf("%s:%d", *Host, *Port)
	fmt.Printf("服务地址: %s\n", serverAddr)
	gRpcEtcd.Register(cli, ServiceName, serverAddr, 5)

	// 关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		gRpcEtcd.UnRegister(cli, ServiceName, serverAddr)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	// 监听服务
	err = srv.Serve(listener)
	if err != nil {
		fmt.Println("gRPC服务的监听异常：", err)
		return
	}
}
