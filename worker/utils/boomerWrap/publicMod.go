package boomerWrap

import (
	"encoding/json"
	"fmt"
	"io"
	proto "locust_hazard/proto"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/kataras/iris/core/errors"
	"github.com/levigross/grequests"
	"github.com/myzhan/boomer"
)

// 执行判断，因为是压力测试，所以判断就简单粗暴一点:状态码或者内容长度
func DoAssertActionChains(resp *grequests.Response, assertActionChain []*proto.InitBommerRequest_AssertAction) (bool, string) {
	respLength := len(resp.Bytes())
	for _, assertAction := range assertActionChain {
		switch assertAction.AssertType {
		case 0:
			{
				if resp.StatusCode != int(assertAction.RuleValue) {
					if resp.StatusCode>=300 && resp.StatusCode<400{
						return false, fmt.Sprintf("请求判断有误,响应码是[%d],期望[%d]" +
							"重定向到[%s]", resp.StatusCode, assertAction.RuleValue,resp.Header.Get("Location"))
					}
					return false, fmt.Sprintf("请求判断有误,响应码是[%d],期望[%d]", resp.StatusCode, assertAction.RuleValue)
				}
				continue
			}
		case 1: // 长度大于
			{
				if respLength <= int(assertAction.RuleValue) {
					return false, fmt.Sprintf("请求判断有误，响应内容长度是[%d],期望大于[%d]", respLength, assertAction.RuleValue)
				}
				continue
			}
		case 2: // 长度等于
			{
				if respLength != int(assertAction.RuleValue) {
					return false, fmt.Sprintf("请求判断有误，响应内容长度是[%d],期望等于[%d]", respLength, assertAction.RuleValue)
				}
				continue
			}
		case 3: // 长度小于
			{
				if respLength >= int(assertAction.RuleValue) {
					return false, fmt.Sprintf("请求判断有误，响应内容长度是[%d],期望小于[%d]", respLength, assertAction.RuleValue)
				}
				continue
			}
		default:
			continue

		}

	}
	return true, ""
}

// 根据json解析语句获取json中想要的字符串，idx需传递一个初始0的参数
func resolve2JosnObj(objI interface{}, savString *string, quotes []string, idx int) {
	if idx >= len(quotes) {
		*savString = fmt.Sprintf("%v", objI)
		return
	}
	switch obj := objI.(type) { // 此处interface{}.(type) 专门用于switch的类型判断
	case string:
		*savString = obj
	case float64:
		*savString = fmt.Sprintf("%v", obj)
	case map[string]interface{}:
		resolve2JosnObj(obj[quotes[idx]], savString, quotes, idx+1)
	case []interface{}:
		i, err := strconv.ParseInt(quotes[idx], 10, 64)
		if err != nil {
			*savString = ""
			return
		}
		if int(i) >= len(obj) {
			*savString = ""
			return
		}
		resolve2JosnObj(obj[i], savString, quotes, idx+1)

	default:
		*savString = fmt.Sprintf("%v", obj)
	}

}

// 执行参数存储
func DoSaveParamActionChains(resp *grequests.Response, savParamsChains []*proto.InitBommerRequest_SaveParamAction, storedParamValues map[string]string) {

	for _, savAction := range savParamsChains {
		switch savAction.GetSaveType() {
		case 0 :
			{ // 响应是html,使用xpath获取想要的内容
				if resp==nil{
					continue
				}
				root, err := htmlquery.Parse(strings.NewReader(resp.String()))
				if err != nil {
					continue
				}
				// fmt.Println(root)
				aimNode := htmlquery.Find(root, savAction.GetRuleValue())
				if len(aimNode) > 0 {
					storedParamValues[savAction.GetParamName()] = htmlquery.InnerText(aimNode[0])
				}
			}
		case 1:
			{ // 响应是json，使用json解析语句去获取想要的内容，比如:"dataList.2.name",
			  // 表示根object中key为“dataList”的列表中获取第“3”个object的key为“name”的value
				if resp==nil{
					continue
				}
				c := 0
				var savStr string
				var jsonObj interface{}
				err := json.Unmarshal([]byte(resp.String()), &jsonObj)
				if err != nil {
					continue
				}
				resolve2JosnObj(jsonObj, &savStr, strings.Split(savAction.GetRuleValue(), "."), c)
				storedParamValues[savAction.GetParamName()] = savStr
			}
		case 2:
			{ // 使用正则表达式去获取想要的内容，支持模式(?iLmsux)开头
				if resp==nil{
					continue
				}
				reg, err := regexp.Compile(savAction.GetRuleValue())
				if err != nil {
					continue
				}
				foundStrSlc := reg.FindStringSubmatch(resp.String())
				if len(foundStrSlc) == 0 {
					continue
				}
				if len(foundStrSlc) == 1 {
					storedParamValues[savAction.GetParamName()] = foundStrSlc[0]
				} else {
					storedParamValues[savAction.GetParamName()] = foundStrSlc[1]
				}
			}
		case 3:
			{ // 直接存储一个固定字符串
				storedParamValues[savAction.GetParamName()] = savAction.GetRuleValue()
			}
		case 4:
			{ // 存储一个指定范围内的随机数字
				rangeInts := strings.Split(savAction.GetRuleValue(), "-")
				if len(rangeInts) != 2 {
					continue
				}
				s, err := strconv.ParseInt(rangeInts[0], 10, 64)
				if err != nil {
					continue
				}
				e, err := strconv.ParseInt(rangeInts[1], 10, 64)
				if err != nil {
					continue
				}
				if s > e {
					s, e = e, s
				}
				storedParamValues[savAction.GetParamName()] = strconv.Itoa(int(s)+rand.Intn(int(e)+1-int(s)))
			}
		case 5:
			{ // 存储一个指定范围内的随机字符串
				if len(strings.TrimSpace(savAction.GetRuleValue()))<=2{
					continue
				}
				sepMark:=string(savAction.GetRuleValue()[0])
				choseStrings := strings.Split(string(savAction.GetRuleValue()[1:]), sepMark)
				if len(choseStrings) == 1 {
					storedParamValues[savAction.GetParamName()]=choseStrings[0]
				}
				storedParamValues[savAction.GetParamName()] = choseStrings[rand.Intn(len(choseStrings))]
			}
		}
	}

}

