#                 **OWL快速部署文档 v5.0.0**

#### 1. Kairosdb

kairosdb可以基于内存、cassandra、hbse等，演示建议使用内存存储，生产建议使用cassandra，安装部署请参考官方文档。

​	cassandra:http://cassandra.apache.org/doc/latest/getting_started/installing.html#installation-from-binary-tarball-files

​	Kairosdb: https://kairosdb.github.io/docs/build/html/GettingStarted.html#using-with-cassandra

#### 2. MySQL

安装步骤略，安装完成初始化数据库

```shell
mysql> CREATE DATABASE `owl` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
#将下面的password替换成自己设置的密码
#如果cfc、api、controller都在本机，可以只授权localhost
#如果部署在其他机器请设置成对应机器的IP地址
mysql> GRANT ALL ON owl.* TO 'root'@'localhost' IDENTIFIED BY 'password';  
mysql> FLUSH PRIVILEGES;
mysql> use owl;
mysql> source owl-5.0.0.sql;
mysql> INSERT INTO `user` (username, password, role, status) VALUES ('admin', '21232f297a57a5a743894a0e4a801fc3', 1, 1);
mysql> exit;
```



#### 3. cfc

```shell
rpm -ivh owl-cfc-5.0.0-1.el7.centos.x86_64.rpm

vim /usr/local/owl-cfc/conf/cfc.conf
#设置数据库访问地址和账号密码
mysql_addr=127.0.0.1:3306
mysql_user=root
mysql_dbname=owl
mysql_password= 

#启动服务
/etc/init.d/owl-cfc start

#检查服务状态
netstat -nltp | grep 10020  
```



####4. repeater

```shell
#安装
rpm -ivh owl-repeater-5.0.0-1.el7.centos.x86_64.rpm

#修改配置文件
vim /usr/local/owl-repeater/conf/repeater.conf
backend=kairosdb
opentsdb_addr=127.0.0.1:4242  #修改为kairosdb安装的主机地址

#启动服务
/etc/init.d/owl-repeater start 

#检查服务
netstat -nltp | grep repeater
```



#### 5. api

```shell
#安装
rpm -ivh owl-api-5.0.0-1.el7.centos.x86_64.rpm

#创建证书目录
mkdir -p /usr/local/owl-api/certs  

#生成rsa密钥对，不要设置密码，直接回车即可
ssh-keygen -t rsa -b 2018 -f /usr/local/owl-api/certs/owl-api.key  
openssl rsa -in /usr/local/owl-api/certs/owl-api.key -pubout -outform PEM -out owl-api.key.pub

#编辑配置文件
vim /usr/local/owl-api/conf/api.conf

public_key=./certs/owl-api.key.pub
private_key =./certs/owl-api.key
#修改为controller机器的IP地址
alarm_health_check_url=http://127.0.0.1:10051  

#配置数据库访问用户密码
mysql_addr=127.0.0.1:3306
mysql_user=root
mysql_dbname=owl
mysql_password=  
timeseirs_storage=kairosdb

#设置kairosdb安装地址和端口
kairosdb_addr=127.0.0.1:8080  

#保存退出


#启动服务
/etc/init.d/owl-api start 

#检查服务端口是否监听
netstat -nltp | grep api
```



#### 6. controller

```shell
#安装
rpm -ivh owl-controller-5.0.0-1.el7.centos.x86_64.rpm 

#修改配置文件
vim /usr/local/owl-controller/conf/controller.conf 
#配置数据库地址和用户密码
mysql_addr=127.0.0.1:3306
mysql_user=root
mysql_dbname=owl
mysql_password=

#启动服务
/etc/init.d/owl-controller start

#检查服务
netstat -nltp | grep controller
```



#### 7. inspector

```shell
#安装
rpm -ivh owl-inspector-5.0.0-1.el7.centos.x86_64.rpm 

#修改配置文件
vim /usr/local/owl-inspector/conf/inspector.conf 
#配置controller地址
controller_addr=127.0.0.1:10050

#配置tsdb类型和地址
backend_tsdb=kairosdb
tsdb_addr=127.0.0.1:8080

#启动服务
/etc/init.d/owl-inspector start

#检查服务，inspector没有启动端口
ps -ef | grep inspector
```



#### 8. client

```shell
#安装
rpm -ivh owl-client-5.0.0-1.el7.centos.x86_64.rpm

#修改配置文件
vim /usr/local/owl-client/conf/client.conf 
#配置cfc和repeater的地址
cfc_addr=127.0.0.1:10020
repeater_addr=127.0.0.1:10040

#启动服务
/etc/init.d/owl-client start

#检查服务，inspector没有启动端口 
ps -ef | grep client

#查看日志是否有错误输出
tail /usr/local/owl-client/logs/client.log  
```





#### 9. frontend

```shell
#frontend 为静态资源文件，需要安装web服务器，这里以nginx为例子
#安装nginx
yum install -y nginx

#编辑配置文件
vim /etc/nginx/nginx.conf
server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;  
        #此处为修改后的root路径，其他保持不变
        root         /usr/share/nginx/html/owl-frontend;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        location / { 
        }   

        error_page 404 /404.html;
            location = /40x.html {
        }   

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }   
    } 

#解压静态文件
tar zxvf owl-frontend-5.0.0.tar.gz -C /usr/share/nginx/html/

#编辑配置
vim /usr/share/nginx/html/owl-frontend/index.html

#修改proUrl为api的安装地址
proUrl:"http://127.0.0.1:10060"

#启动nginx
nginx 

#打开浏览器，检查页面访问，默认用户名密码为 admin/admin
#登录以后进入管理页面即可创建产品线并分配主机和人员
```



