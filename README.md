# 跳板Springboard 
## 后端
### 文件结构
```sh
├── Dockerfile
├── README.md
├── api
│   ├── feedback
│   │   └── route.go
│   ├── init.go
│   ├── oss
│   │   └── route.go
│   └── portfolio
│       └── route.go
├── cmd
│   └── main.go
├── configs
│   └── config.yaml
├── consts
│   └── default.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── conf
│   │   └── conf.go
│   ├── controller
│   │   ├── auth.go
│   │   ├── feedback.go
│   │   ├── oss.go
│   │   ├── portfolio.go
│   │   └── response.go
│   ├── data
│   │   ├── auth.go
│   │   ├── data.go
│   │   ├── feedback.go
│   │   └── portfolio.go
│   └── middleware
│       └── auth.go
└── pkgs
   ├── logger
   │   └── logger.go
   └── oss
       └── oss.go

   
```
