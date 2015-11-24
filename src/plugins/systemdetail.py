#!/usr/bin/python
# conding: utf8

from __future__ import division
import sys
import os
import time
import threading
from datetime import datetime, timedelta
import subprocess

try:
    import simplejson as json
except:
    import json

import pickle

try:
    import psutil
except:
    sys.exit()

def bytes2human(n):
    symbols = ('K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y')
    prefix = {}
    for i, s in enumerate(symbols):
        prefix[s] = 1 << (i + 1) * 10

    for s in reversed(symbols):
        if n >= prefix[s]:
            value = float(n) / prefix[s]
            return '%.2f%s' % (value, s)
    return '%.2fB' % (n)

def cpu():
    cpustates = {'cpu':'', 'logical':'', 'percent':'', 'user':'', 'system':'', 'idle':'', 'nice':'', 'irq':'', 'softirq':'', 'iowait':'', 'steal':'', 'guest':''}
    cpu_count = psutil.cpu_count(logical=False)
    cpu_count_logical = psutil.cpu_count()
    psutil.cpu_percent()
    psutil.cpu_times_percent()
    time.sleep(1)
    cpu_percent = psutil.cpu_percent()
    cpu_times_percent = psutil.cpu_times_percent()
    for metric in cpustates:
        if metric == 'cpu':
            cpustates[metric] = cpu_count
        elif metric == 'logical':
            cpustates[metric] = cpu_count_logical
        elif metric == 'percent':
            cpustates[metric] = cpu_percent
        elif metric == 'user':
            cpustates[metric] = cpu_times_percent.user
        elif metric == 'system':
            cpustates[metric] = cpu_times_percent.system
        elif metric == 'idle':
            cpustates[metric] = cpu_times_percent.idle
        elif metric == 'nice':
            cpustates[metric] = cpu_times_percent.nice
        elif metric == 'irq':
            cpustates[metric] = cpu_times_percent.irq
        elif metric == 'softirq':
            cpustates[metric] = cpu_times_percent.softirq
        elif metric == 'iowait':
            cpustates[metric] = cpu_times_percent.iowait
        elif metric == 'steal':
            cpustates[metric] = cpu_times_percent.steal
        elif metric == 'guest':
            cpustates[metric] = cpu_times_percent.guest

    return json.dumps(cpustates, indent=4)

def load():
    loadstats = {'min1': '', 'min5': '', 'min15': ''}
    load = os.getloadavg()
    loadstats['min1'] = load[0]
    loadstats['min5'] = load[1]
    loadstats['min15'] = load[2]
    return json.dumps(loadstats, indent=4)

def runtime():
    runtime = {'uptime': ''}
    uptime = str(datetime.now() - datetime.fromtimestamp(psutil.boot_time())).split('.')[0]
    runtime['uptime'] = uptime
    return json.dumps(runtime, indent=4)

def disk():
    diskio = psutil.disk_io_counters()
    diskstates = {'iops':'', 'ips':'', 'ops':'', 'bps':'', 'bps_in':'', 'bps_out':'', 'read_time':'', 'write_time':''}
    for metric in diskstates:
	    if metric == 'iops':
		diskstates[metric] = int(diskio.write_count + diskio.read_count)
	    elif metric == 'ips':
		diskstates[metric] = int(diskio.read_count)
	    elif metric == 'ops':
		diskstates[metric] = int(diskio.write_count)
	    elif metric == 'bps':
		diskstates[metric] = int(diskio.write_bytes + diskio.read_bytes)
	    elif metric == 'bps_in':
		diskstates[metric] = int(diskio.read_bytes)
	    elif metric == 'bps_out':
		diskstates[metric] = int(diskio.write_bytes)
	    elif metric == 'read_time':
		diskstates[metric] = int(diskio.read_time)
	    elif metric == 'write_time':
		diskstates[metric] = int(diskio.write_time)

    return json.dumps(diskstates, indent=4)

