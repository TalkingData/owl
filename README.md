# OWL
[![Go Report Card](https://goreportcard.com/badge/github.com/TalkingData/owl)](https://goreportcard.com/report/github.com/TalkingData/owl)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)


​&nbsp;​&nbsp;​&nbsp;​&nbsp;​&nbsp;​&nbsp;OWL 是由国内领先的第三方数据智能服务商 [TalkingData](<https://www.talkingdata.com/>) 开源的一款企业级分布式监控告警系统，目前由 Tech Operation Team 持续开发更新维护。

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;OWL 后台组件全部使用 [Go](https://golang.org/) 语言开发，Go 语言是 Google 开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言，它的并发机制可以充分利用多核，同平台一次编译可以到处运行，运维成本极低，更多的信息可以参考[官方文档](https://golang.org/doc/)。前端页面使用 [iView](<https://github.com/iview/iview>) 开发，iView 同样是由 TalkingData 开源的一套基于 Vue.js 的 UI 组件库，主要服务于 PC 界面的中后台产品。


## Features

- Go语言开发，部署维护简单
- 基于Go-micro框架，实现负载均衡并支持分布式部署方式，支持多机房
- 支持多种报警算法，支持多条件组合、时间范围、报警模板等
- 灵活的插件机制，支持任意语言编写，支持传参，自动同步到客户端
- 丰富的报警渠道，邮件、企业微信、短信、电话以及自定义脚本
- 原始数据永久存储，支持发送到 kairosdb
- 自带 web 管理界面以及强大的自定义图表功能能

## 什么是v6：
- OWL v6是在v5版本上进行一定的组件优化、策略优化、协议优化后的兼容修订版本
- 目前OWL v6只有Cfc、Repeater和Agent组件开发完成，其余组件尚处于开发状态

## OWL V6 Roadmap && ChangeLogs
- [x] Repeater通讯方式改用gRPC，使用context管理协程和服务；  
- [x] CFC通讯方式改用gRPC，使用context管理协程和服务；  
- [x] CFC操作数据库的方法被封装入dao层，同时加入gorm，解决v5版本中存在的数据库注入问题；
- [x] Agent改用gRPC Client与CFC和Repeater通讯，使用context管理协程和服务；
- [x] Agent内部库shirou/gopsutil更换为v3.22.6；
  - [x] 修复对system.cpu.softirq采集逻辑的错误；
  - [x] 不再支持采集内部指标：system.cpu.stolen；
- [x] 移除Netcollect组件（计划后续采用agent的http上报方式替代）；
- [x] 使Proxy可同时具有CFC和Repeater的代理功能；
- [x] Agent支持使用http协议主动上报ts data方式；
- [x] Agent改为只与Proxy通讯，Proxy利用服务注册机制寻找CFC和Repeater并做相应负载均衡策略；
- [ ] Inspector与Controller组件业务逻辑优化（计划整合、支持分布式方式部署和运行）；
- [ ] Api操作数据库的方法被封装入dao层，同时加入gorm，解决v5版本中存在的数据库注入问题；
- [ ] 【Draft】增加对不同OS和Arch的支持；
  - [ ] host数据库表扩充OS和Arch字段，标识Agent所在主机的操作系统和架构；
  - [ ] plugin库表扩充其支持的OS和Arch字段；
  - [ ] 处理Agent发来的插件下载请求时，CFC会根据其OS和Arch做出相应下载返回；
  - [ ] Proxy组件更新并适配此功能；
