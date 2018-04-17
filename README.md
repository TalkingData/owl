# OWL
[![Go Report Card](https://goreportcard.com/badge/github.com/TalkingData/owl)](https://goreportcard.com/report/github.com/TalkingData/owl)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)


distributed monitoring system

OWL是TalkingData公司推出的一款开源分布式监控系统

## Features

- Go语言开发，部署维护简单
- 分布式，支持多机房
- 多维的数据模型，类opentsdb
- 支持多种报警算法，支持多条件组合、时间范围、报警模板等
- 灵活的插件机制，支持任意语言编写，支持传参，自动同步
- 丰富的报警渠道，邮件、微信、短信、电话、自定义
- 原始数据永久存储，支持发送到opentsdb、kairosdb、kafka
- 自带web管理界面以及强大的自定义图表功能

## Architecture
![owl](./arch.png)


## Demo

http://54.223.127.87/

普通用户：demo/demo </br>
管理员：admin/111111 </br>
注: demo 环境数据库每隔 1 个小时会自动恢复

## rpm包地址
https://pan.baidu.com/s/1UTYOOB8YE8nng0guXOXkmg#list/path=%2Fowl

## 前端源码地址
https://github.com/TalkingData/owl-frontend

## QQ Group
492850035
