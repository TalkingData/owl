#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import os
import json
import argparse
import urllib2
import sys

reload(sys)
sys.setdefaultencoding('utf-8')

CORP_ID = ""
SECRET = ""
AGENT_ID = 0
FILE = "/tmp/wechat_token"


def get_access_token(new=False):
    token = ""
    if new:
        api = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={0}&corpsecret={1}"
        url = api.format(CORP_ID, SECRET)
        try:
            response = urllib2.urlopen(url)
            content  = json.loads(response.read())
            if "errcode" in content and content["errcode"] != 0:
                return False, content
            token = content["access_token"]
            with open(FILE, "w") as tmp_file:
                tmp_file.write(token)
        except Exception as e:
            return False, e

        return True, token

    if not os.path.exists(FILE):
        return False, token

    with open(FILE) as tmp_file:
        token = tmp_file.read()

    return True, token


def send(receiver, content, token):
    api = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={0}"
    url = api.format(token)
    data = {}
    data["touser"] = receiver
    data["msgtype"] = "text"
    data["agentid"] = AGENT_ID
    data["text"] = {"content": content}
    try:
        data = json.dumps(data, ensure_ascii=False).encode("utf-8")
        response = urllib2.urlopen(url, data)
        result  = json.loads(response.read())
        if result["errcode"] != 0 or result["invaliduser"] != "" or result["errmsg"] != "ok":
            return False, result
    except Exception as e:
        return False, e

    return True, result

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="script for sending alarm by wechat")
    parser.add_argument("subject", help="the subject of the alarm")
    parser.add_argument("content", help="the content of the alarm")
    parser.add_argument("receiver", help="the alarm to send by wechat")
    args = parser.parse_args()

    token_status, token = get_access_token()
    receiver = json.loads(args.receiver)
    retry = 3
    while retry:
        status, response = send(receiver["wechat"], args.content, token)
        if status:
            break
        retry -= 1
        token_status, token = get_access_token(True)
    if not status:
        sys.exit("{0} {1}".format(status, response))

    print status, response
