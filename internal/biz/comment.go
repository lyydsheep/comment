package biz

import (
	v1 "comment/api/comment/v1"
	"comment/pkg/log"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Comment is a Comment model.
type Comment struct {
	// Module 业务模块表示，用于区分不同业务场景下的评论
	Module int32 `gorm:"column:module;type:tinyint;not null;comment:0:视频,1:文章"`

	// ResourceID 资源唯一标识，表示被评论的资源ID
	ResourceID string `gorm:"column:resource_id;type:varchar(32);not null"`

	// ID 评论唯一标识
	ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`

	// RootCommentID 根评论ID，用于标识评论所属的根评论
	RootCommentID int64 `gorm:"column:root_id;type:varchar(32);not null;comment:根评论"`

	// ParentCommentID 父评论ID，用于构建评论回复关系
	ParentCommentID int64 `gorm:"column:parent_id;type:varchar(32);not null"`

	// UserID 用户唯一标识
	UserID string `gorm:"column:user_id;type:varchar(32);not null"`

	// Username 用户名
	Username string `gorm:"column:username;type:varchar(24);not null"`

	// Avatar 用户头像URL
	Avatar string `gorm:"column:avatar;type:varchar(255);not null;comment:头像 url"`

	// Content 评论内容
	Content string `gorm:"column:content;type:text;not null"`

	// Level 层级
	Level int32 `gorm:"column:level;type:int;not null;default:0"`

	// LikeCount 点赞数
	LikeCount int64 `gorm:"column:like_count;type:int;not null;default:0"`

	// ReplyCount 回复数
	ReplyCount int64 `gorm:"column:reply_count;type:int;not null;default:0"`

	// CreateGmt 创建时间
	CreateGmt time.Time `gorm:"column:create_gmt;type:datetime;not null;default:CURRENT_TIMESTAMP"`

	// UpdateGmt 更新时间
	UpdateGmt time.Time `gorm:"column:update_gmt;type:datetime;not null;default:CURRENT_TIMESTAMP;updateAt"`

	// ReplyComments 回复评论列表（内存中构建，不存储在数据库中）
	ReplyComments []*Comment `gorm:"-"`
}

func (c *Comment) TableName() string {
	return "comment"
}

// CommentRepo is a Comment repo.
type CommentRepo interface {
	// Save saves a Comment.
	Save(context.Context, *Comment) (*Comment, error)
	// Get gets a Comment by ID.
	Get(context.Context, int64) (*Comment, error)
	// Delete deletes a Comment by ID.
	Delete(context.Context, int64) error
	// DeleteBatch deletes Comments by root ID or ID.
	DeleteBatch(context.Context, int64) error
	// ListRootComments 获取根评论列表
	ListRootComments(ctx context.Context, module int32, resourceID string, page, pageSize int32, sortType int32) ([]*Comment, error)
	// ListReplyComments 获取回复评论列表
	ListReplyComments(ctx context.Context, rootIDs []int64, maxDepth int32, sortType int32) ([]*Comment, error)
}

// CommentUsecase is a Comment usecase.
type CommentUsecase struct {
	repo CommentRepo
}

// NewCommentUsecase new a Comment usecase.
func NewCommentUsecase(repo CommentRepo) *CommentUsecase {
	return &CommentUsecase{repo: repo}
}

// CreateComment creates a Comment, and returns the new Comment.
func (uc *CommentUsecase) CreateComment(ctx context.Context, c *Comment) (*v1.Comment, error) {
	log.Debug(ctx, "create comment.", "user_id", c.UserID, "content", c.Content)
	// 落库
	comment, err := uc.repo.Save(ctx, c)
	if err != nil {
		log.Error(ctx, "create comment error.", "err", err)
		return nil, errors.BadRequest(err.Error(), "create comment error.")
	}
	log.Info(ctx, "repo save successful.")

	// 返回参数
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
		ReplyComments: nil,
		CreateTime:    timestamppb.New(comment.CreateGmt),
	}, nil
}

// GetComments gets comments by module and resource id.
func (uc *CommentUsecase) GetComments(ctx context.Context, module int32, resourceID string, maxDepth, page, pageSize int32, sortType int32) ([]*Comment, error) {
	log.Debug(ctx, "get comments.", "module", module, "resource_id", resourceID, "max_depth", maxDepth, "page", page, "page_size", pageSize, "sort_type", sortType)

	// 获取根评论
	comments, err := uc.repo.ListRootComments(ctx, module, resourceID, page, pageSize, sortType)
	if err != nil {
		log.Error(ctx, "get root comments error.", "err", err)
		return nil, errors.BadRequest(err.Error(), "get root comments error.")
	}

	// 如果需要获取回复，则获取回复评论
	if maxDepth > 1 && len(comments) > 0 {
		// 收集所有根评论的ID
		rootIDs := make([]int64, len(comments))
		for i, comment := range comments {
			rootIDs[i] = comment.ID
		}

		// 获取所有回复评论
		replyComments, err := uc.repo.ListReplyComments(ctx, rootIDs, maxDepth-1, sortType)
		if err != nil {
			log.Error(ctx, "get reply comments error.", "err", err)
			return nil, errors.BadRequest(err.Error(), "get reply comments error.")
		}

		// 构建评论树
		uc.buildCommentTree(comments, replyComments, maxDepth)
	}

	log.Info(ctx, "repo get comments successful.")
	return comments, nil
}

// buildCommentTree 构建评论树
func (uc *CommentUsecase) buildCommentTree(rootComments []*Comment, replyComments []*Comment, maxDepth int32) {
	// 创建一个map用于快速查找评论
	commentMap := make(map[int64]*Comment)
	for _, comment := range rootComments {
		commentMap[comment.ID] = comment
	}

	// 将回复评论也加入map
	for _, comment := range replyComments {
		commentMap[comment.ID] = comment
	}

	// 将回复评论挂到对应的父评论下
	for _, reply := range replyComments {
		// 查找父评论
		if parent, exists := commentMap[reply.ParentCommentID]; exists {
			// 计算当前评论的层级
			currentDepth := reply.Level - parent.Level

			// 只在允许的深度范围内添加
			if currentDepth <= maxDepth {
				parent.ReplyComments = append(parent.ReplyComments, reply)
			}
		}
	}
}

// DeleteComment deletes a Comment by ID.
func (uc *CommentUsecase) DeleteComment(ctx context.Context, id int64) error {
	log.Debug(ctx, "delete comment.", "id", id)

	// 首先获取要删除的评论
	comment, err := uc.repo.Get(ctx, id)
	if err != nil {
		log.Error(ctx, "get comment error.", "err", err)
		return errors.BadRequest(err.Error(), "get comment error.")
	}

	// 删除根评论ID为该评论ID的所有回复评论，或者ID为该评论ID的评论
	// 这将删除该评论及其所有回复和回复的回复
	err = uc.deleteCommentAndReplies(ctx, comment)
	if err != nil {
		log.Error(ctx, "delete comment and replies error.", "err", err)
		return errors.BadRequest(err.Error(), "delete comment and replies error.")
	}

	log.Info(ctx, "repo delete successful.")
	return nil
}

// deleteCommentAndReplies 删除评论及其所有回复
func (uc *CommentUsecase) deleteCommentAndReplies(ctx context.Context, comment *Comment) error {
	// 删除所有根评论ID为当前评论ID的评论（即该评论的所有直接回复和间接回复）
	// 或者ID为当前评论ID的评论（即当前评论本身）
	// 这样可以一次性删除整个评论树
	err := uc.repo.DeleteBatch(ctx, comment.ID)
	return err
}
