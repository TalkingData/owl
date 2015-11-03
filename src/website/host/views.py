# coding:utf-8
from django.shortcuts import render_to_response, HttpResponse, get_object_or_404,RequestContext,HttpResponseRedirect
from django.core.serializers.json import DjangoJSONEncoder
from django.contrib.auth.decorators import login_required
from django.http import Http404

from host.models import *
from assest.models import *
#from network.models import *
from django.db.models import Q
from django.contrib.auth.models import  Group
from system.views import paging
from django.conf import settings
import json

symbols= {
    ">": "大于",
    ">=": "大于等于",
    "<": "小于",
    "<=": "小于等于",
    "=": "等于",
    "<>": "不等于"
    }

methods={
    "sum": "加和",
    "avg": "平均值",
    "max": "最大值",
    "min": "最小值",
    "ratio": "环比",
    }


#主机查询视图
def HostListView(request):
    context = {}
    q = request.GET.get("q")
    page = request.GET.get("page")
    if not page:
        page = 1
    queryset=None
    if q:
        queryset = host.objects.filter(
            Q(ip__contains=q)|
            Q(server__sn__contains=q)|
            Q(server__cabinet__name__contains=q)|
            Q(group__name__contains=q)
        )
    else:
        queryset = host.objects.all().order_by('-status', '-alarm', '-c_time')
    context["hosts"] = paging(queryset.distinct(), page, settings.PAGE_SIZE)
    return render_to_response(
        "host/host_list.html",
        context ,
        context_instance=RequestContext(request)
    )

@login_required()
def HostUpdateView(request, pk):
    context = {}
    template_name = "host/host_edit.html"
    context["groups"] = group.objects.all()
    if request.method == "GET":
        context["host"] = host.objects.get(pk=int(pk))
        return render_to_response(template_name, context, context_instance = RequestContext(request))
    else:
        id = request.POST.get("id")
        _host = host.objects.get(pk=int(id))
        _host.idrac = request.POST.get("idrac")
        _host.hostname = request.POST.get("hostname")
        _host.os = request.POST.get("os")
        _host.kernel  = request.POST.get("kernel")
        _host.alarm = request.POST.get("alarm")
        group_id = request.POST.getlist("group")
        _host.group.clear()
        for g_id in group_id:
            _host.group.add(group.objects.get(pk=g_id))
        _host.save()
        context["host"] = _host
        context["message"] = "修改成功"
        return HttpResponseRedirect("/host/host/")

def HostDeleteView(request):
    context = {}
    ids = request.POST.get("ids")
    if request.user.is_superuser is False:
        context["status"] = 1
        context["message"] = "权限拒绝.."
    else:
        for host_id in ids.split(","):
            try:
                _host=host.objects.get(pk=int(host_id))
                if _host.server.has_cabinet():
                    #释放U位
                    _host.server.cabinet.clear_unit(_host.server.location, _host.server.location + _host.server.unit)
                _host.server.delete()
                _host.delete()
                context["status"] = 0
            except:
                context["status"] = 1

    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )
def ServiceListView(request):
    context = {}
    host_id = request.GET.get("host_id")
    template_id = request.GET.get("template_id")
    page = request.GET.get("page")
    q = request.GET.get("q")
    if not page:
        page = 1
    if host_id:
        context["host"] = host.objects.get(pk=int(host_id))
        queryset = context["host"].service_set.all()
    else:
        context["template"] = template.objects.get(pk=int(template_id))
        queryset = context["template"].service_set.all()
    if q:
        queryset = queryset.filter(
            Q(name__contains = q.strip())
        )
    context["services"] = paging(queryset, page, settings.PAGE_SIZE)
    return render_to_response(
        "host/service_list.html",
        context ,
        context_instance=RequestContext(request)
    )

