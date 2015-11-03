from django.db import models
from django.contrib.auth.models import User


class proxy(models.Model):
    name = models.CharField(max_length=255,null=True,blank=True)
    ip = models.IPAddressField()
    status = models.IntegerField(default=0)
    last_check = models.DateTimeField(auto_now_add=True)

    def __unicode__(self):
        return self.name


class userprofile(models.Model):
    user = models.OneToOneField(User)
    realname = models.CharField(max_length=255,null=True,default='')
    phone = models.CharField(max_length=255, null=True, default='')
    weixin = models.CharField(max_length=255, null=True, default='')

    def __unicode__(self):
        return self.user.username


class alarm_history(models.Model):
    ip = models.CharField(max_length=255)
    group = models.CharField(max_length=255)
    type = models.CharField(max_length=255)
    metric = models.CharField(max_length=255)
    exp = models.CharField(max_length=255)
    datetime = models.DateTimeField(auto_now=True)