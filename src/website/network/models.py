#coding:utf-8
from django.db import models
from assest.models import  *

from system.models import  *

from django.contrib.auth.models import Group
import uuid


class device(models.Model):
    uuid = models.CharField(max_length=255, unique=True,default=uuid.uuid1())
    cabinet = models.ForeignKey(cabinet, null=True, blank=True, on_delete=models.SET_NULL)
    group = models.ManyToManyField(Group,null=True,blank=True)
    ip = models.IPAddressField()
    location = models.IntegerField()
    unit = models.IntegerField()
    sn = models.CharField(max_length=255, null=True, blank=True)
    vender = models.CharField(max_length=255, null=True, blank=True)    #品牌
    model = models.CharField(max_length=255, null=True, blank=True)    #型号
    snmp_version = models.CharField(max_length=5, default="2c")
    snmp_community = models.CharField(max_length=255, default="public")
    snmp_port = models.IntegerField(default=161)
    config_update_interval = models.IntegerField(default=60)
    check_interval = models.IntegerField(default=60)
    last_check = models.DateTimeField(auto_now_add=True)
    status = models.CharField(max_length=10, null=True, blank=True, default="Ok")
    proxy = models.ForeignKey(proxy, null=True, blank=True)
    alarm = models.IntegerField(max_length=1, default=0)
    def __unicode__(self):
        return self.ip

class interface(models.Model):
    device = models.ForeignKey(device)
    index = models.IntegerField()
    name = models.CharField(max_length=50)
    mac = models.CharField(max_length=32)
    speed = models.CharField(max_length=20)
    status = models.CharField(max_length=10)
    alarm = models.IntegerField(max_length=1, default=0)
    def __unicode__(self):
	    return self.name


class oid(models.Model):
    device = models.ForeignKey(device)
    name = models.CharField(max_length=255)
    oid = models.CharField(max_length=255)
    def __unicode__(self):
        return self.name
