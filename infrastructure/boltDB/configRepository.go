package boltDB

import (
	"code-fabrik.com/bend/domain/config"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
)

type ConfigRepository struct {
	DB *bolt.DB
}

const configBucket = "config"

func (rr ConfigRepository) Save(config config.Config) error {
	db := rr.DB

	config.Path = configKey(config.Path)

	err := db.Update(func(txn *bolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists([]byte(configBucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		err = b.Put([]byte(configKey(config.Path)), data)
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
			value := bucket.Get([]byte(configKey(path)))
			if value == nil {
				return fmt.Errorf("config not found")
			}

			return json.Unmarshal(value, &result)
		}
		return fmt.Errorf("config not found")
	})

	if err != nil {
		return nil
	}
	return &result
}

func (rr ConfigRepository) FindAll() []config.Config {
	db := rr.DB

	var result []config.Config

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(configBucket))

		if bucket != nil {
			return bucket.ForEach(func(key, value []byte) error {
				v := config.Config{}
				json.Unmarshal(value, &v)

				result = append(result, v)
				return nil
			})
		}
		return nil
	})

	if err != nil {
		return nil
	}
	return result
}

func (rr ConfigRepository) Delete(path string) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(configBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		err := bucket.Delete([]byte(configKey(path)))
		return err
	})
}

func (rr ConfigRepository) DeleteAll() error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		return txn.DeleteBucket([]byte(configBucket))

	})
}
func configKey(path string) string {
	if strings.Index(path, "/") != 0 {
		return "/" + path
	}
	return path
}
