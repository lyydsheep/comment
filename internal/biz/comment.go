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
}

func (c *Comment) TableName() string {
	return "comment"
}

// CommentRepo is a Comment repo.
type CommentRepo interface {
	// Save saves a Comment.
	Save(context.Context, *Comment) (*Comment, error)
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

func (uc *CommentUsecase) constructTree(root *Comment) (*v1.Comment, error) {
	// TODO
	return nil, nil
}