def ServiceUpdateView(request, pk):
    context = {}
    context["groups"] =Group.objects.all()
    template_name = "host/service_edit.html"
    template_id = request.GET.get("template_id")
    host_id = request.GET.get("host_id")
    redirect_url = ""
    if host_id:
        context["host"] = host.objects.get(pk=int(host_id))
        redirect_url = "/host/service/?host_id=%s"%host_id
    if template_id:
        context["template"] = template.objects.get(pk=int(template_id))
        redirect_url = "/host/service/?template_id=%s"%template_id
    if request.method == "GET":
        _service = service.objects.get(pk=int(pk))
        context["service"] = _service
        return render_to_response(
            template_name,
            context,
            context_instance = RequestContext(request)
        )
    else:
        id = request.POST.get("id")
        groups = request.POST.getlist("groups")
        _service = service.objects.get(pk=int(id))
        _service.name = request.POST.get("name")
        _service.plugin = request.POST.get("plugin")
        _service.args = request.POST.get("args")
        _service.exec_interval = request.POST.get("exec_interval")
        _service.alarm = request.POST.get("alarm")
        _service.save()
        _service.group.clear()
        if groups:
            for g_id in groups:
                _service.group.add(Group.objects.get(pk=int(g_id)))
        return HttpResponseRedirect(redirect_url)

def ServiceCreateView(request):
    context = {}
    host_id = request.GET.get("host_id")
    template_id = request.GET.get("template_id")
    redirect_url = ""
    context["groups"] = Group.objects.all()
    if template_id:
        context["template"] = template.objects.get(pk=int(template_id))
    if host_id:
        context["host"] = host.objects.get(pk=int(host_id))

    if request.method == "GET":
        return render_to_response(
            "host/service_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _service=service.objects.create(
            name = request.POST.get("name"),
            plugin = request.POST.get("plugin"),
            args = request.POST.get("args"),
            exec_interval = request.POST.get("exec_interval")
        )
        groups = request.POST.getlist("groups")
        if groups:
            for g in groups:
                _service.group.add(Group.objects.get(pk=int(g)))
        if host_id:
            _service.host = context["host"]
            redirect_url = "/host/service/?host_id=%s"%host_id
        if template_id:
            _service.template = context["template"]
            redirect_url = "/host/service/?template_id=%s"%template_id
        _service.save()
    return HttpResponseRedirect(redirect_url)


def ItemListView(request):
    context = {}
    context["methods"] = methods
    context["symbols"] = symbols
    context["dts"] = ("GAUGE","COUNTER","DERIVE")
    service_id = request.GET.get("service_id")
    host_id = request.GET.get("host_id")
    template_id = request.GET.get("template_id")
    context["service"] = service.objects.get(pk=int(service_id))
    q = request.GET.get("q")
    page = request.GET.get("page")
    if not page:
        page = 1
    if host_id:
        context["host"] = context["service"].host
    if template_id:
        context["template"] = context["service"].template
    context["items"] = context["service"].item_set.all().order_by('key')
    if q:
        context["items"] = context["items"].filter(
            Q(key__contains=q)
        )
    context["items"] = paging(context["items"], page, settings.PAGE_SIZE)
    return render_to_response(
        "host/item_list.html",
        context ,
        context_instance=RequestContext(request)
    )

def ItemUpdateView(request, pk):
    context = {}
    context["methods"] = methods
    context["symbols"] = symbols
    template_name = "host/item_edit.html"
    host_id = request.GET.get("host_id")
    template_id = request.GET.get("template_id")
    service_id = request.GET.get("service_id")
    _item = item.objects.get(pk=int(pk))
    context["service"] = service.objects.get(pk=int(service_id))
    context["item"] = _item
    context["dts"] = ("GAUGE","COUNTER","DERIVE")
    redirect_url = "/host/item/"
    if host_id:
        context["host"] = context["service"].host
        redirect_url += "?host_id=%s&service_id=%s"%(host_id,service_id)
    if template_id:
        context["template"] = context["service"].template
        redirect_url += "?template_id=%s&service_id=%s"%(template_id,service_id)
    if request.method == "GET":
        return render_to_response(
            template_name,
            context,
            context_instance = RequestContext(request)
        )
    else:
        _item.key = request.POST.get("key")
        _item.cycle = request.POST.get("cycle")
        _item.method = request.POST.get("method")
        _item.symbol = request.POST.get("symbol")
        _item.threshold = request.POST.get("threshold")
        _item.attempt = request.POST.get("attempt")
        _item.number = request.POST.get("number")
        _item.dt = request.POST.get("dt")
        _item.floatingthreshold = request.POST.get("floatingthreshold")
        _item.drawing = True if request.POST.get("drawing") == "1" else False
        _item.alarm = request.POST.get("alarm")
        _item.save()
        context["message"] = "修改成功"
        return HttpResponseRedirect(redirect_url)

