package boltDB

import (
	"code-fabrik.com/bend/domain/config"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

type ConfigRepository struct {
	DB *bolt.DB
}

const configBucket = "config"

func (rr ConfigRepository) Save(config config.Config) error {
	db := rr.DB

	err := db.Update(func(txn *bolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists([]byte(configBucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		err = b.Put([]byte(config.Path), data)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (rr ConfigRepository) Find(path string) *config.Config {
	db := rr.DB

	result := config.Config{}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(configBucket))

		if bucket != nil {
			value := bucket.Get([]byte(path))
			if value == nil {
				return nil
			}

			json.Unmarshal(value, &result)
		}
		return fmt.Errorf("config not found")
	})

	if err != nil {
		return nil
	}
	return &result
}

func (rr ConfigRepository) Delete(path string) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(configBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(path))
	})
}
