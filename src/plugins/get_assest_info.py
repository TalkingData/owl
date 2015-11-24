#!/usr/bin/env python

import subprocess
try:
	import simplejson as json
except ImportError:
	import json

import socket
import re

s = {'nics':'', 'disks':'', 'memorys':'', 'cpus':'', 'ip':'', 'model':'', 'sn':'', 'vender':''}

nics = []
disks = []
memorys = []
cpus = []

model = ''
sn = ''
vender = ''

p = subprocess.Popen('dmidecode --type system', shell=True, stdout=subprocess.PIPE)
for line in p.stdout.readlines():
	data = line.strip().split(":",1)
	if len(data) == 2:
		d1 = data[0].strip()
		d2 = data[1].strip()
		if d1 == 'Product Name':
			model = d2
		elif d1 == 'Serial Number':
			sn = d2
		elif d1 == 'Manufacturer':
			vender = d2

s['model'] = model
s['sn'] = sn
s['vender'] = vender

if vender == 'HP':
	name = []
	mac = []

	p = subprocess.Popen("ifconfig -a | grep 'HWaddr' | awk '{print $1,$NF}", shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		data = line.strip().split(":",1)
		d1 = data[0].strip()
		d2 = data[1].strip()
		if d2:
			name.append(d1)
			mac.append(d2)

	for i in xrange(0, len(name)):
		nics.append({'des':'', 'name':name[i], 'vender':'', 'mac':mac[i] })

	s['nics'] = nics


	bus = []
	capacity = []
	production_date = []
	slot = []
	port = []
	box = []
	bay = []
	sn = []
	speed = []
	vender = []
	status = []

	p = subprocess.Popen('hpssacli ctrl slot=0 pd all show detail', shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		data = line.strip().split(":",1)
		if len(data) == 2:
			d1 = data[0].strip()
			d2 = data[1].strip()
			if d1 == 'Interface Type':
				bus.append(d2)
			elif d1 == 'Size':	
				m = re.match(r'(.*GB)', d2)
				n = m.group().replace(',', '')
				capacity.append(n)
			elif d1 == 'Port':
				port.append(d2)
			elif d1 == 'Box':
				box.append(d2)
			elif d1 == 'Bay':
				bay.append(d2)
			elif d1 == 'Serial Number':
				sn.append(d2)
			elif d1 == 'PHY Transfer Rate':
				m = re.match(r'(.*Gbps)', d2)
				n = m.group()
				speed.append(n)
			elif d1 == 'Model':
				vender.append(d2)
			elif d1 == 'Status':
				status.append(d2)

	for i in xrange(0, len(port)):
		number = '%s:%s:%s' % (port[i], box[i], bay[i])  
		slot.append(number)
	
	production_date = ''
	p = subprocess.Popen('dmidecode | grep "Date"', shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines(): 
		data = line.strip().split(":")
		production_date = '/'.join(data[1].strip().split('/')[::-1])

	for i in range(0, len(sn)):
		disks.append({'bus':bus[i], 'capacity':capacity[i], 'media':'', 'product_id':'', 'production_date':production_date, 'slot':slot[i], 'sn':sn[i], 'speed':speed[i], 'vender':vender[i], 'status':status[i]})	

	s['disks'] = disks
else:
	des = []
	name = []
	vender = []
	mac = []

	p1 = subprocess.Popen('/opt/dell/srvadmin/bin/omreport chassis nics', shell=True, stdout=subprocess.PIPE)
	for line in p1.stdout.readlines():
		data = line.strip().split(":")
		if len(data) == 2:
			d1 = data[0].strip()
			d2 = data[1].strip()
			if d1 == 'Interface Name':
				name.append(d2)
			elif d1 == 'Description':
				des.append(d2)
			elif d1  == 'Vendor':
				vender.append(d2)

	for i in xrange(0, len(name)):
		p = subprocess.Popen("ifconfig %s | perl -anle '/(?:ether|HWaddr)\s+(\w+\:\w+\:\w+:\w+\:\w+\:\w+)/ and print $1'" % (name[i],), shell=True, stdout=subprocess.PIPE)
		nics.append({'des':des[i], 'name':name[i], 'vender':vender[i], 'mac': p.stdout.readlines()[0].strip()})

	s['nics'] = nics

	bus = []
	capacity = []
	media = []
	product_id = []
	production_date = []
	production_day = []
	production_week = []
	production_year = []
	slot = []
	sn = []
	speed = []
	vender = []
	status = []

	p = subprocess.Popen('/opt/dell/srvadmin/bin/omreport storage pdisk controller=0', shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		data = line.strip().split(":",1)
		if len(data) == 2:
			d1 = data[0].strip()
			d2 = data[1].strip()
			if d1 == 'Bus Protocol':
				bus.append(d2)
			elif d1 == 'Capacity':	
				m = re.match(r'(.*GB)', d2)
				n = m.group().replace(',', '')
				capacity.append(n)
			elif d1 == 'Media':
				media.append(d2)
			elif d1 == 'Product ID':
				product_id.append(d2)
			elif d1 == 'Manufacture Day':
				production_day.append(d2)
			elif d1 == 'Manufacture Week':
				production_week.append(d2)
			elif d1 == 'Manufacture Year':
				production_year.append(d2)
			elif d1 == 'ID':
				slot.append(d2)
			elif d1 == 'Serial No.':
				sn.append(d2)
			elif d1 == 'Capable Speed':
				speed.append(d2)
			elif d1 == 'Vendor ID':
				vender.append(d2)
			elif d1 == 'Status':
				status.append(d2)

	for i in xrange(0, len(production_year)):
		if production_year[i] == 'Not Available' and production_week[i] == 'Not Available' and production_day[i] == 'Not Available':
			time = ''
		else:
			time = '%s/%s/%s' % (production_year[i], production_week[i], production_day[i])  
		production_date.append(time)

	for i in range(0, len(sn)):
		disks.append({'bus':bus[i], 'capacity':capacity[i], 'media':media[i], 'product_id':product_id[i], 'production_date':production_date[i], 'slot':slot[i], 'sn':sn[i], 'speed':speed[i], 'vender':vender[i], 'status':status[i]})	

	s['disks'] = disks	

locator = []
size = []
sn = []
speed = []

p = subprocess.Popen('dmidecode --type Memory', shell=True, stdout=subprocess.PIPE)
for line in p.stdout.readlines():
	data = line.strip().split(":",1)
	if len(data) == 2:
		d1 = data[0].strip()
		d2 = data[1].strip()
		if d1 == 'Locator':
			locator.append(d2)
		elif d1 == 'Size':
			size.append(d2)
		elif d1 == 'Serial Number':
			sn.append(d2)
		elif d1 == 'Speed':
			speed.append(d2)
for i in xrange(0, len(sn)):
	if size[i] and size[i] != 'No Module Installed':
		memorys.append({'locator':locator[i],'size':size[i],'sn':sn[i],'speed':speed[i]})

s['memorys'] = memorys

sn = []
vender = []
name = []
model = []

p = subprocess.Popen('dmidecode --type processor', shell=True, stdout=subprocess.PIPE)
for line in p.stdout.readlines():
	data = line.strip().split(":",1)
	if len(data) == 2:
		d1 = data[0].strip()
		d2 = data[1].strip()
		if d1 == 'ID':
			sn.append(d2)
		elif d1 == 'Manufacturer':
			vender.append(d2)
		elif d1 == 'Socket Designation':
			name.append(d2)
		elif d1 == 'Version':
			model.append(d2)

for i in xrange(0, len(sn)):
	cpus.append({'sn':sn[i],'vender':vender[i],'name':name[i],'model':model[i]})

s['cpus'] = cpus

soc = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
soc.connect(('10.10.0.1', 0))
ip = soc.getsockname()[0]
s['ip'] = ip

print json.dumps(s, indent=4)
