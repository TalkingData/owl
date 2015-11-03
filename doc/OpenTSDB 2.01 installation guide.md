#OpenTSDB 2.0.1 安装
----

    环境描述：

    System    :   Centos 7.1 x64
    Hbase     :   hbase-1.1.2
    OpenTSDB  :   OpenTSDB-2.0.1
    
    lzo压缩依赖包安装：
    
```python
yum -y install ant ant-nodeps lzo-devel
git clone git://github.com/cloudera/hadoop-lzo.git 
cd hadoop-lzo 
CLASSPATH=path/to/hadoop-core-1.0.4.jar CFLAGS=-m64 CXXFLAGS=-m64 ant compile-native tar 
hbasedir=path/to/hbase 
mkdir -p $hbasedir/lib/native 
cp build/hadoop-lzo-0.4.14/hadoop-lzo-0.4.14.jar $hbasedir/lib
cp -a build/hadoop-lzo-0.4.14/lib/native/* $hbasedir/lib/native
```

## HBase 安装
    OpenTSDB依赖于HBase作为底层存储，所以需要先安装HBase。
    hbase具体安装过程请见hbase集群部署文档。

## OpenTSDB安装
    我们选择源代码安装，需要做好心理准备，非常非常久。
```python
git clone git://github.com/OpenTSDB/opentsdb.git
cd opentsdb
./build.sh
```
    后来发现其实是因为OpenTSDB在build时候去maven中央仓库拉取第三方jar包，它不是通过标准的pom.xml指定依赖，所以即使你有本地setting.xml也没有生效。解决方案是修改maven地址：/Users/argan/tools/opentsdb/third_party /的所有include.mk文件中。
    
    Build成功之后可以将其make install到系统目录，不过这个是可选的。
    然后就要建表： 
```python
    env COMPRESSION=NONE HBASE_HOME=/home/hadoop/apache/hbase     
    ./src/create_table.sh
```
显示信息：
```python
HBase Shell; enter 'help<RETURN>' for list of supported commands.
Type "exit<RETURN>" to leave the HBase Shell
Version 1.0.1.1, re1dbf4df30d214fca14908df71d038081577ea46, Sun May 17 12:34:26 PDT 2015
create 'tsdb-uid',
  {NAME => 'id', COMPRESSION => 'NONE', BLOOMFILTER => 'ROW'},
  {NAME => 'name', COMPRESSION => 'NONE', BLOOMFILTER => 'ROW'}
 row(s) in 0.5970 seconds
Hbase::Table - tsdb-uid
create 'tsdb',
  {NAME => 't', VERSIONS => 1, COMPRESSION => 'NONE', BLOOMFILTER => 'ROW'}
 row(s) in 0.1520 seconds
Hbase::Table - tsdb
create 'tsdb-tree',
  {NAME => 't', VERSIONS => 1, COMPRESSION => 'NONE', BLOOMFILTER => 'ROW'}
 row(s) in 0.1610 seconds
Hbase::Table - tsdb-tree
create 'tsdb-meta',
  {NAME => 'name', COMPRESSION => 'NONE', BLOOMFILTER => 'ROW'}
 row(s) in 0.1540 seconds
Hbase::Table - tsdb-meta
```

    OK，现在我们可以准备启动TSD了。不过在这之前，我们需要配置一下OpenTSDB先 configuration ：
```python
> 将 src/opentsdb.conf 拷贝到如下目录之一： 
•	./opentsdb.conf 
•	/etc/opentsdb.conf 
•	/etc/opentsdb/opentsdb.conf 
•	/opt/opentsdb/opentsdb.conf 
```

---
>然后配置如下四个必须配置项： 
```python
•	tsd.network.port=4242 
•	tsd.http.cachedir=/tmp/tsd - Path to write temporary files to 
•	tsd.http.staticroot=build/staticroot - Path to the static GUI files found in ./build/staticroot 
•	tsd.storage.hbase.zk_quorum=localhost - A comma separated list of Zookeeper hosts to connect to, default is "localhost". If HBase and Zookeeper are not running on the same machine, specify the host and port here. 
•	tsd.core.auto_create_metrics=True - Whether or not to automatically create UIDs for new metric types, default is False. 建议打开。 
```
```python
mkdir –p /tmp/tsd 
```

然后就可以简单启动TSD了： 
```python
./build/tsdb tsd
```
当然，你也可以通过命令行指定（会覆盖opentsdb.conf的配置）：
```python
Usage: tsd --port=PORT --staticroot=PATH --cachedir=PATH
Starts the TSD, the Time Series Daemon
  --async-io=true|false Use async NIO (default true) or traditional blocking io
  --auto-metric         Automatically add metrics to tsdb as they are inserted.  Warning: this may cause unexpected metrics to be tracked
  --backlog=NUM         Size of connection attempt queue (default: 3072 or kernel somaxconn.
  --bind=ADDR           Address to bind to (default: 0.0.0.0).
  --cachedir=PATH       Directory under which to cache result of requests.
  --config=PATH         Path to a configuration file (default: Searches for file see docs).
  --flush-interval=MSEC Maximum time for which a new data point can be buffered (default: 1000).
  --port=NUM            TCP port to listen on.
  --staticroot=PATH     Web root from which to serve static files (/s URLs).
  --table=TABLE         Name of the HBase table where to store the time series (default: tsdb).
  --uidtable=TABLE      Name of the HBase table to use for Unique IDs (default: tsdb-uid).
  --worker-threads=NUM  Number for async io workers (default: cpu * 2).
  --zkbasedir=PATH      Path under which is the znode for the -ROOT- region (default: /hbase).
  --zkquorum=SPEC       Specification of the ZooKeeper quorum to use (default: localhost).
tsdtmp=${TMPDIR-'/tmp'}/tsd    # For best performance, make sure
mkdir -p "$tsdtmp"             # your temporary directory uses tmpfs

./build/tsdb tsd --port=4242 --staticroot=build/staticroot --cachedir="$tsdtmp" --zkquorum=myhost:2181
```

    然后就可以通过 http://127.0.0.1:4242 访问TSD的web界面了

通过http接口 写入数据
```python
[root@master opentsdb]# curl -i  -H "Content-Type: application/json" -X POST -d '{"metric": "sys.cpu.nice", "timestamp": 1433989867597,"value": 18, "tags": { "host": "web01"}}' http://10.10.3.179:4242/api/put/?details               
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 36

{"errors":[],"failed":0,"success":1}[root@master opentsdb]#
```
然后就可以通过 http://127.0.0.1:4242 访问TSD的web界面了
