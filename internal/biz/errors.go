package biz

import "errors"

// 评论相关错误
var (
	// ErrInvalidCommentContent 无效的评论内容
	ErrInvalidCommentContent = errors.New("评论内容不能为空")

	// ErrInvalidUserID 无效的用户ID
	ErrInvalidUserID = errors.New("用户ID无效")

	// ErrInvalidResourceID 无效的资源ID
	ErrInvalidResourceID = errors.New("资源ID无效")

	// ErrCreateCommentFailed 创建评论失败
	ErrCreateCommentFailed = errors.New("创建评论失败")
)
