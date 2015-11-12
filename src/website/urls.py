from django.conf.urls import patterns, include, url

# Uncomment the next two lines to enable the admin:
from django.contrib import admin
admin.autodiscover()

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'website.views.home', name='home'),
    # url(r'^website/', include('website.foo.urls')),

    # Uncomment the admin/doc line below to enable admin documentation:
    # url(r'^admin/doc/', include('django.contrib.admindocs.urls')),

    # Uncomment the next line to enable the admin:
    url(r'^$', 'graph.views.DashBoard'),
    url(r'^admin/', include(admin.site.urls)),
    url(r'^host/', include('host.urls')),
    url(r'^api/', include('api.urls')),
    url(r'^network/', include('network.urls')),
    url(r'^assest/', include('assest.urls')),
    url(r'^system/', include('system.urls')),
    url(r'^graph/',include('graph.urls')),
    url(r'^logout/$', 'system.views.logout'),
    url(r'^login/$', 'system.views.login'),
    url(r'^get_data/$', 'graph.views.get_data'),
    #url(r'^graph/$','graph.views.graph_index'),
    url(r'^\w+/(?P<model>\w+)/delete/$', 'system.views.Delete'),
    url(r'^\w+/(?P<model>\w+)/(?P<method>[a-z]+)/$', 'system.views.ChangeStatus'),
    url(r'^appmonitor/', 'task.views.appmonitor'),
    url('^addAppMonitor', 'task.views.addAppMonitor'),
    url(r'^task/', include('task.urls')),	
)
