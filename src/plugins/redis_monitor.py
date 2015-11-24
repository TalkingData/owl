#!/usr/bin/env python

from __future__ import division
import redis
import sys
import json
import re

def main():
	index = {}
	keys = ('up')
	symbols = ('K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y')
	if len(sys.argv) == 2:
		host = "127.0.0.1"
		port = sys.argv[1]
		r = redis.Redis(host=host, port=int(port))
		index['connected_clients'] = r.info()['connected_clients']
		mem = r.info()['used_memory_human']
		m = re.match(r'(.*)([A-Z])', mem)
		if m:
			val = m.group(1)
			unit = m.group(2)
			if unit == 'K':
				index['used_memory_human'] = float(val) * 1024
			elif unit == 'M':
				index['used_memory_human'] = float(val) * 1024 ** 2
			elif unit == 'G':
				index['used_memory_human'] = float(val) * 1024 ** 3
		
		index['mem_fragmentation_ratio'] = r.info()['mem_fragmentation_ratio']
		index['rdb_bgsave_in_progress'] = r.info()['rdb_bgsave_in_progress']
		index['instantaneous_ops_per_sec'] = r.info()['instantaneous_ops_per_sec']
		index['total_commands_processed'] = r.info()['total_commands_processed']
		index['connected_clients'] = r.info()['connected_clients']
		if r.info().has_key('master_link_status'):
			if r.info()['master_link_status'] == "up":
				index['master_link_status'] = 0
			else:
				index['master_link_status'] = 1

		if r.info().has_key('slave0'):
			datas = r.info()['slave0']	
			if datas['state'] == "online":
				index['slave0'] = 0
			else:
				index['slave0'] = 1

	print json.dumps(index, indent=4)

if __name__ == "__main__":
	main()
