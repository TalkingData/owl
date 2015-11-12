# coding:utf-8
import pycurl
import operator
from notifi import alarm
from views import item_update 
from compare import tsdb_data, tsdb_ratio_data, tsdb_network_data, ratio, bytes2human
from celery import task
from host.models import *
from network.models import *
from round_robin import Round_Robin
from celery import task
from datetime import *
from urllib import quote
from django.shortcuts import get_object_or_404
try:
    import signal
    signal.signal(signal.SIGPIPE, signal.SIG_IGN)
except ImportError:
    pass

tsdb_host = ("127.0.0.1:4242", "127.0.0.1:4243", "127.0.0.1:4244")
num_conn = 12

def Unit(cycle, number):
	start_time = (datetime.now()-timedelta(days=number)-timedelta(minutes=cycle)).strftime("%Y/%m/%d %H:%M:%S")
	stop_time = (datetime.now()-timedelta(days=number)).strftime("%Y/%m/%d %H:%M:%S")
   	return start_time, stop_time

def pycurl_data(filename, url): 
	fp = open(filename, "wb")
	curl = pycurl.Curl()
	curl.setopt(pycurl.URL, url)
	curl.setopt(pycurl.FOLLOWLOCATION, 1)
	curl.setopt(pycurl.MAXREDIRS, 5)
	curl.setopt(pycurl.CONNECTTIMEOUT, 5)
	curl.setopt(pycurl.TIMEOUT, 10)
	curl.setopt(pycurl.NOSIGNAL, 1)
	curl.setopt(pycurl.WRITEDATA, fp)
	curl.perform()
	curl.close()
	fp.close()
	return 0

def tsdb_ratio(h, url, url_ago, name, key, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups, alert):
	filename = "/tmp/%s_ratio_url_%s.%s" % (h.uuid, name, key)
	filename_ago = "/tmp/%s_ratio_url_ago_%s.%s" % (h.uuid, name, key)
	data = pycurl_data(filename, quote(url, ':/=&()?,>.'))
	data_ago = pycurl_data(filename_ago, quote(url_ago, ':/=&()?,>.'))
    	metric = "%s.%s" % (name, key)
	if data == 0 and data_ago == 0:
		dt = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
		symbols = {
			'>' : 'gt',
			'>=' : 'ge',
			'<' : 'lt',
			'<=' : 'le',
			'=' : 'eq',
			'!=' : 'ne',
			'<>' : 'ne'
		}
		comparator = operator.__dict__[symbols[symbol]]
		with open(filename, "rb") as f:
			bfb_val = tsdb_ratio_data(f.readline(), method)

		with open(filename_ago, "rb") as f:
			bfb_val_ago = tsdb_ratio_data(f.readline(), method)
		
		rv = 0
		if bfb_val and bfb_val_ago:
			val = ratio(bfb_val_ago, bfb_val)
	
		if comparator(abs(val), threshold):
			rv = 1

		if rv == 0:
			floatingvalue = 0
			alert = 0
			if counter > 0 and alert == 0:
				counter = 0
				content = "状态:ok %s 主机:%s metric:%s 阀值百分比:%s 方法:环比 %s 结果百分比:|%s|" % (dt, h.ip, metric, bytes2human(threshold), symbol, bytes2human(val))
				print content
				alarm(content, groups)	
				item_update(h, name, key, val, counter, floatingvalue, alert)
			else:
				item_update(h, name, key, val, counter, floatingvalue, alert)

		elif rv == 1:
			if floatingthreshold == 0:
				floatingvalue = 0
				content = "状态:critical %s 主机:%s metric:%s 结果百分比:|%s| 方法:环比 %s 阀值百分比:%s" % (dt, h.ip, metric, bytes2human(val), symbol, bytes2human(threshold))
				counter += 1

				if attempt == 0 and alert== 0:
					alarm(content, groups)
				elif counter < attempt and alert == 0:
					alarm(content, groups)

				item_update(h, name, key, val, counter, floatingvalue, alert)

			else:
				if counter == 0:
					floatingvalue = val + floatingthreshold
					content = "状态:critical %s 主机:%s metric:%s 结果百分比:|%s| 方法:环比 %s 阀值百分比:%s" % (dt, h.ip, metric, bytes2human(val), symbol, bytes2human(threshold))
					print content
					counter += 1
					if attempt == 0 and alert == 0:
						alarm(content, groups)
					elif counter < attempt and alert == 0:
						alarm(content, groups)

					item_update(h, name, key, val, counter, floatingvalue, alert)

				elif counter > 0:
					fv = 0
					if comparator(val, floatingvalue):
						fv = 1

					if fv == 0:
						while not comparator(val, floatingvalue):
							floatingvalue -= floatingthreshold			

						print "floatingvalue:%s" % (floatingvalue)

						floatingvalue += floatingthreshold
						item_update(h, name, key, val, counter, floatingvalue, alert)

					elif fv == 1:
						content = "状态:critical %s 主机:%s metric:%s 结果百分比:|%s| 方法:环比 %s 浮动值百分比:%s" % (dt, h.ip, metric, bytes2human(val), symbol, bytes2human(floatingvalue))
						print content
						if attempt == 0 and alert == 0:
							alarm(content, groups)
						elif counter < attempt and alert == 0:
							alarm(content, groups)

						floatingvalue = val + floatingthreshold
						counter += 1
						item_update(h, name, key, val, counter, floatingvalue, alert)

