#!/usr/bin/env python

import subprocess

try:
	import simplejson as json
except:
	import json


disk_status = {} 

p = subprocess.Popen("dmidecode --type system | perl -alne '/Manufacturer:\s(.*)/ and print $1'", shell=True, stdout=subprocess.PIPE)
vender = p.stdout.readline().strip()

if vender == 'HP':
	slot = []
	port = []
	box = []
	bay = []
	status = []
	interface_type = []

	p = subprocess.Popen("hpssacli ctrl all show | perl -alne '/Slot.*?(\d+)/ and print $1'", shell=True, stdout=subprocess.PIPE)
	ctrl_slot = int(p.stdout.readline().strip())

	p = subprocess.Popen('hpssacli ctrl slot=%d pd all show detail' % (ctrl_slot), shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		data = line.strip().split(":",1)
		if len(data) == 2:
			d1 = data[0].strip()
			d2 = data[1].strip()
			if d1 == 'Port':
				port.append(d2)
			elif d1 == 'Box':
				box.append(d2)
			elif d1 == 'Bay':
				bay.append(d2)	
			elif d1 == 'Status':
				status.append(d2)
			elif d1 == 'Interface Type':
				interface_type.append(d2)

	for i in xrange(0, len(port)):
		number = '%s_%s_%s' % (port[i], box[i], bay[i])  
		slot.append(number)

	for i in range(0, len(slot)):
		if status[i].lower() == 'ok' and interface_type[i].lower() == 'sas':
			disk_status[slot[i]] = 0
		elif interface_type[i].lower() == 'solid state sata':
			disk_status[slot[i]] = 0
		else:
			disk_status[slot[i]] = 1

else:
	slot = []
	status = []
	media = []

	p = subprocess.Popen("/opt/dell/srvadmin/bin/omreport storage controller | perl -alne '/^ID.*?(\d+)/ and print $1'", shell=True, stdout=subprocess.PIPE)
	ctrl_slot = int(p.stdout.readline().strip())

	p = subprocess.Popen('/opt/dell/srvadmin/bin/omreport storage pdisk controller=%d' % (ctrl_slot), shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		data = line.strip().split(":",1)
		if len(data) == 2:
			d1 = data[0].strip()
			d2 = data[1].strip()
			if d1 == 'ID':
				slot.append(d2.replace(':', '_'))
			elif d1 == 'Status':
				status.append(d2)
			elif d1 == 'Media':
				media.append(d2)

	for i in range(0, len(slot)):
		if status[i].lower() == 'ok' and media[i].lower() == 'hdd':
			disk_status[slot[i]] = 0
		elif media[i].lower() == 'ssd':
			disk_status[slot[i]] = 0
		else:
			disk_status[slot[i]] = 1

print json.dumps(disk_status, indent=4)