def ItemCreateView(request):
    context = {}
    context["methods"] = methods
    context["symbols"] = symbols
    context["dts"] = ("GAUGE","COUNTER","DERIVE")
    host_id = request.GET.get("host_id")
    service_id = request.GET.get("service_id")
    template_id = request.GET.get("template_id")
    redirect_url = "/host/item/"
    if service_id:
        context["service"] = service.objects.get(pk=int(service_id))
    if host_id:
        context["host"] = host.objects.get(pk=int(host_id))
        redirect_url += "?host_id=%s&service_id=%s"%(host_id,service_id)
    if template_id:
        context["template"] = template.objects.get(pk=int(template_id))
        redirect_url += "?template_id=%s&service_id=%s"%(template_id,service_id)
    if request.method == "GET":
        return render_to_response(
            "host/item_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        for key in request.POST.get("key").split(","):
            item(
                service = context["service"],
                key = key,
                cycle = request.POST.get("cycle"),
                method = request.POST.get("method"),
                symbol = request.POST.get("symbol"),
                number = request.POST.get("number"),
                threshold = request.POST.get("threshold"),
                floatingthreshold = request.POST.get("floatingthreshold"),
                attempt = True if request.POST.get("attempt") == "1" else False,
                drawing = True if request.POST.get("drawing") == "1" else False
            ).save()
    return HttpResponseRedirect(redirect_url)


def PortListView(request):
    if request.method == "POST":
        return HttpResponse("BAD REQUEST")
    context={}
    template_name = "host/port_list.html"
    host_id = request.GET.get("host_id")
    q = request.GET.get("q")
    context["host"] = host.objects.get(pk=int(host_id))
    context["ports"] = context["host"].port_set.all()
    if q:
        context["ports"] = context["ports"].filter(
            Q(port__contains=q)|
            Q(proc_name__contains=q)
        )
    return render_to_response(
        template_name,
        context,
        context_instance = RequestContext(request)
    )

def PortUpdateView(request,pk):
    context = {}
    template_name = "host/port_edit.html"
    _port = port.objects.get(pk=int(pk))
    context["host"] = _port.host
    if request.method == "GET":
        context["port"]= _port
        return render_to_response(
            template_name,
            context,
            context_instance = RequestContext(request)
        )
    else:
        _port.port = request.POST.get("port")
        _port.proc_name = request.POST.get("proc_name")
        _port.alias = request.POST.get("alias")
        _port.alarm = request.POST.get("alarm")
        _port.save()
        context["message"] = "修改成功"
        return HttpResponseRedirect(
            "/host/port/?host_id=%s"%(context["host"].id)
        )

def PortCreateView(request):
    context = {}
    host_id = request.GET.get("host_id")
    if host_id:context["host"] = host.objects.get(pk=int(host_id))
    if request.method == "GET":
        return render_to_response("host/port_edit.html", context, context_instance = RequestContext(request))
    else:
        port.objects.create(
            host = context["host"],
            port = request.POST.get("port"),
            proc_name = request.POST.get("proc_name"),
            alias = request.POST.get("alias"),
            alarm = request.POST.get("alarm")
        ).save()
    return HttpResponseRedirect("/host/port/?host_id=%s"%(context["host"].id))

def TemplateListView(request):
    context = {}
    context["templates"] = template.objects.all()
    return render_to_response(
        "host/template_list.html",
        context ,
        context_instance=RequestContext(request)
    )

def TemplateUpdateView(request,pk):
    context = {}
    template_name = "host/template_edit.html"
    _template = template.objects.get(pk=int(pk))
    if request.method == "GET":
        context["template"]= _template
        context["groups"] = group.objects.all()
        context["hosts"] = host.objects.all()
        return render_to_response(
            template_name,
            context,
            context_instance = RequestContext(request)
        )
    else:
        _template.name = request.POST.get("name")
        _template.group_set.clear()
        for g_id in request.POST.getlist("groups"):
            _template.group_set.add(group.objects.get(pk=int(g_id)))
        _template.host_set.clear()
        for h_id in request.POST.getlist("hosts"):
            _template.host_set.add(host.objects.get(pk=int(h_id)))
        _template.save()
        context["message"] = "修改成功"
        return HttpResponseRedirect(
            "/host/template/"
        )

def TemplateCreateView(request):
    context = {}
    context["hosts"] = host.objects.all()
    context["groups"] = group.objects.all()
    if request.method == "GET":
        return render_to_response(
            "host/template_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _template = template.objects.create(
            name = request.POST.get("name")
        )
        for g_id in request.POST.getlist("groups"):
            _template.group_set.add(group.objects.get(pk=int(g_id)))
        for h_id in request.POST.getlist("hosts"):
            _template.host_set.add(host.objects.get(pk=int(h_id)))
        _template.save()
    return HttpResponseRedirect("/host/template/")

#主机资产修改
def ServerUpdateView(request, pk):
    context={}
    context["cabinets"] = cabinet.objects.all()
    context["units"] = [1,2,3,4]
    _server = server.objects.get(pk=int(pk))
    host_id = request.GET.get("host_id")
    context["host"] = host.objects.get(pk=int(host_id))
    if request.method == "GET":
        context["server"] = _server
        return render_to_response(
            "host/server_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _server.sn = request.POST.get("sn")
        unit = request.POST.get("unit")
        location = request.POST.get("location")
        _server.vender = request.POST.get("vender")
        _server.model = request.POST.get("model")
        cabinet_id = request.POST.get("cabinet")
        if cabinet_id:
            #获取最新的机柜
            _cabinet = cabinet.objects.get(pk=int(cabinet_id))
            start = int(location)
            end = start + int(unit)
            #过滤是否有已占用U位
            if _cabinet.unit_isused(start, end):
                #有任何一个u位被使用则返回错误消息
                context["message"] = "您选择的U位已有占用，请重新选择"
                _server.location = location
                _server.unit = int(unit)
                _server.cabinet= _cabinet
                context["server"] = _server
                return render_to_response(
                    "host/server_edit.html",
                    context,
                    context_instance = RequestContext(request)
                )
            if _server.has_cabinet():
                #释放之前占用的U位
                _server.cabinet.clear_unit(_server.location, _server.location + _server.unit)
            #设置最新U位占用状态
            _server.cabinet = _cabinet
            _server.location = location
            _server.unit = unit
            _server.cabinet.set_unit(start, end)
        else:
            _server.cabinet.clear_unit(_server.location, _server.location + _server.unit)
            _server.location=0
            _server.unit=0
            _server.cabinet = None
        _server.save()
        return HttpResponseRedirect(
            "/host/host/?host_id=%s"%host_id
        )

def GroupListView(request):
    context = {}
    context["groups"] = group.objects.all()
    return render_to_response(
        "host/group_list.html",
        context,
        context_instance = RequestContext(request)
    )


def GroupUpdateView(request,pk):
    context = {}
    context["group"] = group.objects.get(pk=int(pk))
    context["hosts"] = host.objects.all()
    context["templates"] = template.objects.all()
    if request.method == "GET":
        return render_to_response(
            "host/group_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _group = context["group"]
        _group.name = request.POST.get("name")
        hosts = request.POST.getlist("hosts")
        templates = request.POST.getlist("templates")
        _group.host_set.clear()
        for host_id in hosts:
            _group.host_set.add(host.objects.get(pk=int(host_id)))
        _group.template.clear()
        if templates:
            for template_id in templates:
                _group.template.add(template.objects.get(pk=int(template_id)))
        _group.save()
        return HttpResponseRedirect(
            "/host/group/"
        )

def GroupCreateView(request):
    context = {}
    context["hosts"] = host.objects.all()
    context["templates"] = template.objects.all()
    if request.method=="GET":
        return render_to_response(
            "host/group_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _group=group(
            name = request.POST.get("name")
        )
        _group.save()
        hosts = request.POST.getlist("hosts")
        templates = request.POST.getlist("templates")
        if hosts:
            for host_id in hosts:
                _group.host_set.add(host.objects.get(pk=int(host_id)))
        if templates:
            for template_id in templates:
                _group.template.add(template.objects.get(pk=int(template_id)))
        return HttpResponseRedirect("/host/group/")

def GroupDeleteView(request):
    context = {}
    ids = request.POST.get("ids")
    if request.user.is_superuser is False:
        context["status"] = 1
        context["message"] = "权限拒绝.."
    else:
        for id in ids.split(","):
            group.objects.get(pk=int(id)).delete()
        context["status"] = 0
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )