from django.conf.urls import patterns, include, url

from api.views import *

urlpatterns = patterns('host.views',
    # Examples:
    # url(r'^$', 'website.views.home', name='home'),
    # url(r'^website/', include('website.foo.urls')),

    # Uncomment the admin/doc line below to enable admin documentation:
    # url(r'^admin/doc/', include('django.contrib.admindocs.urls')),

    url(r'GetAllHostIp/$', GetAllHostIP),
    url(r'GetAllDeviceIp/$', GetAllDeviceIP),
    url(r'GetDeviceMetric/(?P<pk>\d+)$', GetDeviceMetric),
    url(r'GetHostMetric/(?P<pk>\d+)$', GetHostMetric),
    url(r'items/mchange/$',MuitlChangeItems),
)