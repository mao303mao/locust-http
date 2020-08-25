# coding:utf-8
import json
import boomerCall_pb2
import traceback

def makeSaveParamChain(saveParamChainObj):
    if not saveParamChainObj:
        return None
    saveParamChain=[]
    for sp in saveParamChainObj:
        saveParamChain.append(
            boomerCall_pb2.InitBommerRequest.SaveParamAction(ParamName=sp.get("ParamName"),
                                                             SaveType=sp.get("SaveType"),
                                                             RuleValue=sp.get("RuleValue")))
    return saveParamChain

def makeAssertChain(assertChainObj):
    if not assertChainObj:
        return None
    assertChain=[]
    for ac in assertChainObj:
        assertChain.append(
            boomerCall_pb2.InitBommerRequest.AssertAction(AssertType=ac.get("AssertType"),
                                                             RuleValue=ac.get("RuleValue")))
    return assertChain


def makeSigleRequest(requestObj):
    return {
        "UrlPath":         requestObj.get("UrlPath"),
        "Method":          requestObj.get("Method"),
        "Headers":         requestObj.get("Headers"),
        "DictData":        requestObj.get("DictData"),
        "Params":          requestObj.get("Params"),
        "RawData":         requestObj.get("RawData"),
        "JsonData":        requestObj.get("JsonData"),
        "SaveParamChain":  makeSaveParamChain(requestObj.get("SaveParamChain")),
        "AssertChain":     makeAssertChain(requestObj.get("AssertChain"))
    }


def makeInitBoomerRequest(jsonFilePath:str,locustMaster):
    try:
        with open(jsonFilePath,mode='rb') as f:
            jsonObj=json.load(f)
            PreTask=None
            preTaskObj=jsonObj.get("PreTask")
            if preTaskObj:
                PreTask = []
                for pt in preTaskObj:
                    PreTask.append(
                        makeSigleRequest(pt)
                    )

            mainTaskObjList=jsonObj.get("MainTask")
            if not mainTaskObjList:
                raise Exception("缺少主测试方法，不能进行测试")
            MainTask = []
            for mt in mainTaskObjList:
                tmpMt={}
                tmpMt["TaskName"]=mt.get("TaskName")
                tmpMt["TaskWeight"]=mt.get("TaskWeight")
                preWork=None
                preWorkObj=mt.get("PreWork")
                if preWorkObj:
                    preWork=[]
                    for pw in preWorkObj:
                        preWork.append(makeSigleRequest(pw))
                tmpMt["PreWork"]=preWork
                if not mt.get("TestWork"):
                    raise Exception("主测试方法中缺少测试方法，不能进行测试")
                tmpMt["TestWork"]=makeSigleRequest(mt.get("TestWork"))
                MainTask.append(tmpMt)
            HttpProxy=jsonObj.get("HttpProxy")
            return boomerCall_pb2.InitBommerRequest(
                isSession=True,LocustMaster=locustMaster,PreTask=PreTask,MainTask=MainTask,HttpProxy=HttpProxy)
    except Exception as e:
        print("%s\n"%(e)+traceback.format_exc())
