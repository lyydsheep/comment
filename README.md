# Comment

## PRD
- 发表评论
  - 可以发送文字和表情
  - 字数 [1, 200] 字
  - 敏感词过滤（黑、灰、白）
  - 通知消息和提醒
- 评论列表
  - 排序：时间、热度（点赞数）
  - 展示评论总数
  - 评论下的回复默认展示 3 条，可以点击“查看更多”展示更多回复
- 删除评论或回复

## API 接口文档
[在线接口文档](https://s.apifox.cn/9b22df33-b9c4-4562-bfba-f1304632aba2)

## DDL
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