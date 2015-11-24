#!/usr/bin/env python
from __future__ import division
import sys
import os
try:
    import simplejson as json
except:
    import json

def getdirsize(dir_path):
    size = 0L
    if os.path.isdir(dir_path):
        for root, dirs, files in os.walk(dir_path):
            for name in files:
                f = os.path.join(root, name)
                if not os.path.islink(f):
                    size += os.path.getsize(f)

    return size

def main():
    if len(sys.argv) > 1:
        dir_path = sys.argv[1]
        size = getdirsize(dir_path)
        print json.dumps(size, indent=4)

if __name__ == '__main__':
    main()
