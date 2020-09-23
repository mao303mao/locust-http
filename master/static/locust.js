$(window).ready(function() {
    if($("#user_count").length > 0) {
        $("#user_count").focus().select();
    }
});

function appearStopped() {
    $(".box_stop").hide();
    $("a.new_test").show();
    $("a.edit_test").hide();
    $(".user_count").hide();
}

$("#box_stop a.stop-button").click(function(event) {
    event.preventDefault();
    $.get($(this).attr("href"),
        (rsp)=>{
            if(rsp.status){
                msgtip(rsp.message)
            }else{
                warntip(rsp.message)
            }
        }
    );
});

$("a.reset-button").click(function(event) {
    event.preventDefault();
    $.get($(this).attr("href"),(d)=>{
        rpsChart.data=[[],[]];
        responseTimeChart.data=[[],[]];
        usersChart.data=[[]];
        window.rpsChart.reset();
        window.responseTimeChart.reset();
        window.usersChart.reset();
        msgtip("已经重置");
        //setTimeout(()=>{window.location.reload()},1000)
    })
    }
);

var covertTimeStr2Second=function(runTimeStr){
    if(!runTimeStr || runTimeStr.trim()===""){
        return 0
    }
	m=runTimeStr.match(/(\d+)h(\d+)m/)
	if(!!m){
		return parseInt(m[1])*3600+parseInt(m[2])*60
	}
	m=runTimeStr.match(/(\d+)m/)
	if(!!m){
		return parseInt(m[1])*60
	}
	m=runTimeStr.match(/(\d+)s/)
	if(!!m){
		return parseInt(m[1])
	}
	return -99
}


$("#new_test").click(function(event) {
    event.preventDefault();
    $("#start").show();
    $("#user_count").focus().select();
});

$(".edit_test").click(function(event) {
    event.preventDefault();
    $("#edit").show();
    $("#new_user_count").focus().select();
});

$(".close_link").click(function(event) {
    event.preventDefault();
    $(this).parent().parent().hide();
});


$("ul.tabs").tabs("div.panes > div").on("onClick", function(event) {
    if (event.target == $(".chart-tab-link")[0]) {
        // trigger resizing of charts
        rpsChart.resize();
        responseTimeChart.resize();
        usersChart.resize();
    }
});

var stats_tpl = $('#stats-template');
var errors_tpl = $('#errors-template');
var exceptions_tpl = $('#exceptions-template');
var workers_tpl = $('#worker-template');
var slaves_tpl = $('#slave-template');

function setHostName(hostname) {
    hostname = hostname || "";
    $('#host_url').text(hostname);
}

$('#swarm_form').submit(function(event) {
    event.preventDefault();
    var runsecs=covertTimeStr2Second($("#run_time1").val());
    if(runsecs<0){
        warntip('请输入正确的时间格式!') //这里content是一个普通的String
	    return
	}
	if(runsecs<30 && runsecs>0){
	    warntip("请输入不小于30s的运行时间")
		return
	}
    $.post($(this).attr("action"), $(this).serialize(),
        function(response) {
            if (response.success) {
                $("body").attr("class", "hatching");
                $("#start").fadeOut();
                $("#status").fadeIn();
                $(".box_running").fadeIn();
                $("a.new_test").fadeOut();
                $("a.edit_test").fadeIn();
                $(".user_count").fadeIn();
                setHostName(response.host);
            }else{
                warntip(response.message)
            }
        }
    );
});

$('#edit_form').submit(function(event) {
    event.preventDefault();
    var runsecs=covertTimeStr2Second($("#run_time2").val());
    if(runsecs<0){
        warntip('请输入正确的时间格式!') //这里content是一个普通的String
	    return
	}
	if(runsecs<30 && runsecs>0){
	    warntip("请输入不小于30s的运行时间")
		return
	}
    $.post($(this).attr("action"), $(this).serialize(),
        function(response) {
            if (response.success) {
                $("body").attr("class", "hatching");
                $("#edit").fadeOut();
                setHostName(response.host);
            }else{
                warntip(response.message)
            }
        }
    );
});

var sortBy = function(field, reverse, primer){
    reverse = (reverse) ? -1 : 1;
    return function(a,b){
        a = a[field];
        b = b[field];
       if (typeof(primer) != 'undefined'){
           a = primer(a);
           b = primer(b);
       }
       if (a<b) return reverse * -1;
       if (a>b) return reverse * 1;
       return 0;
    }
}

// Sorting by column
var alternate = false; //used by jqote2.min.js
var sortAttribute = "name";
var WorkerSortAttribute = "id";
var desc = false;
var WorkerDesc = false;
var report;

function renderTable(report) {
    var totalRow = report.stats.pop();
    totalRow.is_aggregated = true;
    var sortedStats = (report.stats).sort(sortBy(sortAttribute, desc));
    sortedStats.push(totalRow);
    $('#stats tbody').empty();
    $('#errors tbody').empty();

    window.alternate = false;
    $('#stats tbody').jqoteapp(stats_tpl, sortedStats);

    window.alternate = false;
    $('#errors tbody').jqoteapp(errors_tpl, (report.errors).sort(sortBy(sortAttribute, desc)));

    $("#total_rps").html(Math.round(report.total_rps*100)/100);
    $("#fail_ratio").html(Math.round(report.fail_ratio*100));
    $("#status_text").html(report.state);
    $("#userCount").html(report.user_count);
}

