# coding:utf-8
import etcd3
import gevent
class EtcdTooler():
    etcdClient=None
    etcdAddr="127.0.0.1:2379"
    stop_flag=False
    servAddressList = []

    def __init__(self,servAddrPrefix:str,etcdAddr=""):
        self.servAddrPrefix=servAddrPrefix
        if etcdAddr=="":
            self.etcdClient=etcd3.client()
        else:
            ectdAddrReslv=etcdAddr.split(":")
            if len(ectdAddrReslv)!=2:
                self.etcdClient = etcd3.client()
            else:
                self.etcdAddr = etcdAddr
                self.etcdClient=etcd3.client(ectdAddrReslv[0],ectdAddrReslv[1])
        gevent.spawn(self.__checkHeartBeat)



    def __checkHeartBeat(self):
        print("......连接etcd服务:%s,并心跳检查......" % (self.etcdAddr))
        def __resetEtcdClient():
            self.etcdClient.close()
            ectdAddrReslv = self.etcdAddr.split(":")
            self.etcdClient=etcd3.client(ectdAddrReslv[0],ectdAddrReslv[1])
        errCount=1
        while not self.stop_flag:
            self.servAddressList.clear()
            try:
                for kv in self.etcdClient.get_prefix_response(self.servAddrPrefix).kvs:
                    self.servAddressList.append(kv.value.decode('utf-8'))
                errCount=1 # 重置错误次数
                gevent.sleep(2)  # 2s 检查一次
            except Exception as e:
                if errCount>60: # 超过5分钟
                    print(".....etcd尝试次数超过上限，退出....." )
                    exit(2)
                gevent.sleep(5)
                print(".....获取etcd的key：异常：%s。正在尝试第%d次....."%(e,errCount))
                __resetEtcdClient()
                errCount+=1


    def close(self, exc_type, exc_val, exc_tb):
        self.stop_flag=True
        self.etcdClient.close()

