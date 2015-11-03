from django.conf.urls import patterns, include, url

from system.views import *

urlpatterns = patterns('system.views',
    url(r'^user/(?P<pk>\d+)/$', UserUpdateView),
    url(r'^user/$', UserListView),
    url(r'^user/add/$', UserCreateView),

    url(r'^group/(?P<pk>\d+)/$', GroupUpdateView),
    url(r'^group/$', GroupListView),
    url(r'^group/add/$', GroupCreateView),
    #url(r'^group/delete/$', GroupDelteView),
)