@task(bind=True)
def curlmulti_tsdb(self, uuid):
	print uuid
	urls = []
	rr_obj = Round_Robin(tsdb_host)

	h = get_object_or_404(host, uuid=uuid)
	services = h.service_set.filter(alarm=0)
	for s in services:
		groups = s.group.all()
		items = s.item_set.filter(alarm=0) | s.item_set.filter(alarm=2)
		for i in items:
			if not i.symbol:
				continue

			tags = '{uuid=' + uuid + '}'
			metric = "%s.%s" % (s.name, i.key)
			if i.method == 'ratio' and i.number:
				start_time, stop_time = Unit(i.cycle+1, i.number)
				url = 'http://{0}/api/query?start={1}m-ago&m=sum:{2}{3}'.format(rr_obj.get_next()[1], i.cycle+1, metric, tags)
				url_ago = 'http://{0}/api/query?start={1}&end={2}&m=sum:{3}{4}'.format(rr_obj.get_next()[1], start_time, stop_time, metric, tags)
				tsdb_ratio(h, url, url_ago, s.name, i.key, i.method, i.symbol, i.threshold, i.floatingthreshold, i.floatingvalue, i.counter, i.attempt, groups, i.alarm)
			else:
				url = 'http://{0}/api/query?start={1}m-ago&m=sum:{2}{3}'.format(rr_obj.get_next()[1], i.cycle+1, metric, tags)
				urls.append((url, s.name, i.key, i.method, i.symbol, i.threshold, i.floatingthreshold, i.floatingvalue, i.counter, i.attempt, groups, i.alarm))
	queue = []
	for _url in urls:
		url, service, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups, alert = _url
		filename = "/tmp/%s_url_%04d" % (uuid, len(queue)+1)
		queue.append((url, filename, service, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups, alert))

	num_urls = len(urls)

	m = pycurl.CurlMulti()
	m.handles = []

	for i in range(num_conn):
		c = pycurl.Curl()
		c.fp = None
		c.setopt(pycurl.FOLLOWLOCATION, 1)
		c.setopt(pycurl.MAXREDIRS, 5)
		c.setopt(pycurl.CONNECTTIMEOUT, 5)
		c.setopt(pycurl.TIMEOUT, 10)
		c.setopt(pycurl.NOSIGNAL, 1)
		m.handles.append(c)

	#main loop
	freelist = m.handles[:]
	num_processed = 0

	while num_processed < num_urls:
		while queue and freelist:
			url, filename, service, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups, alert = queue.pop()
			c = freelist.pop()
			c.fp = open(filename, "wb")
			c.setopt(pycurl.URL, url)
			c.setopt(pycurl.WRITEDATA, c.fp)
			m.add_handle(c)
			c.filename = filename
			c.url = url
			c.service = service
			c.item = item
			c.method = method
			c.symbol = symbol
			c.threshold = threshold
			c.floatingthreshold = floatingthreshold
			c.floatingvalue = floatingvalue
			c.counter = counter
			c.attempt = attempt
			c.groups = groups
			c.alert = alert

		while 1:
			ret, num_handles = m.perform()
			if ret != pycurl.E_CALL_MULTI_PERFORM:
				break

		while 1:
			num_q, ok_list, err_list = m.info_read()
			for c in ok_list:
				c.fp.close()
				c.fp = None
				with open(c.filename, "rb") as f:
					tsdb_data(f.readline(), c.url, h, c.service, c.item, c.method, c.symbol, c.threshold, c.floatingthreshold, c.floatingvalue, c.counter, c.attempt, c.groups, c.alert)

				m.remove_handle(c)
				freelist.append(c)

			for c, errno, errmsg in err_list:
				c.fp.close()
				c.fp = None
				m.remove_handle(c)
				print "Failed:", c.url, errno, errmsg
				freelist.append(c)

			num_processed = num_processed + len(ok_list) + len(err_list)

			if num_q == 0:
				break

		m.select(1.0)

	for c in m.handles:
		if c.fp is not None:
			c.fp.close()
			c.fp = None

	m.close()
	print self.request

