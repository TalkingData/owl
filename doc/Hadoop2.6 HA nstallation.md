#Hadoop2.6 HA部署
----

    环境描述：

System  ：  Centos 7.1 x64
Hadoop  :   hadoop-2.6.0-cdh5.4.5
Zookeeper   :   zookeeper-3.4.6
JDK     :       jdk1.7.0_80

    添加hosts记录：
```python
10.10.3.179     hadoop0
10.10.3.180     hadoop1
10.10.3.181     hadoop2
```

## 创建root互信

    [root@hadoop0 ~]# ssh-keygen -t rsa
    [root@hadoop0 ~]# cd .ssh/ ; cp id_rsa.pub authorized_keys

 Slave机器创建ssh目录 ： 
        
    mkdir -m 700 /root/.ssh

将公钥复制到Slave机器上：

    scp authorized_keys  hadoop1:/root/.ssh/

    注：由于完全模拟生产环境，故把角色尽量分开，所以需要做2次互信（hadoop1,hadoop2）机器都需要

测试无密码登录是否正常 ！！

    添加环境变量
    [root@hadoop0 src]# more /etc/profile
```python
export JAVA_HOME=/usr/java/jdk1.7.0_80
export HADOOP_HOME=/usr/local/hadoop
export ZOOKEEPER_HOME=/usr/local/zookpeer
PATH=/usr/local/cassandra/bin:$JAVA_HOME/bin:$ZOOKEEPER_HOME/bin:$HADOOP_HOME/bin:$HADOOP_HOME/sbin:$PATH
CLASSPATH=.:$JAVA_HOME/lib/dt.jar:$ZOOKEEPER_HOME/lib:$JAVA_HOME/lib/tools.jar
```
## 安装zookeeper
    
    zookeeper具体安装过程请见hbase集群部署文档。

## 安装hadoop集群

```python
[root@hadoop0 src]# tar zxf hadoop-2.6.0-cdh5.4.5.tar.gz
[root@hadoop0 src]# mv hadoop-2.6.0-cdh5.4.5 ../hadoop
[root@hadoop0 src]# cd ../hadoop/etc/hadoop/
[root@hadoop0 hadoop]# vim hadoop-env.sh 
export JAVA_HOME=/usr/local/jdk1.7.0_80
export HADOOP_OPTS="$HADOOP_OPTS -Djava.library.path=/usr/local/hadoop/lib/"
export HADOOP_COMMON_LIB_NATIVE_DIR="/usr/local/hadoop/lib/native/"
export HADOOP_CONF_DIR=${HADOOP_CONF_DIR:-"/usr/local/hadoop/etc/hadoop"}

[root@hadoop0 hadoop]# vim mapred-env.sh
export JAVA_HOME=/usr/local/jdk1.7.0_80
[root@hadoop0 hadoop]# vim yarn-env.sh
export JAVA_HOME=/usr/local/jdk1.7.0_80

[root@hadoop0 hadoop]# vim core-site.xml
<configuration>
<property>  
	<name>fs.defaultFS</name>	
	<value>hdfs://myhdfs</value>
</property>
<property>  
	<name>hadoop.tmp.dir</name>	
	<value>/usr/local/hadoop/hadoop-tmp-${user.name}</value>
</property>
<property>
	<name>hadoop.native.lib</name>
	<value>true</value>
</property>
<property>
        <name>ha.zookeeper.quorum</name>
        <value>hadoop0,hadoop1,hadoop2</value>
</property>
</configuration>
```
    [root@hadoop0 hadoop]# more hdfs-site.xml 
```python
<configuration>
<property>
	<name>dfs.nameservices</name>
	<value>myhdfs</value>
</property>
<property>
	<name>dfs.replication</name>
	<value>2</value>
</property>
<property>
	<name>dfs.ha.namenodes.myhdfs</name>
	<value>nn1,nn2</value>
</property>
<property>
	<name>dfs.namenode.rpc-address.myhdfs.nn1</name>
	<value>hadoop0:54310</value>
</property>
<property>
	<name>dfs.namenode.rpc-address.myhdfs.nn2</name>
	<value>hadoop1:54310</value>
</property>
<property>
	<name>dfs.namenode.http-address.myhdfs.nn1</name>
	<value>hadoop0:50070</value>
</property>
<property>
	<name>dfs.namenode.http-address.myhdfs.nn2</name>
	<value>hadoop1:50070</value>
</property>
<property>
	<name>dfs.namenode.shared.edits.dir</name>
	<value>qjournal://hadoop0:8485;hadoop1:8485;hadoop2:8485/hadoop-journal</value>
</property>
<property>
	<name>dfs.ha.automatic-failover.enabled</name>
	<value>true</value>
</property>
<property>
	<name>dfs.journalnode.edits.dir</name>
	<value>/usr/local/hadoop/journal</value>
</property>
<property>
	<name>dfs.client.failover.proxy.provider.myhdfs</name>
	<value>org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider</value>
</property>
<property>
	<name>dfs.ha.fencing.methods</name>
	<value>sshfence</value>
	<description>how to communicate in the switch process</description>
</property>
<property>
	<name>dfs.ha.fencing.ssh.private-key-files</name>
	<value>/root/.ssh/id_rsa</value>
	<description>the location stored ssh key</description>
</property>
<property>
	<name>dfs.ha.fencing.ssh.connect-timeout</name>
	<value>5000</value>
</property>
<property>
	<name>dfs.namenode.name.dir</name>
	<value>/usr/local/hadoop/nn</value>
</property>
<property>
	<name>dfs.datanode.data.dir</name>
	<value>/data0/dfs/data,/data1/dfs/data,/data2/dfs/data</value>
</property>
<property>
	<name>dfs.balance.bandwidthPerSec</name>
	<value>26214400</value>
</property>
<property>
	<name>dfs.safemode.threshold.pct</name>
	<value>0.9</value>
</property>
<property>
	<name>dfs.datanode.max.xcievers</name>
	<value>5120</value>
</property>
<property>
	<name>dfs.client.block.write.retries</name>
	<value>3</value>
</property>
<property>
	<name>fs.trash.interval</name>
	<value>600</value>
</property>
<property>  
	<name>dfs.client.failover.proxy.provider.myhdfs</name>                        
	<value>org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider</value>
</property>  
</configuration>
```
    [root@hadoop0 hadoop]# more mapred-site.xml
