from django.db import models
from django.contrib.auth.models import User
# Create your models here.

class graph(models.Model):
    user = models.ForeignKey(User)
    title = models.CharField(max_length=255)
    start = models.CharField(max_length=255)    #convert to second
    def __unicode__(self):
        return self.title

class metric(models.Model):
    graph = models.ForeignKey(graph)
    metric = models.CharField(max_length=255)
    uuid = models.CharField(max_length=255)
    ip = models.CharField(max_length=255)
    def __unicode__(self):
        return self.metric
