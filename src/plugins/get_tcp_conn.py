#!/usr/bin/env python
# coding:utf8
import os
try:
    import simplejson as json
except:
    import json

tcp = "/proc/net/tcp"
tcp6 = "/proc/net/tcp6"

status_code = {	
    '00' : 'ERROR_STATUS',
    '01' : 'TCP_ESTABLISHED',
    '02' : 'TCP_SYN_SENT',
    '03' : 'TCP_SYN_RECV',
    '04' : 'TCP_FIN_WAIT1',
    '05' : 'TCP_FIN_WAIT2',
    '06' : 'TCP_TIME_WAIT',
    '07' : 'TCP_CLOSE',
    '08' : 'TCP_CLOSE_WAIT',
    '09' : 'TCP_LAST_ACK',
    '0A' : 'TCP_LISTEN',
    '0B' : 'TCP_CLOSING'
}

tcp_count = {}

for code in status_code.values():
   tcp_count[code] = 0

if os.path.exists(tcp):
    with open(tcp, 'rb') as f:
        lines = f.readlines()
	for line in lines:
	    st = line.split()[3]
	    for code in status_code.keys(): 
		if st == code:	
		    tcp_count[status_code[code]] += 1

if os.path.exists(tcp6):
    with open(tcp6, 'rb') as f:
        lines = f.readlines()
	for line in lines:
	    st = line.split()[3]
	    for code in status_code.keys(): 
		if st == code:	
		    tcp_count[status_code[code]] += 1

print json.dumps(tcp_count, indent=4)
