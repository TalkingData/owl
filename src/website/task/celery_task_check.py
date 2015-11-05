# coding:utf8
import hashlib
import os
from urllib import quote, urlopen
from datetime import datetime

def alert(content):
    phone = ("18611610862", "13810520844", "15210190629")
    email = ("tong.li", "hierarch.pan")
    alert_url="http://10.10.32.10:9000/sendalart?severity=1&msg={0}&phones={1}&emails={2}&subject=报警".format(content, ','.join(phone), ','.join(email))
    print urlopen(quote(alert_url, ':/=&()?,>.')).read().strip()

def main():
    celery_log = "/deploy/src/website/celery.log"
    dt = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    content = "%s: tsdb 5 minutes don't excute task" % (dt)

    if os.path.exists(celery_log):
        with open(celery_log, 'rb') as f:
            md5_new = hashlib.md5()
            md5_new.update(f.read())
            hash_new = md5_new.hexdigest()
            print "%s: hash_new:" % (dt) + hash_new

        md5_temp_file = "/tmp/celery.txt"

        if os.path.exists(md5_temp_file):
            with open(md5_temp_file, 'r') as f:
                hash_old = f.read()
                print "%s: hash_old:" % (dt) + hash_old
                if hash_new == hash_old: 
                    print "no ok"
                    alert(content) 
                else:
                    print "ok"

            with open(md5_temp_file, 'w') as f:
                f.write(hash_new)

        else:
            with open(md5_temp_file, 'w') as f:
                f.write(hash_new)

    else:
        print "no ok"

if __name__ == "__main__":
    main()
