package db

import (
	"github.com/boltdb/bolt"
)

func OpenDB(dbName string) (*bolt.DB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

/**
 * Functions that writes a pair (key, val) in the bucket "bkt" in the database
 */
func Write(key []byte, val []byte, bkt []byte, database *bolt.DB) error {
	err := database.Update(func(tx *bolt.Tx) error {
		b, err1 := tx.CreateBucketIfNotExists(bkt)
		if err1 != nil {
			return err1
		}
		err2 := b.Put(key, val)
		if err2 != nil {
			return err2
		}
		return nil
	})
	return err
}

/**
 * Functions that transforms a nested bucket of the db (of members, events, etc.) into a list
 * (of byte array currently, tbd according to the definite field of the structs)
 * ça va surement tej'
 */
func NestedToList(b *bolt.Bucket) [][]byte {
	var list [][]byte
	b.ForEach(func(k, v []byte) error {
		list = append(list, v)
		return nil
	})

	return list
}
