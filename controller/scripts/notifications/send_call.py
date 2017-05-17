#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import json
import argparse
import urllib2
import sys
reload(sys)
sys.setdefaultencoding('utf-8')

SERVICE_TOKEN = ""


def send(subject, content, receiver):
    api = "http://www.linkedsee.com/alarm/channel"
    data = dict()
    data["receiver"] = receiver
    data["type"] = "phone"
    data["title"] = subject
    data["content"] = content
    try:
        data = json.dumps(data, ensure_ascii=False).encode("utf-8")
        req = urllib2.Request(api, data)
        req.add_header("Servicetoken", SERVICE_TOKEN)
        res = urllib2.urlopen(req)
        result = json.loads(res.read(), encoding="utf-8")
        if "status" in result and not result["status"]:
            return False, json.dumps(result, ensure_ascii=False).encode("utf-8")
    except Exception as e:
        return False, e

    return True, json.dumps(result, ensure_ascii=False).encode("utf-8")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="script for sending alarm by linkedsee")
    parser.add_argument("subject", help="the subject of the alarm call")
    parser.add_argument("content", help="the content of the alarm call")
    parser.add_argument("receiver", help="the alarm to send by linkedsee")
    args = parser.parse_args()

    retry = 3
    while retry:
        status, response = send(args.subject, args.content, args.receiver)
        if status:
            break
        retry -= 1
    if not status:
        sys.exit("{0} {1}".format(status, response))

    print status, response
