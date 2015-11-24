#!/usr/bin/env python

import sys
import socket
import time
import re
import copy

from StringIO import StringIO
import json

ZK_METRICS = {
    'time' : 0,
    'data' : {}
}

ZK_LAST_METRICS = copy.deepcopy(ZK_METRICS)

class ZooKeeperServer(object):
    def __init__(self, host='127.0.0.1', port='2181', timeout=1):
        self._address = (host, int(port))
        self._timeout = timeout

    def get_stats(self):
        global ZK_METRICS, ZK_LAST_METRICS
        ZK_METRICS = {
          'time' : time.time(),
          'data' : {}
        }
        data = self._send_cmd('mntr')
        if data:
            parsed_data =  self._parse(data)
        else:
            data = self._send_cmd('stat')
            parsed_data = self._parse_stat(data)

        ZK_METRICS['data'] = parsed_data
        ZK_LAST_METRICS = copy.deepcopy(ZK_METRICS)

	del parsed_data["zk_version"]
	del parsed_data["zk_server_state"]
        return parsed_data

    def _create_socket(self):
        return socket.socket()

    def _send_cmd(self, cmd):
        s = self._create_socket()
        s.settimeout(self._timeout)

        s.connect(self._address)
        s.send(cmd)

        data = ""
        newdata = s.recv(2048)
        while newdata:
            data += newdata
            newdata = s.recv(2048)

        s.close()

        return data

    def _parse(self, data):
        h = StringIO(data)

        result = {}
        for line in h.readlines():
            try:
                key, value = self._parse_line(line)
                result[key] = value
            except ValueError:
                pass

        return result

    def _parse_stat(self, data):
        global ZK_METRICS, ZK_LAST_METRICS

        h = StringIO(data)

        result = {}

        while h.readline().strip(): pass

        for line in h.readlines():
            m = re.match('Latency min/avg/max: (\d+)/(\d+)/(\d+)', line)
            if m is not None:
                result['zk_min_latency'] = int(m.group(1))
                result['zk_avg_latency'] = int(m.group(2))
                result['zk_max_latency'] = int(m.group(3))
                continue

            m = re.match('Received: (\d+)', line)
            if m is not None:
                cur_packets = int(m.group(1))
                packet_delta = cur_packets - ZK_LAST_METRICS['data'].get('zk_packets_received_total', cur_packets)
                time_delta = ZK_METRICS['time'] - ZK_LAST_METRICS['time']
                try:
                    result['zk_packets_received_total'] = cur_packets
                    result['zk_packets_received'] = packet_delta / float(time_delta)
                except ZeroDivisionError:
                    result['zk_packets_received'] = 0
                continue

            m = re.match('Sent: (\d+)', line)
            if m is not None:
                cur_packets = int(m.group(1))
                packet_delta = cur_packets - ZK_LAST_METRICS['data'].get('zk_packets_sent_total', cur_packets)
                time_delta = ZK_METRICS['time'] - ZK_LAST_METRICS['time']
                try:
                    result['zk_packets_sent_total'] = cur_packets
                    result['zk_packets_sent'] = packet_delta / float(time_delta)
                except ZeroDivisionError:
                    result['zk_packets_sent'] = 0
                continue

            m = re.match('Outstanding: (\d+)', line)
            if m is not None:
                result['zk_outstanding_requests'] = int(m.group(1))
                continue

#            m = re.match('Mode: (.*)', line)
#            if m is not None:
#                result['zk_server_state'] = m.group(1)
#                continue
#
            m = re.match('Node count: (\d+)', line)
            if m is not None:
                result['zk_znode_count'] = int(m.group(1))
                continue

        return result

    def _parse_line(self, line):
        try:
            key, value = map(str.strip, line.split('\t'))
        except ValueError:
            raise ValueError('Found invalid line: %s' % line)

        if not key:
            raise ValueError('The key is mandatory and should not be empty')

        try:
            value = int(value)
        except (TypeError, ValueError):
            pass

        return key, value

if __name__ == "__main__":
	zk = ZooKeeperServer()
	print json.dumps(zk.get_stats(), indent=4)
