## Users


#### 1、获取用户列表
```
GET /api/v1/users
```

__Example request:__
```
GET  /api/v1/users HTTP/1.1
```

__Example response:__

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "users": [
    {
      "id": 1,
      "username": "chao.ma",
      "role": 0,
      "phone": "18611744661",
      "mail": "chao.ma@tendcloud.com",
      "weixin": "18611744661",
      "status": 0,
      "groups": [
        {
          "id": 2,
          "name": "appcpa"
        }
      ]
    },
    {
      "id": 48,
      "username": "yingsong.wu",
      "role": 1,
      "phone": "15210190629",
      "mail": "yingsong.wu@tenddata.com",
      "weixin": "15210190629",
      "status": 1,
      "groups": [
        {
          "id": 2,
          "name": "appcpa"
        }
      ]
    },
    {
      "id": 50,
      "username": "hierarch.pan",
      "role": 1,
      "phone": "13810520844",
      "mail": "hierarch.pan@tendcloud.com",
      "weixin": "13810520844",
      "status": 1,
      "groups": []
    }
  ]
}
```

__Query parameters:__

+ `q`:查询字符串,支持用户名名,手机号,微信号,邮箱地址过滤
+ `group_id`: 查询指定组下的用户, 不指定则为所有用户
+ `status`: 根据状态查询, 不指定表示所有状态
+  `hasNull`: 非零值表示过滤出手机号,邮箱地址,微信有空值的用户,默认不过滤

<br/>


#### 2、创建用户
```
PUT /api/v1/users HTTP/1.1
```

__Example request:__
```
PUT  /api/v1/users  HTTP/1.1
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
