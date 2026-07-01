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

// countsBucket holds an incrementally maintained request count per path so the
// dashboard never has to scan the whole (multi-GB) database to build the
// endpoint list. flagsBucket records one-off migration flags.
const (
	countsBucket = "__meta_pathCounts"
	flagsBucket  = "__meta_flags"
	backfillFlag = "countsBackfilled"
)

type RequestRepository struct {
	DB *bolt.DB
}

func bucketName(path string) string {
	return bucketPrefix + path
}

// adjustPathCount changes the stored request count for a path by delta within an
// existing write transaction. The entry is removed once it reaches zero.
func adjustPathCount(tx *bolt.Tx, path string, delta int) error {
	b, err := tx.CreateBucketIfNotExists([]byte(countsBucket))
	if err != nil {
		return err
	}
	count := 0
	if v := b.Get([]byte(path)); v != nil {
		count, _ = strconv.Atoi(string(v))
	}
	count += delta
	if count <= 0 {
		return b.Delete([]byte(path))
	}
	return b.Put([]byte(path), []byte(strconv.Itoa(count)))
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
		return adjustPathCount(txn, req.Path, 1)
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

// GetSummariesPage returns a page of request summaries for a path, newest first.
// Requests are keyed by a monotonically increasing sequence, so descending
// sequence order is also newest-first chronological order. before is an
// exclusive upper bound on the sequence id (pass 0 for the newest page); at most
// limit summaries are returned. The bool reports whether older requests remain,
// enabling incremental (infinite-scroll) loading without reading the whole
// bucket into memory.
func (rr RequestRepository) GetSummariesPage(path string, before int, limit int) ([]request.Summary, bool) {
	db := rr.DB

	if limit <= 0 {
		return nil, false
	}

	var summaries []request.Summary
	hasMore := false

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName(path)))
		if bucket == nil {
			return nil
		}

		start := before - 1
		if before <= 0 {
			start = int(bucket.Sequence())
		}

		for id := start; id >= 1; id-- {
			value := bucket.Get([]byte(strconv.Itoa(id)))
			if value == nil {
				continue
			}
			if len(summaries) >= limit {
				hasMore = true
				break
			}
			summary := request.Summary{}
			_ = json.Unmarshal(value, &summary)
			summaries = append(summaries, summary)
		}
		return nil
	})

	return summaries, hasMore
}

// GetPathCounts returns the request count per path. Counts are read from the
// incrementally maintained countsBucket, so this stays cheap even for a database
// with hundreds of thousands of requests across tens of thousands of paths (a
// full scan would need to read the entire multi-GB file from disk).
func (rr RequestRepository) GetPathCounts() map[string]int {
	db := rr.DB

	counts := make(map[string]int)

	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(countsBucket))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			n, _ := strconv.Atoi(string(v))
			counts[string(k)] = n
			return nil
		})
	})

	return counts
}

// BackfillPathCounts populates countsBucket from a one-time full scan of the
// database. It is a no-op once the backfill flag is set, so it only pays the
// expensive scan on the first start after this feature is deployed. Intended to
// run in a background goroutine at start-up.
func (rr RequestRepository) BackfillPathCounts() error {
	db := rr.DB

	already := false
	_ = db.View(func(tx *bolt.Tx) error {
		if fb := tx.Bucket([]byte(flagsBucket)); fb != nil && fb.Get([]byte(backfillFlag)) != nil {
			already = true
		}
		return nil
	})
	if already {
		return nil
	}

	// Snapshot the per-path counts with a read transaction (does not block
	// request tracking). Requests added during the scan keep incrementing
	// countsBucket via Add and are reconciled by the absolute write below.
	snapshot := make(map[string]int)
	if err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			n := string(name)
			if !strings.HasPrefix(n, bucketPrefix) {
				return nil
			}
			path := strings.TrimPrefix(n, bucketPrefix)
			count := 0
			cursor := b.Cursor()
			for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
				count++
			}
			snapshot[path] = count
			return nil
		})
	}); err != nil {
		return err
	}

	// Write the scanned counts as absolute values and set the flag atomically so
	// an interrupted backfill simply re-runs cleanly on the next start.
	return db.Update(func(tx *bolt.Tx) error {
		cb, err := tx.CreateBucketIfNotExists([]byte(countsBucket))
		if err != nil {
			return err
		}
		for path, count := range snapshot {
			if count <= 0 {
				continue
			}
			if err := cb.Put([]byte(path), []byte(strconv.Itoa(count))); err != nil {
				return err
			}
		}
		fb, err := tx.CreateBucketIfNotExists([]byte(flagsBucket))
		if err != nil {
			return err
		}
		return fb.Put([]byte(backfillFlag), []byte("1"))
	})
}

// GetRequestCountForPath returns the number of stored requests for a path. It
// reads the incrementally maintained countsBucket so it stays cheap regardless
// of how many requests the path holds.
func (rr RequestRepository) GetRequestCountForPath(path string) int {
	db := rr.DB
	result := 0

	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(countsBucket))
		if b == nil {
			return nil
		}
		if v := b.Get([]byte(path)); v != nil {
			result, _ = strconv.Atoi(string(v))
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
		if err := txn.DeleteBucket([]byte(bucketName(path))); err != nil {
			return err
		}
		if cb := txn.Bucket([]byte(countsBucket)); cb != nil {
			return cb.Delete([]byte(path))
		}
		return nil
	})
}

func (rr RequestRepository) DeleteRequestForPath(path string, id string) error {
	db := rr.DB

	return db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(bucketName(path)))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		if bucket.Get([]byte(id)) == nil {
			return nil
		}
		if err := bucket.Delete([]byte(id)); err != nil {
			return err
		}
		return adjustPathCount(txn, path, -1)
	})
}
