# coding:utf-8
# Create your views here.
from django.views.generic import ListView, UpdateView, CreateView
from django.shortcuts import HttpResponse, Http404,HttpResponseRedirect,RequestContext
from django.shortcuts import render_to_response
from django.contrib.auth.models import Group
from django.db.models import Q

from django.conf import settings
from network.models import *
from assest.models import *
from system.models import proxy
from host.models import item
from host.views import symbols,methods
from system.views import paging
import json


keys =  (
        "inOctets",
        "outOctets",
        "inUcastPkts",
        "outUcastPkts",
        "inDiscards",
        "outDiscards",
        "inErrors",
        "outErrors",
        "inUnknownProtos",
    )

def getAvailableUintsByID(request):
    cabinet_id = request.POST.get("cabinet_id")
    try:
        units = cabinet.objects.get(pk=int(cabinet_id)).unit_set.all()
    except:
        raise Http404
    context = {}
    for u in units:
        context[u.number] = u.used
    return HttpResponse(json.dumps(context), content_type="application/json")

def DeviceListView(request):
    context = {}
    context["proxys"] = proxy.objects.all()
    context["cabinets"] = cabinet.objects.all()
    context["units"] = [1,2,3,4]
    context["devices"] = device.objects.all()
    q =request.GET.get("q")
    page = request.GET.get("page")
    if not page:
        page =1
    if q:
        context["devices"]=context["devices"].filter(
            Q(ip__contains=q)|
            Q(sn__contains=q)|
            Q(cabinet__name__contains=q)|
            Q(model__contains=q)
        )
    context["devices"] = paging(context["devices"], page, settings.PAGE_SIZE)
    return render_to_response(
        "network/device_list.html",
        context,
        context_instance = RequestContext(request)
    )