// 解析请求参数数据, 提供从保存的参数字典中取值
func ResolvingHttpItem(httpItem string, storedParamValues map[string]string) string {

	reg, err := regexp.Compile("(?m){@([^@]+)@}")
	if err != nil {
		return httpItem
	}
	matchArgs := reg.FindAllString(httpItem,-1)
	if len(matchArgs)==0{
		return httpItem
	}
	tmpMap := map[string]byte{}
	argSlc:=[]string{}
	for _,argMt := range matchArgs{
		l:=len(tmpMap)
		tmpMap[argMt[2:len(argMt)-2]]=0
		if len(tmpMap)>l{
			argSlc=append(argSlc, argMt[2:len(argMt)-2])
		}
	}

	var resolvedHttpItem =httpItem
	for i,_:=range argSlc{
		resolvedHttpItem = strings.ReplaceAll(resolvedHttpItem, "{@"+argSlc[i]+"@}", storedParamValues[argSlc[i]])
	}
	return resolvedHttpItem

}

// 解析Map类型的请求参数数据
func ResolvingMapItem(mapItem map[string]string, storedParamValues map[string]string) map[string]string {
	if mapItem == nil {
		return nil
	}
	reslvdMapItem := map[string]string{}
	for k, v := range mapItem {
		reslvdMapItem[k] = ResolvingHttpItem(v, storedParamValues)
	}
	return reslvdMapItem
}

// 进行前置http请求，一般用于前置条件判断或进行关联参数的处理
func DoPreTask(baseReqOptions *grequests.RequestOptions, reqSession *grequests.Session, perTaskRequests []*proto.InitBommerRequest_HttpRequest, storedParamValues map[string]string) error {
	reqOptions := *baseReqOptions // 复制一个新的ReqOption
	for _, preTaskReq := range perTaskRequests {
		if preTaskReq.GetUrlPath() != "" {
			var resp *grequests.Response
			var rawData io.Reader
			reslvdUrlPath := ResolvingHttpItem(preTaskReq.GetUrlPath(), storedParamValues)
			reslvdRawData := ResolvingHttpItem(preTaskReq.GetRawData(), storedParamValues)
			if reslvdRawData == "" {
				rawData = nil
			} else {
				rawData = strings.NewReader(reslvdRawData)
			}
			var jsonData interface{}
			reslvJsonData := ResolvingHttpItem(preTaskReq.GetJsonData(), storedParamValues)
			if reslvJsonData == "" {
				jsonData = nil
			} else {
				jsonData = reslvJsonData
			}
			reqOptions.Headers = ResolvingMapItem(preTaskReq.GetHeaders(), storedParamValues)
			if reqOptions.Headers["User-Agent"]!=""{ // 此处为了兼容grequest的User-Agent的处理
				reqOptions.UserAgent=reqOptions.Headers["User-Agent"]
			}
			reqOptions.Params = ResolvingMapItem(preTaskReq.GetParams(), storedParamValues)
			reqOptions.Data = ResolvingMapItem(preTaskReq.GetDictData(), storedParamValues)
			reqOptions.JSON = jsonData
			reqOptions.RequestBody = rawData
			var err error
			if reqSession == nil { // 不使用session
				resp, err = grequests.DoRegularRequest(strings.ToUpper(preTaskReq.GetMethod()), reslvdUrlPath, &reqOptions)
			} else {
				if strings.ToUpper(preTaskReq.GetMethod()) == "GET" {
					resp, err = reqSession.Get(reslvdUrlPath, &reqOptions)
				} else {
					resp, err = reqSession.Post(reslvdUrlPath, &reqOptions)
				}
			}
			if err != nil {
				return err
			}
			if resp.StatusCode >= 400 {
				resp.ClearInternalBuffer()
				resp.Close() // defer 用在loop中可能造成资源无法回收，所以改成程序的出口处都加上清理语句
				return errors.New(fmt.Sprintf("请求[%s]的状态码有误，为[%d]", reslvdUrlPath, resp.StatusCode))
			}
			if preTaskReq.GetAssertChain() != nil {
				if assertstate, msg := DoAssertActionChains(resp, preTaskReq.GetAssertChain()); !assertstate {
					resp.ClearInternalBuffer()
					resp.Close()
					return errors.New(fmt.Sprintf("请求[%s]的判断有失败，为[%s]", reslvdUrlPath, msg))
				}
			}
			if preTaskReq.GetSaveParamChain() != nil {
				DoSaveParamActionChains(resp, preTaskReq.GetSaveParamChain(), storedParamValues)
			}
			resp.ClearInternalBuffer()
			resp.Close()
			continue
		}
		if preTaskReq.GetSaveParamChain() != nil {
			DoSaveParamActionChains(nil, preTaskReq.GetSaveParamChain(), storedParamValues)
		}
	}
	return nil
}

