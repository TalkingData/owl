#coding:utf8
# Create your views here.

from django.shortcuts import render_to_response,RequestContext
from django.http import HttpResponse,HttpResponseRedirect
from django.conf import settings
from django.db.models import Q
from host.models import *
from network.models import *
from .models import *
from system.views import paging
import time
import json

NetMertic = [
    "inOctets",
    "outOctets",
    "inUcastPkts",
    "outUcastPkts",
    "inDiscards",
    "outDiscards",
    "inErrors",
    "outErrors",
    "inUnknownProtos"
]

def GetIndexView(request):
    uuid = request.GET.get("uuid")
    service_id = request.GET.get("service_id")
    interface_id = request.GET.get("interface_id")
    oid_id = request.GET.get("oid_id")
    context={}
    context["metrics"] = []
    template_name = "graph/graph.html"
    if service_id:
        context["host"] = host.objects.get(uuid=uuid)
        context["service"] = service.objects.get(pk=int(service_id))
    elif interface_id:
        context["device"] = device.objects.get(uuid=uuid)
        context["interface"] = interface.objects.get(pk=int(interface_id))
        context["metric"] = NetMertic
    elif oid_id:
        _oid=oid.objects.get(pk=int(oid_id))
        context["device"] = device.objects.get(uuid=uuid)
        context["oid"] = _oid

    return render_to_response(
        template_name,
        context,
        context_instance = RequestContext(request)
    )


def get_data(request):
    import urllib2,json
    metric = request.GET.get("metric")
    uuid = request.GET.get("uuid")
    start = request.GET.get("start")
    end = request.GET.get("end")
    if not end:
        end = int(time.time())
    if not start:
        start = end - (1 * 24  * 60 * 60) #one week
    url = "http://%s/api/query?start=%s&end=%s&m=sum:%s{uuid=%s}" % (settings.OPENTSDB_ADDR, start, end, metric, uuid)
    req = urllib2.Request(url)
    data = []
    try:
        resp=urllib2.urlopen(req)
        result = json.load(resp)
        data = [[int(k)*1000,v] for k,v in result[0]["dps"].items()]
        data = sorted(data, key= lambda data:data[0])
    except urllib2.URLError,e:
        print e.reason

    return HttpResponse(
        json.dumps(data),
        content_type="application/json"
    )

def DashBoard(request):
    context = {}
    page = request.GET.get("page")
    q = request.GET.get("q")
    if not page:
        page = 1
    queryset = graph.objects.filter(user__username__exact=request.user.username)
    if q:
        for qs in q.split(" "):
            queryset = queryset.filter(
                 Q(title__contains=qs)
            )
    context["graphs"] = paging(queryset, page, 3)

    return render_to_response(
        "graph/dashboard.html",
        context,
        context_instance = RequestContext(request)
    )

def GraphListView(request):
    context = {}
    page = request.GET.get("page")
    q = request.GET.get("q")
    if not page:
        page = 1
    if request.user.is_superuser:
        queryset = graph.objects.all()
    else:
        queryset = graph.objects.filter(user__username__exact=request.user.username)
    if q:
        queryset=queryset.filter(
            Q(title__contains=q)
        )
    context["graphs"] = paging(queryset, page, 5)
    return render_to_response(
        "graph/graph_list.html",
        context,
        context_instance = RequestContext(request)
    )

def GraphCreateView(request):
    context = {}
    if request.method == "GET":
        return  render_to_response(
            "graph/graph_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _graph = graph(
            title=request.POST.get("title"),
            user =request.user
        )
        _graph.save()
        return HttpResponseRedirect(
            "/graph/metric/?graph_id=%s"%_graph.id
        )

def GraphDeleteView(request):
    context = {}
    ids = request.POST.get("ids")
    for id in ids.split(","):
        _g = graph.objects.get(pk=int(id))
        _g.delete()
        context["status"] = 0
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )


def MetricListView(request):
    context = {}
    graph_id = request.GET.get("graph_id")
    context["graph"] = graph.objects.get(pk=int(graph_id))
    context["metrics"] = context["graph"].metric_set.all()
    if not request.user.is_superuser and context["graph"].user != request.user:
        return HttpResponse(
            "权限拒绝"
        )
    page = request.GET.get("page")
    q = request.GET.get("q")
    if not page:
        page = 1
    if q:
        context["metrics"] = context["metrics"].filter(
            Q(metric__contains=q)
        )
    context["metrics"] = paging(context["metrics"],page, settings.PAGE_SIZE)
    return  render_to_response(
        "graph/metric_list.html",
        context,
        context_instance = RequestContext(request)
    )

def MetricCreateView(request):
    context = {}
    graph_id = request.GET.get("graph_id")
    context["graph"] = graph.objects.get(pk=int(graph_id))
    if request.method =="GET":
        return render_to_response(
            "graph/metric_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        t = request.POST.get("type")
        uuid = ""
        if t == "host":
            uuid = host.objects.get(ip=request.POST.get("ip")).uuid
        if t == "network":
            uuid = device.objects.get(ip=request.POST.get("ip")).uuid
        metrics = request.POST.getlist("metric")
        ip = request.POST.get("ip")
        for m in metrics:
            metric(
                graph=context["graph"],
                metric = m,
                uuid = uuid,
                ip = ip
            ).save()
        return HttpResponseRedirect(
            "/graph/metric/?graph_id=%s" % graph_id
        )


def MetricDeleteView(request):
    context = {}
    ids = request.POST.get("ids")
    context["status"] = 0
    for id in ids.split(","):
        _metric = metric.objects.get(pk=int(id))
        _metric.delete()
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )