# Comment
# 评论服务系统

## 项目简介
这是一个基于 Go 语言开发的评论服务系统，提供完整的评论功能，包括发表评论、查询评论列表、删除评论以及点赞评论等功能。

## 技术栈
- **开发语言**: Go 1.25.0
- **后端框架**: [Kratos v2](https://github.com/go-kratos/kratos)
- **数据库**: MySQL (GORM 1.30.1)
- **缓存**: Redis
- **API框架**: gRPC + RESTful API
- **依赖注入**: Wire
- **配置管理**: YAML
- **日志**: tint

## 核心功能

### 1. 发表评论
- 支持发送文字和表情
- 字数限制：1-2000字
- 支持多级评论回复
- 支持不同业务模块（如文章、视频等）

### 2. 评论列表查询
- 支持按点赞数或创建时间降序排序
- 支持分页查询
- 支持层级展示评论（默认展示3条回复，可通过参数控制）

### 3. 删除评论
- 支持删除指定评论
- 支持批量删除关联回复

### 4. 评论互动
- 支持点赞和取消点赞评论
- 实时更新点赞数量

## 项目结构

```
├── api/                    # API 定义目录
│   └── comment/v1/         # 评论服务 v1 版本 API
├── cmd/                    # 命令行入口
│   └── comment/            # 评论服务主程序
│       ├── main.go         # 程序入口
│       ├── wire.go         # 依赖注入配置
│       └── wire_gen.go     # 自动生成的依赖注入代码
├── configs/                # 配置文件目录
│   └── config.yaml         # 主配置文件
├── internal/               # 内部代码，不对外暴露
│   ├── biz/                # 业务逻辑层
│   ├── data/               # 数据访问层
│   ├── middleware/         # 中间件
│   ├── server/             # 服务器实现
│   └── service/            # 服务实现层
├── pkg/                    # 公共包
│   └── log/                # 日志工具
└── third_party/            # 第三方依赖
```

## 快速开始

### 环境要求
- Go 1.25.0 或更高版本
- MySQL 数据库
- Redis

### 安装依赖
```bash
go mod tidy
```

### 数据库准备
1. 创建数据库：
```sql
CREATE DATABASE comment;
```

2. 导入数据表：
```bash
mysql -u root -p comment < init_comments.sql
```

### 配置修改
修改 `configs/config.yaml` 文件中的数据库和Redis连接信息：
```yaml
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:33060)/comment?parseTime=True&loc=UTC
  redis:
    addr: 127.0.0.1:6379
```

### 运行项目
```bash
# 生成依赖注入代码
go generate ./cmd/comment

# 启动服务
make run
```

## API 接口文档
[在线接口文档](https://s.apifox.cn/9b22df33-b9c4-4562-bfba-f1304632aba2)

## 数据库设计

### 评论表 (comment)
```sql
create table comment
(
  id          bigint auto_increment
        primary key,
  module      tinyint                            not null comment '0：视频，1：文章',
  resource_id varchar(32)                        not null,
  root_id     varchar(32)                        not null comment '根评论',
  parent_id   varchar(32)                        not null,
  level       int      default 0                 not null,
  user_id     varchar(32)                        not null,
  username    varchar(24)                        not null,
  avatar      varchar(255)                       not null comment '头像 url',
  content     text                               not null,
  like_num    int      default 0                 not null,
  reply_count int      default 0                 not null,
  create_gmt  datetime default CURRENT_TIMESTAMP not null,
  update_gmt  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);
```

## 配置说明

### 服务配置
```yaml
http:
  addr: 0.0.0.0:8000        # HTTP 服务监听地址
  timeout: 1s               # HTTP 请求超时时间
grpc:
  addr: 0.0.0.0:9000        # gRPC 服务监听地址
  timeout: 1s               # gRPC 请求超时时间
```

### 数据库配置
```yaml
data:
  database:
    driver: mysql                         # 数据库驱动
    source: root:root@tcp(127.0.0.1:33060)/comment?parseTime=True&loc=UTC  # 数据库连接字符串
    ConnMaxLifeTime: 300s                 # 连接最大生命周期
    ConnMaxIdleTime: 120s                 # 连接空闲超时时间
    IdleConns: 10                         # 最小空闲连接数
    MaxOpenConns: 50                      # 最大打开连接数
```

### Redis 配置
```yaml
data:
  redis:
    addr: 127.0.0.1:6379      # Redis 服务器地址
    read_timeout: 0.2s        # 读取超时时间
    write_timeout: 0.2s       # 写入超时时间
```

## 核心 API

### CommentService 服务

#### 创建评论
```protobuf
rpc CreateComment (CreateCommentRequest) returns (Comment)
```

#### 获取评论列表
```protobuf
rpc GetComment (GetCommentRequest) returns (CommentTree)
```

#### 删除评论
```protobuf
rpc DeleteComment (DeleteCommentRequest) returns (DeleteResponse)
```

#### 点赞评论
```protobuf
rpc LikeComment (LikeCommentRequest) returns (LikeResponse)
```

#### 取消点赞评论
```protobuf
rpc UnlikeComment (UnlikeCommentRequest) returns (UnlikeResponse)
```

## 开发指南

### 目录说明
- `api/`: 存放 API 定义文件（protobuf）
- `internal/biz/`: 存放业务实体和业务逻辑
- `internal/data/`: 存放数据访问层代码
- `internal/service/`: 存放服务实现代码
- `cmd/`: 存放程序入口代码

### 开发流程
1. 定义或修改 protobuf 文件
2. 生成对应的 Go 代码
3. 实现业务逻辑（biz 层）
4. 实现数据访问（data 层）
5. 实现服务（service 层）
6. 配置依赖注入（wire.go）
7. 运行和测试

## License
[MIT](LICENSE)