#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import argparse
import sys
import smtplib
import json
from email.mime.text import MIMEText

mail_host = ""
mail_user = ""
mail_pass = ""


def send(sub, content, receiver):
    me = "<"+mail_user+">"
    msg = MIMEText(content,_subtype="plain",_charset="utf-8")
    msg["Subject"] = sub
    msg["From"] = me
    msg["To"] = receiver
    try:
        server = smtplib.SMTP()
        server.connect(mail_host)
        server.starttls()
        server.login(mail_user, mail_pass)
        server.sendmail(me, receiver, msg.as_string())
        server.quit()
        return True, "200 OK" 
    except Exception, e:
        return False, str(e)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="script for sending alarm by email")
    parser.add_argument("subject", help="the subject of the alarm")
    parser.add_argument("content", help="the content of the alarm")
    parser.add_argument("receiver", help="the alarm to send by mail")
    args = parser.parse_args()

    receiver = json.loads(args.receiver)
    status, response = send(args.subject, args.content, receiver["mail"])
    if not status:
        sys.exit("{0} {1}".format(status, response))

    print status, response
