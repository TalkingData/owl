#!/usr/bin/env python

import subprocess
try:
    from vertica_python import connect
except:
    p = subprocess.call("easy_install vertica_python", shell=True)
import json
import sys

class Vertica(object):

    def __init__(self):
	self.host = "127.0.0.1" 
	if len(sys.argv) == 2:
        	self.port = sys.argv[1]
	else:
        	self.port = 5433
	
        self.user = "user"
        self.password = "password"
	self.database = ""

	self.conn = connect(host=self.host, port=int(self.port), user=self.user, password=self.password, database=self.database)
	self.cur = self.conn.cursor()

    def getValue(self, query):
	self.cur.execute(query)
	result = self.cur.fetchone()
	return result

    def getAll(self, query):
	self.cur.execute(query)
	result = self.cur.fetchall()
	return result
	
    def closeConn(self):
        self.conn.close()

if __name__ == "__main__":
	vertica = Vertica()
	performance_index = {}
	performance_index["session_percent"] = vertica.getValue("select a.sess, round((a.sess / b.current_value * 100), 2) session_percent from (select count(*) sess from v_monitor.sessions) a,(select current_value From v_monitor.CONFIGURATION_PARAMETERS where parameter_name = 'MaxClientSessions') b;")[1]
	performance_index["usage_percent"] = float(vertica.getValue("select GET_COMPLIANCE_STATUS();")[0].split('\n')[2].split(':')[1].strip().replace('%', ''))
	
	status = vertica.getAll("select node_address,node_state from nodes;")
	performance_index["node_state"] = 0
	for s in status:
		if s[1] != "UP":	
			performance_index["node_state"] = 1
		
	vertica.closeConn()
	print json.dumps(performance_index, indent=4)
