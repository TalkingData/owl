## Groups


#### 1、获取主机组列表
```
GET /api/v1/groups
```

__Example request:__
```
GET  /api/v1/groups?q=data HTTP/1.1
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "groups": [
    {
      "id": 3,
      "name": "data-cloud"
    }
  ]
}
```

__Query parameters:__

+ `q`:查询字符串,支持主机组名过滤

<br/>


#### 2、创建主机组
```
PUT /api/v1/groups HTTP/1.1
```

__Example request:__
```
PUT  /api/v1/hosts/839787f0772e1a27/disable  HTTP/1.1
Content-Type: application/json

{
  "name": "test_group",
  "hosts":[
    {
        "id": "03381a1a3ea35d5a",
        "name": "",
        "ip": "172.28.5.36",
        "sn": "ba111677-a504-4481-962e-40bcc972b800",
        "hostname": "owl-agent-12.novalocal",
        "agent_version": "0.1",
        "status": "1"
    }
  ]
}
```

__Example response:__
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "group": {
    "id": 8,
    "name": "test_group",
    "hosts": [
      {
        "id": "03381a1a3ea35d5a",
        "name": "",
        "ip": "172.28.5.36",
        "sn": "ba111677-a504-4481-962e-40bcc972b800",
        "hostname": "owl-agent-12.novalocal",
        "agent_version": "0.1",
        "status": "1"
      }
    ]
  },
  "message": "group create successful"
}
```

#### 3、更新主机组
```
POST /api/v1/groups  HTTP/1.1
```
__Example request:__
```
POST  /api/v1/groups  HTTP/1.1
Content-Type: application/json

{
  "id": 8,
  "name": "test_group_new",
  "hosts":[
    {
        "id": "03381a1a3ea35d5a",
        "name": "",
        "ip": "172.28.5.36",
        "sn": "ba111677-a504-4481-962e-40bcc972b800",
        "hostname": "owl-agent-12.novalocal",
        "agent_version": "0.1",
        "status": "1"
    }
  ]
}
```

__Example response:__
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "group": {
    "id": 8,
    "name": "test_group_new",
    "hosts": [
      {
        "id": "03381a1a3ea35d5a",
        "name": "",
        "ip": "",
        "sn": "",
        "hostname": "",
        "agent_version": "",
        "status": ""
      }
    ]
  },
  "message": "group update successful"
}
```

#### 4、删除主机组
```
DELETE /api/v1/groups/:id  HTTP/1.1
```
__Example request:__
```
DELETE  /api/v1/groups/8  HTTP/1.1
Content-Type: application/json

{}
```

__Example response:__
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "message": "test_group_new contains some host, is not allowed to delete"
}
```

备注:当主机组包含有主机时,不允许删除
