#! /usr/bin/env python2
# encoding:utf-8
# python 2.7 测试通过
# python 3 更换适当的开发库就能使用，在此我们不额外提供

import httplib
import json
import hashlib
import random
import time
import argparse

class SmsSingleSender:
    """ 单发类定义"""
    appid = 0
    appkey = ""
    url = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms"
    template = "短信报警:"

    def __init__(self, appid, appkey):
        self.appid = appid
        self.appkey = appkey
        self.util = SmsSenderUtil()

    def send(self, sms_type, nation_code, phone_number, msg, extend, ext):
        """ 普通群发接口
        明确指定内容，如果有多个签名，请在内容中以【】的方式添加到信息内容中，否则系统将使用默认签名

        Args:
            sms_type: 短信类型，0 为普通短信，1 为营销短信
            nation_code: 国家码，如 86 为中国
            phone_number: 不带国家码的手机号
            msg: 信息内容，必须与申请的模板格式一致，否则将返回错误
            extend: 扩展码，可填空串
            ext: 服务端原样返回的参数，可填空串

        Returns:
            json string { "result": xxxx, "errmsg": "xxxxx" ... }，被省略的内容参见协议文档
            请求包体
            {
                "tel": {
                    "nationcode": "86",
                    "mobile": "13788888888"
                },
                "type": 0,
                "msg": "你的验证码是1234",
                "sig": "fdba654e05bc0d15796713a1a1a2318c",
                "time": 1479888540,
                "extend": "",
                "ext": ""
            }
            应答包体
            {
                "result": 0,
                "errmsg": "OK",
                "ext": "",
                "sid": "xxxxxxx",
                "fee": 1
            }
        """
        rnd = self.util.get_random()
        cur_time = self.util.get_cur_time()

        data = {}

        tel = {"nationcode": nation_code, "mobile": phone_number}
        data["tel"] = tel
        data["type"] = sms_type
        data["msg"] = msg
        data["sig"] = hashlib.sha256("appkey=" + self.appkey + "&random=" + str(rnd)
                                     + "&time=" + str(cur_time) + "&mobile=" + phone_number).hexdigest()
        data["time"] = cur_time
        data["extend"] = extend
        data["ext"] = ext

        whole_url = self.url + "?sdkappid=" + str(self.appid) + "&random=" + str(rnd)
        result  = self.util.send_post_request("yun.tim.qq.com", whole_url, data)
        obj = json.loads(result)
        if obj["result"] == 0 and obj["errmsg"] == "OK":
            return True, result
        else:
            return False, result


class SmsSenderUtil:
    """ 工具类定义 """

    def get_random(self):
        return random.randint(100000, 999999)

    def get_cur_time(self):
        return long(time.time())

    def calculate_sig(self, appkey, rnd, cur_time, phone_numbers):
        phone_numbers_string = phone_numbers[0]
        for i in range(1, len(phone_numbers)):
            phone_numbers_string += "," + phone_numbers[i]
        return hashlib.sha256("appkey=" + appkey + "&random=" + str(rnd) + "&time=" + str(cur_time)
                              + "&mobile=" + phone_numbers_string).hexdigest()

    # def calculate_sig_for_templ_phone_numbers(self, appkey, rnd, cur_time, phone_numbers):
    #     """ 计算带模板和手机号列表的 sig """
    #     phone_numbers_string = phone_numbers[0]
    #     for i in range(1, len(phone_numbers)):
    #         phone_numbers_string += "," + phone_numbers[i]
    #     return hashlib.sha256("appkey=" + appkey + "&random=" + str(rnd) + "&time="
    #                           + str(cur_time) + "&mobile=" + phone_numbers_string).hexdigest()

    # def calculate_sig_for_templ(self, appkey, rnd, cur_time, phone_number):
    #     phone_numbers = [phone_number]
    #     return self.calculate_sig_for_templ_phone_numbers(appkey, rnd, cur_time, phone_numbers)

    # def phone_numbers_to_list(self, nation_code, phone_numbers):
    #     tel = []
    #     for phone_number in phone_numbers:
    #         tel.append({"nationcode": nation_code, "mobile":phone_number})
    #     return tel

    def send_post_request(self, host, url, data):
        con = None
        try:
            con = httplib.HTTPSConnection(host)
            con.request('POST', url, json.dumps(data))
            response = con.getresponse()
            if '200' != str(response.status):
                obj = {}
                obj["result"] = -1
                obj["errmsg"] = "connect failed:\t"+str(response.status) + " " + response.reason
                result = json.dumps(obj)
            else:
                result = response.read()
        except Exception,e:
            obj = {}
            obj["result"] = -2
            obj["errmsg"] = "connect failed:\t" + str(e)
            result = json.dumps(obj)
        finally:
            if con:
                con.close()
        return result

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="script for sending alarm sms_type")
    parser.add_argument("subject", help="the subject of the alarm sms")
    parser.add_argument("content", help="the content of the alarm sms")
    parser.add_argument("receiver", help="the phone number who receive the sms")
    args = parser.parse_args()

    ss = SmsSingleSender(SmsSingleSender.appid, SmsSingleSender.appkey)
    receiver = json.loads(args.receiver)
    status, response = ss.send(0, 86, receiver["phone"], "{0}{1}".format(SmsSingleSender.template, args.content), "", "")
    print status, response
