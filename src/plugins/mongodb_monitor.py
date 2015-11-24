#!/usr/bin/env python
import subprocess

try:
	from pymongo import MongoClient
except:
	p = subprocess.call("pip install pymongo", shell=True)

import json
import sys

class MongoDB(object):

    def __init__(self):
        self.mongo_host = "127.0.0.1"
	if len(sys.argv) == 2:
        	self.mongo_port = int(sys.argv[1])
	else:
        	self.mongo_port = 27017
		
        self.mongo_db = ["admin", ]
        self.mongo_user = None
        self.mongo_password = None

    def do_server_status(self):
        conn = MongoClient(host=self.mongo_host, port=self.mongo_port)
        db = conn[self.mongo_db[0]]
	performance_index = {}
        if self.mongo_user and self.mongo_password:
            db.authenticate(self.mongo_user, self.mongo_password)

        server_status = db.command('serverStatus')

        # operations
        for k, v in server_status['opcounters'].items():
	    performance_index[k] = v

        # memory
        for t in ['resident', 'virtual']:
            performance_index[t] = server_status['mem'][t]

        # connections
	if 'current' in server_status['connections']:
            performance_index['current'] = server_status['connections']['available']

	if 'available' in server_status['connections']:
            performance_index['available'] = server_status['connections']['available']

	if 'totalCreated' in server_status['connections']:
            performance_index['totalCreated'] = server_status['connections']['totalCreated']

	# network
	if 'network' in server_status:
	    for t in ['bytesIn', 'bytesOut', 'numRequests']:
                performance_index[t] = server_status['network'][t]

	if 'ismaster' in server_status['repl']:
		if server_status['repl']['ismaster']:
			performance_index['ismaster'] = 0
		else:
			performance_index['ismaster'] = 1
		
	if 'secondary' in server_status['repl']:
		if server_status['repl']['secondary']:
			performance_index['secondary'] = 0
		else:
			performance_index['secondary'] = 1
		
	return json.dumps(performance_index, indent=4)

if __name__ == "__main__":
	mongodb = MongoDB()
	print mongodb.do_server_status()
