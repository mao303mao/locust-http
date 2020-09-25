layui.use([ 'layer', 'form','upload'], function () {
    var $ = layui.jquery;
    var form = layui.form;
    var layer = layui.layer;
    var upload = layui.upload;
    form.on("select(saveType)",function(val){
        var val_input=$(this).parents('.saveParam-component').find("input.fake-ruleValue");
        switch (val.value){
            case "0":
                val_input.attr("placeholder","例://input[@name='a']/@value或//p[@name='b']/text()");
                val_input.attr("lay-verify","");
                break;
            case "1":
                val_input.attr("placeholder","例:针对{\"data\":[{\"id\":1},{\"id\":2}]},通过data.0.id可以获取1");
                val_input.attr("lay-verify","");
                break;
            case "2":
                val_input.attr("placeholder","获取正则表达式第1个子匹配，支持模式(?iLmsux)开头");
                val_input.attr("lay-verify","");
                break;
            case "3":
                val_input.attr("placeholder","保存一个固定的字符串");
                val_input.attr("lay-verify","");
                break;
            case "4":
                val_input.attr("placeholder","获取范围内随机整树,例:1-100");
                val_input.attr("lay-verify","number_range");
                break;
            case "5":
                val_input.attr("placeholder","注：第1个字符是分割符,例:|aaaa|bbbb|cccc");
                val_input.attr("lay-verify","");
                break;
            default:
                val_input.attr("placeholder","");
                val_input.attr("lay-verify","");
        }
    });
    form.verify({
        number_range:function(value, item){ //value：表单的值、item：表单的DOM对象
            if(!/^\d+-\d+$/.test(value)){
                return "请输入'数字-数字'格式";
            }
        },
        ipv4:function(value,item){
            if(value.trim()===""){
                return
            }
            if(!/^\d{1,3}\.\d{1,3}.\d{1,3}.\d{1,3}:\d+$/.test(value)){
                return "请输入正确的ipv4格式";
            }
            for(var item of $('[name=HttpProxy]').val().split(":")[0].split(".")){
                if(parseInt(item)>255){
                    return "请输入正确的ipv4格式";
                }
            }
        }

    });
    form.on("submit(save)",function(e){
        $.ajax({
            url : "/saveTrans",
            type : 'POST',
            data : $(e.form).serialize(),
            dataType : "json",
            processData:false,
            traditional:false,
            success : function(jsonData) {
                if(jsonData.success){
                    layer.alert(jsonData.message);
                }else{
                    layer.alert(jsonData.message);
                }
            }
        });

    });
    form.on("submit(backup)",function(e){
        $(e.form).attr("action",'/backupTrans');
        $(e.form).submit();
    });
    upload.render({
        elem: '#btn-import'
        ,url: '/importTrans' //改成您自己的上传接口
        ,accept: 'file' //普通文件
        ,done: function(res){
            if(res.success){
                window.location="/importedTrans"
            }else{
                layer.alert(res.message)
            }
        }
  });
});

function makeKeyvalueComponent(x) {
    var html = `
            <div style="padding: 10px 0 0 80px;" class="keyvalue-component">
                <div class="layui-inline">
                    <div class="layui-input-inline" style="width: 200px;">
                        <input type="text" lay-verify="required" name="`
                        +x+"-key"+
                        `" placeholder="key" autocomplete="off" class="layui-input">
                    </div>
                    <div class="layui-input-inline" style="width: 400px;">
                        <input type="text" name="`
                        +x+"-value"+
                        `" placeholder="value,参数用{@参数名@}代替" autocomplete="off" class="layui-input">
                    </div>
                    <button type="button" class="layui-btn layui-btn-warm layui-btn-sm" title="删除" onclick="delKeyValue(this)"><i class="layui-icon">&#xe67e;</i></button>
                </div>
            </div>
        `
    return html;
}

function delKeyValue(obj) {
    $(obj).parents(".keyvalue-component").remove();
}

function addKeyValue(obj, x) {
    var html = makeKeyvalueComponent(x);
    var parentElem=$(obj).parents(".layui-tab-item");
    var keyElems=parentElem.find("input[name='"+x+"-key"+"']")
    if(keyElems && keyElems.size()>0){ // 前一个key没有完成，不能添加新的
        for(i=0;i<keyElems.size();i++){
            if($(keyElems[i]).val().trim()===""){
                return
            }
        }
    }
    parentElem.append(html);
    layui.form.render();
}

