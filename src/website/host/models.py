#coding:utf-8
from django.db import models
from assest.models import *
from django.contrib.auth.models import Group
from network.models import oid,interface
from system.models import *
class host(models.Model):
    server = models.ForeignKey(server, null=True, blank=True, default=None)
    uuid = models.CharField(max_length=255, unique=True)
    ip = models.IPAddressField(blank=True, null=True)
    idrac = models.IPAddressField(null=True, blank=True)
    group = models.ManyToManyField('group', null=True, blank=True)
    template = models.ManyToManyField('template', blank=True, null=True)
    hostname = models.CharField(max_length=255, null=True, blank=True)
    os = models.CharField(max_length=255, null=True, blank=True)
    kernel = models.CharField(max_length=255, null=True, blank=True)
    last_check = models.DateTimeField(auto_now_add=True)
    proxy = models.ForeignKey(proxy, blank=True, null=True)
    status = models.IntegerField(default=3, blank=True, null=True)
    alarm = models.IntegerField(max_length=1, default=0)
    c_time = models.DateTimeField(auto_now_add=True)
    
    def __unicode__(self):
        return self.uuid

class template(models.Model):
    name = models.CharField(max_length=255)
    def __unicode__(self):
        return self.name

class group(models.Model):
    name = models.CharField(max_length=255)
    template = models.ManyToManyField('template', null=True, blank=True)
    def __unicode__(self):
        return self.name


class service(models.Model):
    host = models.ForeignKey(host, null=True, blank=True)
    group = models.ManyToManyField(Group, null=True, blank=True)
    name = models.CharField(max_length=255)
    plugin = models.CharField(max_length=255)
    args = models.CharField(max_length=255, blank=True, null=True)
    exec_interval = models.IntegerField()
    status = models.IntegerField(max_length=1, null=True, blank=True)  #delete
    template = models.ForeignKey("template", null=True, blank=True)
    alarm = models.IntegerField(max_length=1, default=0)
    def __unicode__(self):
        return self.name

class item(models.Model):
    service = models.ForeignKey(service, null=True, blank=True)
    interface = models.ForeignKey(interface, null=True, blank=True)
    oid  = models.ForeignKey(oid, null=True, blank=True)
    key = models.CharField(max_length=255)                          #指标
    alarm = models.IntegerField(default=0)                          #报警开关 0报警 1 不报警
    last_check = models.DateTimeField(auto_now_add=True)
    dt = models.CharField(max_length=10,default="GAUGE")           #数据类型 COUNTER,GAUGE,AVG
    last_check = models.DateTimeField(auto_now=True)                #最后一次检查时间
    duration = models.DateTimeField(auto_now_add=True)              #状态变更时间
    attempt = models.IntegerField(default=5)                        #最多报警次数
    counter = models.IntegerField(default=0)                        #错误计数
    symbol = models.CharField(max_length=255)                       #比较运算符
    method = models.CharField(max_length=255)                       #运算方法
    threshold = models.BigIntegerField(default=0)                   #报警阈值
    units = models.CharField(max_length=10)                         #单位
    current = models.BigIntegerField(default=0)                     #当前值
    cycle = models.IntegerField(default=5)                          #统计周期
    drawing = models.BooleanField(default=False)
    floatingthreshold = models.BigIntegerField(default=0)           #浮动阀值
    floatingvalue = models.BigIntegerField(default=0)               #浮动值
    number = models.IntegerField(default=0)
    def __unicode__(self):
        return self.key



class port(models.Model):
    host = models.ForeignKey(host)
    port = models.IntegerField()
    proc_name = models.CharField(max_length=255)
    alias = models.CharField(max_length=255,null=True,blank=True)
    alarm = models.IntegerField(max_length=1, default=1)
    status = models.IntegerField(max_length=1, default=0, blank=True)
    def __unicode__(self):
        return self.proc_name
		
class agent(models.Model):
    version = models.CharField(max_length=10)
    timestramp = models.DateTimeField()
