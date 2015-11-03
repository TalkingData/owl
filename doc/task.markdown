								安装文档

+ 依赖包安装

```
yum install libcurl libcurl-devel
pip install django==1.7.5 MySQL-python gevent celery django-celery numpy
```
pycurl install:

```
export PYCURL_SSL_LIBRARY=nss
easy_install pycurl
```
+ 首先需要任务数据库添加一条自动添加主机的定时任务
insert into djcelery_intervalschedule(every, period) values (5, 'minutes');

insert into djcelery_periodictask(name, task, interval_id) select 'add_task', 'task.run.add_task', id from djcelery_intervalschedule where every=5 and period='minutes';

+ 把定期任务加入到queue

```
	python manager.py celery beat
```

+ 执行定期任务

```
	python manager.py celery worker -P gevent
```

