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
	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 创建评论
	if err := tx.Model(&biz.Comment{}).Create(c).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果是回复评论，更新父评论的回复数
	if c.ParentCommentID > 0 {
		if err := tx.Model(&biz.Comment{}).Where("id = ?", c.ParentCommentID).
			UpdateColumn("reply_count", tx.Model(&biz.Comment{}).Select("reply_count + ?", 1).Where("id = ?", c.ParentCommentID)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return c, nil
}

func (r *commentRepo) Get(ctx context.Context, id int64) (*biz.Comment, error) {
	var comment biz.Comment
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) Delete(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Where("id = ?", id).Delete(&biz.Comment{}).Error
}

// DeleteBatch 删除指定评论ID的所有相关评论（包括该评论本身及其所有回复）
func (r *commentRepo) DeleteBatch(ctx context.Context, id int64) error {
	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 找出所有要删除的评论ID
	var commentIDs []int64
	if err := tx.Model(&biz.Comment{}).Where("id = ? OR root_id = ?", id, id).Pluck("id", &commentIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果有要删除的评论
	if len(commentIDs) > 0 {
		// 删除所有相关的点赞记录
		if err := tx.Where("comment_id IN ?", commentIDs).Delete(&CommentLike{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 找出所有这些评论的父评论ID
		var parentIDs []int64
		if err := tx.Model(&biz.Comment{}).Where("id IN ? AND parent_id > 0", commentIDs).Pluck("parent_id", &parentIDs).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 删除所有指定的评论
		if err := tx.Where("id IN ?", commentIDs).Delete(&biz.Comment{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 更新所有父评论的回复数
		for _, parentID := range parentIDs {
			// 重新计算父评论的回复数
			var replyCount int64
			if err := tx.Model(&biz.Comment{}).Where("parent_id = ?", parentID).Count(&replyCount).Error; err != nil {
				tx.Rollback()
				return err
			}

			// 更新父评论的回复数
			if err := tx.Model(&biz.Comment{}).Where("id = ?", parentID).UpdateColumn("reply_count", replyCount).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

// ListRootComments 获取根评论列表
func (r *commentRepo) ListRootComments(ctx context.Context, module int32, resourceID string, page, pageSize int32, sortType int32) ([]*biz.Comment, error) {
	var comments []*biz.Comment

	// 计算偏移量
	offset := (page - 1) * pageSize

	query := r.data.db.WithContext(ctx).Model(&biz.Comment{}).
		Where("module = ? AND resource_id = ? AND level = 0", module, resourceID)

	// 根据排序类型添加排序条件
	switch sortType {
	case 1: // CREATE_TIME_DESC 按创建时间降序
		query = query.Order("create_gmt DESC")
	default: // LIKE_COUNT_DESC 按点赞数降序（默认）
		query = query.Order("like_count DESC, create_gmt DESC")
	}

	err := query.Limit(int(pageSize)).Offset(int(offset)).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// ListReplyComments 获取回复评论列表
func (r *commentRepo) ListReplyComments(ctx context.Context, rootIDs []int64, replyLimit int32, sortType int32) ([]*biz.Comment, error) {
	var comments []*biz.Comment

	query := r.data.db.WithContext(ctx).Model(&biz.Comment{}).
		Where("root_id IN ?", rootIDs)

	// 回复评论永远按照点赞数由高到低排序
	query = query.Order("like_count DESC, create_gmt DESC")

	err := query.Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}