@task(bind=True)
def network_curlmulti_tsdb(self, uuid):
	print uuid
	urls = []
	rr_obj = Round_Robin(tsdb_host)

	n = get_object_or_404(device, uuid=uuid)		
	groups = n.group.all()
	interfaces = n.interface_set.filter(alarm=0)

	for i in interfaces:
		items = i.item_set.filter(alarm=0)
		for _i in items:
			if not _i.symbol:
				continue
			tags = '{uuid=' + uuid + '}'
			metric = "%s.%s" % (i.name, _i.key)
			url = 'http://{0}/api/query?start={1}m-ago&m=sum:{2}{3}'.format(rr_obj.get_next()[1], _i.cycle+1, metric, tags)
			urls.append((url, i.name, _i.key, _i.method, _i.symbol, _i.threshold, _i.floatingthreshold, _i.floatingvalue, _i.counter, _i.attempt, groups))

	queue = []
	for _url in urls:
		url, interface, item,  method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups = _url
		filename = "/tmp/%s_url_%04d" % (uuid, len(queue)+1)
		queue.append((url, filename, interface, item,  method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups))

	num_urls = len(urls)

	m = pycurl.CurlMulti()
	m.handles = []

	for i in range(num_conn):
		c = pycurl.Curl()
		c.fp = None
		c.setopt(pycurl.FOLLOWLOCATION, 1)
		c.setopt(pycurl.MAXREDIRS, 5)
		c.setopt(pycurl.CONNECTTIMEOUT, 5)
		c.setopt(pycurl.TIMEOUT, 10)
		c.setopt(pycurl.NOSIGNAL, 1)
		m.handles.append(c)

	#main loop
	freelist = m.handles[:]
	num_processed = 0

	while num_processed < num_urls:
		while queue and freelist:
			url, filename, interface, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups = queue.pop()
			c = freelist.pop()
			c.fp = open(filename, "wb")
			c.setopt(pycurl.URL, url)
			c.setopt(pycurl.WRITEDATA, c.fp)
			m.add_handle(c)
			c.filename = filename
			c.url = url
			c.interface = interface
			c.item = item
			c.method = method
			c.symbol = symbol
			c.threshold = threshold
			c.floatingthreshold = floatingthreshold
			c.floatingvalue = floatingvalue
			c.counter = counter
			c.attempt = attempt
			c.groups = groups

		while 1:
			ret, num_handles = m.perform()
			if ret != pycurl.E_CALL_MULTI_PERFORM:
				break

		while 1:
			num_q, ok_list, err_list = m.info_read()
			for c in ok_list:
				c.fp.close()
				c.fp = None
				with open(c.filename, "rb") as f:
					tsdb_network_data(f.readline(), c.url, n, c.interface, c.item, c.method, c.symbol, c.threshold, c.floatingthreshold, c.floatingvalue, c.counter, c.attempt, c.groups)
				m.remove_handle(c)
				freelist.append(c)

			for c, errno, errmsg in err_list:
				c.fp.close()
				c.fp = None
				m.remove_handle(c)
				print "Failed:", c.url, errno, errmsg
				freelist.append(c)

			num_processed = num_processed + len(ok_list) + len(err_list)

			if num_q == 0:
				break

		m.select(1.0)

	for c in m.handles:
		if c.fp is not None:
			c.fp.close()
			c.fp = None

	m.close()
	print self.request
