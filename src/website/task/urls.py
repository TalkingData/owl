from django.conf.urls import patterns, include, url
from task.views import *

urlpatterns = patterns('task.views',
    url(r'acknowledged/all_acknowledged/$', all_acknowledged),
    url(r'acknowledged/disable', acknowledged),
    url(r'acknowledged', alert_data),
)
