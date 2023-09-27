package migrations

import (
	bbolt2 "github.com/scriptscat/cloudcat/pkg/bbolt"
	"go.etcd.io/bbolt"
)

func T20230210() *bbolt2.Migration {
	return &bbolt2.Migration{
		ID: "20230904",
		Migrate: func(db *bbolt.DB) error {
			return db.Update(func(tx *bbolt.Tx) error {
				if _, err := tx.CreateBucketIfNotExists([]byte("value")); err != nil {
					return err
				}
				if _, err := tx.CreateBucketIfNotExists([]byte("script")); err != nil {
					return err
				}
				if _, err := tx.CreateBucketIfNotExists([]byte("cookie")); err != nil {
					return err
				}
				if _, err := tx.CreateBucketIfNotExists([]byte("token")); err != nil {
					return err
				}
				return nil
			})
		},
		Rollback: func(tx *bbolt.DB) error {
			return nil
		},
	}
}
