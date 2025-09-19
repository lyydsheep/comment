-- 初始化评论数据的SQL语句
-- 初始化评论数据的SQL语句
-- 包含视频和文章两种模块，以及根评论和回复评论

-- 清空评论表（可选，仅在需要重新初始化时使用）
-- TRUNCATE TABLE comment;

-- 插入视频模块的根评论 (module=2)
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (2, 'video123', 0, 0, 'user1', '张三', 'https://example.com/avatars/user1.png', '这个视频太精彩了，学到了很多！', 0, 45, 2, NOW(), NOW()),
    (2, 'video123', 0, 0, 'user2', '李四', 'https://example.com/avatars/user2.png', '讲解很清晰，推荐给大家', 0, 28, 1, DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_SUB(NOW(), INTERVAL 1 HOUR)),
    (2, 'video123', 0, 0, 'user3', '王五', 'https://example.com/avatars/user3.png', '希望能出更多类似的内容', 0, 15, 0, DATE_SUB(NOW(), INTERVAL 2 HOUR), DATE_SUB(NOW(), INTERVAL 2 HOUR)),
    (2, 'video456', 0, 0, 'user4', '赵六', 'https://example.com/avatars/user4.png', '这个教程对新手很友好', 0, 32, 0, DATE_SUB(NOW(), INTERVAL 3 HOUR), DATE_SUB(NOW(), INTERVAL 3 HOUR));

-- 插入文章模块的根评论
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (1, 'article789', 0, 0, 'user5', '孙七', 'https://example.com/avatars/user5.png', '文章分析得很到位，有深度', 0, 67, 2, DATE_SUB(NOW(), INTERVAL 4 HOUR), DATE_SUB(NOW(), INTERVAL 4 HOUR)),
    (1, 'article789', 0, 0, 'user6', '周八', 'https://example.com/avatars/user6.png', '学到了新知识，感谢分享', 0, 41, 0, DATE_SUB(NOW(), INTERVAL 5 HOUR), DATE_SUB(NOW(), INTERVAL 5 HOUR)),
    (1, 'article101', 0, 0, 'user7', '吴九', 'https://example.com/avatars/user7.png', '观点很独特，值得思考', 0, 29, 0, DATE_SUB(NOW(), INTERVAL 6 HOUR), DATE_SUB(NOW(), INTERVAL 6 HOUR));

-- 假设上面的插入语句执行后，视频模块的根评论ID为1,2,3,4，文章模块的根评论ID为5,6,7
-- 插入回复评论
-- 注意：实际使用时需要根据真实的ID进行调整

-- 回复视频评论1的评论 (module=2)
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (2, 'video123', 1, 1, 'user8', '郑十', 'https://example.com/avatars/user8.png', '完全同意，我也是受益匪浅', 1, 12, 0, DATE_SUB(NOW(), INTERVAL 30 MINUTE), DATE_SUB(NOW(), INTERVAL 30 MINUTE)),
    (2, 'video123', 1, 1, 'user9', '钱十一', 'https://example.com/avatars/user9.png', '请问哪里可以找到相关资料？', 1, 5, 1, DATE_SUB(NOW(), INTERVAL 20 MINUTE), DATE_SUB(NOW(), INTERVAL 20 MINUTE));

-- 回复视频评论2的评论 (module=2)
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (2, 'video123', 2, 2, 'user10', '孙十二', 'https://example.com/avatars/user10.png', '确实不错，已收藏', 1, 8, 0, DATE_SUB(NOW(), INTERVAL 90 MINUTE), DATE_SUB(NOW(), INTERVAL 90 MINUTE));

-- 回复文章评论5的评论
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (1, 'article789', 5, 5, 'user11', '李十三', 'https://example.com/avatars/user11.png', '补充一点，还有另外一个角度可以考虑', 1, 15, 2, DATE_SUB(NOW(), INTERVAL 210 MINUTE), DATE_SUB(NOW(), INTERVAL 210 MINUTE)),
    (1, 'article789', 5, 5, 'user12', '张十四', 'https://example.com/avatars/user12.png', '期待后续更新', 1, 7, 0, DATE_SUB(NOW(), INTERVAL 225 MINUTE), DATE_SUB(NOW(), INTERVAL 225 MINUTE));

-- 嵌套回复（回复评论的评论）
-- 回复视频评论1的回复（假设ID为8） (module=2)
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (2, 'video123', 1, 8, 'user1', '张三', 'https://example.com/avatars/user1.png', '可以去官方文档看看，有详细说明', 2, 3, 0, DATE_SUB(NOW(), INTERVAL 15 MINUTE), DATE_SUB(NOW(), INTERVAL 15 MINUTE));

-- 回复文章评论5的回复（假设ID为11）
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (1, 'article789', 5, 11, 'user5', '孙七', 'https://example.com/avatars/user5.png', '感谢补充，下次我会考虑这个角度', 2, 6, 0, DATE_SUB(NOW(), INTERVAL 3 HOUR), DATE_SUB(NOW(), INTERVAL 3 HOUR)),
    (1, 'article789', 5, 11, 'user13', '王十五', 'https://example.com/avatars/user13.png', '这个补充很有价值', 2, 4, 0, DATE_SUB(NOW(), INTERVAL 150 MINUTE), DATE_SUB(NOW(), INTERVAL 150 MINUTE));

-- 生成更多不同场景的评论数据
-- 视频模块的热门评论（高点赞数） (module=2)
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (2, 'video123', 0, 0, 'user14', '刘十六', 'https://example.com/avatars/user14.png', '干货满满，强烈推荐！', 0, 128, 0, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY));

-- 文章模块的热门评论（高点赞数）
INSERT INTO comment (module, resource_id, root_id, parent_id, user_id, username, avatar, content, level, like_count, reply_count, create_gmt, update_gmt)
VALUES
    (1, 'article789', 0, 0, 'user15', '陈十七', 'https://example.com/avatars/user15.png', '写得太好了，转发给朋友了', 0, 215, 0, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY));