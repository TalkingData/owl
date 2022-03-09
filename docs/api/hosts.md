## Hosts

#### 1、获取主机列表

```
GET /api/v1/hosts
```

__Example request:__

```
GET  /api/v1/hosts?page=1&pageSize=15&status=1&group_id=2&q=10.10.32.10  HTTP/1.1
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "id":"e736e4bc76a95088",
        "name":"app01",
        "sn":"F0CHDC2",
        "ip":"127.0.0.1",
        "hostname":"mysql-01",
        "agent_version":"1.0",
        "status":0,
        "create_at":"2016-07-26 22:26:54",
        "update_at":"2016-07-26 22:26:54",
        "groups":[
            {
                "id":1,
                "name":"appcpa"
            },{
                "id":2,
                "name":"game"
            }
        ],
        "metric_cnt":5,
        "strategy":[1,2,3]
    },{
        "id":"abe91c6115a4abf0",
        "name":"",
        "sn":"F1CBDA2",
        "ip":"127.0.0.1",
        "hostname":"mysql-02",
        "agent_version":"1.0",
        "status":0,
        "create_at":"2016-07-26 22:26:54",
        "update_at":"2016-07-26 22:26:54",
        "groups":[
            {
                "id":1,
                "name":"appcpa"
            },{
                "id":2,
                "name":"game"
            }
        ],
        "metric_cnt":10,
        "strategy":[1,2,3]
    },
]
```

__Query parameters:__

+ `page`: 当前页码
+ `pageSize`: 每页显示条数
+ `status`: 状态字段,有效值 0 正常, 1 故障, 2 禁用, 3 新增
+ `group_id`:主机组ID,根据主机组过滤主机
+ `q`:查询字符串,支持ip地址、主机名、自定义名称

<br/>

#### 2、禁用主机

```
POST /api/v1/hosts/host_id/disable  HTTP/1.1
```

__Example request:__

```
POST  /api/v1/hosts/839787f0772e1a27/disable  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "message": "disable the host success"
}
```

#### 3、启用主机

```
POST /api/v1/hosts/host_id/enable  HTTP/1.1
```

__Example request:__

```
POST  /api/v1/hosts/839787f0772e1a27/enable  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "message": "enable the host success"
}
```

#### 4、重命名主机

```
POST /api/v1/hosts/host_id/rename  HTTP/1.1
```

__Example request:__

```
POST  /api/v1/hosts/839787f0772e1a27/enable  HTTP/1.1
Content-Type: application/json

{
    "name":"host-01"
}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "message": "rename the host success"
}
```

#### 4、删除主机

```
DELETE /api/v1/hosts/host_id  HTTP/1.1
```

__Example request:__

```
DELETE  /api/v1/hosts/839787f0772e1a27  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "message":"delete success"
}
```

#### 5、获取统计信息

```
GET /api/v1/hosts/info  HTTP/1.1
```

> 包含主机数、主机组数、metric数

__Example request:__

```
GET  /api/v1/hosts/info  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "groups": 3,  //当有group_id参数时不返回该项
  "hosts": 16,
  "metrics": 1005
}
```

__Query parameters:__

+ `group_id`: 获取主机组下的统计信息

#### 6、获取主机状态统计信息

```
GET /api/v1/hosts/status  HTTP/1.1
```

> 包含正常、故障、禁用

__Example request:__

```
GET  /api/v1/hosts/status?group_id=2  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "disable": 1,
  "failed": 0,
  "normal": 15,
  "pending": 0
}
```

__Query parameters:__

+ `group_id`: 获取主机组下的主机状态信息

#### 6、获取主机下策略计数

```
GET /api/v1/hosts/strategy/host_id  HTTP/1.1
```

> 包含全局策略数量、主机组策略数量、主机策略数量

__Example request:__

```
GET  /api/v1/hosts/strategy/839787f0772e1a27  HTTP/1.1
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "response": {
    "global_strategy": [
      {
        "id": 2,
        "name": "global host status"
      }
    ],
    "group_strategy": [],
    "host_strategy": []
  }
}
```

__Query parameters:__

+ `group_id`: 获取主机组下的主机状态信息

#### 7、获取主机下metric列表

```
GET /api/v1/hosts/metric/host_id  HTTP/1.1
```

__Example request:__

```
GET  /api/v1/hosts/metric/839787f0772e1a27  HTTP/1.1
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "metrics": [
    {
      "id": 1,
      "name": "agent.alive",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:32:48+08:00",
      "update_at": "2016-09-02T17:32:48+08:00"
    },
    {
      "id": 2,
      "name": "mem.total",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 3,
      "name": "mem.usedprecent",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 4,
      "name": "mem.active",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 5,
      "name": "mem.buffers",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 6,
      "name": "mem.free",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 7,
      "name": "mem.used",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 8,
      "name": "swap.total",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 9,
      "name": "swap.usedprecent",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    },
    {
      "id": 10,
      "name": "swap.free",
      "dt": "GAUGE",
      "cycle": 30,
      "creat_at": "2016-09-02T17:33:14+08:00",
      "update_at": "2016-09-02T17:33:14+08:00"
    }
  ]
}
```

__Query parameters:__

+ `group_id`: 获取主机组下的主机状态信息
+ `page`: 当前页,默认为1
+ `pageSize`: 每页的记录数,默认为10
