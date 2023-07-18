package script_repo

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"

	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"
)

type ScriptRepo interface {
	Find(ctx context.Context, id string) (*script_entity.Script, error)
	FindByPrefixID(ctx context.Context, prefix string) (*script_entity.Script, error)
	FindPage(ctx context.Context) ([]*script_entity.Script, error)
	Create(ctx context.Context, script *script_entity.Script) error
	Update(ctx context.Context, script *script_entity.Script) error
	Delete(ctx context.Context, id string) error

	FindByStorage(ctx context.Context, storageName string) ([]*script_entity.Script, error)
	StorageList(ctx context.Context) ([]*script_entity.Storage, error)
}

var defaultScript ScriptRepo

func Script() ScriptRepo {
	return defaultScript
}

func RegisterScript(i ScriptRepo) {
	defaultScript = i
}

type scriptRepo struct {
}

func NewScript() ScriptRepo {
	return &scriptRepo{}
}

func (u *scriptRepo) Find(ctx context.Context, id string) (*script_entity.Script, error) {
	script := &script_entity.Script{}
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("script"))
		data := b.Get([]byte(id))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, script)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return script, nil
}

func (u *scriptRepo) FindByPrefixID(ctx context.Context, prefix string) (*script_entity.Script, error) {
	script := &script_entity.Script{}
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("script"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.HasPrefix(string(k), prefix) {
				if err := json.Unmarshal(v, script); err != nil {
					return err
				}
				return nil
			}
		}
		return bbolt.ErrNil
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return script, nil
}

func (u *scriptRepo) Create(ctx context.Context, script *script_entity.Script) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(script)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("script")).Put([]byte(script.ID), data)
	})
}

func (u *scriptRepo) Update(ctx context.Context, script *script_entity.Script) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(script)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("script")).Put([]byte(script.ID), data)
	})
}

func (u *scriptRepo) Delete(ctx context.Context, id string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("script")).Delete([]byte(id))
	})
}

func (u *scriptRepo) FindPage(ctx context.Context) ([]*script_entity.Script, error) {
	list := make([]*script_entity.Script, 0)
	if err := bbolt.Default().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("script"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			script := &script_entity.Script{}
			if err := json.Unmarshal(v, script); err != nil {
				return err
			}
			list = append(list, script)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

func (u *scriptRepo) FindByStorage(ctx context.Context, storageName string) ([]*script_entity.Script, error) {
	list, err := u.FindPage(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*script_entity.Script, 0)
	for _, v := range list {
		if v.StorageName() == storageName {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (u *scriptRepo) StorageList(ctx context.Context) ([]*script_entity.Storage, error) {
	list, err := u.FindPage(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	for _, v := range list {
		_, ok := m[v.StorageName()]
		if !ok {
			m[v.StorageName()] = make([]string, 0)
		}
		m[v.StorageName()] = append(m[v.StorageName()], v.ID)
	}
	ret := make([]*script_entity.Storage, 0)
	for k, v := range m {
		ret = append(ret, &script_entity.Storage{
			Name:         k,
			LinkScriptID: v,
		})
	}
	return ret, nil
}
