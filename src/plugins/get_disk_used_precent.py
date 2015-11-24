#!/usr/bin/python
# conding: utf8
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
def diskusage():
    disk = {}
    p = subprocess.Popen("df -lh", shell=True, stdout=subprocess.PIPE)
    p1 = subprocess.Popen("sed '1d;/ /!N;s/\\n//;s/ \+/ /'", shell=True, stdin=p.stdout, stdout=subprocess.PIPE)
    p2 = subprocess.Popen("awk '{print $NF,$(NF-1)}'", shell=True, stdin=p1.stdout, stdout=subprocess.PIPE)
    p3 = subprocess.Popen("grep -v '/dev/shm'", shell=True, stdin=p2.stdout, stdout=subprocess.PIPE)
    _usage = {}

    for line in p3.stdout.readlines():
        l = line.split()
        if l[0] == '/':
            _usage['root'] = int(l[1].replace('%',''))
        else:
           key = "%s" % (l[0].strip('/').replace('/','.'))
           _usage[key] = int(l[1].replace('%',''))

    for part in psutil.disk_partitions(all=False):
        diskmetric = ['percent']
        usage = psutil.disk_usage(part.mountpoint)

        for metric in diskmetric:
                if part.mountpoint == '/':
                        key = "root"
                        if metric == 'total':
                                disk[key] = usage.total
                        elif metric == 'used':
                                disk[key] = usage.used
                        elif metric == 'free':
                                disk[key] = usage.free
                        elif metric == 'percent':
                                disk[key] = _usage[key]
                else:
                        mount = part.mountpoint.strip('/').replace('/','.')
                        key = mount
                        if metric == 'total':
                                disk[key] = usage.total
                        elif metric == 'used':
                                disk[key] = usage.used
                        elif metric == 'free':
                                disk[key] = usage.free
                        elif metric == 'percent':
                                if key in _usage.keys():
                                    disk[key] = _usage[key]
    return json.dumps(disk, indent=4)

if __name__ == "__main__":
	print diskusage()
