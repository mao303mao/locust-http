package gRpcEtcd

// 此处选摘“疯一样的狼人”的文档
// https://www.cnblogs.com/wujuntian/p/12838041.html
import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
)

const schema = "ns"

// 创建Etcd客户端
func NewEtcdClient(etcdAddr string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ";"),
		DialTimeout: 5 * time.Second,
	})

}

func GetLocalNetIpV4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return "127.0.0.1"
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.Index(ipnet.IP.String(), "192.168.") == 0 { // 优先C段私有地址
					return ipnet.IP.String()
				}
				if strings.Index(ipnet.IP.String(), "10.") == 0 { // A段私有地址
					return ipnet.IP.String()
				}
				if strings.Index(ipnet.IP.String(), "172.") == 0 { // B段私有地址
					ipv4number := strings.Split(ipnet.IP.String(), ".")
					ipv4_1, _ := strconv.Atoi(ipv4number[1])
					if ipv4_1 >= 16 && ipv4_1 <= 31 {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "127.0.0.1"
}

// 将服务地址注册到etcd中
func Register(cli *clientv3.Client, serviceName, serverAddr string, ttl int64) error {

	// 进行心跳检测
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		key := "/" + schema + "/" + serviceName + "/" + serverAddr
		log.Println("注册boomerhazard服务",key)
		for {
			resp, err := cli.Get(context.Background(), key)
			// fmt.Printf("resp:%+v\n", resp)
			if err != nil {
				log.Printf("获取服务地址失败：%s\n", err)
			} else if resp.Count == 0 { // 未注册
				err = KeepAlive(cli, serviceName, serverAddr, ttl)
				if err != nil {
					log.Printf("保持连接失败：%s\n", err)
				}
			}
			<-ticker.C
		}
	}()

	return nil
}

// 保持服务器与etcd的长连接
func KeepAlive(cli *clientv3.Client, serviceName, serverAddr string, ttl int64) error {
	// 创建租约
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		log.Printf("创建租期失败：%s\n", err)
		return err
	}

	// 将服务地址注册到etcd中
	key := "/" + schema + "/" + serviceName + "/" + serverAddr
	_, err = cli.Put(context.Background(), key, serverAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("注册服务失败：%s\n", err)
		return err
	}

	// 建立长连接
	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("建立长连接失败：%s\n", err)
		return err
	}

	// 清空keepAlive返回的channel
	go func() {
		for {
			<-ch
		}
	}()
	return nil
}

// 取消注册
func UnRegister(cli *clientv3.Client, serviceName, serverAddr string) {
	if cli != nil {
		key := "/" + schema + "/" + serviceName + "/" + serverAddr
		cli.Delete(context.Background(), key)
		cli.Close()
	}
}

// 获取分布式锁，本应用用不到
func NewTaskMutex(cli *clientv3.Client, serviceName string) (*concurrency.Session, *concurrency.Mutex, error) {
	prefix := "/myLock/" + serviceName
	mutexSession, err := concurrency.NewSession(cli)
	if err != nil {
		return nil, nil, err
	}
	return mutexSession, concurrency.NewMutex(mutexSession, prefix), nil
}

// 获取key值
func GetKey(cli *clientv3.Client, keyName string) []map[string]string {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := cli.Get(ctx, keyName)
	cancelFunc()
	if err != nil {
		return nil
	}
	var kvMap []map[string]string
	for k, v := range res.Kvs {
		log.Println("获取到Kvs的值", k, string(v.Key), string(v.Value))
		kvMap = append(kvMap, map[string]string{"key": string(v.Key), "value": string(v.Value)})
	}
	return kvMap
}

// 设置key值
func SetKey(cli *clientv3.Client, keyName string, keyValue string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := cli.Put(ctx, keyName, keyValue)
	cancelFunc()
	return err
}
