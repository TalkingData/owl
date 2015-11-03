from django.conf.urls import patterns, include, url

from network.views import *

urlpatterns = patterns('',
    url(r'^getAvailableUintsByID/$', 'network.views.getAvailableUintsByID'),
    url(r'^device/(?P<pk>\d+)/$', DeviceUpdateView),
    url(r'^device/add/$', DeviceCreateView),
    url(r'^device/$', DeviceListView),
    url(r'^interface/$', InterfaceListView),
    url(r'^oid/(?P<pk>\d+)/$', CustomOIDUpdateView),
    url(r'^oid/add/$', CustomOIDCreateView),
    url(r'^oid/$', CustomOIDListView),
    url(r'^item/(?P<pk>\d+)/$', ItemUpdateView),
    url(r'^item/add/$', ItemCreateView),
    url(r'^item/$', ItemListView),
)