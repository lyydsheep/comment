package biz

import (
	"container/heap"
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewCommentUsecase)

// TxnManager 事务管理
type TxnManager interface {
	// Txn 开启一个事务执行 fn 函数
	Txn(ctx context.Context, fn func(ctx context.Context) error) error
	heap.Interface
}