function renderWorkerTable(report) {
    if (report.workers) {
        var workers = (report.workers).sort(sortBy(WorkerSortAttribute, WorkerDesc));
        $("#workers tbody").empty();
        window.alternate = false;
        $("#workers tbody").jqoteapp(workers_tpl, workers);
        $("#workerCount").html(workers.length);
    }
    if (report.slaves) {
        var slaves = report.slaves;
        $("#slaves tbody").empty();
        window.alternate = false;
        $("#slaves tbody").jqoteapp(slaves_tpl, slaves);
        $("#slaveCount").html(slaves.length);
    }
}


$("#stats .stats_label").click(function(event) {
    event.preventDefault();
    sortAttribute = $(this).attr("data-sortkey");
    desc = !desc;
    renderTable(window.report);
});

$("#workers .stats_label").click(function(event) {
    event.preventDefault();
    WorkerSortAttribute = $(this).attr("data-sortkey");
    WorkerDesc = !WorkerDesc;
    renderWorkerTable(window.report);
});

// init charts
var rpsChart = new LocustLineChart($(".charts-container"), "Total Requests per Second", ["RPS", "Failures/s"], "reqs/s", ['#00ca5a', '#ff6d6d']);
var responseTimeChart = new LocustLineChart($(".charts-container"), "Response Times (ms)", ["Median Response Time", "95% percentile"], "ms");
var usersChart = new LocustLineChart($(".charts-container"), "Number of Users", ["Users"], "users");

function updateStats() {
    $.get('./stats/requests', function (report) {
        window.report = report;
        renderTable(report);
        renderWorkerTable(report);
        // console.log((report.state);
        if (report.state !== "stopped" && report.state !== "ready"){
            // get total stats row
            var total = report.stats[report.stats.length-1];
            // update charts
            rpsChart.addValue([total.current_rps, total.current_fail_per_sec]);
            responseTimeChart.addValue([report.current_response_time_percentile_50, report.current_response_time_percentile_95]);
            usersChart.addValue([report.user_count]);
        } else {
            appearStopped();
        }
        setTimeout(updateStats, 2000);
    });
}
updateStats();

function updateExceptions() {
    $.get('./exceptions', function (data) {
        $('#exceptions tbody').empty();
        $('#exceptions tbody').jqoteapp(exceptions_tpl, data.exceptions);
        setTimeout(updateExceptions, 5000);
    });
}
updateExceptions();

$("#initBoomer").click(function(event){
      var elem_slaves=$("#slaves tr td:nth-child(1)")
      if(elem_slaves.size()==0){
        warntip("当前没有可用的压测机!")
        return
      }
      currentList=elem_slaves.toArray().map(x=>x.innerHTML)
      layer.open({
        type: 1
        ,title: false //不显示标题栏
        ,closeBtn: false
        ,area: '300px;'
        ,shade: 0.8
        ,id: 'LAY_layuipro' //设定一个id，防止重复弹出
        ,btn: ['确定', '取消']
        ,btnAlign: 'c'
        ,moveType: 1 //拖拽模式，0或者1
        ,content: `<div style="color:black; font:normal bold 20px sans-serif; text-align: center;">请选择压测机</div>
                   <form id="initBoomer">
                   <div style="padding: 20px 0 15px 10px; line-height: 22px; background-color: #393D49; color: #fff; font-size: 17px;">`
                    +currentList.map(x=>'<input type="checkbox" name="servAddr[]" value='+x+'>'+x).join('<br/>')
                   +`</div></form>`
        ,success: function(layero){
          var btn = layero.find('.layui-layer-btn');
          btn.find('.layui-layer-btn0').click(
            function(event){
                event.preventDefault();
                if($("input[name='servAddr[]']:checked").size()==0){
                    warntip("请选择对应压测机！");
                    return
                }
                $.post('./initBoomer',$("form#initBoomer").serialize(),function(jsonData){
                     msgtip(jsonData.message)
                })
            }
          );
        }
      });

})
$("#stopBoomer").click(function(event){
    $.post('./shutdownBoomer',function(jsonData){
        if(jsonData.success){
             msgtip(jsonData.message)
        }else{
            warntip(jsonData.message)
        }

    })
})

// 自定义提示
msgtip=function(msg){
    layui.layer.msg(
        `<div style="color: black; font-size: 15px; text-align: center">`+msg+`</div>`,
        {
          offset: '15px',
          icon: 1,
          time: 2000
        })
}
// 自定义提示
warntip=function(msg){
    layui.layer.msg(
        `<div style="color: black; font-size: 15px; text-align: center">`+msg+`</div>`,
        {
          offset: '15px',
          icon: 5,
          time: 2000
        })
}

function reinitIframe(){
    var iframe = document.getElementById("ifm_transcation");
    try{
        var bHeight = iframe.contentWindow.document.body.scrollHeight;
        var dHeight = iframe.contentWindow.document.documentElement.scrollHeight;
        var height = Math.min(bHeight, dHeight);
        iframe.height = height;
    }catch (ex){}
}
window.setInterval("reinitIframe()", 500);