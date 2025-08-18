package service

import (
	"comment/pkg/log"
	"context"
	"time"

	v1 "comment/api/comment/v1"
	"comment/internal/biz"
)

// CommentService is a comment service.
type CommentService struct {
	v1.UnimplementedCommentServiceServer

	uc *biz.CommentUsecase
}

// NewCommentService new a comment service.
func NewCommentService(uc *biz.CommentUsecase) *CommentService {
	return &CommentService{uc: uc}
}

// CreateComment 实现评论创建接口
// ctx - 请求上下文
// in - 创建评论请求参数
// 返回 - 创建的评论信息和可能的错误
func (s *CommentService) CreateComment(ctx context.Context, in *v1.CreateCommentRequest) (*v1.Comment, error) {
	log.Info(ctx, "create comment")
	log.Debug(ctx, "CreateComment", "user_id", in.UserId, "content", in.Content)
	// 转换请求参数为业务模型
	comment := &biz.Comment{
		Module:          in.Module,
		ResourceID:      in.ResourceId,
		RootCommentID:   in.RootCommentId,
		ParentCommentID: in.ParentCommentId,
		UserID:          in.UserId,
		Username:        in.Username,
		Avatar:          in.Avatar,
		Content:         in.Content,
		Level:           in.Level,
		LikeCount:       0,
		ReplyCount:      0,
		CreateGmt:       time.Now().UTC(),
		UpdateGmt:       time.Now().UTC(),
	}

	// 调用业务层创建评论
	createdComment, err := s.uc.CreateComment(ctx, comment)
	if err != nil {
		log.Error(ctx, "create comment failed.", "error", err)
		return nil, err
	}

	// 返回 API 响应
	log.Info(ctx, "create comment successful.")
	return createdComment, nil
}
