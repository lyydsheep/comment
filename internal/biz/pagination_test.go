package biz

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCommentRepo 是 CommentRepo 的一个模拟实现
type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) ListRootComments(ctx context.Context, module int32, resourceID string, page, pageSize, sortType int32) ([]*Comment, error) {
	args := m.Called(ctx, module, resourceID, page, pageSize, sortType)
	return args.Get(0).([]*Comment), args.Error(1)
}

func (m *MockCommentRepo) ListReplyComments(ctx context.Context, rootIDs []int64, replyLimit int32, sortType int32) ([]*Comment, error) {
	args := m.Called(ctx, rootIDs, replyLimit, sortType)
	return args.Get(0).([]*Comment), args.Error(1)
}

// 这里需要实现MockCommentRepo的其他方法，仅包含测试所需的实现
func (m *MockCommentRepo) Save(ctx context.Context, c *Comment) (*Comment, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*Comment), args.Error(1)
}

func (m *MockCommentRepo) Get(ctx context.Context, id int64) (*Comment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Comment), args.Error(1)
}

func (m *MockCommentRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepo) DeleteBatch(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepo) LikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	args := m.Called(ctx, commentID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCommentRepo) UnlikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	args := m.Called(ctx, commentID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCommentUsecase_GetComments_Pagination(t *testing.T) {
	// 创建测试用例
	t.Run("Should return correct page when using pagination", func(t *testing.T) {
		// 准备模拟数据
		mockRepo := new(MockCommentRepo)
		uc := NewCommentUsecase(mockRepo)

		// 创建预期的评论数据
		comments := []*Comment{
			{
				ID:           1,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user1",
				Username:     "user1",
				Avatar:       "avatar1",
				Content:      "Comment 1",
				Level:        0,
				LikeCount:    10,
				ReplyCount:   2,
				CreateGmt:    time.Now(),
			},
			{
				ID:           2,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user2",
				Username:     "user2",
				Avatar:       "avatar2",
				Content:      "Comment 2",
				Level:        0,
				LikeCount:    5,
				ReplyCount:   1,
				CreateGmt:    time.Now().Add(-1 * time.Hour),
			},
		}

		// 设置模拟对象的行为
		mockRepo.On("ListRootComments", mock.Anything, int32(1), "article1", int32(1), int32(10), int32(0)).Return(comments, nil)
		mockRepo.On("ListReplyComments", mock.Anything, []int64{1, 2}, int32(5), int32(0)).Return([]*Comment{}, nil)

		// 执行测试
		result, err := uc.GetComments(context.Background(), 1, "article1", 5, 1, 10, 0)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		mockRepo.AssertExpectations(t)
	})
}

func TestCommentUsecase_GetComments_Sorting(t *testing.T) {
	// 创建测试用例
	t.Run("Should sort comments by like count desc when sort type is 0", func(t *testing.T) {
		// 准备模拟数据
		mockRepo := new(MockCommentRepo)
		uc := NewCommentUsecase(mockRepo)

		// 创建预期的评论数据 - 按点赞数降序排列
		comments := []*Comment{
			{
				ID:           3,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user3",
				Username:     "user3",
				Avatar:       "avatar3",
				Content:      "Comment with most likes",
				Level:        0,
				LikeCount:    100,
				ReplyCount:   5,
				CreateGmt:    time.Now().Add(-2 * time.Hour),
			},
			{
				ID:           4,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user4",
				Username:     "user4",
				Avatar:       "avatar4",
				Content:      "Comment with fewer likes",
				Level:        0,
				LikeCount:    50,
				ReplyCount:   3,
				CreateGmt:    time.Now(),
			},
		}

		// 设置模拟对象的行为 - 使用默认排序类型(0: 按点赞数降序)
		mockRepo.On("ListRootComments", mock.Anything, int32(1), "article1", int32(1), int32(10), int32(0)).Return(comments, nil)
		mockRepo.On("ListReplyComments", mock.Anything, []int64{3, 4}, int32(5), int32(0)).Return([]*Comment{}, nil)

		// 执行测试
		result, err := uc.GetComments(context.Background(), 1, "article1", 5, 1, 10, 0)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(3), result[0].ID) // 点赞数最多的应该在第一个
		assert.Equal(t, int64(4), result[1].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should sort comments by create time desc when sort type is 1", func(t *testing.T) {
		// 准备模拟数据
		mockRepo := new(MockCommentRepo)
		uc := NewCommentUsecase(mockRepo)

		// 创建预期的评论数据 - 按创建时间降序排列
		comments := []*Comment{
			{
				ID:           5,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user5",
				Username:     "user5",
				Avatar:       "avatar5",
				Content:      "Newest comment",
				Level:        0,
				LikeCount:    10,
				ReplyCount:   1,
				CreateGmt:    time.Now(),
			},
			{
				ID:           6,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user6",
				Username:     "user6",
				Avatar:       "avatar6",
				Content:      "Older comment",
				Level:        0,
				LikeCount:    100,
				ReplyCount:   5,
				CreateGmt:    time.Now().Add(-2 * time.Hour),
			},
		}

		// 设置模拟对象的行为 - 使用创建时间排序类型(1: 按创建时间降序)
		mockRepo.On("ListRootComments", mock.Anything, int32(1), "article1", int32(1), int32(10), int32(1)).Return(comments, nil)
		mockRepo.On("ListReplyComments", mock.Anything, []int64{5, 6}, int32(5), int32(1)).Return([]*Comment{}, nil)

		// 执行测试
		result, err := uc.GetComments(context.Background(), 1, "article1", 5, 1, 10, 1)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(5), result[0].ID) // 最新的评论应该在第一个
		assert.Equal(t, int64(6), result[1].ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestCommentUsecase_GetComments_Error(t *testing.T) {
	// 创建测试用例
	t.Run("Should return error when repo returns error", func(t *testing.T) {
		// 准备模拟数据
		mockRepo := new(MockCommentRepo)
		uc := NewCommentUsecase(mockRepo)

		// 设置模拟对象的行为 - 返回错误
		mockRepo.On("ListRootComments", mock.Anything, int32(1), "article1", int32(1), int32(10), int32(0)).Return([]*Comment{}, gorm.ErrRecordNotFound)

		// 执行测试
		result, err := uc.GetComments(context.Background(), 1, "article1", 5, 1, 10, 0)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error when reply comments query fails", func(t *testing.T) {
		// 准备模拟数据
		mockRepo := new(MockCommentRepo)
		uc := NewCommentUsecase(mockRepo)

		// 创建预期的评论数据
		comments := []*Comment{
			{
				ID:           7,
				Module:       1,
				ResourceID:   "article1",
				UserID:       "user7",
				Username:     "user7",
				Avatar:       "avatar7",
				Content:      "Comment with replies",
				Level:        0,
				LikeCount:    15,
				ReplyCount:   3,
				CreateGmt:    time.Now(),
			},
		}

		// 设置模拟对象的行为 - 根评论成功，回复评论失败
		mockRepo.On("ListRootComments", mock.Anything, int32(1), "article1", int32(1), int32(10), int32(0)).Return(comments, nil)
		mockRepo.On("ListReplyComments", mock.Anything, []int64{7}, int32(5), int32(0)).Return([]*Comment{}, gorm.ErrRecordNotFound)

		// 执行测试
		result, err := uc.GetComments(context.Background(), 1, "article1", 5, 1, 10, 0)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}