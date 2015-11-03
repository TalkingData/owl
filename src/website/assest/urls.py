from django.conf.urls import patterns, include, url

from assest.views import *

urlpatterns = patterns('assest.views',

    url(r'^idc/(?P<pk>\d+)/$', IdcUpdateView.as_view()),
    url(r'^cabinet/(?P<pk>\d+)/$', CabinetUpdateView.as_view()),

    url(r'^assest/$', AssestListView.as_view()),
    url(r'^idc/$', IdcListView.as_view()),
    url(r'^cabinet/$', CabinetListView.as_view()),

    url(r'idc/add/$', IdcCreateView.as_view()),
    url(r'cabinet/add/$', "CabinetCreateView"),
)