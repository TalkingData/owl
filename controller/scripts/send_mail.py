#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import smtplib
import sys
from email.mime.text import MIMEText

mail_host=""
mail_user=""
mail_pass=""

def send(sub, content, receiver):
    me="<"+mail_user+">"
    msg = MIMEText(content,_subtype='plain',_charset='utf-8')
    msg['Subject'] = sub
    msg['From'] = me
    msg['To'] = receiver
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
    script_name, subject, content, receiver = sys.argv[:]
    result, response = send(subject, content, receiver)
    if not result:
        sys.exit("{0} {1}".format(result, response))
    print result, response
