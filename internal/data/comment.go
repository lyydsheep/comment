package data

import (
	"comment/internal/biz"
	"context"
)

type commentRepo struct {
	data *Data
}

// NewCommentRepo .
func NewCommentRepo(data *Data) biz.CommentRepo {
	return &commentRepo{
		data: data,
	}
}

func (r *commentRepo) Save(ctx context.Context, c *biz.Comment) (*biz.Comment, error) {
	err := r.data.db.WithContext(ctx).Model(&biz.Comment{}).Create(c).Error
	return c, err
}
