#!/usr/bin/env python

import subprocess

try:
	import simplejson as json
except:
	import json

memory_status = {} 

p = subprocess.Popen("dmidecode --type system | perl -alne '/Manufacturer:\s(.*)/ and print $1'", shell=True, stdout=subprocess.PIPE)
vender = p.stdout.readline().strip()

if vender == 'HP':
	mem_slot = []
	mem_status = []
	p = subprocess.call("rpm -qa | grep hp-health", shell=True, stdout=subprocess.PIPE)
	if p != 0:
		subprocess.call("rpm -ivh http://10.10.32.35/game-configuration/hp-health-9.40-1602.44.rhel6.x86_64.rpm", shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

	p = subprocess.Popen("hpasmcli -s 'show dimm'", shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		a = {}
		data = line.strip().split(":")
		if len(data) == 2:
			a[data[0].strip()] = data[1].strip()
		
		if 'Status' in a.keys():
			status = a['Status']
                        if status == "Ok":
                                status = 0
                                mem_status.append(status)

	length = len(mem_status)
	if length:
		for i in range(0, length):
			slot = "slot_%d"  % (i)
			mem_slot.append(slot)

	memory_status = dict(zip(mem_slot, mem_status))
else:
	mem_slot = []
	mem_status = []
	p = subprocess.Popen("omreport chassis Memory", shell=True, stdout=subprocess.PIPE)
	for line in p.stdout.readlines():
		a = {}
		data = line.strip().split(":")
		if len(data) == 2:
			a[data[0].strip()] = data[1].strip()

		if 'Index' in a.keys():
			slot = "slot_" + a['Index']
			if slot:
				mem_slot.append(slot)		
		elif 'Status' in a.keys():
			status = a['Status']
			if status != "Unknown" and status == "Ok":
				status = 0
				mem_status.append(status)

	memory_status = dict(zip(mem_slot, mem_status))
			
print json.dumps(memory_status, indent=4)
