package biz

import (
	v1 "comment/api/comment/v1"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// CommentRepoMock 是CommentRepo接口的mock实现
type CommentRepoMock struct {
	mock.Mock
}

func (m *CommentRepoMock) Save(ctx context.Context, comment *Comment) (*Comment, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(*Comment), args.Error(1)
}

func (m *CommentRepoMock) Get(ctx context.Context, id int64) (*Comment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Comment), args.Error(1)
}

func (m *CommentRepoMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *CommentRepoMock) DeleteBatch(ctx context.Context, rootID int64) error {
	args := m.Called(ctx, rootID)
	return args.Error(0)
}

func (m *CommentRepoMock) ListRootComments(ctx context.Context, module int32, resourceID string, page, pageSize int32, sortType int32) ([]*Comment, error) {
	args := m.Called(ctx, module, resourceID, page, pageSize, sortType)
	return args.Get(0).([]*Comment), args.Error(1)
}

func (m *CommentRepoMock) ListReplyComments(ctx context.Context, rootIDs []int64, replyLimit int32, sortType int32) ([]*Comment, error) {
	args := m.Called(ctx, rootIDs, replyLimit, sortType)
	return args.Get(0).([]*Comment), args.Error(1)
}

func (m *CommentRepoMock) LikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	args := m.Called(ctx, commentID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *CommentRepoMock) UnlikeComment(ctx context.Context, commentID int64, userID string) (int64, error) {
	args := m.Called(ctx, commentID, userID)
	return args.Get(0).(int64), args.Error(1)
}

// CommentTestSuite 是测试套件
type CommentTestSuite struct {
	suite.Suite
	repoMock *CommentRepoMock
	usecase  *CommentUsecase
}