```python
<configuration>
<property>
	<name>mapreduce.framework.name</name>
	<value>yarn</value>
</property>
<property>
	<name>mapred.system.dir</name>
	<value>${hadoop.tmp.dir}/system</value>
</property>
<property>
	<name>mapred.local.dir</name>
	<value>${hadoop.tmp.dir}/mapred-local</value>
</property>
<property>
	<name>mapreduce.jobhistory.address</name>
	<value>hadoop1:54311</value>
</property>
<property>
	<name>mapreduce.jobhistory.webapp.address</name>
	<value>hadoop1:19888</value>
</property>
</configuration>
```
    [root@hadoop0 hadoop]# more yarn-site.xml
```python
<configuration>
 <property>
        <name>yarn.nodemanager.aux-services</name>
        <value>mapreduce_shuffle</value>
 </property>
 <property>                                                                
	<name>yarn.nodemanager.aux-services.mapreduce.shuffle.class</name>
        <value>org.apache.hadoop.mapred.ShuffleHandler</value>
 </property>
 <property>
        <name>yarn.resourcemanager.address</name>
        <value>hadoop1:8032</value>
</property>
<property>
        <name>yarn.resourcemanager.scheduler.address</name>
        <value>hadoop1:8030</value>
</property>
<property>
	<name>yarn.resourcemanager.resource-tracker.address</name>
	<value>hadoop1:8031</value>
</property>
<property>
        <name>yarn.resourcemanager.admin.address</name>
        <value>hadoop1:8033</value>
</property>
<property>
        <name>yarn.resourcemanager.webapp.address</name>
        <value>hadoop1:50030</value>
</property>
</configuration>
```
    [root@hadoop0 hadoop]# more slaves
```python
hadoop1
hadoop2
```
    创建数据存储目录：
```python
[root@hadoop0 hadoop]# mkdir -p /usr/local/hadoop/{nn，journal}
```
    将修改好的hadoop目录拷贝到各个节点：[root@hadoop0 local]# scp -r hadoop hadoop1:/usr/local/
    为datanode节点创建数据目录：
```python
    [root@hadoop0 sbin]# ./slaves.sh mkdir /data{0,1,2}/dfs/
```
    拷贝系统/环境变量到各个服务器：
```python
    [root@hadoop0 ~]# scp /etc/hosts /etc/profile hadoop1:/etc/
    [root@hadoop0 ~]# hdfs zkfc -formatZK
```
    第一次启动格式化HDFS
```python
    hdfs namenode -format
```
    启动hdfs服务：
    通过start-dfs.sh 直接启动所有服务:
```python
    [root@hadoop0 sbin]# ./start-dfs.sh 
    [root@hadoop0 sbin]# jps
    20422 Jps
    17813 JournalNode
    18174 DFSZKFailoverController
    17595 NameNode
```
    访问 http://hadoop0:50070 会看到该节点已经成为active
    下面需要同步一次元数据：
```python
    [root@hadoop0 sbin]#hdfs namenode -bootstrapStandby
```    
    访问 http://hadoop1:50070/dfshealth.html#tab-overview  会看到该节点已经成为standby。
    
    然后kill掉hadoop0上的active NN进程，standby NN会成为active
    注意：手动切换时，会提示下面警告。所以一般在启动zkfc的情况下也无需进行切换。
```python
[root@hadoop0 sbin]#hdfs haadmin -transitionToActive nn1
```

    启动yarn服务：
```python
[root@hadoop1 sbin]# ./start-yarn.sh 

[root@hadoop1 sbin]# jps
6994 QuorumPeerMain
13657 ResourceManager
13537 JournalNode
14174 Jps
```
