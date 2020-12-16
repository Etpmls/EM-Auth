# EM-Auth

Englsih | [简体中文](./README_zh-CN.md)

## Introduction
This project is developed based on Etpmls-Micro.

This project is a general control center service that integrates RBAC0 authentication of users, roles, permissions, custom menus, clear cache, disk cleanup and other functions.

![Process](docs/images/Process.jpg)

## Run

MySQL/MariaDB
```go
go run -tags=mysql main.go
```
PostgreSQL
```go
go run -tags=postgresql main.go
```