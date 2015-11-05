# coding:utf8
from django.shortcuts import render, render_to_response,RequestContext
from django.contrib.auth.models import Group
from django.views.decorators.csrf import csrf_exempt
from django.http import HttpResponse 
from host.models import *
from network.models import *
import re
import json

def item_update(h, service, item, val, counter, floatingvalue, alert):
	s = h.service_set.get(name=service)
	i = s.item_set.get(key=item)
	i.current = val
	i.counter = counter
	i.floatingvalue = floatingvalue
	i.alarm = alert
	i.save(update_fields=['current', 'counter', 'floatingvalue', 'alarm'])

def item_network_update(n, interface, item, val, counter, floatingvalue):
	i = n.interface_set.get(name=interface)
	_i = i.item_set.get(key=item)
	_i.current = val
	_i.counter = counter
	_i.floatingvalue = floatingvalue
	_i.save(update_fields=['current', 'counter', 'floatingvalue'])

def appmonitor(request):
	groups = Group.objects.all()
	return render_to_response('appmonitor.html', {"groups": groups},context_instance = RequestContext(request))

def check_url(url):
    match = re.match(r'http:/{2}(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):\d+/.*', url)
    if match:
        return 1
    else:
        return 0

def Unit(interval, unit):
    second = 0
    try:
        if type(interval) == int:
            if unit == 'D':
                second = interval * 60 * 60 * 24
            elif unit == 'H':
                second = interval * 60 * 60
            elif unit == 'M':
                second = interval * 60
            elif unit == 'S':
                second = interval

        return second
    except:
        return "interval type error"

@csrf_exempt
def addAppMonitor(request):
    if request.method == 'POST':
        url = request.POST.get('url', '')
        content = request.POST.get('content', '')
	groups = request.POST.getlist('groups', '')
	keys = request.POST.getlist('key', '')
	cycles = request.POST.getlist('cycle', '')
	methods = request.POST.getlist('method', '')
	symbols = request.POST.getlist('symbol', '')
	thresholds = request.POST.getlist('threshold', '')
	drawings = request.POST.getlist('drawing', '')
	alarms = request.POST.getlist('alarm', '')
        interval = request.POST.get('interval', '')
        unit = request.POST.get('unit', '')
        if check_url(url):
            s = Unit(int(interval), unit)
            status_info = auto_add_metric(url, content, s, groups, keys, cycles, methods, symbols, thresholds, drawings, alarms)
            if status_info:
                return HttpResponse("%s" % (status_info))
            else:
                return HttpResponse("success")
        else:
            return HttpResponse("the url don't match")

    else:
        return HttpResponse("the request method don't post")

def auto_add_metric(url, content, task_interval, groups, keys, cycles, methods, symbols, thresholds, drawings, alarms):
    try:
        datas = eval(content)
        if type(datas) == dict:
            if "host" in datas.keys():
                ip = datas['host']
                _h = host.objects.filter(ip=ip)
		if _h:
			h = host.objects.get(ip=ip)
			if "app_name" in datas.keys() and "domain" in datas.keys():
			    service_name = "%s.%s" % (datas["app_name"], datas["domain"])
			    if h.service_set.filter(name=service_name):
				s = h.service_set.get(name=service_name)
				if keys:
					for k, v in enumerate(keys):
						if s.item_set.filter(key=v):
							i = s.item_set.get(key=v)
							i.cycle = cycles[k]
							i.method = methods[k]
							i.symbol = symbols[k]
							i.threshold = thresholds[k]
							if drawings[k] == "1":
								i.drawing = False
							else:
								i.drawing = True
							i.alarm = alarms[k]
							i.save(update_fields=['cycle', 'method', 'symbol', 'threshold', 'drawing', 'alarm'])
						else:
							i = item.objects.create(key=v, cycle=cycles[k], method=methods[k], symbol=symbols[k], threshold=thresholds[k], drawing=drawings[k], alarm=alarms[k])
							s.item_set.add(i)

				if groups:
					for g in groups:
						group = Group.objects.get(name=g)
						s.group.add(group)
			    else:
				s = service.objects.create(name=service_name, plugin="appmonitor.py", args=url, exec_interval=int(task_interval), status=0)
				h.service_set.add(s)
				if keys:
					for k, v in enumerate(keys):
						if s.item_set.filter(key=v):
							i = s.item_set.get(key=v)
							i.cycle = cycles[k]
							i.method = methods[k]
							i.symbol = symbols[k]
							i.threshold = thresholds[k]
							i.drawing = drawings[k]
							i.alarm = alarms[k]
							i.save(update_fields=['cycle', 'method', 'symbol', 'threshold', 'drawing', 'alarm'])
						else:
							i = item.objects.create(key=v, cycle=cycles[k], method=methods[k], symbol=symbols[k], threshold=thresholds[k], drawing=drawings[k], alarm=alarms[k])
							s.item_set.add(i)

				if groups:
					for g in groups:
						if not s.group.filter(name=g):
							group = Group.objects.get(name=g)
							s.group.add(group)
		else:
			return "the host don't exists"
					
            else:
                    return "the json don't exists key (app_name) and (domain)"
        else:
            return "the json data type is wrong"
    except:
            return "the json data type is wrong"

