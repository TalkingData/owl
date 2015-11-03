#coding:utf-8
from django.db import models

'''
name：idc名称
address：idc地址
tel:联系方式
remark：备注
'''
class idc(models.Model):
    name = models.CharField(max_length=255)
    address = models.CharField(max_length=255, null=True, blank=True)
    tel = models.CharField(max_length=255,null=True, blank=True)
    remark = models.CharField(max_length=255, null=True, blank=True)

    def __unicode__(self):
        return self.name

'''
idc：机房信息
sn：机柜编号
'''
class cabinet(models.Model):
    idc = models.ForeignKey(idc)
    name = models.CharField(max_length=255)
    capacity = models.IntegerField() #机柜容量, 单位U

    def get_used_unit(self):
       return [ int(unit.unit) for unit in self.units_set.filter(used=True)]

    def unit_isused(self, start, end):
        return True if self.units_set.filter(
            unit__gte=start,
            unit__lte=end,
            used = True
        ) else False

    def set_unit(self,start, end):
        self.units_set.filter(
            unit__gte=start,
            unit__lte=end
        ).update(used=True)

    def clear_unit(self,start, end):
        self.units_set.filter(
                unit__gte=start,
                unit__lte=end
            ).update(used=False)

    def __unicode__(self):
        return self.name

class units(models.Model):
    cabinet = models.ForeignKey(cabinet)
    unit = models.IntegerField()
    used = models.BooleanField(default=False)

'''
name：名称，唯一
sn：序列号
model：型号
vender：品牌
'''
class cpu(models.Model):
    server = models.ForeignKey('server')
    name = models.CharField(max_length=255)
    sn = models.CharField(max_length=255, blank=True, null=True)
    model = models.CharField(max_length=255, blank=True, null=True)
    vender = models.CharField(max_length=255, blank=True, null=True)

    def __unicode__(self):
        return self.sn

'''
sn：序列号
vender：品牌
product_id：id号
slot：槽位
'''
class disk(models.Model):
    server = models.ForeignKey('server')
    sn = models.CharField(max_length=255)                                       #序列号
    vender = models.CharField(max_length=255, blank=True, null=True)            #供应商
    product_id = models.CharField(max_length=255, blank=True, null=True)        #产品id
    slot = models.CharField(max_length=255, blank=True, null=True)              #硬盘槽位
    production_date = models.CharField(max_length=255, blank=True, null=True)   #生产日期
    capacity = models.CharField(max_length=255, blank=True, null=True)          #容量
    speed = models.CharField(max_length=255, blank=True, null=True)             #速率
    bus = models.CharField(max_length=255, blank=True, null=True)               #总线类型
    media = models.CharField(max_length=255, blank=True, null=True)             #介质类型HDD,SSD
    def __unicode__(self):
        return self.sn

class nic(models.Model):
    server = models.ForeignKey('server')
    name = models.CharField(max_length=255, blank=True, null=True)          #网卡名称
    mac = models.CharField(max_length=255, blank=True, null=True)           #mac地址
    vender = models.CharField(max_length=255, blank=True, null=True)        #供应商
    description = models.CharField(max_length=255, blank=True, null=True)   #详细信息
    def __unicode__(self):
        return self.name

'''
sn：序列号
size：容量
speed：速率
locator：槽位
vender：品牌
'''
class memory(models.Model):
    server = models.ForeignKey('server')
    sn = models.CharField(max_length=40)
    size = models.CharField(max_length=30, blank=True, null=True)
    speed = models.CharField(max_length=30, blank=True, null=True)
    locator = models.CharField(max_length=30, blank=False, null=False)
    vender = models.CharField(max_length=40, blank=True, null=True)

    def __unicode__(self):
        return self.sn

'''
sn：序列号
cabinet:机柜号
vender：品牌
model：型号
cpus：cpu信息
memorys：内存信息
disks：磁盘信息
nics：网卡信息
'''
class server(models.Model):
    sn = models.CharField(max_length=255)
    cabinet = models.ForeignKey(cabinet, null=True, blank=True, on_delete=models.SET_NULL)
    location = models.IntegerField()
    unit = models.IntegerField()
    vender = models.CharField(max_length=255, null=True, blank=True)
    model = models.CharField(max_length=255, null=True, blank=True)

    def has_cabinet(self):
        if self.cabinet:
            return True
        return False


    def __unicode__(self):
        return self.sn

