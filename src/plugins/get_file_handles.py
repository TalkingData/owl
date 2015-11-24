#!/usr/bin/python

try:
	import simplejson as json
except:
	import json

def num_fds():
	fd = {}
	with open('/proc/sys/fs/file-nr', 'rb') as f:
		fds = f.readline().strip().split('\t')
	fd['allocated'] = int(fds[0])
	fd['unused'] = int(fds[1])
	fd['max'] = int(fds[2])
	return json.dumps(fd, indent=4)

def main():
	print num_fds()

if __name__ == '__main__':
	main()



