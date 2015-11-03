# coding:utf-8
from __future__ import unicode_literals
from notifi import alarm
from host.models import *
from network.models import *
from django.contrib.auth.models import *
from round_robin import Round_Robin
from django.http import HttpResponse
from datetime import datetime
from curlmulti import curlmulti_tsdb
from celery import task
from djcelery.models import *
from datetime import datetime
from gevent import monkey
monkey.patch_all()

@task(bind=True)
def add_task(self):
	dt = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
	hosts = host.objects.all()
	networks = device.objects.all()

	t_uuid = [p.name for p in PeriodicTask.objects.all()]

	i = IntervalSchedule.objects.filter(every=5, period="minutes")
	if not i:
		i = IntervalSchedule.objects.create(every=5, period="minutes")
	
	interval = IntervalSchedule.objects.get(every=5, period="minutes")

	if hosts:
		for h in hosts:
			#主机状态0, 报警状态0
			if h.status == 0 and h.alarm == 0:
				# 报警
				if h.uuid not in t_uuid:
					PeriodicTask.objects.create(name=h.uuid, task="task.curlmulti.curlmulti_tsdb", interval=interval, args="[\"%s\"]" % (h.uuid))

				ports = h.port_set.filter(status=1, alarm=0)
				if ports:
					for p in ports:	
						content = "critical: %s 主机:%s 进程名:%s 端口:%s down"  % (dt, h.ip, p.alias if p.alias else p.proc_name, p.port)
						for g in h.group.all():
							alarm(content, Group.objects.filter(name=g.name))
				
			elif h.status == 0  and h.alarm == 1:
				if h.uuid in t_uuid:
					PeriodicTask.objects.get(name=h.uuid).delete()

			elif h.status == 1 and h.alarm == 0:
				#报警
				if h.uuid in t_uuid:
					PeriodicTask.objects.get(name=h.uuid).delete()
				content = "critical: %s 主机:%s down" % (dt, h.ip)
				print content
				for g in h.group.all():
					alarm(content, Group.objects.filter(name=g.name))
				
			elif h.status == 1 and h.alarm == 1:
				if h.uuid in t_uuid:
					PeriodicTask.objects.get(name=h.uuid).delete()

	if networks:
		for n in networks:
			if n.alarm == 0:
				if n.uuid not in t_uuid:	
					PeriodicTask.objects.create(name=n.uuid, task="task.curlmulti.network_curlmulti_tsdb", interval=interval, args="[\"%s\"]" % (n.uuid))
			else:
				if n.uuid in t_uuid:	
					PeriodicTask.objects.get(name=n.uuid).delete()
	print self.request

@task(bind=True)
def check_service_key(self):
        hosts = host.objects.all()
        for h in hosts:
		if "基础信息采集模板" not in [t.name for t in h.template.all()]:
			_t = template.objects.get(name="基础信息采集模板")
			h.template.add(_t)

                hw_mems = h.service_set.filter(name="hw.memory.status")
                if hw_mems:
                        for i in hw_mems:
                                mems = i.item_set.all()
                                for m in mems:
					if not m.method:
						m.cycle = 5
						m.method = "max"
						m.symbol = "="
						m.threshold = 1
						m.save()

                disk_mems = h.service_set.filter(name="hw.disk.status")
                if disk_mems:
                        for i in disk_mems:
                                disks = i.item_set.all()
                                for d in disks:
					if not d.method:
						d.cycle = 5
						d.method = "max"
						d.symbol = "="
						d.threshold = 1
						d.save()

		disk_percent = h.service_set.filter(name="disk.used.precent")
                if disk_percent:
                        for i in disk_percent:
                                dirs = i.item_set.all()
                                for d in dirs:
                                        if not d.method:
						d.cycle = 5
                                                d.method = "max"
                                                d.symbol = ">"
                                                d.threshold = 93 
                                                d.save()

		for s in h.service_set.all():
			g = Group.objects.get(name="duty")
			if g not in s.group.all():
				s.group.add(g)
	
	g = group.objects.get(name="10GNic")
	g_hosts = g.host_set.all()
	for h in g_hosts:
	    for s in h.service_set.all():
		if s.name == 'net.eth0':
		    for i in s.item_set.all():
			if i.key == 'bytes_recv' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
			elif i.key == 'bytes_sent' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
		elif s.name == 'net.eth1':
		    for i in s.item_set.all():
			if i.key == 'bytes_recv' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
			elif i.key == 'bytes_sent' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
		elif s.name == 'net.em1':
		    for i in s.item_set.all():
			if i.key == 'bytes_recv' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
			elif i.key == 'bytes_sent' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
		elif s.name == 'net.em2':
		    for i in s.item_set.all():
			if i.key == 'bytes_recv' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()
			elif i.key == 'bytes_sent' and i.threshold != 1000000000:
			    i.threshold = 1000000000
			    i.save()

	print self.request
