# -*- coding: utf-8 -*-

import csv
import logging
import os.path
from functools import wraps
from html import escape
from io import StringIO
from itertools import chain
from time import time
import traceback
import gevent
from flask import Flask,make_response,jsonify,render_template,request,send_from_directory, send_file
from flask_basicauth import BasicAuth
from gevent import pywsgi
import datetime
from locust import __version__ as version
from locust import runners
from locust.exception import AuthCredentialsError
from locust.runners import MasterRunner
from locust.log import greenlet_exception_logger
from locust.stats import sort_stats,StatsCSV
import locust.stats as stats_module
from locust.util.cache import memoize
from locust.util.rounding import proper_round
from locust.util.timespan import parse_timespan
import boomerCall_pb2
import boomerCall_pb2_grpc
import grpc
from utils.etcdTool import EtcdTooler
from utils.requestMaker import makeInitBoomerRequest
import json
# from gevent.lock import RLock
# servListLock=RLock()

logger = logging.getLogger(__name__)
greenlet_exception_handler = greenlet_exception_logger(logger)

DEFAULT_CACHE_TIME = 2.0


def parse2Int4saveTrans(oriStr:str)->int:
    """
    字符串转数字
    :param oriStr:
    :return:
    """
    if not oriStr.isnumeric():
        return 0
    try:
        return int(oriStr)
    except:
        return 0

