package token_repo

import (
	"context"
	"encoding/json"

	"github.com/scriptscat/cloudcat/internal/model/entity/token_entity"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"
)

type TokenRepo interface {
	Find(ctx context.Context, id string) (*token_entity.Token, error)
	FindPage(ctx context.Context) ([]*token_entity.Token, error)
	Create(ctx context.Context, token *token_entity.Token) error
	Update(ctx context.Context, token *token_entity.Token) error
	Delete(ctx context.Context, id string) error
}

var defaultToken TokenRepo

func Token() TokenRepo {
	return defaultToken
}

func RegisterToken(i TokenRepo) {
	defaultToken = i
}

type tokenRepo struct {
}

func NewToken() TokenRepo {
	return &tokenRepo{}
}

func (u *tokenRepo) Find(ctx context.Context, id string) (*token_entity.Token, error) {
	token := &token_entity.Token{}
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("token"))
		data := b.Get([]byte(id))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, token)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}

func (u *tokenRepo) Create(ctx context.Context, token *token_entity.Token) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(token)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("token")).Put([]byte(token.ID), data)
	})
}

func (u *tokenRepo) Update(ctx context.Context, token *token_entity.Token) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(token)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("token")).Put([]byte(token.ID), data)
	})
}

func (u *tokenRepo) Delete(ctx context.Context, id string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("token")).Delete([]byte(id))
	})
}

func (u *tokenRepo) FindPage(ctx context.Context) ([]*token_entity.Token, error) {
	list := make([]*token_entity.Token, 0)
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("token"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			token := &token_entity.Token{}
			if err := json.Unmarshal(v, token); err != nil {
				return err
			}
			list = append(list, token)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}
