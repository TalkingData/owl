# coding:utf-8
import json
import operator
from compute import Stats
from notifi import alarm
from datetime import datetime
from views import item_update, item_network_update

def ratio(data1, data2):
    r = 0
    if data1 != 0:
        r = (data2 - data1)/data1
	return int(r*100)
    return 0

def bytes2human(n):
    symbols = ('K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y')
    prefix = {}

    for i, s in enumerate(symbols):
        prefix[s] = 1 << (i + 1) * 10

    for s in reversed(symbols):
        if n >= prefix[s]:
            value = float(n) / prefix[s]
            return "%.1f%s" % (value, s)    

    return "%.1f" % (n) 

def tsdb_ratio_data(data, method):
	_data = json.loads(data)
	val = 0
	if len(_data) and type(_data) == list:
                data_dict = _data[0]
                keys = data_dict.keys()

                if 'dps' in keys:
                        if len(data_dict['dps']):
                                points = data_dict['dps'].values()
                                stats = Stats(points)
                                if method == "max":
                                        val = stats.max()
                                elif method == "min":
                                        val = stats.min()
                                elif method == "avg":
                                        val = stats.avg()
                                elif method == "sum":
                                        val = stats.sum()
                                elif method == "count":
                                        val = stats.count()
                                else:
                                        val = stats.avg()

	return val

def tsdb_data(data, url, h, service, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups):
	dt = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
	_data = json.loads(data)
	symbols = {
		'>' : 'gt',
		'>=' : 'ge',
		'<' : 'lt',
		'<=' : 'le',
		'=' : 'eq',
		'!=' : 'ne',
		'<>' : 'ne'
	}
	if len(_data) and type(_data) == list:
		data_dict = _data[0]
		metric = "%s.%s" % (service, item)
		keys = data_dict.keys()

		comparator = operator.__dict__[symbols[symbol]]
		val = 0
		rv = 0
		if 'dps' in keys:
			if len(data_dict['dps']):
				points = data_dict['dps'].values()
				stats = Stats(points)
				if method == "max":
					val = stats.max()	
				elif method == "min":
					val = stats.min()	
				elif method == "avg":
					val = stats.avg()	
				elif method == "sum":
					val = stats.sum()	
				elif method == "count":
					val = stats.count()	
				else:
					val = stats.avg()	

		if comparator(val, threshold):
			rv = 1

		if rv == 0:
			floatingvalue = 0
			if counter > 0:
				counter = 0
				content = "状态:ok %s 主机:%s metric:%s 阀值:%s 方法:%s %s 结果:%s" % (dt, h.ip, metric, bytes2human(threshold), method, symbol, bytes2human(val))
				alarm(content, groups)	
				item_update(h, service, item, val, counter, floatingvalue)
			else:
				item_update(h, service, item, val, counter, floatingvalue)

		elif rv == 1:
			if floatingthreshold == 0:
				floatingvalue = 0
				content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 阀值:%s" % (dt, h.ip, metric, bytes2human(val), method, symbol, bytes2human(threshold))
				counter += 1

				if attempt == 0:
					alarm(content, groups)
				elif counter < attempt:
					alarm(content, groups)

				item_update(h, service, item, val, counter, floatingvalue)

			else:
				if counter == 0:
					floatingvalue = val + floatingthreshold
					content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 阀值:%s" % (dt, h.ip, metric, bytes2human(val), method, symbol, bytes2human(threshold))
					counter += 1
					if attempt == 0:
						alarm(content, groups)
					elif counter < attempt:
						alarm(content, groups)

					item_update(h, service, item, val, counter, floatingvalue)

				elif counter > 0:
					fv = 0
					if comparator(val, floatingvalue):
						fv = 1

					if fv == 0:
						while not comparator(val, floatingvalue):
							floatingvalue -= floatingthreshold			

						floatingvalue += floatingthreshold
						item_update(h, service, item, val, counter, floatingvalue)

					elif fv == 1:
						content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 浮动值:%s" % (dt, h.ip, metric, bytes2human(val), method, symbol, bytes2human(floatingvalue))
						if attempt == 0:
							alarm(content, groups)
						elif counter < attempt:
							alarm(content, groups)

						floatingvalue = val + floatingthreshold
						counter += 1
						item_update(h, service, item, val, counter, floatingvalue)

def tsdb_network_data(data, url, n, interface, item, method, symbol, threshold, floatingthreshold, floatingvalue, counter, attempt, groups):
	dt = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
	_data = json.loads(data)
	symbols = {
		'>' : 'gt',
		'>=' : 'ge',
		'<' : 'lt',
		'<=' : 'le',
		'=' : 'eq',
		'!=' : 'ne',
		'<>' : 'ne'
	}
	if len(_data) and type(_data) == list:
		data_dict = _data[0]
		metric = "%s.%s" % (interface, item)
		keys = data_dict.keys()

		comparator = operator.__dict__[symbols[symbol]]
		val = 0
		rv = 0
		if 'dps' in keys:
			if len(data_dict['dps']):
				points = data_dict['dps'].values()
				stats = Stats(points)
				if method == "max":
					val = stats.max()	
				elif method == "min":
					val = stats.min()	
				elif method == "avg":
					val = stats.avg()	
				elif method == "sum":
					val = stats.sum()	
				elif method == "count":
					val = stats.count()	
				else:
					val = stats.avg()	

		if comparator(val, threshold):
			rv = 1

		if rv == 0:
			floatingvalue = 0
			if counter > 0:
				counter = 0
				content = "状态:ok %s 主机:%s metric:%s 阀值:%s 方法:%s %s 结果:%s" % (dt, n.ip, metric, bytes2human(threshold), method, symbol, bytes2human(val))
				alarm(content, groups)	
				item_network_update(n, interface, item, val, counter, floatingvalue)
			else:
				item_network_update(n, interface, item, val, counter, floatingvalue)
		elif rv == 1:
			if floatingthreshold == 0:
				floatingvalue = 0
				content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 阀值:%s" % (dt, n.ip, metric, bytes2human(val), method, symbol, bytes2human(threshold))
				counter += 1

				if attempt == 0:
					alarm(content, groups)
				elif counter < attempt:
					alarm(content, groups)

				item_network_update(n, interface, item, val, counter, floatingvalue)
			else:
				if counter == 0:
					floatingvalue = val + floatingthreshold
					content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 阀值:%s" % (dt, n.ip, metric, bytes2human(val), method, symbol, bytes2human(threshold))
					counter += 1
					if attempt == 0:
						alarm(content, groups)
					elif counter < attempt:
						alarm(content, groups)

					item_network_update(n, interface, item, val, counter, floatingvalue)

				elif counter > 0:
					fv = 0
					if comparator(val, floatingvalue):
						fv = 1

					if fv == 0:
						while not comparator(val, floatingvalue):
							floatingvalue -= floatingthreshold			

						floatingvalue += floatingthreshold
						item_network_update(n, interface, item, val, counter, floatingvalue)

					elif fv == 1:
						content = "状态:critical %s 主机:%s metric:%s 结果:%s 方法:%s %s 浮动值:%s" % (dt, n.ip, metric, bytes2human(val), method, symbol, bytes2human(floatingvalue))
						if attempt == 0:
							alarm(content, groups)
						elif counter < attempt:
							alarm(content, groups)

						floatingvalue = val + floatingthreshold
						counter += 1
						item_network_update(n, interface, item, val, counter, floatingvalue)	
