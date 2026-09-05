# ProjectA

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

基于 Go + Gin + GORM + Redis 的 Web 后端项目

## 技术栈
- Go 1.26
- Gin Web 框架
- GORM ORM
- MySQL 8.0
- Redis 7.0
- JWT 认证
- Zap 日志
- Docker + Docker Compose

## 快速启动

```bash
# 克隆项目
git clone https://github.com/lsnb6666/ProjectA.git

# 进入项目目录
cd ProjectA

# 用 Docker 启动所有服务
docker compose up -d
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 注册 |
| POST | /api/auth/login | 登录 |
| POST | /api/articles | 创建文章 |
| GET | /api/articles | 获取文章列表 |
| POST | /api/article/:id/like | 点赞 |

## 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。