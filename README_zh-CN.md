# EM-Auth

[English](./README.md) | 简体中文

## 介绍
本项目基于Etpmls-Micro开发。

本项目是一个总控制中心服务，集成用户、角色、权限的RBAC0的鉴权、自定义菜单、清除缓存、磁盘清理等功能。

## 运行
MySQL/MariaDB
```go
go run -tags=mysql main.go
```
PostgreSQL
```go
go run -tags=postgresql main.go
```