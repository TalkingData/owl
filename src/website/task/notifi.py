# coding:utf-8

import pycurl
from round_robin import Round_Robin
from urllib import urlopen, quote, quote_plus

import sys
reload(sys)
sys.setdefaultencoding('utf8')

addrs = ['192.168.1.1']
rr_obj = Round_Robin(addrs)

def alarm(content, groups):
	names = []
	phones = []
	weixins = []
	
	if groups:
		for g in groups:	
			users = g.user_set.all()
			if users:
				for u in g.user_set.all():
					names.append(u.userprofile.realname)
					phones.append(u.userprofile.phone)
					weixins.append(u.userprofile.weixin)
			
	url = "http://{0}/sendalart?severity=1&msg={1}&phones={2}&subject=报警".format(rr_obj.get_next()[1], content, ','.join(phones))
	if phones:
		urlopen(quote(url, ':/=&()?,>.')).read().strip()
	else:
		pass