// 构造boomer的任务
func MakeTestTask(baseReqOptions *grequests.RequestOptions, reqSession *grequests.Session, testTask *proto.InitBommerRequest_TestTask, gBoomer *boomer.Boomer, storedParamValues map[string]string) func() {
	return func() {
		// log.Printf("请求%v\n", *HttpRequestEntry)
		localStoredParamValues := map[string]string{}
		for k, v := range storedParamValues { // 复制全局变量存储器到局部变量中
			localStoredParamValues[k] = v
		}
		start := time.Now() // 整个事务开始
		if testTask.GetPreWork() != nil {
			// 有PreWork时，保证事务内的请求一致性，强制使用每个协程独立一个TestTask一个session。
			// 如果不使用全局reqSession，就需要使用新建session
			if reqSession==nil{
				reqSession=grequests.NewSession(baseReqOptions)
			}
			err := DoPreTask(baseReqOptions, reqSession, testTask.GetPreWork(), localStoredParamValues)
			if err != nil {
				responseTime := time.Since(start).Nanoseconds() / int64(time.Millisecond)
				gBoomer.RecordFailure("http", testTask.GetTaskName(), responseTime, fmt.Sprintf("前置任务呢遇到异常：%s", err.Error()))
				return
			}
		}
		reqOptions := *baseReqOptions // 复制一个新的ReqOption
		var rawData io.Reader
		taskWork := testTask.GetTestWork()
		reslvdUrlPath := ResolvingHttpItem(taskWork.GetUrlPath(), localStoredParamValues)
		reslvdRawData := ResolvingHttpItem(taskWork.GetRawData(), localStoredParamValues)
		if reslvdRawData == "" {
			rawData = nil
		} else {
			rawData = strings.NewReader(reslvdRawData)
		}
		var jsonData interface{}
		reslvJsonData := ResolvingHttpItem(taskWork.GetJsonData(), localStoredParamValues)
		if reslvJsonData == "" {
			jsonData = nil
		} else {
			jsonData = reslvJsonData
		}
		reqOptions.Headers = ResolvingMapItem(taskWork.GetHeaders(), localStoredParamValues)
		if reqOptions.Headers["User-Agent"]!=""{ // 此处为了兼容grequest的User-Agent的处理
			reqOptions.UserAgent=reqOptions.Headers["User-Agent"]
		}
		reqOptions.Params = ResolvingMapItem(taskWork.GetParams(), localStoredParamValues)
		reqOptions.Data = ResolvingMapItem(taskWork.GetDictData(), localStoredParamValues)
		reqOptions.JSON = jsonData
		reqOptions.RequestBody = rawData
		var resp *grequests.Response
		var err error
		if reqSession == nil { // 不使用session
			resp, err = grequests.DoRegularRequest(strings.ToUpper(taskWork.GetMethod()), reslvdUrlPath, &reqOptions)
		} else {
			if strings.ToUpper(strings.ToUpper(taskWork.GetMethod())) == "GET" {
				resp, err = reqSession.Get(reslvdUrlPath, &reqOptions)
			} else {
				resp, err = reqSession.Post(reslvdUrlPath, &reqOptions)
			}
		}
		responseTime := time.Since(start).Nanoseconds() / int64(time.Millisecond)
		if err != nil {
			log.Println(err.Error())
			gBoomer.RecordFailure("http", testTask.GetTaskName(), responseTime, err.Error())
			return
		}
		defer resp.Close()
		contentLength := len(resp.Bytes())
		defer resp.ClearInternalBuffer()
		log.Println(resp.StatusCode, contentLength, responseTime)
		if resp.StatusCode >= 400 {
			gBoomer.RecordFailure("http", testTask.GetTaskName(), responseTime, fmt.Sprintf("状态码有误，为[%d]",resp.StatusCode))
			return
		}
		if taskWork.GetAssertChain() != nil {
			if assertstate, msg := DoAssertActionChains(resp, taskWork.GetAssertChain()); !assertstate {
				gBoomer.RecordFailure("http", testTask.GetTaskName(), responseTime, fmt.Sprintf("请求[%s]的判断有失败，为[%s]", reslvdUrlPath, msg))
				return
			}
		}
		gBoomer.RecordSuccess("http", testTask.GetTaskName(), responseTime, int64(contentLength))

	}

}


