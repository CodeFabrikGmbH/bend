package boltDB

import (
	"code-fabrik.com/bend/domain/request"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
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

func (rr RequestRepository) Add(req request.Request) (request.Request, error) {
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
	return req, err
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

// GetSummariesForPath returns a lightweight projection (ID + timestamp) of the
// requests for a path. It unmarshals only the summary fields, so bodies, headers
// and responses are never read into memory for list views.
func (rr RequestRepository) GetSummariesForPath(path string) []request.Summary {
	db := rr.DB

	var summaries []request.Summary

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName(path)))
		if bucket != nil {
			return bucket.ForEach(func(name []byte, value []byte) error {
				summary := request.Summary{}
				_ = json.Unmarshal(value, &summary)
				summaries = append(summaries, summary)
				return nil
			})
		}
		return nil
	})

	return summaries
}

// GetPathCounts returns the request count per path, computed parse-free in a
// single read transaction instead of one transaction and full scan per path.
func (rr RequestRepository) GetPathCounts() map[string]int {
	db := rr.DB

	counts := make(map[string]int)

	_ = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucketName := string(name)
			if strings.Index(bucketName, bucketPrefix) != 0 {
				return nil
			}
			path := strings.TrimPrefix(bucketName, bucketPrefix)

			count := 0
			cursor := b.Cursor()
			for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
				count++
			}
			counts[path] = count
			return nil
		})
	})

	return counts
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

func (rr RequestRepository) DeletePath(path string) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		return txn.DeleteBucket([]byte(bucketName(path)))

	})
}

func (rr RequestRepository) DeleteRequestForPath(path string, id string) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(bucketName(path)))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(id))
	})
}
