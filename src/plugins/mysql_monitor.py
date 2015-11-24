#!/usr/bin/python
# coding:utf8

from __future__ import division
import MySQLdb
import string
import os
import sys
import optparse
from datetime import datetime, timedelta
try:
    import simplejson as json
except:
    import json

class Mysql:
    def __init__(self, host, user, passwd, port):
        self.host = host
        self.user = user
        self.passwd = passwd
        self.port = port
        try:
            self.conn = MySQLdb.connect(host=host, user=user, passwd=passwd, port=port)
            self.cur = self.conn.cursor()
        except MySQLdb.Error, e:
            print e.args[0], e.args[1]
            sys.exit(1)

    def getValue(self, query):
        self.cur.execute(query)
        result = self.cur.fetchone()
        return int(result[1])

    def getQuery(self, query):
        self.cur.execute(query)
        result = self.cur.fetchall()
        return result

    def executeQuery(self, query):
        self.cur.execute(query)
        self.conn.commit()
        print "%s execute success" % (query)

    def closeConn(self):
        self.cur.close()
        self.conn.close()

def quantity(name, count):
    count = str(count)
    _f = '/tmp/%s.txt' % (name)
    c = 0 
    if os.path.isfile(_f):
        with open(_f, 'rb') as f:
            _count = f.readline()
            c = int(count) - int(_count)

        with open(_f, 'wb') as f:
            f.write(count)
    else:
        with open(_f, 'wb') as f:
            f.write(count)

    return c

