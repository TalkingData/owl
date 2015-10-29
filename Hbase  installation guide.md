#  Hbase 基础与集群部署
------
    环境描述：
| System| Ipaddress|  Hostname  |Service|Version|
| --------   | -----:| :----:|
| Centos 7.1 x64|10.10.3.179|hadoop0|Hmaster+zookeeper|hbase-1.1.2，zookeeper-3.4.6|
| Centos 7.1 x64|10.10.3.180|hadoop1|HRegionServer|hbase-1.1.2，zookeeper-3.4.6|
| Centos 7.1 x64|10.10.3.181|hadoop2|HRegionServer|hbase-1.1.2，zookeeper-3.4.6|
## 1.ZooKeeper

### 1. zookeeper作用
    可以用来保证数据在zk集群之间的数据的事务性一致。

###  2.如何搭建ZooKeeper服务器集群

#### 2.1 zk服务器集群规模不小于3个节点，要求各服务器之间系统时间要保持一致
#### 2.2 在hadoop0的/usr/local目录下，解压缩zk....tar.gz，设置环境变量
#### 2.3 在conf目录下，修改文件 vi zoo_sample.cfg  zoo.cfg
#### 2.4 编辑该文件，执行vi zoo.cfg
#### 修改dataDir=/usr/local/zookpeer/zk_data
新增
```python
   server.0=hadoop0:2888:3888
   server.1=hadoop1:2888:3888
   server.2=hadoop2:2888:3888
```
#### 2.5 创建文件夹mkdir /usr/local/zookpeer/zk_data
#### 2.6 在data目录下，创建文件myid，值为0
#### 2.7 把zk目录复制到hadoop1和hadoop2中
#### 2.8 把hadoop1中相应的myid的值改为1
   把hadoop2中相应的myid的值改为2
#### 2.9 启动，在三个节点上分别执行命令zkServer.sh start
#### 2.10 检验，在三个节点上分别执行命令zkServer.sh status

## 2.Hbase 基础

### 1.HBase(NoSQL)的数据模型
#### 1.1 表(table)，是存储管理数据的。
#### 1.2 行键(row key)，类似于MySQL中的主键。
      行键是HBase表天然自带的。 
#### 1.3 列族(column family)，列的集合。
    HBase中列族是需要在定义表时指定的，列是在插入记录时动态增加的。
    HBase表中的数据，每个列族单独一个文件。
#### 1.4 时间戳(timestamp)，列(也称作标签、修饰符)的一个属性。
    行键和列确定的单元格，可以存储多个数据，每个数据含有时间戳属性，数据具有版本特性。 
    如果不指定时间戳或者版本，默认取最新的数据。
#### 1.5 存储的数据都是字节数组。
#### 1.6 表中的数据是按照行键的顺序物理存储的。


### 2.HBase的物理模型

#### 2.1 HBase是适合海量数据(如20PB)的秒级简单查询的数据库。
#### 2.2 HBase表中的记录，按照行键进行拆分， 拆分成一个个的region。
    许多个region存储在region server(单独的物理机器)中的。
    这样，对表的操作转化为对多台region server的并行查询。

### 3.HBase的体系结构

    HBase是主从式结构，HMaster、HRegionServer

### 4.HBase伪分布安装
#### 4.1 解压缩、重命名、设置环境变量
#### 4.2 修改$HBASE_HOME/conf/hbase-env.sh，修改内容如下：
```python
    export JAVA_HOME=/usr/local/jdk
    export HBASE_MANAGES_ZK=true
```
#### 4.3 修改$HBASE_HOME/conf/hbase-site.xml，修改内容如下：
```python
<property>
        <name>hbase.rootdir</name>
        <value>hdfs:// hadoop0:9000/hbase</value>
</property>

<property>
        <name>hbase.cluster.distributed</name>
        <value>true</value>
</property>

<property>
        <name>hbase.zookeeper.quorum</name>
        <value> hadoop0 </value>
</property>
 
<property>
        <name> dfs.replication </name>
        <value> 1</value>
</property>
```
#### 4.3 (可选)文件regionservers的内容为hadoop0
#### 4.4 启动hbase，执行命令start-hbase.sh
    ******启动hbase之前，确保hadoop是运行正常的，并且可以写入文件*******
#### 4.5 验证
    (1)执行jps，发现新增加了3个java进程，分别是HMaster、HRegionServer、HQuorumPeer 
    (2)使用浏览器访问http://hadoop0:16010



## 3.Hbase集群
### 1.hbase的机群搭建过程(在原来的hadoop0上的hbase伪分布基础上进行搭建)
#### 1.1 集群结构，主节点(hmaster)是hadoop0，从节点(region server)是hadoop1和hadoop2
#### 1.2 修改hadoop0上的hbase的几个文件
    (1)修改hbase-env.sh的最后一行export HBASE_MANAGES_ZK=false
    (2)修改hbase-site.xml文件的hbase.zookeeper.quorum的值为hadoop0,hadoop1,hadoop2
    (3)修改regionservers文件(存放的region server的hostname)，内容修改为hadoop1、hadoop2
#### 1.3 复制hadoop0中的hbase文件夹到hadoop1、hadoop2中 
    复制hadoop0中的/etc/profile到hadoop1、hadoop2中，在hadoop1、hadoop2上执行source /etc/profile
#### 1.4 首先启动hadoop，然后启动zookeeper集群。
    最后在hadoop0上启动hbase集群。
```python
[hadoop@master ~]$ start-hbase.sh
```
