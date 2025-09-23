package data

import (
	"comment/internal/biz"
	"context"
	"time"
)

// CommentLike 评论点赞记录模型
type CommentLike struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CommentID  int64     `gorm:"column:comment_id;type:bigint;not null;index:idx_comment_user,unique"`
	UserID     string    `gorm:"column:user_id;type:varchar(32);not null;index:idx_comment_user,unique"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP"`
}

func (cl *CommentLike) TableName() string {
	return "comment_like"
}

// LikeComment 点赞评论
func (r *commentRepo) LikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 检查是否已经点赞
	var existingLike CommentLike
	error := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&existingLike).Error
	if error == nil {
		// 已经点赞过，直接返回当前点赞数
		var likeCount int64
		if err := tx.Model(&biz.Comment{}).Where("id = ?", commentID).Select("like_count").Scan(&likeCount).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		return likeCount, nil
	}

	// 添加点赞记录
	like := &CommentLike{
		CommentID:  commentID,
		UserID:     userID,
		CreateTime: time.Now(),
	}
	if err := tx.Create(like).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 更新评论的点赞数
	var likeCount int64
	if err := tx.Model(&biz.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", tx.Model(&biz.Comment{}).Select("like_count + ?", 1).Where("id = ?", commentID)).
		Select("like_count").Scan(&likeCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return likeCount, nil
}

// UnlikeComment 取消点赞评论
func (r *commentRepo) UnlikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 删除点赞记录
	if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&CommentLike{}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 如果没有删除任何记录，说明用户没有点赞过，直接返回当前点赞数
	if tx.RowsAffected == 0 {
		var likeCount int64
		if err := tx.Model(&biz.Comment{}).Where("id = ?", commentID).Select("like_count").Scan(&likeCount).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		return likeCount, nil
	}

	// 更新评论的点赞数
	var likeCount int64
	if err := tx.Model(&biz.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", tx.Model(&biz.Comment{}).Select("like_count - ?", 1).Where("id = ?", commentID)).
		Select("like_count").Scan(&likeCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return likeCount, nil
}