def main():
    performance = {} 

    Questions="show global status like 'Questions'"
    Uptime = "show global status like 'Uptime'"
    Com_commit = "show global status like 'Com_commit'"
    Com_rollback = "show global status like 'Com_rollback'"
    Key_reads = "show global status like 'Key_reads'"
    Key_read_requests = "show global status like 'Key_read_requests'"
    Key_writes = "show global status like 'Key_writes'" 
    Key_write_requests = "show global status like 'Key_write_requests'"

    Have_innodb = "show global variables like 'have_innodb'"
    Innodb_buffer_pool_reads = "show global status like 'Innodb_buffer_pool_reads'"
    Innodb_buffer_pool_read_requests = "show global status like 'Innodb_buffer_pool_read_requests'"

    Qcache_hits = "show global status like 'Qcache_hits'"
    Qcache_inserts = "show global status like 'Qcache_inserts'"
    Open_tables = "show global status like 'Open_tables'"
    Opened_tables = "show global status like 'Opened_tables'"
    Threads_created = "show global status like 'Threads_created'"
    Threads_running = "show global status like 'Threads_running'"
    Threads_connected = "show global status like 'Threads_connected'"
    Aborted_connects = "show global status like 'Aborted_connects'"
    Connections = "show global status like 'Connections'"
    Com_begin = "show global status like 'Com_begin'"
    Com_select = "show global status like 'Com_select'"
    Com_insert = "show global status like 'Com_insert'"
    Com_update = "show global status like 'Com_update'"
    Com_delete = "show global status like 'Com_delete'"
    Com_replace = "show global status like 'Com_replace'"
    Table_locks_waited = "show global status like 'Table_locks_waited'"
    Table_locks_immediate = "show global status like 'Table_locks_immediate'"
    Created_tmp_tables = "show global status like 'Created_tmp_tables'"
    Created_tmp_disk_tables = "show global status like 'Created_tmp_disk_tables'"
    Slow_queries = "show global status like 'Slow_queries'"
    Select_full_join = "show global status like 'Select_full_join'"

    Bytes_received = "show global status like 'Bytes_received'"
    Bytes_sent = "show global status like 'Bytes_sent'"

    slave_status = "show slave status"

    p = optparse.OptionParser()
    p.add_option('--host', '-H', default='127.0.0.1')
    p.add_option('--user', '-u', default='root')
    p.add_option('--passwd', '-p', default='')
    p.add_option('--port', '-P', default=3306)

    options, arguments = p.parse_args()
    port = int(options.port)
    conn = Mysql(options.host, options.user, options.passwd, port)

    Uptime = conn.getValue(Uptime)
    s = quantity("Uptime_%s" % (port), Uptime)
    
    slave_status = conn.getQuery(slave_status)
    if slave_status:
	if slave_status[0][10] == "Yes":
		performance['Slave_IO_Running'] = 0
	else:
		performance['Slave_IO_Running'] = 1

	if slave_status[0][11] == "Yes":
		performance['Slave_SQL_Running'] = 0
	else:
		performance['Slave_SQL_Running'] = 1

	if slave_status[0][32]:
		performance['Seconds_Behind_Master']= slave_status[0][32]
	else:
		performance['Seconds_Behind_Master']= 0

    Questions = conn.getValue(Questions)
    _Questions = quantity("Questions_%s" % (port), Questions)

    Com_commit = conn.getValue(Com_commit)
    _Com_commit = quantity("Com_commit_%s" % (port), Com_commit)
    Com_rollback = conn.getValue(Com_rollback)
    _Com_rollback = quantity("Com_rollback_%s" % (port), Com_rollback)

    Key_reads = conn.getValue(Key_reads)
    _Key_reads = quantity("Key_reads_%s" % (port), Key_reads)
    Key_read_requests = conn.getValue(Key_read_requests)
    _Key_read_requests = quantity("Key_read_requests_%s" % (port), Key_read_requests)
    Key_writes = conn.getValue(Key_writes)
    _Key_writes = quantity("Key_writes_%s" % (port), Key_writes)
    Key_write_requests = conn.getValue(Key_write_requests)
    _Key_write_requests = quantity("Key_write_requests_%s" % (port), Key_write_requests)

    Qcache_hits = conn.getValue(Qcache_hits)
    _Qcache_hits = quantity("Qcache_hits_%s" % (port), Qcache_hits)
    Qcache_inserts = conn.getValue(Qcache_inserts)
    _Qcache_inserts = quantity("Qcache_inserts_%s" % (port), Qcache_inserts)

    Open_tables = conn.getValue(Open_tables)

    Opened_tables = conn.getValue(Opened_tables)
    _Opened_tables = quantity("Opened_tables_%s" % (port), Opened_tables)

    Threads_created = conn.getValue(Threads_created)
    Connections = conn.getValue(Connections)

    Threads_running = conn.getValue(Threads_running)
    Threads_connected = conn.getValue(Threads_connected) 
    Aborted_connects = conn.getValue(Aborted_connects)
    _Aborted_connects = quantity("Aborted_connects_%s" % (port), Aborted_connects)
    
    Com_begin = conn.getValue(Com_begin)
    _Com_begin = quantity("Com_begin_%s" % (port), Com_begin)
    Com_select = conn.getValue(Com_select)
    _Com_select = quantity("Com_select_%s" % (port), Com_select)
    Com_insert = conn.getValue(Com_insert)
    _Com_insert = quantity("Com_insert_%s" % (port), Com_insert)
    Com_update = conn.getValue(Com_update)
    _Com_update = quantity("Com_update_%s" % (port), Com_update)
    Com_delete = conn.getValue(Com_delete)
    _Com_delete = quantity("Com_delete_%s" % (port), Com_delete)
    Com_replace = conn.getValue(Com_replace) 
    _Com_replace = quantity("Com_replace_%s" % (port), Com_replace)

    Table_locks_immediate = conn.getValue(Table_locks_immediate)  
    Table_locks_waited = conn.getValue(Table_locks_waited)

    Created_tmp_tables = conn.getValue(Created_tmp_tables)
    _Created_tmp_tables = quantity("Created_tmp_tables_%s" % (port), Created_tmp_tables)
    Created_tmp_disk_tables = conn.getValue(Created_tmp_disk_tables)
    _Created_tmp_disk_tables = quantity("Created_tmp_disk_tables_%s" % (port), Created_tmp_disk_tables)

    Slow_queries = conn.getValue(Slow_queries)
    _Slow_queries = quantity("Slow_queries_%s" % (port), Slow_queries)
    Select_full_join = conn.getValue(Select_full_join)

    Bytes_received = conn.getValue(Bytes_received)
    _Bytes_received = quantity("Bytes_received_%s" % (port), Bytes_received)
    Bytes_sent = conn.getValue(Bytes_sent)
    _Bytes_send = quantity("Bytes_sent_%s" % (port), Bytes_sent)

    if s:
        performance['qps'] = int(_Questions / s)
        performance['tps'] = int((_Com_commit + _Com_rollback) / s)

        performance['begin'] = int(_Com_begin / s)
        performance['select'] = int(_Com_select / s)
        performance['insert'] = int(_Com_insert / s)
        performance['update'] = int(_Com_update / s)
        performance['delete'] = int(_Com_delete / s)

        performance['commit'] = int(_Com_commit / s)
        performance['rollback'] = int(_Com_rollback / s)
        performance['replace'] = int(_Com_replace / s)

        performance['bytes_received'] = int(_Bytes_received / s)
        performance['bytes_send'] = int(_Bytes_send / s)

        performance['threads_running'] = int(Threads_running)
        performance['threads_connected'] = int(Threads_connected)
        performance['aborted_connects'] = int(_Aborted_connects)

        performance['open_tables'] = int(Open_tables)

    conn.closeConn()
    print json.dumps(performance, indent=4)

if __name__ == '__main__':
    main()
