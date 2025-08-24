package service

import (
	"comment/pkg/log"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// GetComment 实现获取评论接口
// ctx - 请求上下文
// in - 获取评论请求参数
// 返回 - 评论信息树和可能的错误
func (s *CommentService) GetComment(ctx context.Context, in *v1.GetCommentRequest) (*v1.CommentTree, error) {
	log.Info(ctx, "get comment")
	log.Debug(ctx, "GetComment", "module", in.Module, "resource_id", in.ResourceId, "max_depth", in.MaxDepth)

	// 设置默认值
	page := in.GetPage()
	if page <= 0 {
		page = 1
	}

	pageSize := in.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用业务层获取评论
	comments, err := s.uc.GetComments(ctx, in.Module, in.ResourceId, in.MaxDepth, page, pageSize, int32(in.GetSortType()))
	if err != nil {
		log.Error(ctx, "get comments failed.", "error", err)
		return nil, err
	}

	// 转换为API响应格式
	apiComments := make([]*v1.Comment, len(comments))
	for i, comment := range comments {
		apiComments[i] = s.convertToAPIComment(comment)
	}

	// 返回 API 响应
	log.Info(ctx, "get comment successful.")
	return &v1.CommentTree{
		Comments: apiComments,
	}, nil
}

// convertToAPIComment 将biz.Comment转换为v1.Comment
func (s *CommentService) convertToAPIComment(comment *biz.Comment) *v1.Comment {
	// 递归转换回复评论
	var replyComments []*v1.Comment
	if comment.ReplyComments != nil && len(comment.ReplyComments) > 0 {
		replyComments = make([]*v1.Comment, len(comment.ReplyComments))
		for i, reply := range comment.ReplyComments {
			replyComments[i] = s.convertToAPIComment(reply)
		}
	}

	return &v1.Comment{
		Module:        comment.Module,
		ResourceId:    comment.ResourceID,
		CommentId:     comment.ID,
		UserId:        comment.UserID,
		Username:      comment.Username,
		Avatar:        comment.Avatar,
		Content:       comment.Content,
		Level:         comment.Level,
		LikeCount:     comment.LikeCount,
		ReplyCount:    comment.ReplyCount,
		ReplyComments: replyComments,
		CreateTime:    timestamppb.New(comment.CreateGmt),
	}
}

// DeleteComment 实现删除评论接口
// ctx - 请求上下文
// in - 删除评论请求参数
// 返回 - 删除结果和可能的错误
func (s *CommentService) DeleteComment(ctx context.Context, in *v1.DeleteCommentRequest) (*v1.DeleteResponse, error) {
	log.Info(ctx, "delete comment")
	log.Debug(ctx, "DeleteComment", "module", in.Module, "resource_id", in.ResourceId, "comment_id", in.CommentId, "user_id", in.UserId)

	// 调用业务层删除评论
	err := s.uc.DeleteComment(ctx, in.CommentId)
	if err != nil {
		log.Error(ctx, "delete comment failed.", "error", err)
		return nil, err
	}

	// 返回 API 响应
	log.Info(ctx, "delete comment successful.")
	return &v1.DeleteResponse{
		Success: true,
	}, nil
}
