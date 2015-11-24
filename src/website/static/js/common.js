/**
* Created by wuyingsong on 2015/3/27.
*/


//$(document).ready(function(){

function prev_page(){
    window.location.href=document.referrer;
}

$("#return").click(function(){
     prev_page();
});

$("#select_all").click(function(){
    if($("#select_all").is(":checked")){
        $("[name='delete']").prop("checked", true);
        $("[name='delete']").parent().parent().addClass("info");
    }else{
        $("[name='delete']").prop("checked", false);
        $("[name='delete']").parent().parent().removeClass("info");
    }
});

function getChecked(){
    var ids = "";
    $("[name='delete']").each(function(){
        if($(this).is(":checked")){
            ids += $(this).val() + ",";
        }
    });
    console.log(ids);
    return ids.substring(0, ids.length-1)
}

//全选或者取消单个选择，取消全选框
$("[name='delete']").change(function(){
    var counter = 0;
     $("[name='delete']").each(function(){
        if($(this).is(":checked")){
            $(this).parent().parent().addClass("info");
            counter ++
        }else{
            $(this).parent().parent().removeClass("info");
        }
    });
    var sum = $("[name='delete']").length;
    if(counter == sum ){
        $("#select_all").prop("checked", true);
    }else{
        $("#select_all").prop("checked", false);
    }
});

function Deletes(ids){
    var url = window.location.pathname + "delete/";
    console.log(url);
    $.ajax({
        type: "post",
        url: url ,
        data:{"ids":ids, "csrfmiddlewaretoken":$("[name='csrfmiddlewaretoken']").val()},
        success: function (result) {
            if(result.status ==0) {
                location.reload();
            }else{
                alert(result.message);
            }
        }
    });
}

function Disables(ids){
    var url = window.location.pathname + "disable/";
    $.ajax({
        type:"post",
        url:url,
        data:{"ids":ids, "csrfmiddlewaretoken":$("[name='csrfmiddlewaretoken']").val()},
        success:function(result){
            if(result.status ==0) {
                location.reload();
            }else{
                alert(result.message);
            }
        }
    });
}
function Enables(ids){
    var url = window.location.pathname + "enable/";
    $.ajax({
        type:"post",
        url:url,
        data:{"ids":ids, "csrfmiddlewaretoken":$("[name='csrfmiddlewaretoken']").val()},
        success:function(result){
            if(result.status ==0) {
                location.reload();
            }else{
                alert(result.message);
            }
        }
    });
}
//批量删除
$("#mDelete").click(function(){
    var ids = getChecked();
    if(ids.length == 0){
        alert("至少选中一条记录进行删除!");
        return false;
    }
    if(confirm("确定要删除" + ids.split(",").length + "条记录吗?")){
        Deletes(ids);
    }
});

//批量禁用
$("#mDisable").click(function(){
    var ids = getChecked();
    if(ids.length == 0){
        alert("至少选择一条记录");
        return false;
    }
    if(confirm("确定要禁用" + ids.split(",").length + "条记录吗?")){
        Disables(ids);
    }
});

$("a[name='enable']").click(function(){
    if(confirm("确定启用吗?")){
        Enables($(this).attr("title"));
    }
});


$("a[name='disable']").click(function(){
    if(confirm("确定禁用吗?")){
        Disables($(this).attr("title"));
    }
});

$("a[name='delete']").click(function(){
    if(confirm("确定删除吗?")){
        Deletes($(this).attr("title"));
    }
});
//批量启用

$("#mEnable").click(function(){
    var ids = getChecked();
    if(ids.length == 0){
        alert("至少选择一条记录");
        return false;
    }
    if(confirm("确定要启用" + ids.split(",").length + "条记录吗?")){
        Enables(ids);
    }
});


function getCsrfToken(){
    return $("[name='csrfmiddlewaretoken']").val()
}


//格式化时间格式为数组
//2015-06-17 14:46:10 -> ["2015", "06", "17", "14", "46", "10"]
function parse_datetime(datetime){
        date_arr = datetime.split(" ")[0].split("-");
        time_arr = datetime.split(" ")[1].split(":");
        return $.merge(date_arr, time_arr);
}

function get_data(url){
    var data;
    $.ajax({
        type:"GET",
        url:url,
        async:false,
        dataType:"json",
        success:function(result){
            data = result;
        }
    });
    return data;
}

function getUrlParam(name){
    var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r!=null) return r[2]; return null;
}

$("#acknowledged").click(function(){
            var ids = getChecked();
            if(ids.length == 0){
                alert("至少选择一条记录");
                return false;
            }
            if(confirm("确定要启用" + ids.split(",").length + "条记录吗?")){
                all_acknowledged(ids);
            }
        });

function all_acknowledged(ids){
    var url = window.location.pathname + "/all_acknowledged/";
    $.ajax({
		type:"post",
		url:url,
		data:{"ids":ids, "csrfmiddlewaretoken":$("[name='csrfmiddlewaretoken']").val()},
		success:function(result){
		//alert(result);
		    if(result.status ==0) {
			location.reload();
		    }else{
			alert(result.message);
		    }
		}
	});
}
