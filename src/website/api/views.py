# coding:utf-8
from django.shortcuts import render
from django.http import HttpResponse
# Create your views here.

from host.models import host, item
from network.models import device

from network.views import keys

import json
def GetAllHostIP(request):
    context = {}
    context["ips"] = []
    for _host in host.objects.all():
        context["ips"].append(
            {
                "id" : _host.id,
                "uuid" : _host.uuid,
                "ip" : _host.ip
            }
        )
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )

def GetAllDeviceIP(request):
    context = {}
    context["ips"] = []
    for _device in device.objects.all():
        context["ips"].append({
            "id": _device.id,
            "uuid" : _device.uuid,
            "ip" : _device.ip
        })


    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )

def GetHostMetric(request,pk):
    context = {}
    context["metrics"] = []
    try:
        _host = host.objects.get(pk=int(pk))
        for _service in _host.service_set.all():
            for _item in _service.item_set.all():
                context["metrics"].append(_service.name + "." + _item.key)
    except:
        pass
    context["metrics"].sort()
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )

def GetDeviceMetric(request,pk):
    context = {}
    context["metrics"] = []
    try:
        _device = device.objects.get(pk=int(pk))
        print _device
        for _ift in _device.interface_set.all():
            for _key in keys:
                context["metrics"].append(_ift.name + "." + _key)
        for _oid in _device.oid_set.all():
            context["metrics"].append(_oid.name)
    except:
        pass

    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )


def MuitlChangeItems(request):
    context = {}
    ids = request.POST.get("ids")
    method = request.POST.get("method")
    symbol = request.POST.get("symbol")
    threshold = request.POST.get("threshold")
    attempt = request.POST.get("attempt")
    dt = request.POST.get("dt")
    drawing = request.POST.get("drawing")
    cycle = request.POST.get("cycle")
    floatingthreshold = request.POST.get("floatingthreshold")
    if len(ids) == 0:
        context["status"] = 0
        context["message"] =  "未选择要修改的内容"
    else:
        for id in ids.split(","):
            _item = item.objects.get(pk=int(id.strip()))
            if method:_item.method = method
            if symbol:_item.symbol = symbol
            if threshold:_item.threshold = threshold
            if attempt:_item.attempt = attempt
            if dt:_item.dt = dt
            if cycle:_item.cycle = cycle
            if floatingthreshold:_item.floatingthreshold = floatingthreshold
            if drawing:
                _item.drawing = True if drawing == "1" else False
            _item.save()
        context["status"] = 1
        context["message"] = "批量修改完成"
    return HttpResponse(
        json.dumps(context),
        content_type= "application/json"
    )