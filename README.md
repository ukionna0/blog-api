# 简易个人博客 API

https://img.shields.io/badge/RESTful%20API-blog--api-red

## 项目简介
一个基于 Golang 的 RESTful Blog-api，技术栈：Golang + Gin + GORM + JWT。
实现功能：登录、注册、发表文章和评论以及对应的 CRUD。

## 快速开始
### 环境要求
* Go 1.24.5
* MySQL 9.4.0

### 安装与运行
1.  克隆项目：
    ```bash
    git clone https://github.com/ukionna0/blog-api.git
    cd blog-api
    ```
2.  安装依赖：
    ```bash
    go mod download
    ```
3.  启动项目：
    ```bash
    go run main.go
    ```
