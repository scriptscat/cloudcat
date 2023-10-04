package resource_repo

import (
	"context"
	"encoding/json"

	"github.com/scriptscat/cloudcat/internal/model/entity/resource_entity"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"
)

type ResourceRepo interface {
	Find(ctx context.Context, url string) (*resource_entity.Resource, error)
	FindPage(ctx context.Context) ([]*resource_entity.Resource, error)
	Create(ctx context.Context, resource *resource_entity.Resource) error
	Update(ctx context.Context, resource *resource_entity.Resource) error
	Delete(ctx context.Context, url string) error
}

var defaultResource ResourceRepo

func Resource() ResourceRepo {
	return defaultResource
}

func RegisterResource(i ResourceRepo) {
	defaultResource = i
}

type resourceRepo struct {
}

func NewResource() ResourceRepo {
	return &resourceRepo{}
}

func (u *resourceRepo) Find(ctx context.Context, url string) (*resource_entity.Resource, error) {
	resource := &resource_entity.Resource{}
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resource"))
		data := b.Get([]byte(url))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, resource)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return resource, nil
}

func (u *resourceRepo) Create(ctx context.Context, resource *resource_entity.Resource) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(resource)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("resource")).Put([]byte(resource.URL), data)
	})
}

func (u *resourceRepo) Update(ctx context.Context, resource *resource_entity.Resource) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(resource)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("resource")).Put([]byte(resource.URL), data)
	})
}

func (u *resourceRepo) Delete(ctx context.Context, url string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resource")).Delete([]byte(url))
	})
}

func (u *resourceRepo) FindPage(ctx context.Context) ([]*resource_entity.Resource, error) {
	list := make([]*resource_entity.Resource, 0)
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resource"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			resource := &resource_entity.Resource{}
			if err := json.Unmarshal(v, resource); err != nil {
				return err
			}
			list = append(list, resource)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}