def alert_data(request):
	if request.method == 'POST':
		g = request.POST.get('group', '')
		h = request.POST.get('host', '')
		datas = []
		if g and h:
			_g = group.objects.get(name=g)
			for _h in _g.host_set.all():
				if _h.ip == h:
					for _s in _h.service_set.all():
						for _i in _s.item_set.all():
							metric = "%s.%s" % (_s.name, _i.key)
							if abs(_i.threshold) > abs(_i.val) or abs(_i.floatingvalue) > abs(_i.val):
								datas.append({'id': _i.id, 'ip': _h.ip, 'metric': metric, 'val': val, 'threshold': _i.threshold, 'floatvalue': _i.floatingvalue})
		elif g and not h:
			_g = group.objects.get(name=g)
			for _h in _g.host_set.all():
				for _s in _h.service_set.all():
					for _i in _s.item_set.all():
						metric = "%s.%s" % (_s.name, _i.key)
						if abs(_i.threshold > _i.val) or abs(_i.floatingvalue) > abs(_i.val):
							datas.append({'id':_i.id, 'ip': _h.ip, 'metric': metric, 'val': val, 'threshold': _i.threshold, 'floatvalue': _i.floatingvalue})
		elif not g and h:
			_h = host.objects.get(ip=h)
			for _s in _h.service_set.all():
				for _i in _s.item_set.all():
					metric = "%s.%s" % (_s.name, _i.key)
					if abs(_i.threshold) > abs(_i.val) or abs(_i.floatingvalue) > abs(_i.val):
						datas.append({'id':_i.id, 'ip': _h.ip, 'metric': metric, 'val': val, 'threshold': _i.threshold, 'floatvalue': _i.floatingvalue})
		else:
			for _h in host.objects.all():
				for _s in _h.service_set.all():
					for _i in _s.item_set.all():
						metric = "%s.%s" % (_s.name, _i.key)
						if abs(_i.threshold) > abs(_i.val) or abs(_i.floatingvalue) > abs(_i.val):
							datas.append({'id':_i.id, 'ip': _h.ip, 'metric': metric, 'val': val, 'threshold': _i.threshold, 'floatvalue': _i.floatingvalue})

        	return HttpResponse(json.dumps(datas, indent=4))
	else:
		return HttpResponse(json.dumps({'status':1}, indent=4))

def acknowledged(request):
	if request.method == 'POST':
		ids = request.POST.getlist('id', '')
		if ids:
			for _id in ids:
				i = item.objects.get(id=_id)
				i.alarm = 2
				i.save()

		return HttpResponse(json.dumps({'status':0}, indent=4))
	else:
		return HttpResponse(json.dumps({'status':1}, indent=4))
