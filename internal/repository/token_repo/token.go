package token_repo

import (
	"context"
	"encoding/json"

	"github.com/scriptscat/cloudcat/internal/model/entity/token_entity"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"
)

type SecretRepo interface {
	Find(ctx context.Context, id string) (*token_entity.Token, error)
	FindPage(ctx context.Context) ([]*token_entity.Token, error)
	Create(ctx context.Context, secret *token_entity.Token) error
	Update(ctx context.Context, secret *token_entity.Token) error
	Delete(ctx context.Context, id string) error
}

var defaultSecret SecretRepo

func Secret() SecretRepo {
	return defaultSecret
}

func RegisterSecret(i SecretRepo) {
	defaultSecret = i
}

type secretRepo struct {
}

func NewSecret() SecretRepo {
	return &secretRepo{}
}

func (u *secretRepo) Find(ctx context.Context, id string) (*token_entity.Token, error) {
	secret := &token_entity.Token{}
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("token"))
		data := b.Get([]byte(id))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, secret)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return secret, nil
}

func (u *secretRepo) Create(ctx context.Context, secret *token_entity.Token) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(secret)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("token")).Put([]byte(secret.ID), data)
	})
}

func (u *secretRepo) Update(ctx context.Context, secret *token_entity.Token) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(secret)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("token")).Put([]byte(secret.ID), data)
	})
}

func (u *secretRepo) Delete(ctx context.Context, id string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("token")).Delete([]byte(id))
	})
}

func (u *secretRepo) FindPage(ctx context.Context) ([]*token_entity.Token, error) {
	list := make([]*token_entity.Token, 0)
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("token"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			secret := &token_entity.Token{}
			if err := json.Unmarshal(v, secret); err != nil {
				return err
			}
			list = append(list, secret)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}