func TestComment_TableName(t *testing.T) {
	type fields struct {
		Module          int32
		ResourceID      string
		ID              int64
		RootCommentID   int64
		ParentCommentID int64
		UserID          string
		Username        string
		Avatar          string
		Content         string
		Level           int32
		LikeCount       int64
		ReplyCount      int64
		CreateGmt       time.Time
		UpdateGmt       time.Time
		ReplyComments   []*Comment
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Comment{
				Module:          tt.fields.Module,
				ResourceID:      tt.fields.ResourceID,
				ID:              tt.fields.ID,
				RootCommentID:   tt.fields.RootCommentID,
				ParentCommentID: tt.fields.ParentCommentID,
				UserID:          tt.fields.UserID,
				Username:        tt.fields.Username,
				Avatar:          tt.fields.Avatar,
				Content:         tt.fields.Content,
				Level:           tt.fields.Level,
				LikeCount:       tt.fields.LikeCount,
				ReplyCount:      tt.fields.ReplyCount,
				CreateGmt:       tt.fields.CreateGmt,
				UpdateGmt:       tt.fields.UpdateGmt,
				ReplyComments:   tt.fields.ReplyComments,
			}
			if got := c.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *CommentTestSuite) SetupTest() {
	s.repoMock = new(CommentRepoMock)
	s.usecase = NewCommentUsecase(s.repoMock)
}

// TestNewCommentUsecase 测试创建CommentUsecase
func (s *CommentTestSuite) TestNewCommentUsecase() {
	tests := []struct {
		name string
		repo CommentRepo
		want *CommentUsecase
	}{
		{
			name: "正常创建CommentUsecase",
			repo: new(CommentRepoMock),
			want: &CommentUsecase{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := NewCommentUsecase(tt.repo)
			s.Assert().NotNil(got)
			s.Assert().Equal(tt.repo, got.repo)
		})
	}
}

// TestCommentUsecase_CreateComment 测试创建评论
func (s *CommentTestSuite) TestCommentUsecase_CreateComment() {
	tests := []struct {
		name    string
		prepare func()
		args    *Comment
		want    *v1.Comment
		wantErr bool
	}{
		{
			name: "正常创建评论",
			prepare: func() {
				s.repoMock.On("Save", mock.Anything, mock.MatchedBy(func(c *Comment) bool {
					return c.Content == "这是一条测试评论"
				})).Return(&Comment{
					Module:          1,
					ResourceID:      "resource_123",
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					UserID:          "user_123",
					Username:        "test_user",
					Avatar:          "avatar_url",
					Content:         "这是一条测试评论",
					Level:           0,
					LikeCount:       0,
					ReplyCount:      0,
					CreateGmt:       time.Now(),
					UpdateGmt:       time.Now(),
				}, nil).Once()
			},
			args: &Comment{
				Module:     1,
				ResourceID: "resource_123",
				UserID:     "user_123",
				Username:   "test_user",
				Avatar:     "avatar_url",
				Content:    "这是一条测试评论",
				Level:      0,
			},
			want: &v1.Comment{
				Module:     1,
				ResourceId: "resource_123",
				CommentId:  1,
				UserId:     "user_123",
				Username:   "test_user",
				Avatar:     "avatar_url",
				Content:    "这是一条测试评论",
				Level:      0,
				LikeCount:  0,
				ReplyCount: 0,
			},
			wantErr: false,
		},
		{
			name: "创建评论时数据库错误",
			prepare: func() {
				s.repoMock.On("Save", mock.Anything, mock.MatchedBy(func(c *Comment) bool {
					return c.Content == "数据库错误测试"
				})).Return((*Comment)(nil), errors.New("数据库保存失败")).Once()
			},
			args: &Comment{
				Content: "数据库错误测试",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			got, err := s.usecase.CreateComment(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil && tt.want != nil {
				// 只比较关心的字段，忽略时间等动态字段
				s.Assert().Equal(tt.want.Module, got.Module)
				s.Assert().Equal(tt.want.ResourceId, got.ResourceId)
				s.Assert().Equal(tt.want.CommentId, got.CommentId)
				s.Assert().Equal(tt.want.UserId, got.UserId)
				s.Assert().Equal(tt.want.Username, got.Username)
				s.Assert().Equal(tt.want.Avatar, got.Avatar)
				s.Assert().Equal(tt.want.Content, got.Content)
				s.Assert().Equal(tt.want.Level, got.Level)
				s.Assert().Equal(tt.want.LikeCount, got.LikeCount)
				s.Assert().Equal(tt.want.ReplyCount, got.ReplyCount)
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// TestCommentUsecase_DeleteComment 测试删除评论
func (s *CommentTestSuite) TestCommentUsecase_DeleteComment() {
	tests := []struct {
		name    string
		prepare func()
		id      int64
		wantErr bool
	}{
		{
			name: "正常删除评论",
			prepare: func() {
				// 模拟获取评论成功
				s.repoMock.On("Get", mock.Anything, int64(1)).Return(&Comment{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "要删除的评论",
				}, nil).Once()

				// 模拟删除评论及回复成功
				s.repoMock.On("DeleteBatch", mock.Anything, int64(1)).Return(nil).Once()
			},
			id:      1,
			wantErr: false,
		},
		{
			name: "删除不存在的评论",
			prepare: func() {
				// 模拟获取评论失败
				s.repoMock.On("Get", mock.Anything, int64(999)).Return((*Comment)(nil), errors.New("评论不存在")).Once()
			},
			id:      999,
			wantErr: true,
		},
		{
			name: "删除评论时数据库错误",
			prepare: func() {
				// 模拟获取评论成功
				s.repoMock.On("Get", mock.Anything, int64(2)).Return(&Comment{
					ID:              2,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "要删除的评论",
				}, nil).Once()

				// 模拟删除评论及回复失败
				s.repoMock.On("DeleteBatch", mock.Anything, int64(2)).Return(errors.New("数据库删除失败")).Once()
			},
			id:      2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.usecase.DeleteComment(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// TestCommentUsecase_GetComments 测试获取评论
func (s *CommentTestSuite) TestCommentUsecase_GetComments() {
	tests := []struct {
		name       string
		prepare    func()
		module     int32
		resourceID string
		replyLimit int32
		page       int32
		pageSize   int32
		sortType   int32
		want       []*Comment
		wantErr    bool
	}{
		{
			name: "正常获取根评论",
			prepare: func() {
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return([]*Comment{
						{
							Module:          1,
							ResourceID:      "resource_123",
							ID:              1,
							RootCommentID:   0,
							ParentCommentID: 0,
							UserID:          "user_1",
							Username:        "user1",
							Avatar:          "avatar1",
							Content:         "第一条根评论",
							Level:           0,
							LikeCount:       10,
							ReplyCount:      2,
							CreateGmt:       time.Now(),
							UpdateGmt:       time.Now(),
						},
						{
							Module:          1,
							ResourceID:      "resource_123",
							ID:              2,
							RootCommentID:   0,
							ParentCommentID: 0,
							UserID:          "user_2",
							Username:        "user2",
							Avatar:          "avatar2",
							Content:         "第二条根评论",
							Level:           0,
							LikeCount:       5,
							ReplyCount:      1,
							CreateGmt:       time.Now(),
							UpdateGmt:       time.Now(),
						},
					}, nil).Once()
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 0, // 无限制
			page:       1,
			pageSize:   10,
			sortType:   0,
			want: []*Comment{
				{
					ID:         1,
					Content:    "第一条根评论",
					Level:      0,
					LikeCount:  10,
					ReplyCount: 2,
				},
				{
					ID:         2,
					Content:    "第二条根评论",
					Level:      0,
					LikeCount:  5,
					ReplyCount: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "获取根评论和回复评论",
			prepare: func() {
				// 模拟获取根评论
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return([]*Comment{
						{
							ID:              1,
							RootCommentID:   0,
							ParentCommentID: 0,
							Content:         "根评论",
							Level:           0,
						},
					}, nil).Once()

				// 模拟获取回复评论
				s.repoMock.On("ListReplyComments", mock.Anything, []int64{1}, int32(2), int32(0)).
					Return([]*Comment{
						{
							ID:              2,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复评论",
							Level:           1,
						},
					}, nil).Once()
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 2,
			page:       1,
			pageSize:   10,
			sortType:   0,
			want: []*Comment{
				{
					ID:      1,
					Content: "根评论",
					Level:   0,
					ReplyComments: []*Comment{
						{
							ID:            2,
							Content:       "回复评论",
							Level:         1,
							ReplyComments: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "获取根评论时数据库错误",
			prepare: func() {
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return(([]*Comment)(nil), errors.New("数据库查询失败")).Once()
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 1,
			page:       1,
			pageSize:   10,
			sortType:   0,
			want:       nil,
			wantErr:    true,
		},
		{
			name: "获取回复评论时数据库错误",
			prepare: func() {
				// 模拟获取根评论成功
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return([]*Comment{
						{
							ID:              1,
							RootCommentID:   0,
							ParentCommentID: 0,
							Content:         "根评论",
							Level:           0,
						},
					}, nil).Once()

				// 模拟获取回复评论失败
				s.repoMock.On("ListReplyComments", mock.Anything, []int64{1}, int32(2), int32(0)).
					Return(([]*Comment)(nil), errors.New("数据库查询失败")).Once()
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 2,
			page:       1,
			pageSize:   10,
			sortType:   0,
			want:       nil,
			wantErr:    true,
		},
		{
			name: "没有根评论的情况",
			prepare: func() {
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return([]*Comment{}, nil).Once()
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 2,
			page:       1,
			pageSize:   10,
			sortType:   0,
			want:       []*Comment{},
			wantErr:    false,
		},
		{
			name: "replyLimit为0时不获取回复评论",
			prepare: func() {
				// 模拟获取根评论
				s.repoMock.On("ListRootComments", mock.Anything, int32(1), "resource_123", int32(1), int32(10), int32(0)).
					Return([]*Comment{
						{
							ID:              1,
							RootCommentID:   0,
							ParentCommentID: 0,
							Content:         "根评论",
							Level:           0,
						},
					}, nil).Once()

				// replyLimit为0时不应该调用ListReplyComments
			},
			module:     1,
			resourceID: "resource_123",
			replyLimit: 0, // 不获取回复
			page:       1,
			pageSize:   10,
			sortType:   0,
			want: []*Comment{
				{
					ID:            1,
					Content:       "根评论",
					Level:         0,
					ReplyComments: []*Comment{}, // 应该是空的回复列表
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			got, err := s.usecase.GetComments(context.Background(), tt.module, tt.resourceID, tt.replyLimit, tt.page, tt.pageSize, tt.sortType)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				s.Assert().Equal(len(tt.want), len(got))
				for i := range tt.want {
					s.assertCommentEqual(tt.want[i], got[i])
				}
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// assertCommentEqual 比较两个评论是否相等
func (s *CommentTestSuite) assertCommentEqual(expected, actual *Comment) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		s.T().Errorf("Expected or actual comment is nil")
		return
	}

	s.Assert().Equal(expected.ID, actual.ID)
	s.Assert().Equal(expected.Content, actual.Content)
	s.Assert().Equal(expected.Level, actual.Level)
	s.Assert().Equal(expected.LikeCount, actual.LikeCount)
	s.Assert().Equal(expected.ReplyCount, actual.ReplyCount)
	s.Assert().Equal(len(expected.ReplyComments), len(actual.ReplyComments))

	for i := range expected.ReplyComments {
		if i < len(actual.ReplyComments) {
			s.assertCommentEqual(expected.ReplyComments[i], actual.ReplyComments[i])
		}
	}
}

// TestCommentUsecase_buildCommentTree 测试构建评论树
func (s *CommentTestSuite) TestCommentUsecase_buildCommentTree() {
	tests := []struct {
		name          string
		rootComments  []*Comment
		replyComments []*Comment
		replyLimit    int32
		want          []*Comment // 期望的结果
	}{
		{
			name: "构建简单的评论树",
			rootComments: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论1",
					Level:           0,
				},
				{
					ID:              2,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论2",
					Level:           0,
				},
			},
			replyComments: []*Comment{
				{
					ID:              3,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "回复根评论1-1",
					Level:           1,
				},
				{
					ID:              4,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "回复根评论1-2",
					Level:           1,
				},
				{
					ID:              5,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "回复根评论1-3",
					Level:           1,
				},
			},
			replyLimit: 0, // 无限制
			want: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论1",
					Level:           0,
					ReplyComments: []*Comment{
						{
							ID:              3,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复根评论1-1",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
						{
							ID:              4,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复根评论1-2",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
						{
							ID:              5,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复根评论1-3",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
					},
				},
				{
					ID:              2,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论2",
					Level:           0,
					ReplyComments:   []*Comment{},
				},
			},
		},
		{
			name: "限制回复个数",
			rootComments: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
				},
			},
			replyComments: []*Comment{
				{
					ID:              2,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "第一条评论",
					Level:           1,
				},
				{
					ID:              3,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "第二条评论",
					Level:           1,
				},
				{
					ID:              4,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "第三条评论",
					Level:           1,
				},
			},
			replyLimit: 2, // 限制为2个回复
			want: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
					ReplyComments: []*Comment{
						{
							ID:              2,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "第一条评论",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
						{
							ID:              3,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "第二条评论",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
						// 第三条评论应该被限制，不会出现在结果中
					},
				},
			},
		},
		{
			name:          "空评论列表",
			rootComments:  []*Comment{},
			replyComments: []*Comment{},
			replyLimit:    0,
			want:          []*Comment{},
		},
		{
			name: "只有根评论没有回复",
			rootComments: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
				},
			},
			replyComments: []*Comment{},
			replyLimit:    0,
			want: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
					ReplyComments:   []*Comment{},
				},
			},
		},
		{
			name: "replyLimit为负数时视为无限制",
			rootComments: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
				},
			},
			replyComments: []*Comment{
				{
					ID:              2,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "回复1",
					Level:           1,
				},
				{
					ID:              3,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "回复2",
					Level:           1,
				},
			},
			replyLimit: -1, // 负数视为无限制
			want: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
					ReplyComments: []*Comment{
						{
							ID:              2,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复1",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
						{
							ID:              3,
							RootCommentID:   1,
							ParentCommentID: 1,
							Content:         "回复2",
							Level:           1,
							ReplyComments:   []*Comment{},
						},
					},
				},
			},
		},
		{
			name: "复杂的嵌套回复结构",
			rootComments: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
				},
			},
			replyComments: []*Comment{
				{
					ID:              2,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "直接回复根评论",
					Level:           1,
				},
				{
					ID:              3,
					RootCommentID:   1,
					ParentCommentID: 2,
					Content:         "回复ID为2的评论",
					Level:           2,
				},
				{
					ID:              4,
					RootCommentID:   1,
					ParentCommentID: 1,
					Content:         "另一个直接回复根评论",
					Level:           1,
				},
			},
			replyLimit: 2, // 限制为2个直接回复
			want: []*Comment{
				{
					ID:              1,
					RootCommentID:   0,
					ParentCommentID: 0,
					Content:         "根评论",
					Level:           0,
					ReplyComments: []*Comment{
						{
							ID:      2,
							Content: "直接回复根评论",
							Level:   1,
							ReplyComments: []*Comment{
								{
									ID:            3,
									Content:       "回复ID为2的评论",
									Level:         2,
									ReplyComments: []*Comment{},
								},
							},
						},
						{
							ID:            4,
							Content:       "另一个直接回复根评论",
							Level:         1,
							ReplyComments: []*Comment{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// 执行构建评论树操作
			s.usecase.buildCommentTree(tt.rootComments, tt.replyComments, tt.replyLimit)

			// 验证结果
			s.Assert().Equal(len(tt.want), len(tt.rootComments))
			for i := range tt.want {
				s.assertCommentEqual(tt.want[i], tt.rootComments[i])
			}
		})
	}
}

// TestCommentUsecase_deleteCommentAndReplies 测试删除评论及回复
func (s *CommentTestSuite) TestCommentUsecase_deleteCommentAndReplies() {
	tests := []struct {
		name    string
		prepare func()
		comment *Comment
		wantErr bool
	}{
		{
			name: "正常删除评论及回复",
			prepare: func() {
				s.repoMock.On("DeleteBatch", mock.Anything, int64(1)).Return(nil).Once()
			},
			comment: &Comment{
				ID:              1,
				RootCommentID:   0,
				ParentCommentID: 0,
				Content:         "要删除的评论",
			},
			wantErr: false,
		},
		{
			name: "删除评论时数据库错误",
			prepare: func() {
				s.repoMock.On("DeleteBatch", mock.Anything, int64(1)).Return(errors.New("数据库删除失败")).Once()
			},
			comment: &Comment{
				ID:              1,
				RootCommentID:   0,
				ParentCommentID: 0,
				Content:         "要删除的评论",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.usecase.deleteCommentAndReplies(context.Background(), tt.comment)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("deleteCommentAndReplies() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// TestCommentUsecase_LikeComment 测试点赞评论
func (s *CommentTestSuite) TestCommentUsecase_LikeComment() {
	tests := []struct {
		name    string
		prepare func()
		commentID int64
		userID string
		wantCount int64
		wantErr bool
	}{{
		name: "NormalLike",
		prepare: func() {
			s.repoMock.On("LikeComment", mock.Anything, int64(1), "user_123").Return(int64(1), nil).Once()
		},
		commentID: 1,
		userID: "user_123",
		wantCount: 1,
		wantErr: false,
	}, {
		name: "LikeWithDatabaseError",
		prepare: func() {
			s.repoMock.On("LikeComment", mock.Anything, int64(2), "user_123").Return(int64(0), errors.New("database error")).Once()
		},
		commentID: 2,
		userID: "user_123",
		wantCount: 0,
		wantErr: true,
	}, {
		name: "DuplicateLike",
		prepare: func() {
			s.repoMock.On("LikeComment", mock.Anything, int64(3), "user_123").Return(int64(1), nil).Once()
		},
		commentID: 3,
		userID: "user_123",
		wantCount: 1,
		wantErr: false,
	}}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			gotCount, err := s.usecase.LikeComment(context.Background(), tt.commentID, tt.userID)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("LikeComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				s.T().Errorf("LikeComment() gotCount = %v, want %v", gotCount, tt.wantCount)
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// TestCommentUsecase_UnlikeComment 测试取消点赞评论
func (s *CommentTestSuite) TestCommentUsecase_UnlikeComment() {
	tests := []struct {
		name    string
		prepare func()
		commentID int64
		userID string
		wantCount int64
		wantErr bool
	}{{
		name: "NormalUnlike",
		prepare: func() {
			s.repoMock.On("UnlikeComment", mock.Anything, int64(1), "user_123").Return(int64(0), nil).Once()
		},
		commentID: 1,
		userID: "user_123",
		wantCount: 0,
		wantErr: false,
	}, {
		name: "UnlikeWithDatabaseError",
		prepare: func() {
			s.repoMock.On("UnlikeComment", mock.Anything, int64(2), "user_123").Return(int64(1), errors.New("database error")).Once()
		},
		commentID: 2,
		userID: "user_123",
		wantCount: 0,
		wantErr: true,
	}, {
		name: "UnlikeNotLikedComment",
		prepare: func() {
			s.repoMock.On("UnlikeComment", mock.Anything, int64(3), "user_123").Return(int64(0), nil).Once()
		},
		commentID: 3,
		userID: "user_123",
		wantCount: 0,
		wantErr: false,
	}}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // 重新初始化mock
			if tt.prepare != nil {
				tt.prepare()
			}
			gotCount, err := s.usecase.UnlikeComment(context.Background(), tt.commentID, tt.userID)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("UnlikeComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				s.T().Errorf("UnlikeComment() gotCount = %v, want %v", gotCount, tt.wantCount)
			}

			// 确保预期的调用都发生了
			s.repoMock.AssertExpectations(s.T())
		})
	}
}

// TestCommentSuite 启动测试套件
func TestCommentSuite(t *testing.T) {
	suite.Run(t, new(CommentTestSuite))
}