def DeviceCreateView(request):
    context = {}
    context["proxys"] = proxy.objects.all()
    context["cabinets"] = cabinet.objects.all()
    context["units"] = [1,2,3,4]
    context["groups"] = Group.objects.all()
    if request.method == "GET":
        return render_to_response(
            "network/device_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        import uuid
        cabinet_id = request.POST.get("cabinet")
        proxy_id = request.POST.get("proxy")
        _device=device(
            uuid=uuid.uuid1(),
            sn = request.POST.get("sn"),
            vender = request.POST.get("vender"),
            model = request.POST.get("model"),
            location = int(request.POST.get("location")),
            unit = int(request.POST.get("unit")),
            ip = request.POST.get("ip"),
            snmp_community = request.POST.get("snmp_community"),
            snmp_port = request.POST.get("snmp_port"),
            config_update_interval = request.POST.get("config_update_interval"),
            check_interval = request.POST.get("check_interval")
        )
        if cabinet_id:
            _device.cabinet = cabinet.objects.get(pk=int(cabinet_id))
            if _device.cabinet.unit_isused(
                    _device.location,
                    _device.unit + _device.location
            ):
                context["message"] = "您选择的U位已有占用，请重新选择"
                context["device"] = _device
                return render_to_response(
                    "network/device_edit.html",
                    context,
                    context_instance = RequestContext(request)
                )
            _device.cabinet.set_unit(_device.location, _device.location + _device.unit)

        else:
            _device.location = 0
            _device.unit = 0
        if proxy_id:
            _device.proxy = proxy.objects.get(pk=int(proxy_id))
        groups = request.POST.getlist("groups")
        if groups:
            for g_id in  groups:
                _device.group.add(Group.objects.get(pk=int(g_id)))
        _device.save()

        return HttpResponseRedirect(
            "/network/device/"
        )



def DeviceUpdateView(request, pk):
    context = {}
    context["proxys"] = proxy.objects.all()
    context["cabinets"] = cabinet.objects.all()
    context["units"] = [1,2,3,4]
    context["groups"] = Group.objects.all()
    _device = device.objects.get(pk=int(pk))
    if request.method == "GET":
        context["device"]=_device
        return  render_to_response(
            "network/device_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        location = int(request.POST.get("location"))
        unit = int(request.POST.get("unit"))
        _device.sn = request.POST.get("sn")
        _device.vender = request.POST.get("vender")
        _device.model = request.POST.get("model")
        _device.ip = request.POST.get("ip")
        _device.snmp_community = request.POST.get("snmp_community")
        _device.snmp_port = request.POST.get("snmp_port")
        _device.config_update_interval = request.POST.get("config_update_interval")
        _device.check_interval = request.POST.get("check_interval")
        cabinet_id = request.POST.get("cabinet")
        proxy_id = request.POST.get("proxy")
        if cabinet_id:
            _cabinet = cabinet.objects.get(pk=int(cabinet_id))
            #机柜和U位都没有变更，则不进行验证
            if _device.cabinet == _cabinet and _device.location == int(location) and _device.unit == int(unit):
                pass
            elif _cabinet.unit_isused(
                location,
                unit + location
            ):
                context["message"] = "您选择的U位已有占用，请重新选择"
                _device.location = int(location)
                _device.unit = int(unit)
                context["device"] = _device
                return render_to_response(
                    "network/device_edit.html",
                    context,
                    context_instance = RequestContext(request)
                )
            if _device.cabinet:
                _device.cabinet.clear_unit(_device.location, _device.location + _device.unit)
                _device.location=0
                _device.unit = 0
            _device.cabinet = _cabinet
            _device.location  = location
            _device.unit = unit
            _device.cabinet.set_unit(location, unit+location)
        else:
            _device.cabinet.clear_unit(_device.location, _device.location + _device.unit)
            _device.location = 0
            _device.unit=0
            _device.cabinet=None
        if proxy_id:
            _device.proxy = proxy.objects.get(pk=int(proxy_id))
        else:
            _device.proxy = None
        _device.group.clear()
        groups = request.POST.getlist("groups")
        if groups:
            for g_id in  groups:
                _device.group.add(Group.objects.get(pk=int(g_id)))
        _device.save()
        return HttpResponseRedirect(
            "/network/device/"
        )

def InterfaceListView(request):
    context = {}
    device_id = request.GET.get("device_id")
    q = request.GET.get("q")
    page = request.GET.get("page")
    if not page:
        page = 1
    if request.method == "GET":
        context["device"] = device.objects.get(pk=int(device_id))
        context["interfaces"] = context["device"].interface_set.all().order_by("-status", "index")
        if q:
            context["interfaces"] = context["interfaces"].filter(
                Q(name__contains=q)
            )
        context["interfaces"] = paging(context["interfaces"], page, settings.PAGE_SIZE)
        return render_to_response(
            "network/interfaces_list.html",
            context,
            context_instance = RequestContext(request)
        )

def CustomOIDListView(request):
    context = {}
    device_id = request.GET.get("device_id")
    context["device"] = device.objects.get(pk=int(device_id))
    context["oids"] = context["device"].oid_set.all()
    return render_to_response(
        "network/custom_oid_list.html",
        context,
        context_instance = RequestContext(request)
    )

def CustomOIDCreateView(request):
    context = {}
    device_id = request.GET.get("device_id")
    context["device"] = device.objects.get(pk=int(device_id))
    if request.method=="GET":
        return render_to_response(
            "network/custom_oid_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        oid(
            device=context["device"],
            name=request.POST.get("name"),
            oid=request.POST.get("oid")
        ).save()
        return HttpResponseRedirect(
            "/network/oid/?device_id=%s"%device_id
        )
def CustomOIDUpdateView(request, pk):
    context = {}
    _oid = oid.objects.get(pk=int(pk))
    device_id = request.GET.get("device_id")
    if request.method == "GET":
        context["oid"] = _oid
        context["device"] = device.objects.get(pk=int(device_id))
        return render_to_response(
            "network/custom_oid_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _oid.name = request.POST.get("name")
        _oid.oid = request.POST.get("oid")
        _oid.save()
        return HttpResponseRedirect(
            "/network/oid/?device_id=%s" % device_id
        )

def ItemListView(request):
    context = {}
    device_id = request.GET.get("device_id")
    interface_id = request.GET.get("interface_id")
    oid_id = request.GET.get("oid_id")
    context["methods"] = methods
    context["symbols"] = symbols
    context["device"] = device.objects.get(pk=int(device_id))
    if interface_id:
        context["interface"] = interface.objects.get(pk=int(interface_id))
        context["items"]= context["interface"].item_set.all()
    if oid_id:
        context["oid"] = oid.objects.get(pk=int(oid_id))
        context["items"] = context["oid"].item_set.all()
    #context["items"] = item.objects.all()
    return render_to_response(
        "network/item_list.html",
        context,
        context_instance = RequestContext(request)
    )

def ItemCreateView(request):
    context = {}
    context["methods"] = methods
    context["symbols"] = symbols
    context["keys"] = keys
    interface_id = request.GET.get("interface_id")
    oid_id = request.GET.get("oid_id")
    device_id = request.GET.get("device_id")
    context["device"] = device.objects.get(pk=int(device_id))
    if interface_id:
        context["interface"] = interface.objects.get(pk=int(interface_id))
    if oid_id:
        context["oid"] = oid.objects.get(pk=int(oid_id))
    if request.method == "GET":
        return render_to_response(
            "network/item_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        redirect_url = "/network/item/?device_id=%s&"%device_id
        _item = item(
            #key = request.POST.get("key"),
            cycle = request.POST.get("cycle"),
            method = request.POST.get("method"),
            symbol = request.POST.get("symbol"),
            threshold = request.POST.get("threshold"),
            attempt = request.POST.get("attempt"),
            alarm = request.POST.get("alarm")
        )
        if interface_id:
            _item.key = request.POST.get("key")
            _item.interface = context["interface"]
            if item.objects.filter(interface=_item.interface, key=_item.key).count()>0:
                context["item"] = _item
                context["message"] = u"重复的key: %s"%_item.key
                return render_to_response(
                    "network/item_edit.html",
                    context,
                    context_instance = RequestContext(request)
                )
            redirect_url += "interface_id=%s" % interface_id
        if oid_id:
            _item.oid = context["oid"]
            _item.key = _item.oid.name
            redirect_url += "oid_id=%s" % oid_id
        _item.save()
        return HttpResponseRedirect(
            redirect_url
        )

def ItemUpdateView(request, pk=None):
    context = {}
    context["methods"] = methods
    context["symbols"] = symbols
    context["keys"] = keys
    interface_id = request.GET.get("interface_id")
    oid_id = request.GET.get("oid_id")
    device_id = request.GET.get("device_id")
    redirect_url = "/network/item/?device_id=%s&"%device_id
    context["device"] = device.objects.get(pk=int(device_id))
    if interface_id:
        context["interface"] = interface.objects.get(pk=int(interface_id))
        redirect_url += "interface_id=%s" % interface_id
    if oid_id:
        context["oid"] = oid.objects.get(pk=int(oid_id))
        redirect_url += "oid_id=%s" % oid_id
    _item = item.objects.get(pk=int(pk))
    if request.method == "GET":
        context["item"] = _item
        return render_to_response(
            "network/item_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _item.cycle = request.POST.get("cycle")
        _item.method = request.POST.get("method")
        _item.symbol = request.POST.get("symbol")
        _item.threshold = request.POST.get("threshold")
        _item.attempt = request.POST.get("attempt")
        _item.alarm = request.POST.get("alarm")
        if interface_id:
            _item.key = request.POST.get("key")
        _item.save()
        return HttpResponseRedirect(
            redirect_url
        )
