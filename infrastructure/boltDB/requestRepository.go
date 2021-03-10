package boltDB

import (
	"code-fabrik.com/bend/domain/request"
	"encoding/json"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

const bucketPrefix = "requestPath"

type RequestRepository struct {
	DB *bolt.DB
}

func bucketName(path string) string {
	return bucketPrefix + path
}

func (rr RequestRepository) Save(req request.Request) error {
	db := rr.DB

	err := db.Update(func(txn *bolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists([]byte(bucketName(req.Path)))
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		req.ID = strconv.Itoa(int(id))

		data, err := json.Marshal(req)
		if err != nil {
			return err
		}
		err = b.Put([]byte(req.ID), data)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (rr RequestRepository) GetRequest(path string, id string) request.Request {
	db := rr.DB

	result := request.Request{}

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName(path)))

		if bucket != nil {
			value := bucket.Get([]byte(id))
			if value == nil {
				return nil
			}

			json.Unmarshal(value, &result)
		}
		return nil
	})

	return result
}

func (rr RequestRepository) GetRequestsForPath(path string) []request.Request {
	db := rr.DB

	var requests []request.Request

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName(path)))
		if bucket != nil {
			bucket.ForEach(func(name []byte, value []byte) error {
				req := request.Request{}
				json.Unmarshal(value, &req)
				requests = append(requests, req)
				return nil
			})
		}

		return nil
	})

	return requests
}

func (rr RequestRepository) GetRequestCountForPath(path string) int {
	db := rr.DB
	var result int

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName(path)))

		if bucket != nil {
			bucket.ForEach(func(name []byte, value []byte) error {
				result++
				return nil
			})
		}

		return nil
	})

	return result
}

func (rr RequestRepository) GetPaths() []string {
	db := rr.DB

	var result []string

	err := db.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucketName := string(name)
			if strings.Index(bucketName, bucketPrefix) == 0 {
				path := strings.TrimPrefix(bucketName, bucketPrefix)
				result = append(result, path)
			}
			return nil
		})

		return err
	})

	if err != nil {
		panic(err)
	}
	return result
}
