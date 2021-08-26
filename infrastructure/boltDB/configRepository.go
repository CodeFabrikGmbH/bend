package boltDB

import (
	"code-fabrik.com/bend/domain/config"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
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
		err = b.Put([]byte(config.Id.String()), data)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (rr ConfigRepository) Find(id uuid.UUID) *config.Config {
	db := rr.DB

	result := config.Config{}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(configBucket))

		if bucket != nil {
			value := bucket.Get([]byte(id.String()))
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
				_ = json.Unmarshal(value, &v)

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

func (rr ConfigRepository) Delete(id uuid.UUID) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(configBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		err := bucket.Delete([]byte(id.String()))
		return err
	})
}

func (rr ConfigRepository) DeleteAll() error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		return txn.DeleteBucket([]byte(configBucket))

	})
}