function makeAssertComponent(x){
    html=`
        <div style="padding: 10px 0 0 80px;" class="assert-component">
            <div class="layui-inline">
                <div class="layui-input-inline" style="width: 200px;">
                    <select name="`+x+"-assertType" +
                    `" lay-verify="required">
                        <option value="0">响应状态码=</option>
                        <option value="1">响应字节长度></option>
                        <option value="2">响应字节长度=</option>
                        <option value="3">响应字节长度<</option>
                    </select>
                </div>
                <div class="layui-input-inline" style="width: 400px;">
                    <input type="text" name="`+x+"-assertValue"+
                    `" lay-verify="required|number" placeholder="比较值" autocomplete="off"
                           class="layui-input">
                </div>
                <button type="button" class="layui-btn layui-btn-warm layui-btn-sm" title="删除"
                        onclick="delAssert(this)"><i class="layui-icon">&#xe67e;</i></button>
            </div>
        </div>
    `
    return html;
}
function delAssert(obj) {
    $(obj).parents(".assert-component").remove();
}

function addAssert(obj, x) {
    var html = makeAssertComponent(x);
    var parentElem=$(obj).parents(".layui-tab-item");
    var assertElems=parentElem.find("input[name='"+x+"-assertValue"+"']")
    if(assertElems && assertElems.size()>0){ // 前一个key没有完成，不能添加新的
        for(i=0;i<assertElems.size();i++){
            if($(assertElems[i]).val().trim()===""){
                return
            }
        }
    }
    parentElem.append(html);
    layui.form.render();
}
function makeSaveParamComponent(x){
    html=`
        <div style="padding: 10px 0 0 80px;" class="saveParam-component">
             <div class="layui-inline">
                 <div class="layui-input-inline" style="width: 200px;">
                     <select name="`+x+`-saveType" lay-verify="required" lay-filter="saveType">
                         <option value="0">XPATH解析HTML</option>
                         <option value="1">点分割解析JSON</option>
                         <option value="2">文本正则匹配</option>
                         <option value="3">保存固定字符串</option>
                         <option value="4">范围内随机整数</option>
                         <option value="5">随机字符串</option>
                     </select>
                 </div>
                 <div class="layui-input-inline" style="width: 200px;">
                     <input type="text" name="`+x+`-paramName" lay-verify="required" placeholder="参数名称" autocomplete="off"
                            class="layui-input">
                 </div>
                 <div class="layui-input-inline" style="width: 400px;">
                     <input type="text" name="`+x+`-ruleValue" lay-verify="required" placeholder="例://input[@name='a']/@value或//p[@name='b']/text()" autocomplete="off"
                            class="layui-input fake-ruleValue">
                 </div>
                 <button type="button" class="layui-btn layui-btn-warm layui-btn-sm" title="删除"
                         onclick="delSaveParam(this)"><i class="layui-icon">&#xe67e;</i></button>
             </div>
        </div>
    `
    return html;
}
function delSaveParam(obj) {
    $(obj).parents(".saveParam-component").remove();
}

function addSaveParam(obj, x) {
    var html = makeSaveParamComponent(x);
    var parentElem=$(obj).parents(".layui-tab-item");
    var nameElems=parentElem.find("input[name='"+x+"-paramName"+"']")
    var ruleElems=parentElem.find("input[name='"+x+"-ruleValue"+"']")
    if(nameElems && nameElems.size()>0){ // 前一个key没有完成，不能添加新的
        for(i=0;i<nameElems.size();i++){
            if($(nameElems[i]).val().trim()==="" || $(ruleElems[i]).val().trim()===""){
                return
            }
        }
    }
    parentElem.append(html);
    layui.form.render();
}

