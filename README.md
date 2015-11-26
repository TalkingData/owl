# OWL-v3

distributed monitoring system

OWL-v3是TalkingData公司推出的一款开源分布式监控系统，感兴趣的DevOps可以下载使用，我们也希望更多的人员能够加入我们的

OWL-V3目前用到的技术主要以下几种：

1.Golang和python.

2.hadoop,hbase,zookeeper.

3.opentsdb 2.0.1具体的安装可以参考官方的安装文档.

具体的环境搭建参考相关的文档

功能介绍：
目前OWL采用独立的agent常驻OS，监控指标的采集支持自己写插件的方式，采集不同的监控指标，采集后采用socket的私有协议，上行server端，服务器端采用批量socket的接口，采用batch的方式插入到opentsdb里面，后端报警处理模块实现不同的算法报警。

目前报警算法上支持

1.单个Metrics的固定阀值报警，这是比较常规的报警规则，指标采集之后到达一定的阀值，触发报警。

2.浮动阀值报警，就是说阀值可以设置一个固定，还有一个浮动的参数可以设置，数值采集到达固定阀值之后，系统会根据浮动的数据设定，增加一个数据，故障恢复之后，系统会自动恢复固定阀值。

3.环比报警，OWL目前可以设置一个历史的周期，在算法报警上自动计算当前的数值和之前指定的历史周期的数值的差异，触发报警，这样有效的降低报警设置的合理性。




###Docker Image Address: 
http://pan.baidu.com/s/1eQlXNjk

下载完成解压: 
```Bash
tar zxf  owl-docker.tar.gz
```
导入镜像：
```Bash
docker load < owl.tar 
```
查看镜像：
```Bash
docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
<none>              <none>              6353ef07ffe1        2 weeks ago         1.994 GB
```
启动一个容器：
```Bash
docker run -d -p 8080:80 6353ef07ffe1 /etc/init.d/start-all
```
打开浏览器访问即可
默认用户名密码为：root/123456

###QQ交流群
492850035
