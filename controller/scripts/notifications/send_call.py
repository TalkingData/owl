#! /usr/bin/env python2
# encoding:utf-8

import httplib
import json
import random
import time
import hashlib
import argparse

'''语音通知发送'''
class VoicePromptSender:
    appid = 0
    appkey = ""
    url = "/v5/tlsvoicesvr/sendvoiceprompt"
    template = "电话报警:"
    def __init__(self, appid, appkey):
        self.appid = appid
        self.appkey = appkey
        self.util = SmsSenderUtil()

    """ 语音验证码发送
    Returns:
        请求包体
        {
            "tel": {
                "nationcode": "86", //国家码
                "mobile": "13788888888" //手机号码
            },
            "prompttype": 2, //语音类型，目前固定为2
            "promptfile": "语音内容文本", //通知内容，utf8编码，支持中文英文、数字及组合，需要和语音内容模版相匹配
            "playtimes": 2, //播放次数，可选，最多3次，默认2次
            "sig": "30db206bfd3fea7ef0db929998642c8ea54cc7042a779c5a0d9897358f6e9505", //app凭证，具体计算方式见下注
            "time": 1457336869, //unix时间戳，请求发起时间，如果和系统时间相差超过10分钟则会返回失败
            "ext": "" //用户的session内容，腾讯server回包中会原样返回，可选字段，不需要就填空。
        }
        应答包体
        {
            "result": 0, //0表示成功，非0表示失败
            "errmsg": "OK", //result非0时的具体错误信息
            "ext": "", //用户的session内容，腾讯server回包中会原样返回
            "callid": "xxxx" //标识本次发送id，标识一次下发记录
        }
    参数说明:
        nation_code: 国家码，如 86 为中国
        phone_number: 不带国家码的手机号
        msg: 信息内容，必须与申请的模板格式一致，否则将返回错误
        ext: 服务端原样返回的参数，可填空串
    """
    def send(self,nation_code, phone_number,playtimes,msg, ext):
        rnd = self.util.get_random()
        cur_time = self.util.get_cur_time()

        data = {}
        tel = {"nationcode": nation_code, "mobile": phone_number}
        data["tel"] = tel
        data["prompttype"] = 2
        data["promptfile"] = msg
        data["playtimes"] = playtimes
        data["sig"] = hashlib.sha256("appkey=" + self.appkey + "&random=" + str(rnd)
                                     + "&time=" + str(cur_time) + "&mobile=" + phone_number).hexdigest()
        data["time"] = cur_time
        data["ext"] = ext

        whole_url = self.url + "?sdkappid=" + str(self.appid) + "&random=" + str(rnd)
        result = self.util.send_post_request("yun.tim.qq.com", whole_url, data)
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

    def calculate_sig_for_templ_phone_numbers(self, appkey, rnd, cur_time, phone_numbers):
        """ 计算带模板和手机号列表的 sig """
        phone_numbers_string = phone_numbers[0]
        for i in range(1, len(phone_numbers)):
            phone_numbers_string += "," + phone_numbers[i]
        return hashlib.sha256("appkey=" + appkey + "&random=" + str(rnd) + "&time="
                              + str(cur_time) + "&mobile=" + phone_numbers_string).hexdigest()

    def calculate_sig_for_templ(self, appkey, rnd, cur_time, phone_number):
        phone_numbers = [phone_number]
        return self.calculate_sig_for_templ_phone_numbers(appkey, rnd, cur_time, phone_numbers)

    def phone_numbers_to_list(self, nation_code, phone_numbers):
        tel = []
        for phone_number in phone_numbers:
            tel.append({"nationcode": nation_code, "mobile":phone_number})
        return tel

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
    parser = argparse.ArgumentParser(description="script for sending alarm call")
    parser.add_argument("subject", help="the subject of the alarm call")
    parser.add_argument("content", help="the content of the alarm call")
    parser.add_argument("receiver", help="the phone number who receive the call")
    args = parser.parse_args()
    receiver = json.loads(args.receiver)
    vps = VoicePromptSender(VoicePromptSender.appid, VoicePromptSender.appkey)
    status, response = vps.send(86, receiver["phone"], 2, "{0}{1}".format(VoicePromptSender.template, args.content), "")
    if not status:
        sys.exit("{0} {1}".format(status, response))
    print status, response
