from django.conf.urls import patterns, include, url

from host.views import *

urlpatterns = patterns('host.views',
    # Examples:
    # url(r'^$', 'website.views.home', name='home'),
    # url(r'^website/', include('website.foo.urls')),

    # Uncomment the admin/doc line below to enable admin documentation:
    # url(r'^admin/doc/', include('django.contrib.admindocs.urls')),

    url(r'host/(?P<pk>\d+)/$', HostUpdateView),
    url(r'host/$', HostListView),
    url(r'host/delete/$', HostDeleteView),

    url(r'group/(?P<pk>\d+)/$',GroupUpdateView),
    url(r'group/add/$',GroupCreateView),
    url(r'^group/delete/$', GroupDeleteView),
    url(r'group/$',GroupListView),

    url(r'service/(?P<pk>\d+)/$', ServiceUpdateView),
    url(r'service/$', ServiceListView),
    url(r'service/add/$', ServiceCreateView),

    url(r'item/(?P<pk>\d+)/$', ItemUpdateView),
    url(r'item/add/$', ItemCreateView),
    url(r'item/$', "ItemListView"),

    url(r'port/(?P<pk>\d+)/$', PortUpdateView),
    url(r'port/add/$', PortCreateView),
    url(r'port/$', PortListView),

    url(r'template/(?P<pk>\d+)/$', TemplateUpdateView),
    url(r'template/add/$', TemplateCreateView),
    url(r'template/$', TemplateListView),



    url(r'server/(?P<pk>\d+)/$', ServerUpdateView),
)