from django.conf.urls import patterns, include, url

from graph.views import *

urlpatterns = patterns('graph.views',

    url(r'graph/(?P<pk>\d+)/$', GraphListView),
    url(r'graph/add/$', GraphCreateView),
    url(r'graph/delete/$', GraphDeleteView),
    url(r'graph/$', GraphListView),
    #url(r'metric/(?P<pk>\d+)/$', MetricUpdateView),
    url(r'metric/add/$', MetricCreateView),
    url(r'metric/delete/$', MetricDeleteView),
    url(r'metric/$', MetricListView),
    url(r'draw/$', GetIndexView),
    )