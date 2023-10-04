package value_repo

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/scriptscat/cloudcat/internal/model/entity/value_entity"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"
)

type ValueRepo interface {
	Find(ctx context.Context, storageName, key string) (*value_entity.Value, error)
	FindPage(ctx context.Context, storageName string) ([]*value_entity.Value, int64, error)
	Create(ctx context.Context, value *value_entity.Value) error
	Update(ctx context.Context, value *value_entity.Value) error
	Delete(ctx context.Context, storageName, key string) error

	DeleteByStorage(ctx context.Context, storageName string) error
}

var defaultValue ValueRepo

func Value() ValueRepo {
	return defaultValue
}

func RegisterValue(i ValueRepo) {
	defaultValue = i
}

type valueRepo struct {
}

func NewValue() ValueRepo {
	return &valueRepo{}
}

func (u *valueRepo) Find(ctx context.Context, storageName, key string) (*value_entity.Value, error) {
	value := &value_entity.Value{}
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("value")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		data := b.Get([]byte(key))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, value)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

func (u *valueRepo) Create(ctx context.Context, value *value_entity.Value) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		b, err := tx.Bucket([]byte("value")).CreateBucketIfNotExists([]byte(value.StorageName))
		if err != nil {
			return err
		}
		return b.Put([]byte(value.Key), data)
	})
}

func (u *valueRepo) Update(ctx context.Context, value *value_entity.Value) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		b, err := tx.Bucket([]byte("value")).CreateBucketIfNotExists([]byte(value.StorageName))
		if err != nil {
			return err
		}
		return b.Put([]byte(value.Key), data)
	})
}

func (u *valueRepo) Delete(ctx context.Context, storageName, key string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("value")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

func (u *valueRepo) DeleteByStorage(ctx context.Context, storageName string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("value")).DeleteBucket([]byte(storageName))
		if errors.Is(err, bolt.ErrBucketNotFound) {
			return nil
		}
		return err
	})
}

func (u *valueRepo) FindPage(ctx context.Context, storageName string) ([]*value_entity.Value, int64, error) {
	values := make([]*value_entity.Value, 0)
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("value")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			value := &value_entity.Value{}
			if err := json.Unmarshal(v, value); err != nil {
				return err
			}
			values = append(values, value)
		}
		return nil
	}); err != nil {
		return nil, 0, err
	}
	return values, int64(len(values)), nil
}
