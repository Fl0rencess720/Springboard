# Springboard（跳板） 后端
## 技术栈（暂时）：*Gin*+*GORM*+*JWT*+*MySQL*+*redis*+*OSS*
## 快速部署：
```bash
docker compose up -d --build
```
## CI/CD 
* 要执行流水线，需要为新版本打上tag，例如：`git tag v1.0.0`，然后执行`git push origin v1.0.0`，然后会自动部署到服务器
* 服务器使用traefik作为反向代理，提供https访问能力