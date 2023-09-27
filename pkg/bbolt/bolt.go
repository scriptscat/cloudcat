package bbolt

import (
	"context"
	"errors"
	"time"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/gogo"
	"github.com/codfrm/cago/pkg/logger"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var (
	ErrNil = errors.New("nil")
)

func IsNil(err error) bool {
	return errors.Is(err, ErrNil)
}

var db *bolt.DB

type Config struct {
	Path string `yaml:"path"`
}

func Bolt() cago.FuncComponent {
	return func(ctx context.Context, cfg *configs.Config) error {
		var err error
		config := &Config{}
		if err := cfg.Scan("db", config); err != nil {
			return err
		}
		db, err = bolt.Open(config.Path, 0600, &bolt.Options{
			Timeout: 5 * time.Second,
		})
		if err != nil {
			return err
		}
		if err := gogo.Go(func(ctx context.Context) error {
			<-ctx.Done()
			err := db.Close()
			if err != nil {
				logger.Ctx(ctx).Error("close leveldb err: %v", zap.Error(err))
			}
			return nil
		}, gogo.WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}
}

func Default() *bolt.DB {
	return db
}

type contextKey int

const (
	transactionKey contextKey = iota + 1
)

func TxCtx(ctx context.Context) *bolt.Tx {
	tx, ok := ctx.Value(transactionKey).(*bolt.Tx)
	if ok {
		return tx
	}
	return nil
}

func Transaction(ctx context.Context, fn func(ctx context.Context, tx *bolt.Tx) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		ctx = context.WithValue(ctx, transactionKey, tx)
		return fn(ctx, tx)
	})
}
