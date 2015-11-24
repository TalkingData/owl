#!/usr/bin/env python
# coding:utf8
import subprocess
try:
	import simplejson as json
except:
	import json

dfi1 = subprocess.Popen(["df -ih | grep -v tmpfs"], stdout=subprocess.PIPE, shell=True)
dfi2 = subprocess.Popen(["sed '1d;/ /!N;s/\\n//;s/ \+/ /'"], stdin=dfi1.stdout, stdout=subprocess.PIPE, shell=True)
dfi3 = subprocess.Popen(["awk '{print $NF,$(NF-1)}'"], stdin=dfi2.stdout, stdout=subprocess.PIPE, shell=True)
dfi4 = subprocess.Popen(["grep -v '/dev/shm'"], stdin=dfi3.stdout, stdout=subprocess.PIPE, shell=True)
outputi = dfi4.communicate()[0].strip()
disklisti = outputi.split("\n")

fillratesi = {}

for lines in disklisti:	
	part = lines.split()
	usage_int = int(part[1].replace("%", ""))
	if usage_int > 0:
		if part[0] == '/':
			fillratesi['root'] = usage_int
		else:
			part = "%s" % (part[0].strip('/').replace('/','.'))   
			fillratesi[part] = usage_int

print json.dumps(fillratesi, indent=4)
	