def net(dev):
	networkmetric = {}
	networkio = psutil.net_io_counters(pernic=True)
    	second_now = datetime.now().strftime("%s")
	metrics = ['bytes_recv', 'bytes_sent', 'packets_recv', 'packets_sent', 'err_in', 'err_out', 'drop_in', 'drop_out']
	netcounters = networkio[dev]
	for metric in metrics:
		key = "%s" % (metric)
		if metric == 'bytes_recv':		
			networkmetric[key] = int(netcounters.bytes_recv)
		elif metric == 'bytes_sent':
			networkmetric[key] = int(netcounters.bytes_sent)
		elif metric == 'packets_recv':
			networkmetric[key] = int(netcounters.packets_recv)
		elif metric == 'packets_sent':
			networkmetric[key] = int(netcounters.packets_sent)
		elif metric == 'err_in':
			networkmetric[key] = int(netcounters.errin)
		elif metric == 'err_out':
			networkmetric[key] = int(netcounters.errout)
		elif metric == 'drop_in':
			networkmetric[key] = int(netcounters.dropin)
		elif metric == 'drop_out':
			networkmetric[key] = int(netcounters.dropout)

	return json.dumps(networkmetric, indent=4)
		
def mem():
    memoryio = psutil.virtual_memory()
    memstates = {'total':'', 'available':'', 'percent':'', 'used':'', 'free':'', 'active':'', 'inactive':'', 'buffers':'', 'cached':''}
    for metric in memstates:
        if metric == 'total':
            memstates[metric] = memoryio.total
        elif metric == 'available':
            memstates[metric] = memoryio.available
        elif metric == 'percent':
            memstates[metric] = memoryio.percent
        elif metric == 'used':
            memstates[metric] = memoryio.used
        elif metric == 'free':
            memstates[metric] = memoryio.free
        elif metric == 'active':
            memstates[metric] = memoryio.active
        elif metric == 'inactive':
            memstates[metric] = memoryio.inactive
        elif metric == 'buffers':
            memstates[metric] = memoryio.buffers
        elif metric == 'cached':
            memstates[metric] = memoryio.cached
    return json.dumps(memstates, indent=4)

def memswap():
    swapmemory = psutil.swap_memory()
    swapmemorystates = {'total':'', 'used':'', 'free':'', 'percent':'', 'sin':'', 'sout':''}
    for metric in swapmemorystates:
        if metric == 'total':
            swapmemorystates[metric] = swapmemory.total
        elif metric == 'used':
            swapmemorystates[metric] = swapmemory.used
        elif metric == 'free':
            swapmemorystates[metric] = swapmemory.free
        elif metric == 'percent':
            swapmemorystates[metric] = swapmemory.percent
        elif metric == 'sin':
            swapmemorystates[metric] = swapmemory.sin
        elif metric == 'sout':
            swapmemorystates[metric] = swapmemory.sout
    return json.dumps(swapmemorystates, indent=4)

def user():
    users = psutil.users()
    userstates = []
    for user in users:
        usermetric = {'name':'', 'terminal':'', 'host':'', 'started':''}
        for metric in usermetric:
            if metric == 'name':
                usermetric[metric] = user.name
            elif metric == 'terminal':
                usermetric[metric] = user.terminal
            elif metric == 'host':
                usermetric[metric] = user.host
            elif metric == 'started':
                usermetric[metric] = str(datetime.fromtimestamp(user.started))
        userstates.append(usermetric)
    return json.dumps(userstates, indent=4)

procstates = []

def main():
	if len(sys.argv) >= 2:
	    argvs = sys.argv
	    if argvs[1] == 'disk':
		print disk()
	    elif argvs[1] == 'cpu':
		print cpu()
	    elif argvs[1] == 'load':
		print load()
	    elif argvs[1] == 'runtime':
		print runtime()
	    elif argvs[1] == 'net':
		try:
			dev = argvs[2]	
			print net(dev)
		except:
			pass
	    elif argvs[1] == 'mem':
		print mem()
	    elif argvs[1] == 'memswap':
		print memswap()

if __name__ == '__main__':
	main()
