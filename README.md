# EM-Auth

Englsih | [简体中文](./README_zh-CN.md)

## Introduction
This project is developed based on Etpmls-Micro.

This project is a general control center service that integrates RBAC0 authentication of users, roles, permissions, custom menus, clear cache, disk cleanup and other functions.

![Process](docs/images/Process.jpg)

## Configuration

### EM configuration

Refer to Etpmls-Micro manual [EM Configuration](https://github.com/Etpmls/Etpmls-Micro#em-configuration)

### Gateway configuration

Configure in the configuration file service-discovery=>service=>http=>tag

```yaml
      tag: [
        "em.http.routers.em-AuthHttpService.entrypoints=web,websecure",
        "em.http.routers.em-AuthHttpService.rule=Host(`[YOUR_DOMAIN]`) && PathPrefix(`/api/auth/`)",
        "em.http.routers.em-AuthHttpService.tls.certresolver=myresolver",
        "em.http.routers.em-AuthHttpService.middlewares=forwardAuth@file,circuitBreaker_em-auth@file",
        "em.http.routers.em-AuthHttpService.service=em-AuthHttpService",

        "em.http.routers.em-AuthHttpService-checkAuth.entrypoints=web,websecure",
        "em.http.routers.em-AuthHttpService-checkAuth.rule=Host(`[YOUR_DOMAIN]`,`[YOUR_TRAEFIK_ADDRESS]`) && Path(`/api/checkAuth`)",
        "em.http.routers.em-AuthHttpService-checkAuth.tls.certresolver=myresolver",
        "em.http.routers.em-AuthHttpService-checkAuth.middlewares=circuitBreaker_em-auth@file",
        "em.http.routers.em-AuthHttpService-checkAuth.service=em-AuthHttpService",

        "em.http.services.em-AuthHttpService.loadbalancer.passhostheader=true",
      ]
```

Replace [YOUR_DOMAIN] with your domain name, and [YOUR_TRAEFIK_ADDRESS] with Traefik address

## Run

MySQL/MariaDB
```go
go run -tags=mysql main.go
```
PostgreSQL
```go
go run -tags=postgresql main.go
```