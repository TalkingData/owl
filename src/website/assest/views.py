from django.views.generic import ListView, CreateView, UpdateView
from .models import  *
from host.models import host
from network.models import *
from django.shortcuts import render_to_response, HttpResponse, get_object_or_404,RequestContext,HttpResponseRedirect
def AssestListView(request):
    context = {}


class AssestListView(ListView):
    model = idc
    template_name = "assest/assest_list.html"

    def get_context_data(self, **kwargs):
        context = super(AssestListView, self).get_context_data(**kwargs)
        context["cabinets"] = cabinet.objects.all()
        for c in context["cabinets"]:
            c.obj = {i:None for i in range(1, int(c.capacity)+1)}
            for device in c.device_set.all():
                if device.location + device.unit > c.capacity:continue
                c.obj[device.location] = device
                if device.unit == 2:del c.obj[device.location+1]
            for server in c.server_set.all():
                if server.location + server.unit > c.capacity:continue
                c.obj[server.location] = server
                if server.unit == 2:del c.obj[server.location+1]
        return context

class IdcListView(ListView):
    model = idc
    template_name = "assest/idc_list.html"
    context_object_name = "idcs"

class CabinetListView(ListView):
    model = cabinet
    template_name = "assest/cabinet_list.html"
    context_object_name = "cabinets"

class IdcUpdateView(UpdateView):
    model = idc
    template_name = "assest/idc_edit.html"
    success_url = "/assest/idc/"


class CabinetUpdateView(UpdateView):
    model = cabinet
    template_name = "assest/cabinet_edit.html"
    success_url = "/assest/cabinet/"

    def get_context_data(self, **kwargs):
        context = super(CabinetUpdateView, self).get_context_data(**kwargs)
        context["idcs"] = idc.objects.all()
        return context


class IdcCreateView(CreateView):
    model = idc
    template_name = "assest/idc_edit.html"
    success_url = "/assest/idc/"


def CabinetCreateView(request):
    context = {}
    if request.method == "GET":
        context["idcs"] = idc.objects.all()
        return render_to_response(
            "assest/cabinet_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    idc_id = request.POST.get("idc")
    name = request.POST.get("name")
    capacity = request.POST.get("capacity")
    _idc = None
    if idc_id:
        _idc = idc.objects.get(pk=int(idc_id))
    _cabinet = cabinet(
        idc=_idc,
        name=name,
        capacity=int(capacity)
    )
    _cabinet.save()
    for unit in range(1, _cabinet.capacity+1):
        _cabinet.units_set.create(unit=unit)
    return HttpResponseRedirect(
        "/assest/cabinet/"
    )
