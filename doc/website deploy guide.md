# website部署指南

---
### MySQL安装
（略）

### Python 依赖包安装
- Django 1.7.5
- MySQL-python
- gevent
- celery
- django-celery
- numpy
- pycurl

### 数据库设置
- 创建数据库
```
CREATE DATABASE `owl`  /*!40100 DEFAULT CHARACTER SET utf8 */ ;
```
- 数据库授权
```
#将语句中的ip地址替换成部署owl机器的ip地址
#password替换成自己设置的密码
GRANT ALL ON owl.* to owl@'ip' identified by 'password';
```
- 修改website配置文件

```
cd OWL-v3/src/website
编辑settings.py，配置好数据库设置以及opentsdb地址，保存即可

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.mysql',
        'NAME': 'owl',                      
        'USER': 'owl',
        'PASSWORD': 'owl',
        'HOST':'localhost',
        'PORT': '', 
    }
}
OPENTSDB_ADDR＝"127.0.0.1:4242"
```
- 生成数据库表结构
```
cd OWL-v3/src/website
python manage.py syncdb
#根据提示建立超级用户，可用于登陆管理界面
```

### 启动website
```
python manage.py runserver 0.0.0.0:80
```
打开浏览器http://ip
即可进入登陆页面

注：实际使用建议通过nginx或apache来发布