function makeHttpRequest(z){
    html=`
        <div class="http-request-pre">
                <div class="layui-form-item">
                    <div class="layui-inline">
                        <label class="layui-form-label">Method</label>
                        <div class="layui-input-inline" style="width: 80px;">
                            <select name="`+z+`-method">
                                <option value="Get">Get</option>
                                <option value="Post">Post</option>
                            </select>
                        </div>
                        <label class="layui-form-mid">URL</label>
                        <div class="layui-input-inline" style="width: 600px;">
                            <input type="text" name="`+z+`-urlPath" placeholder="参数用{@参数名@}代替" autocomplete="off"
                                   class="layui-input">
                        </div>
                    </div>
                </div>
                <div class="layui-tab layui-tab-brief">
                    <ul class="layui-tab-title">
                        <li class="layui-this">请求头设置</li>
                        <li>URL参数</li>
                        <li>Form(字段名不同)</li>
                        <li>Json格式</li>
                        <li>Raw内容</li>
                        <li><i class="layui-icon" style="color:green; ">&#x1005;</i>响应断言</li>
                        <li><i class="layui-icon" style="color:blue; ">&#xe641;</i>保存参数</li>
                    </ul>
                    <div class="layui-tab-content">
                        <div class="layui-tab-item layui-show">
                            <label class="layui-form-mid">请求头格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-warm layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'`+z+`-headers')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">Params格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-warm layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'`+z+`-params')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">FORM字段格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-warm layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'`+z+`-dictdata')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <div class="layui-form-item layui-form-text">
                                <label class="layui-form-label">Json内容</label>
                                <div class="layui-input-block">
                                    <textarea name="`+z+`-jsondata" placeholder="请输入Json内容,参数用{@参数名@}代替" class="layui-textarea"></textarea>
                                </div>
                            </div>
                        </div>
                        <div class="layui-tab-item">
                            <div class="layui-form-item layui-form-text">
                                <label class="layui-form-label">Raw内容</label>
                                <div class="layui-input-block">
                                    <textarea name="`+z+`-rawdata" placeholder="请输入Raw内容,参数用{@参数名@}代替" class="layui-textarea"></textarea>
                                </div>
                            </div>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">断言格式：判断类型+比较值</label>
                            <button type="button" class="layui-btn layui-btn-warm layui-btn-sm " title="添加"
                                    onclick="addAssert(this,'`+z+`')"><i class="layui-icon">&#xe654;</i></button>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">保存参数：保存类型+参数名称+取值规则</label>
                            <button type="button" class="layui-btn layui-btn-warm layui-btn-sm " title="添加"
                                    onclick="addSaveParam(this,'`+z+`')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
    `
    return html;
}

function addPreTask(){
    var idx=window.PreTaskMark;
    var pre_pretask_url=$("input[name=pretask-"+idx+"-urlPath")
    if(pre_pretask_url.length>1){
        layui.layer.alert("全局前置处理出现乱序，请删除重来！")
        return
    }
    if(pre_pretask_url.size()==1){
        var pre_pretask_saveParam=$("input[name='pretask-"+idx+"-paramName']");
        if(pre_pretask_url.val().trim()===""){
            if(pre_pretask_saveParam.length==0 ||
                pre_pretask_saveParam.eq(pre_pretask_saveParam.length-1).val().trim()===""
            ){
                layui.layer.alert("请完善全局前置处理：请求URL或保存参数！")
                return
            }
        }
    }
    var z="pretask-"+(idx+1)
    var html = makeHttpRequest(z);
    var parentElem=$("div.fakeclass-pretasks");
    parentElem.append(html);
    window.PreTaskMark=window.PreTaskMark+1;
    console.log('PreTaskMark', window.PreTaskMark)
    layui.form.render();
}

function delLastPreTask(){
    var preTasksElem=$("div.fakeclass-pretasks div.http-request-pre");
    preTasksElem.eq(preTasksElem.length-1).remove();
    if(window.PreTaskMark<=0){
        window.PreTaskMark=0
    }else{
        window.PreTaskMark=window.PreTaskMark-1;
    }
    console.log('PreTaskMark', window.PreTaskMark)
}


function addPreWork(testtask_idx){
    if(window.TestTaskMark.length<1 || window.TestTaskIdMark.length<1 || window.TestTaskMark.length!= window.TestTaskIdMark.length){
        layui.layer.alert("出现测试事务乱序，请重置重来！")
        return
    }
    idIndx=window.TestTaskIdMark.indexOf(testtask_idx) // 对应任务所在index
    if(idIndx<0){
        layui.layer.alert("事务的id号有误，请重置重来！")
        return
    }
    var idx=window.TestTaskMark[idIndx]; // 事务的前置处理的序号
    var pre_pretask_url=$("input[name="+"testtask-"+(testtask_idx)+"-prework-"+idx+"-urlPath")
    if(pre_pretask_url.length>1){
        layui.layer.alert("事务前置处理出现乱序，请删除重来！")
        return
    }
    if(pre_pretask_url.size()==1){
        var pre_pretask_saveParam=$("input[name='"+"testtask-"+(testtask_idx)+"-prework-"+idx+"-paramName']");
        if(pre_pretask_url.val().trim()===""){
            if(pre_pretask_saveParam.length==0 ||
                pre_pretask_saveParam.eq(pre_pretask_saveParam.length-1).val().trim()===""
            ){
                layui.layer.alert("请完善事务前置处理：请求URL或保存参数！")
                return
            }
        }
    }
    var z="testtask-"+(testtask_idx)+"-prework-"+(idx+1)
    var html = makeHttpRequest(z);
    var parentElem=$("div.fakeclass-prework-"+testtask_idx);
    parentElem.append(html);
    window.TestTaskMark[idIndx]=window.TestTaskMark[idIndx]+1; // 对应事务的前置处理的序号+1
    console.log('TestTaskMark', window.TestTaskMark)
    layui.form.render();
}

