package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	client *Client
}

// NewTransactionManager 创建新的事务管理器
func NewTransactionManager(client *Client) *TransactionManager {
	return &TransactionManager{
		client: client,
	}
}

// TransactionFunc 事务函数类型
type TransactionFunc func(sessCtx mongo.SessionContext) error

// WithTransaction 执行事务
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	session, err := tm.client.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// 设置事务选项
	txnOpts := options.Transaction().
		SetReadPreference(nil).
		SetWriteConcern(nil).
		SetReadConcern(nil)

	// 执行事务
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

// WithSession 使用会话执行操作
func (tm *TransactionManager) WithSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	session, err := tm.client.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, fn)
}

// TransactionalRepository 支持事务的仓储
type TransactionalRepository struct {
	*Collection
	session mongo.Session
}

// NewTransactionalRepository 创建支持事务的仓储
func NewTransactionalRepository(client *Client, collectionName string) *TransactionalRepository {
	return &TransactionalRepository{
		Collection: NewCollection(client, collectionName),
	}
}

// WithTransaction 在事务中执行操作
func (tr *TransactionalRepository) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext, repo *Collection) error) error {
	session, err := tr.cli.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	txnOpts := options.Transaction()

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 创建一个在事务上下文中的仓储
		txnRepo := &Collection{
			cli:        tr.cli,
			collection: tr.collection,
		}
		return nil, fn(sessCtx, txnRepo)
	}, txnOpts)

	return err
}
