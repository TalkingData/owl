#coding:utf8
from assest.models import *
from host.models import *
from network.models import *
from .models import *
from graph.models import *
from django.shortcuts import HttpResponse, RequestContext
from django.shortcuts import render_to_response
from django.http import HttpResponseRedirect,Http404
from django.contrib import auth
from django.contrib.auth.models import User,Group
from django.contrib.auth.decorators import login_required
from django.core.paginator import Paginator,EmptyPage,PageNotAnInteger,InvalidPage
from django.db.models import Q
from django.conf import settings
import json
def paging(queryset, page, pagesize):
    paginator = Paginator(queryset,pagesize)
    try:
        newpage=int(page)
    except ValueError:
        newpage=1
    try:
        contacts = paginator.page(newpage)
    except PageNotAnInteger:
        contacts = paginator.page(1)
    except EmptyPage:
        contacts = paginator.page(paginator.num_pages)
    return contacts



def ChangeStatus(request, model=None, method=None):
    status = -1
    context = {}
    ids = request.POST.get("ids")
    if model and method:
        if method == "enable":status = 0
        if method == "disable":status = 1
    if status == -1 or ids is None:
        context["status"] = 1
    else:
        obj = []
        flag = -1
        for id in ids.split(","):
            try:
                qs = eval(model).objects.get(pk=int(id))
                qs.alarm = status
                obj.append(qs)
            except:
                flag = 0
        if flag == -1:
            for qs in obj:
                qs.save()
            context["status"] = 0
        else:
            context["status"] = 1
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )



def Delete(request, model=None):
    context = {}
    ids = request.POST.get("ids")
    if request.user.is_superuser is False:
        context["status"] = 1
        context["message"] = "权限拒绝.."
    else:
        for id in ids.split(","):
            if model == "user":
                User.objects.get(pk=int(id)).delete()
            elif model == "group":
                print int(id)
                Group.objects.get(pk=int(id)).delete()
            else:
                eval(model).objects.get(pk=int(id)).delete()
        context["status"] = 0
    return HttpResponse(
        json.dumps(context),
        content_type="application/json"
    )

def UserListView(request):
    if not request.user.is_superuser:
        raise Http404
    context={}
    context["users"] = User.objects.all()
    q = request.GET.get("q")
    page = request.GET.get("page")
    if not page:
        page = 1
    if q:
        context["users"] = context["users"].filter(
            Q(username__contains=q)
        )
    context["users"] = paging(context["users"], page, settings.PAGE_SIZE)
    return render_to_response(
        "system/user_list.html",
        context,
        context_instance = RequestContext(request)
    )

def UserCreateView(request):
    if not request.user.is_superuser:
        raise Http404
    context = {}
    context["groups"] = Group.objects.all()
    if request.method == "GET":
        return render_to_response(
            "system/user_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        user=User(
            username=request.POST.get("username"),
            email=request.POST.get("email"),
            is_active = False if request.POST.get("is_active") == "0" else True,
            is_superuser = False if request.POST.get("is_superuser") == "0" else True,
        )
        password = request.POST.get("password")
        if password:
            user.set_password(password)
        user.save()
        groups = request.POST.getlist("groups")
        if groups:
            for g_id in groups:
                user.groups.add(Group.objects.get(pk=int(g_id)))
        userprofile(
            user=user,
            realname=request.POST.get("realname"),
            weixin=request.POST.get("weixin"),
            phone=request.POST.get("phone")
        ).save()
        return HttpResponseRedirect(
            "/system/user/"
        )

def UserUpdateView(request,pk):
    context = {}
    if not request.user.is_superuser:
        if request.user.id != int(pk):
            raise Http404
    _user = User.objects.get(pk=int(pk))
    if request.method == "GET":
        context["cuser"] = _user
        context["groups"] = Group.objects.all()
        return render_to_response(
            "system/user_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _user.email = request.POST.get("email")
        password = request.POST.get("password")
        if password:
            _user.set_password(password)
        if request.user.is_superuser:
            is_active = request.POST.get("is_active")
            is_superuser = request.POST.get("is_superuser")
            _user.is_active = True if is_active == "1" else False
            _user.is_superuser = True if is_superuser == "1" else False
        groups = request.POST.getlist("groups")
        _user.groups.clear()
        if groups:
            for g_id in groups:
                _user.groups.add(Group.objects.get(pk=int(g_id)))
        try:
            _user.userprofile.phone = request.POST.get("phone")
            _user.userprofile.weixin = request.POST.get("weixin")
            _user.userprofile.realname = request.POST.get("realname")
            _user.userprofile.save()
        except :
            userprofile(
                user=_user,
                realname=request.POST.get("realname"),
                weixin = request.POST.get("weixin"),
                phone = request.POST.get("phone")
            ).save()
        _user.save()
        if request.user.id == int(pk):
            if len(password) != 0:
                return logout(request)
            return HttpResponseRedirect(
                "/system/user/%s/" % pk
            )
        return HttpResponseRedirect(
            "/system/user/"
        )


def login(request):
    if request.method == "GET":
        return render_to_response(
            "system/login.html",
            context_instance = RequestContext(request)
        )
    context = {}
    username = request.POST.get("username")
    password = request.POST.get("password")
    user = auth.authenticate(username=username, password=password)
    if user is not None and user.is_active:
        auth.login(request,user)
        return HttpResponseRedirect("/host/host/")
    else:
        context["message"] = u"登陆失败,请检查用户名或密码"
        return render_to_response(
            "system/login.html",
            context,
            context_instance = RequestContext(request)
        )


# @login_required(login_url='/login/')
def logout(request):
    auth.logout(request)
    return HttpResponseRedirect("/login/")


def GroupListView(request):
    context={}
    context["groups"] = Group.objects.all()
    return render_to_response(
        "system/group_list.html",
        context,
        context_instance = RequestContext(request)
    )

def GroupCreateView(request):
    context = {}
    context["users"] = User.objects.all()
    if request.method == "GET":
        return render_to_response(
            "system/group_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        group = Group(
            name = request.POST.get("name")
        )
        group.save()
        users_id = request.POST.getlist("users")
        if users_id:
            for id in users_id:
                group.user_set.add(
                    User.objects.get(pk=int(id))
                )
        return HttpResponseRedirect(
            "/system/group/"
        )


def GroupUpdateView(request,pk):
    context = {}
    _group = Group.objects.get(pk=int(pk))
    if request.method == "GET":
        context["group"] = _group
        context["users"] = User.objects.all()
        return render_to_response(
            "system/group_edit.html",
            context,
            context_instance = RequestContext(request)
        )
    else:
        _group.name = request.POST.get("name")
        _group.save()
        _group.user_set.clear()
        users_id = request.POST.getlist("users")
        if users_id:
            for id in users_id:
                _group.user_set.add(
                    User.objects.get(pk=int(id))
                )

        return HttpResponseRedirect(
            "/system/group/"
        )