function delLastPreWork(testtask_idx){
    if(window.TestTaskMark.length<1 || window.TestTaskIdMark.length<1 || window.TestTaskMark.length!= window.TestTaskIdMark.length){
        layui.layer.alert("出现测试事务乱序，请重置重来！")
        return
    }
    idIndx=window.TestTaskIdMark.indexOf(testtask_idx)
    if(idIndx<0){
        layui.layer.alert("要删除的事务的id号有误！")
        return
    }
    var preTasksElem=$("div.fakeclass-prework-"+testtask_idx+" div.http-request-pre");
    preTasksElem.eq(preTasksElem.length-1).remove();
    window.TestTaskMark[idIndx]=window.TestTaskMark[idIndx]-1; // 对应事务的前置处理的序号-1
    console.log('TestTaskMark', window.TestTaskMark)
}

function makeTestTask(t){
    var html=`
    <div class="layui-card">
        <div class="layui-card-header" >
                <button type="button" style="margin-left: 10px; font-weight: bold;"
                        class="layui-btn layui-btn-xs layui-btn-normal"
                        onclick="delTestTask(this,`+t+`)" title="删除该事务">
                    <i class="layui-icon">&#xe67e;</i>
                </button>
            <b style="margin: 0 20px 0 15px; color:blue ; font-size:15px;">事务-`+t+`</b>
        </div>
        <div class="layui-card-body">
            <div class="layui-form-item">
                <div class="layui-inline">
                    <label class="layui-form-label">权重</label>
                    <input placeholder="数字" name="testtask-`+t+`-taskWeight" lay-verify="required|number"
                           autocomplete="off"
                           style="width:140px" class="layui-input" type="text">
                </div>
                <div class="layui-inline">
                    <label class="layui-form-label">事务名称</label>
                    <input placeholder="" name="testtask-`+t+`-TaskName" lay-verify="required" autocomplete="off"
                           style="width:200px" class="layui-input" type="text">
                </div>
                <div class="layui-inline">
                    <b style="margin-left:15px">事务内前置处理</b>
                    <div class="layui-btn-group" style="margin-top:-2px">
                        <button type="button" style="margin-left: 10px; font-weight: bold;"
                                class="layui-btn layui-btn-xs layui-btn-primary"
                                onclick="addPreWork(`+t+`)" title="添加新的事务前置处理">
                            <i class="layui-icon">&#xe654;</i>
                        </button>
                        <button type="button" style="margin-left: 10px; font-weight: bold;"
                                class="layui-btn layui-btn-xs layui-btn-primary"
                                onclick="delLastPreWork(`+t+`)" title="删除最后的事务前置处理">
                            <i class="layui-icon">&#xe67e;</i>
                        </button>
                    </div>
                </div>
            </div>
            <div class="fakeclass-prework-`+t+`">
            </div>
            <div class="http-request">
                <div class="layui-form-item">
                    <div class="layui-inline">
                        <label class="layui-form-label">Method</label>
                        <div class="layui-input-inline" style="width: 80px;">
                            <select name="testwork-`+t+`-method">
                                <option value="Get">Get</option>
                                <option value="Post">Post</option>
                            </select>
                        </div>
                        <label class="layui-form-mid">URL</label>
                        <div class="layui-input-inline" style="width: 600px;">
                            <input type="text" lay-verify="required" name="testwork-`+t+`-urlPath" placeholder="参数用{@参数名@}代替" autocomplete="off"
                                   class="layui-input">
                        </div>
                    </div>
                </div>
                <div class="layui-tab layui-tab-brief">
                    <ul class="layui-tab-title">
                        <li class="layui-this">请求头设置</li>
                        <li>URL参数</li>
                        <li>Form(字段名不同)</li>
                        <li>Json格式</li>
                        <li>Raw内容</li>
                        <li><i class="layui-icon" style="color:green; ">&#x1005;</i>响应断言</li>
                    </ul>
                    <div class="layui-tab-content">
                        <div class="layui-tab-item layui-show">
                            <label class="layui-form-mid">请求头格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'testwork-`+t+`-headers')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">Params格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'testwork-`+t+`-params')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">FORM字段格式：名称+对应值</label>
                            <button type="button" class="layui-btn layui-btn-sm " title="添加"
                                    onclick="addKeyValue(this,'testwork-`+t+`-dictdata')"><i class="layui-icon">&#xe654;</i>
                            </button>
                        </div>
                        <div class="layui-tab-item">
                            <div class="layui-form-item layui-form-text">
                                <label class="layui-form-label">Json内容</label>
                                <div class="layui-input-block">
                                    <textarea name="testwork-`+t+`-jsondata" placeholder="请输入Json内容,参数用{@参数名@}代替" class="layui-textarea"></textarea>
                                </div>
                            </div>
                        </div>
                        <div class="layui-tab-item">
                            <div class="layui-form-item layui-form-text">
                                <label class="layui-form-label">Raw内容</label>
                                <div class="layui-input-block">
                                    <textarea name="testwork-`+t+`-rawdata" placeholder="请输入Raw内容,参数用{@参数名@}代替" class="layui-textarea"></textarea>
                                </div>
                            </div>
                        </div>
                        <div class="layui-tab-item">
                            <label class="layui-form-mid">断言格式：判断类型+比较值</label>
                            <button type="button" class="layui-btn layui-btn-sm " title="添加"
                                    onclick="addAssert(this,'testwork-`+t+`')"><i class="layui-icon">&#xe654;</i></button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
    return html;
}

function addTestTask(){
    if(window.TestTaskMark.length<1 || window.TestTaskIdMark.length<1 || window.TestTaskMark.length!= window.TestTaskIdMark.length){
        layui.layer.alert("出现测试事务乱序，请重置重来！")
        return
    }
    window.TestTaskId=window.TestTaskId+1;
    var html=makeTestTask(window.TestTaskId);
    var parentElem=$("div.fakeclass-testtasks");
    parentElem.append(html);
    window.TestTaskIdMark.push(window.TestTaskId)
    window.TestTaskMark.push(0)
    console.log('TestTaskIdMark', window.TestTaskIdMark)
    console.log('TestTaskMark', window.TestTaskMark)
    layui.form.render();
}

function delTestTask(obj,id){
    if(window.TestTaskMark.length<1 || window.TestTaskIdMark.length<1 || window.TestTaskMark.length!= window.TestTaskIdMark.length){
        layui.layer.alert("出现测试事务乱序，请重置重来！")
        return
    }
    if(window.TestTaskIdMark.length<2){
        layui.layer.alert("至少要保留一个事务！")
        return
    }
    testtasksElem=$(obj).parents(".layui-card").remove();
    idIndx=window.TestTaskIdMark.indexOf(id)
    if(idIndx<0){
        layui.layer.alert("要删除的事务的id号有误！")
        return
    }
    window.TestTaskIdMark.splice(idIndx,1)
    window.TestTaskMark.splice(idIndx,1)
    console.log('TestTaskIdMark', window.TestTaskIdMark)
    console.log('TestTaskMark', window.TestTaskMark)
}


$("#btn-save-trans").click(function(event){
    event.preventDefault();
    $("input[name=PreTaskMark]").val(window.PreTaskMark)
    $("input[name=TestTaskIdMark]").val(window.TestTaskIdMark)
    $("input[name=TestTaskMark]").val(window.TestTaskMark)
});

$("#btn-backup-trans").click(function(event){
    event.preventDefault();
    $("input[name=PreTaskMark]").val(window.PreTaskMark)
    $("input[name=TestTaskIdMark]").val(window.TestTaskIdMark)
    $("input[name=TestTaskMark]").val(window.TestTaskMark)
});

$("#btn-reset").click(function(event){
    event.preventDefault();
    $.post("/resetTrans",rsp=>{
        layui.layer.msg("已经重置！")
        setTimeout(()=>{
            window.location="/transaction";
        },1500)

    })

});