class WebUI:
    """
    Sets up and runs a Flask web app that can start and stop load tests using the
    :attr:`environment.runner <locust.env.Environment.runner>` as well as show the load test statistics
    in :attr:`environment.stats <locust.env.Environment.stats>`
    """

    app = None
    greenlet = None
    server = None
    etcdt = None
    reporter_running_status=False
    """Reference to the :class:`pyqsgi.WSGIServer` instance"""

    def  __init__(self,environment,host,port,masterHost,auth_credentials=None,tls_cert=None,tls_key=None,stats_csv_writer=None):
        """
        Create WebUI instance and start running the web server in a separate greenlet (self.greenlet)

        Arguments:
        environment: Reference to the curren Locust Environment
        host: Host/interface that the web server should accept connections to
        port: Port that the web server should listen to
        auth_credentials:  If provided, it will enable basic auth with all the routes protected by default.
                           Should be supplied in the format: "user:pass".
        tls_cert: A path to a TLS certificate
        tls_key: A path to a TLS private key
        """
        environment.web_ui = self
        self.stats_csv_writer = stats_csv_writer or StatsCSV(environment,stats_module.PERCENTILES_TO_REPORT)
        self.environment = environment
        self.host = host
        self.port = port
        self.tls_cert = tls_cert
        self.tls_key = tls_key
        app = Flask(__name__)
        self.app = app
        app.debug = True
        app.root_path = os.path.dirname(os.path.abspath(__file__))
        self.app.config["BASIC_AUTH_ENABLED"] = False
        self.auth = None
        self.greenlet = None
        self.masterHost = masterHost
        self.etcdt=EtcdTooler("/ns/boomer_service/")
        self.workedServser={}
        self.recvMesg={}

        if auth_credentials is not None:
            credentials = auth_credentials.split(':')
            if len(credentials) == 2:
                self.app.config["BASIC_AUTH_USERNAME"] = credentials[0]
                self.app.config["BASIC_AUTH_PASSWORD"] = credentials[1]
                self.app.config["BASIC_AUTH_ENABLED"] = True
                self.auth = BasicAuth()
                self.auth.init_app(self.app)
            else:
                raise AuthCredentialsError(
                    "Invalid auth_credentials. It should be a string in the following format: 'user.pass'")

        def stats_history(self):
            """Save current stats info to history for charts of report."""
            while self.reporter_running_status:
                stats = environment.runner.stats
                if not stats.total.use_response_times_cache:
                    break
                r = {
                    "time":                        datetime.datetime.now().strftime("%H:%M:%S"),
                    "current_rps":                 stats.total.current_rps or 0,
                    "current_fail_per_sec":        stats.total.current_fail_per_sec or 0,
                    "response_time_percentile_95": stats.total.get_current_response_time_percentile(0.95) or 0,
                    "response_time_percentile_50": stats.total.get_current_response_time_percentile(0.5) or 0,
                    "user_count":                  environment.runner.user_count or 0,
                }
                stats.history.append(r)
                gevent.sleep(5)
            print("结束了指标历史记录任务....")

        def initTask(self,rpcServAddr):
            """
            调用gRPC，通知初始化Boomer，改成同步处理
            :param rpcServAddr:
            :return:
            """
            channel = grpc.insecure_channel(rpcServAddr)  # 连接 rpc 服务器
            # 调用 rpc 服务
            stub = boomerCall_pb2_grpc.BoomerCallServiceStub(channel)
            try:
                initBommerRequest=makeInitBoomerRequest("jsons/tmp.json",self.masterHost)
            except Exception as e:
                print(e)
                return
            beforeClientIds=[ worker.id  for worker in  environment.runner.clients.values()]
            response = stub.InitBommer(initBommerRequest)
            self.recvMesg[rpcServAddr] = response.message
            if response.status:
                tryCount=0
                newId=""
                while tryCount<=50: # 5s的超时
                    afterClientIds = [worker.id for worker in environment.runner.clients.values()]
                    dfIdSet = set(afterClientIds).difference(set(beforeClientIds))
                    if not dfIdSet:
                        gevent.sleep(0.1)
                        tryCount+=1
                    else:
                        newId=dfIdSet.pop()
                        break
                self.workedServser[rpcServAddr]=newId
            try:
                channel.close()
            except:
                pass

        def shutTask(self,rpcServAddr):
            """
            调用gRPC，通知关闭Boomer
            :param rpcServAddr:
            :return:
            """
            channel = grpc.insecure_channel(rpcServAddr)  # 连接 rpc 服务器
            # 调用 rpc 服务
            stub = boomerCall_pb2_grpc.BoomerCallServiceStub(channel)
            response = stub.EndBommer(boomerCall_pb2.EndBommerRequest())
            print("client received: %s, from %s" % (response.message,rpcServAddr))
            self.recvMesg[rpcServAddr] = response.message
            if response.status:
                if self.workedServser.__contains__(rpcServAddr):
                    del self.workedServser[rpcServAddr]
                if self.recvMesg.__contains__(rpcServAddr):
                    # 成功情况下，不需要获取消息
                    del self.recvMesg[rpcServAddr]
            try:
                channel.close()
            except:
                pass

        @app.route('/')
        @self.auth_required_if_enabled
        def index():
            if not environment.runner:
                return make_response("Error: Locust Environment does not have any runner",500)

            is_distributed = isinstance(environment.runner,MasterRunner)
            if is_distributed:
                worker_count = environment.runner.worker_count
                slave_count = len(self.etcdt.servAddressList)
            else:
                worker_count = 0
                slave_count = 0

            override_host_warning = False
            if environment.host:
                host = environment.host
            elif environment.runner.user_classes:
                all_hosts = set([l.host for l in environment.runner.user_classes])
                if len(all_hosts) == 1:
                    host = list(all_hosts)[0]
                else:
                    # since we have mulitple User classes with different host attributes, we'll
                    # inform that specifying host will override the host for all User classes
                    override_host_warning = True
                    host = None
            else:
                host = None

            options = environment.parsed_options
            return render_template("index.html",
                                   state=environment.runner.state,
                                   is_distributed=is_distributed,
                                   user_count=environment.runner.user_count,
                                   version=version,
                                   host=host,
                                   override_host_warning=override_host_warning,
                                   num_users=options and options.num_users,
                                   hatch_rate=options and options.hatch_rate,
                                   step_users=options and options.step_users,
                                   step_time=options and options.step_time,
                                   worker_count=worker_count,
                                   slave_count= slave_count,
                                   is_step_load=environment.step_load,
                                   )




        @app.route('/swarm',methods=["POST"])
        @self.auth_required_if_enabled
        def swarm():
            if environment.runner.state not in (runners.STATE_STOPPED,runners.STATE_INIT):
                return jsonify({'success': False,'message': '当前有任务正在执行，先停止测试再尝试'})
            if not  environment.runner.worker_count:
                 return jsonify({'success': False,'message': '没有可用的work, 不能运行测试'})
            assert request.method == "POST"
            # 开启协程写入指标历史记录
            self.reporter_running_status=True
            gevent.spawn(stats_history,self)
            # 开始压测任务
            user_count = int(request.form["user_count"])
            hatch_rate = float(request.form["hatch_rate"])
            run_seconds=None
            if request.form["run_time"]:
                try:
                    run_seconds = parse_timespan(request.form["run_time"])
                except ValueError:
                    run_seconds = None
            if (request.form.get("host")):
                environment.host = str(request.form["host"])
            if environment.step_load:
                step_user_count = int(request.form["step_user_count"])
                step_duration = parse_timespan(str(request.form["step_duration"]))
                environment.runner.start_stepload(user_count,hatch_rate,step_user_count,step_duration)
                return jsonify(
                    {'success': True,'message': 'Swarming started in Step Load Mode','host': environment.host})

            def stopRunAfterSecs(x):
                if x>3600*6:
                    x=3600*6
                count=0
                while count<=x: # 这种方式避免长时间休眠造成问题
                    gevent.sleep(1)
                    count+=1
                self.reporter_running_status = False  # 结束指标历史记录
                environment.runner.stop()
            if run_seconds and run_seconds>=30:
                gevent.spawn(stopRunAfterSecs,run_seconds)
            environment.runner.start(user_count,hatch_rate)
            return jsonify({'success': True,'message': 'Swarming started','host': environment.host})

        @app.route('/stop')
        @self.auth_required_if_enabled
        def stop():
            self.reporter_running_status=False #结束指标历史记录
            if environment.runner.state not in (runners.STATE_STOPPING,runners.STATE_STOPPED,runners.STATE_INIT):
                environment.runner.stop()
            return jsonify({'success': True,'message': 'Test stopped'})

        @app.route("/stats/reset")
        @self.auth_required_if_enabled
        def reset_stats():
            environment.runner.stats.reset_all()
            environment.runner.exceptions = {}
            return "ok"

        def makeHttpItem(request,prefixStr,i):
            """
            构造任务
            :return:
            """
            xPretask = {}
            pretask_x_method = request.form.get("%s-%d-method" % (prefixStr,i),type=str,default="Get")
            xPretask["Method"] = pretask_x_method
            pretask_x_urlPath = request.form.get("%s-%d-urlPath" % (prefixStr,i),type=str,default="")
            xPretask["UrlPath"] = pretask_x_urlPath
            # Headers 处理
            pretask_x_headers_key = request.form.getlist("%s-%d-headers-key" % (prefixStr,i))
            pretask_x_headers_value = request.form.getlist("%s-%d-headers-value" % (prefixStr,i))
            if pretask_x_headers_key and pretask_x_headers_value and \
                    len(pretask_x_headers_key) == len(pretask_x_headers_value):
                xPretask["Headers"]={}
                for j in range(len(pretask_x_headers_key)):
                    xPretask["Headers"][pretask_x_headers_key[j]] = pretask_x_headers_value[j]
            else:
                xPretask["Headers"] = None
            # Params 处理
            pretask_x_params_key = request.form.getlist("%s-%d-params-key" % (prefixStr,i))
            pretask_x_params_value = request.form.getlist("%s-%d-params-value" % (prefixStr,i))
            if pretask_x_params_key and pretask_x_params_value and \
                    len(pretask_x_params_key) == len(pretask_x_params_value):
                xPretask["Params"]={}
                for j in range(len(pretask_x_params_key)):
                    xPretask["Params"][pretask_x_params_key[j]] = pretask_x_params_value[j]
            else:
                xPretask["Params"] = None
            # DictData 处理
            pretask_x_dictdata_key = request.form.getlist("%s-%d-dictdata-key" % (prefixStr,i))
            pretask_x_dictdata_value = request.form.getlist("%s-%d-dictdata-value" % (prefixStr,i))
            if pretask_x_dictdata_key and pretask_x_dictdata_value and \
                    len(pretask_x_dictdata_key) == len(pretask_x_dictdata_value):
                xPretask["DictData"]={}
                for j in range(len(pretask_x_dictdata_key)):
                    xPretask["DictData"][pretask_x_dictdata_key[j]] = pretask_x_dictdata_value[j]
            else:
                xPretask["DictData"] = None
            # RawData 处理
            pretask_x_rawdata = request.form.get("%s-%d-rawdata" % (prefixStr,i),type=str,default="")
            xPretask["RawData"] = pretask_x_rawdata
            # JsonData 处理
            pretask_x_jsondata = request.form.get("%s-%d-jsondata" % (prefixStr,i),type=str,default="")
            xPretask["JsonData"] = pretask_x_jsondata
            # AssertChain 处理
            pretask_x_assertType = request.form.getlist("%s-%d-assertType" % (prefixStr,i))
            pretask_x_assertValue = request.form.getlist("%s-%d-assertValue" % (prefixStr,i))
            if pretask_x_assertType and pretask_x_assertValue and \
                    len(pretask_x_assertType) == len(pretask_x_assertValue):
                xPretask["AssertChain"] = []
                for j in range(len(pretask_x_assertType)):
                    xPretask["AssertChain"].append({
                        "AssertType": parse2Int4saveTrans(pretask_x_assertType[j]),
                        "RuleValue":  parse2Int4saveTrans(pretask_x_assertValue[j])
                    })
            else:
                xPretask["AssertChain"] = None
            # SaveParamAction 处理
            pretask_x_saveType = request.form.getlist("%s-%d-saveType" % (prefixStr,i))
            pretask_x_paramName = request.form.getlist("%s-%d-paramName" % (prefixStr,i))
            pretask_x_ruleValue = request.form.getlist("%s-%d-ruleValue" % (prefixStr,i))
            if pretask_x_saveType and pretask_x_paramName and pretask_x_ruleValue and \
                    len(pretask_x_saveType) == len(pretask_x_paramName) and \
                    len(pretask_x_saveType) == len(pretask_x_ruleValue):
                xPretask["SaveParamChain"] = []
                for j in range(len(pretask_x_saveType)):
                    xPretask["SaveParamChain"].append({
                        "SaveType":  parse2Int4saveTrans(pretask_x_saveType[j]),
                        "ParamName": pretask_x_paramName[j],
                        "RuleValue": pretask_x_ruleValue[j]
                    })
            else:
                xPretask["SaveParamChain"] = None
            return xPretask

        @app.route('/saveTrans',methods=["POST"])
        @self.auth_required_if_enabled
        def saveTransation():
            if environment.runner.state not in (runners.STATE_STOPPED,runners.STATE_INIT):
                return jsonify({'success': False,'message': '当前有任务正在执行，先停止测试再尝试'})
            transObj={}
            isSession=request.form.get("isSession",type=int,default=0)
            transObj["isSession"]=isSession==1
            httpProxy=request.form.get("HttpProxy",type=str,default="")
            transObj["HttpProxy"] = httpProxy
            preTaskMark=request.form.get("PreTaskMark",type=int,default=0)
            # PreTask 处理
            if preTaskMark==0:
                transObj["PreTask"] = None
            else:
                preTask=[]
                for i in range(1,preTaskMark+1):
                    preTask.append(makeHttpItem(request,"pretask",i))
                transObj["PreTask"]=preTask
            # TestTask 处理
            testTask=[]
            TestTaskIdMark = request.form.get("TestTaskIdMark",type=str,default="1")
            TestTaskMark=request.form.get("TestTaskMark",type=str,default="0")
            if(len(TestTaskIdMark.split(","))!=len(TestTaskMark.split(",")) or len(TestTaskIdMark.split(","))<1):
                return jsonify({'success': False,'message': '测试事务的id号出现错乱，可能需要重新来！'})
            for i,ttm in enumerate(TestTaskMark.split(",")):
                taskId=parse2Int4saveTrans(TestTaskIdMark.split(",")[i])
                if taskId==0:
                    return jsonify({'success': False,'message': '测试事务的id号出现错误，可能需要重新来！'})
                xTestTask={}
                xtm=parse2Int4saveTrans(ttm)
                xTaskWeight=request.form.get("testtask-%d-taskWeight"%(taskId),type=int,default=0)
                xTestTask["TaskWeight"]=xTaskWeight
                xTaskName=request.form.get("testtask-%d-TaskName"%(taskId),type=str,default="测试事务%d"%(i+1))
                xTestTask["TaskName"] = xTaskName
                if xtm==0:
                    xTestTask["PreWork"]=None
                else:
                    preWork=[]
                    for j in range(1,xtm + 1):
                        preWork.append(makeHttpItem(request,"testtask-%d-prework"%(taskId),j))
                    xTestTask["PreWork"]=preWork
                xTestTask["TestWork"]=makeHttpItem(request,"testwork",taskId)
                testTask.append(xTestTask)
            transObj["MainTask"]=testTask
            try:
                with open("./jsons/main.json",mode="w+") as f:
                    f.write(json.dumps(transObj,indent=4))
                return jsonify({'success': True,'message': '保存成功'})
            except Exception as e:
                return jsonify({'success': False,'message': e})

        @app.route('/importTrans',methods=['Post'])
        @self.auth_required_if_enabled
        def importTrans():
            try:
                f=request.files["file"]
                f.save("./jsons/tmp.json")
                return jsonify({'success': True,'message': ''})
            except Exception as e:
                return jsonify({'success': False,'message': str(e)+"\n"+traceback.format_exc()})

        @app.route('/transaction',methods=["GET"])
        @self.auth_required_if_enabled
        def transation():
            return render_template("transaction.html")

        @app.route('/importedTrans',methods=['Get'])
        @self.auth_required_if_enabled
        def importedTrans():
            try:
                with open("./jsons/tmp.json",mode='rb') as f:
                    transation=json.load(f)
                    if not transation.get("PreTask"):
                        PreTaskMark=0
                    else:
                        PreTaskMark=len(transation.get("PreTask"))
                    if not transation.get("MainTask"):
                        print("json中没有MainTask")
                        return render_template("transaction.html")
                    TestTaskMark=[]
                    TestTaskIdMark=[]
                    tid=0
                    for mt in transation.get("MainTask"):
                        TestTaskIdMark.append(tid+1)
                        tid+=1
                        if not mt.get("PreWork"):
                            TestTaskMark.append(0)
                        else:
                            TestTaskMark.append(len(mt.get("PreWork")))
                    transMark = {
                        "PreTaskMark":  PreTaskMark,
                        "TestTaskMark": TestTaskMark,
                        "TestTaskIdMark": TestTaskIdMark,
                        "TestTaskId":tid
                    }
                    return render_template("importedTransaction.html",transMark=transMark,transation=transation)
            except Exception as e:
                print(e)
                return render_template("transaction.html")


        @app.route('/backupTrans',methods=['POST'])
        @self.auth_required_if_enabled
        def backupTrans():
            transObj = {}
            isSession = request.form.get("isSession",type=int,default=0)
            transObj["isSession"] = isSession == 1
            httpProxy = request.form.get("HttpProxy",type=str,default="")
            transObj["HttpProxy"] = httpProxy
            preTaskMark = request.form.get("PreTaskMark",type=int,default=0)
            # PreTask 处理
            if preTaskMark == 0:
                transObj["PreTask"] = None
            else:
                preTask = []
                for i in range(1,preTaskMark + 1):
                    preTask.append(makeHttpItem(request,"pretask",i))
                transObj["PreTask"] = preTask
            # TestTask 处理
            testTask = []
            TestTaskIdMark = request.form.get("TestTaskIdMark",type=str,default="1")
            TestTaskMark = request.form.get("TestTaskMark",type=str,default="0")
            if (len(TestTaskIdMark.split(",")) != len(TestTaskMark.split(",")) or len(TestTaskIdMark.split(",")) < 1):
                return render_template("transaction.html")
            for i,ttm in enumerate(TestTaskMark.split(",")):
                taskId = parse2Int4saveTrans(TestTaskIdMark.split(",")[i])
                if taskId == 0:
                    return render_template("transaction.html")
                xTestTask = {}
                xtm = parse2Int4saveTrans(ttm)
                xTaskWeight = request.form.get("testtask-%d-taskWeight" % (taskId),type=int,default=0)
                xTestTask["TaskWeight"] = xTaskWeight
                xTaskName = request.form.get("testtask-%d-TaskName" % (taskId),type=str,
                                             default="测试事务%d" % (i + 1))
                xTestTask["TaskName"] = xTaskName
                if xtm == 0:
                    xTestTask["PreWork"] = None
                else:
                    preWork = []
                    for j in range(1,xtm + 1):
                        preWork.append(makeHttpItem(request,"testtask-%d-prework" % (taskId),j))
                    xTestTask["PreWork"] = preWork
                xTestTask["TestWork"] = makeHttpItem(request,"testwork",taskId)
                testTask.append(xTestTask)
            transObj["MainTask"] = testTask
            rawData = StringIO()
            rawData.write(json.dumps(transObj,indent=4))
            rawData.seek(0)
            response = make_response(rawData.getvalue())
            rawData.close()
            file_name = "backup_{0}.json".format(time())
            disposition = "attachment;filename={0}".format(file_name)
            response.headers["Content-type"] = "text/json"
            response.headers["Content-disposition"] = disposition
            return response


        @app.route('/download_boomer',methods=['GET'])
        @self.auth_required_if_enabled
        def downloadBoomer():
            if os.path.isfile(os.path.join("./slaveEXE","boomerHazardServer.exe")):
                return send_from_directory("./slaveEXE","boomerHazardServer.exe",as_attachment=True)

        @app.route('/download_boomer_linux',methods=['GET'])
        @self.auth_required_if_enabled
        def downloadBoomerLinux():
            if os.path.isfile(os.path.join("./slaveEXE","boomerHazardServer")):
                return send_from_directory("./slaveEXE","boomerHazardServer",as_attachment=True)


        @app.route('/initBoomer',methods=["POST"])
        @self.auth_required_if_enabled
        def initBoomer():
            if not self.etcdt.servAddressList:
                return jsonify({'success': False,'message': '没有可用的压力机'})
            selectServAddrList = request.form.getlist("servAddr[]")
            if not selectServAddrList:
                return jsonify({'success': False,'message': '请选择压力机'})
            for rpcServAddr in set(self.etcdt.servAddressList).intersection(set(selectServAddrList)):
                initTask(self,rpcServAddr)
            return jsonify({'success': True,'message': '已通知压测机初始化，请检查Workers中各压力机的最新消息'})

        @app.route('/shutdownBoomer',methods=["POST"])
        @self.auth_required_if_enabled
        def shutdownBoomer():
            if  environment.runner.state  not in (runners.STATE_STOPPED,runners.STATE_INIT):
                return jsonify({'success': False,'message': '当前有任务正在执行，先停止测试再尝试'})
            if not self.etcdt.servAddressList:
                return jsonify({'success': False,'message': '没有可用的压力机'})
            for rpcServAddr in self.etcdt.servAddressList: # 全部关闭
                shutTask(self,rpcServAddr)
            return jsonify({'success': True,'message': '已通知压测机停止，请检查Workers中各压力机的最新情况'})


        def _download_csv_suggest_file_name(suggest_filename_prefix):
            """Generate csv file download attachment filename suggestion.

            Arguments:
            suggest_filename_prefix: Prefix of the filename to suggest for saving the download. Will be appended with timestamp.
            """

            return f"{suggest_filename_prefix}_{time()}.csv"

        def _download_csv_response(csv_data, filename_prefix):
            """Generate csv file download response with 'csv_data'.

            Arguments:
            csv_data: CSV header and data rows.
            filename_prefix: Prefix of the filename to suggest for saving the download. Will be appended with timestamp.
            """

            response = make_response(csv_data)
            response.headers["Content-type"] = "text/csv"
            response.headers[
                "Content-disposition"
            ] = f"attachment;filename={_download_csv_suggest_file_name(filename_prefix)}"
            return response

        @app.route("/stats/requests/csv")
        @self.auth_required_if_enabled
        def request_stats_csv():
            data = StringIO()
            writer = csv.writer(data)
            self.stats_csv_writer.requests_csv(writer)
            return _download_csv_response(data.getvalue(),"requests")

        @app.route("/stats/requests_full_history/csv")
        @self.auth_required_if_enabled
        def request_stats_full_history_csv():
            options = self.environment.parsed_options
            if options and options.stats_history_enabled:
                return send_file(
                    os.path.abspath(self.stats_csv_writer.stats_history_file_name()),
                    mimetype="text/csv",
                    as_attachment=True,
                    attachment_filename=_download_csv_suggest_file_name("requests_full_history"),
                    add_etags=True,
                    cache_timeout=None,
                    conditional=True,
                    last_modified=None,
                )

            return make_response("Error: Server was not started with option to generate full history.",404)

        @app.route("/stats/failures/csv")
        @self.auth_required_if_enabled
        def failures_stats_csv():
            data = StringIO()
            writer = csv.writer(data)
            self.stats_csv_writer.failures_csv(writer)
            return _download_csv_response(data.getvalue(),"failures")

        @app.route('/stats/requests')
        @self.auth_required_if_enabled
        @memoize(timeout=DEFAULT_CACHE_TIME,dynamic_timeout=True)
        def request_stats():
            stats = []

            for s in chain(sort_stats(self.environment.runner.stats.entries),[environment.runner.stats.total]):
                stats.append({
                    "method":                  s.method,
                    "name":                    s.name,
                    "safe_name":               escape(s.name,quote=False),
                    "num_requests":            s.num_requests,
                    "num_failures":            s.num_failures,
                    "avg_response_time":       s.avg_response_time,
                    "min_response_time":       0 if s.min_response_time is None else proper_round(s.min_response_time),
                    "max_response_time":       proper_round(s.max_response_time),
                    "current_rps":             s.current_rps,
                    "current_fail_per_sec":    s.current_fail_per_sec,
                    "median_response_time":    s.median_response_time,
                    "ninetieth_response_time": s.get_response_time_percentile(0.9),
                    "avg_content_length":      s.avg_content_length,
                })

            errors = []
            for e in environment.runner.errors.values():
                err_dict = e.to_dict()
                err_dict["name"] = escape(err_dict["name"])
                err_dict["error"] = escape(err_dict["error"])
                errors.append(err_dict)

            # Truncate the total number of stats and errors displayed since a large number of rows will cause the app
            # to render extremely slowly. Aggregate stats should be preserved.
            report = {"stats": stats[:500],"errors": errors[:500]}
            if len(stats) > 500:
                report["stats"] += [stats[-1]]

            if stats:
                report["total_rps"] = stats[len(stats) - 1]["current_rps"]
                report["fail_ratio"] = environment.runner.stats.total.fail_ratio
                report[
                    "current_response_time_percentile_95"] = environment.runner.stats.total.get_current_response_time_percentile(
                    0.95)
                report[
                    "current_response_time_percentile_50"] = environment.runner.stats.total.get_current_response_time_percentile(
                    0.5)

            is_distributed = isinstance(environment.runner,MasterRunner)
            if is_distributed:
                workers = []
                for worker in environment.runner.clients.values():
                    worker.state!='missing' and workers.append({"id":        worker.id,"state": worker.state,"user_count": worker.user_count,
                                    "cpu_usage": worker.cpu_usage})

                report["workers"] = workers
                report["slaves"] = [{
                    "slave":x,
                    "clientId": (lambda t:"-" if not t or not t in [w.get("id") for w in workers] else t)(self.workedServser.get(x)),
                    "rectMsg": (lambda t: "" if not t else t)(self.recvMesg.get(x))
                } for x in self.etcdt.servAddressList ]
            # print("environment.runner.state",environment.runner.state)
            report["state"] = environment.runner.state
            report["user_count"] = environment.runner.user_count

            return jsonify(report)

        @app.route("/exceptions")
        @self.auth_required_if_enabled
        def exceptions():
            return jsonify({
                'exceptions': [
                    {
                        "count":     row["count"],
                        "msg":       row["msg"],
                        "traceback": row["traceback"],
                        "nodes":     ", ".join(row["nodes"])
                    } for row in environment.runner.exceptions.values()
                ]
            })

        @app.route("/exceptions/csv")
        @self.auth_required_if_enabled
        def exceptions_csv():
            data = StringIO()
            writer = csv.writer(data)
            writer.writerow(["Count","Message","Traceback","Nodes"])
            for exc in environment.runner.exceptions.values():
                nodes = ", ".join(exc["nodes"])
                writer.writerow([exc["count"],exc["msg"],exc["traceback"],nodes])

            response = make_response(data.getvalue())
            file_name = "exceptions_{0}.csv".format(time())
            disposition = "attachment;filename={0}".format(file_name)
            response.headers["Content-type"] = "text/csv"
            response.headers["Content-disposition"] = disposition
            return response

        # start the web server
        self.greenlet = gevent.spawn(self.start)
        self.greenlet.link_exception(greenlet_exception_handler)

        @app.route("/stats/report")
        @self.auth_required_if_enabled
        def stats_report():
            stats = self.environment.runner.stats
            if not stats or not stats.start_time or not stats.last_request_timestamp or not stats.entries:
                return  render_template(
                "report.html")

            start_ts = stats.start_time
            start_time = datetime.datetime.fromtimestamp(start_ts)
            start_time = start_time.strftime("%Y-%m-%d %H:%M:%S")

            end_ts = stats.last_request_timestamp
            end_time = datetime.datetime.fromtimestamp(end_ts)
            end_time = end_time.strftime("%Y-%m-%d %H:%M:%S")

            host = None
            if environment.host:
                host = environment.host
            elif environment.runner.user_classes:
                all_hosts = set([l.host for l in environment.runner.user_classes])
                if len(all_hosts) == 1:
                    host = list(all_hosts)[0]

            requests_statistics = list(chain(sort_stats(stats.entries),[stats.total]))
            failures_statistics = sort_stats(stats.errors)
            exceptions_statistics = []
            for exc in environment.runner.exceptions.values():
                exc["nodes"] = ", ".join(exc["nodes"])
                exceptions_statistics.append(exc)

            history = stats.history

            static_js = ""
            js_files = ["jquery-1.11.3.min.js","echarts.common.min.js","vintage.js","chart.js"]
            for js_file in js_files:
                path = os.path.join(os.path.dirname(__file__),"static",js_file)
                with open(path,encoding="utf8") as f:
                    content = f.read()
                static_js += "// " + js_file + "\n"
                static_js += content
                static_js += "\n\n\n"

            res = render_template(
                "report.html",
                int=int,
                round=round,
                requests_statistics=requests_statistics,
                failures_statistics=failures_statistics,
                exceptions_statistics=exceptions_statistics,
                start_time=start_time,
                end_time=end_time,
                host=host,
                history=history,
                static_js=static_js,
            )
            if request.args.get("download"):
                res = app.make_response(res)
                res.headers["Content-Disposition"] = "attachment;filename=report_%s.html" % time()
            return res

    def start(self):
        if self.tls_cert and self.tls_key:
            self.server = pywsgi.WSGIServer((self.host,self.port),self.app,log=None,keyfile=self.tls_key,
                                            certfile=self.tls_cert)
        else:
            self.server = pywsgi.WSGIServer((self.host,self.port),self.app,log=None)
        self.server.serve_forever()

    def stop(self):
        """
        Stop the running web server
        """
        self.server.stop()

    def auth_required_if_enabled(self,view_func):
        """
        Decorator that can be used on custom route methods that will turn on Basic Auth
        authentication if the ``--web-auth`` flag is used. Example::

            @web_ui.app.route("/my_custom_route")
            @web_ui.auth_required_if_enabled
            def my_custom_route():
                return "custom response"
        """

        @wraps(view_func)
        def wrapper(*args,**kwargs):
            if self.app.config["BASIC_AUTH_ENABLED"]:
                if self.auth.authenticate():
                    return view_func(*args,**kwargs)
                else:
                    return self.auth.challenge()
            else:
                return view_func(*args,**kwargs)

        return wrapper
