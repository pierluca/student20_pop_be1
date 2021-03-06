package db

// This file contains functions to manage subscribers of one channel.
// THIS FILE IS NOW OBSOLETE AND UNUSED. DONT REMOVE UNTIL FULLY TESTED

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	"student20_pop/lib"
)

const bucketSubscribers = "sub"
const SubscribeDB = "sub.db"

// Subscribe is a function that subscribe a user to a channel. ONLY AT THE PUB/SUB LAYER
// DEPRECATED : Subscribers are not stored in a database anymore
// if user was already subscribed, does nothing
// does not change LAO's member field
func Subscribe(userId int, channelId []byte) error {

	db, err := OpenDB(SubscribeDB)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err1 := tx.CreateBucketIfNotExists([]byte(bucketSubscribers))
		if err1 != nil {
			return err1
		}
		var ints []int
		//gets the list of subscribers if exists, converts it to a list of int
		data := b.Get(channelId)
		if data != nil {
			err1 = json.Unmarshal(data, &ints)
			if err1 != nil {
				return err1
			}
		}

		//check if was already susbscribed
		if _, found := lib.Find(ints, userId); found {
			fmt.Println("user was already subscribed")
			return lib.ErrResourceAlreadyExists
		}
		ints = append(ints, userId)
		//converts []int to string to []byte
		data = []byte(strings.Trim(strings.Join(strings.Split(fmt.Sprint(ints), " "), ","), ""))
		//push values back
		err1 = b.Put(channelId, data)
		return err1
	})

	return err
}

// Unsubscribe is a function that unsubscribes a user from a channel. ONLY AT THE PUB/SUB LAYER
// DEPRECATED : Subscribers are not stored in a database anymore. Will be removed after full tests
// does nothing if that user was not already subscribed
// does not change LAO's member field
func Unsubscribe(userId int, channelId []byte) error {

	db, err := OpenDB(SubscribeDB)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err1 := tx.CreateBucketIfNotExists([]byte(bucketSubscribers))
		if err1 != nil {
			return err1
		}
		var ints []int
		//gets the list of subscribers if exists, converts it to a list of int
		data := b.Get(channelId)
		if data != nil {
			err1 = json.Unmarshal(data, &ints)
			if err1 != nil {
				return err1
			}
		}

		//check if was already susbscribed
		i, found := lib.Find(ints, userId)
		if !found {
			fmt.Println("this user was not subscribed to this channel")
			return lib.ErrInvalidResource
		}
		//remove elem from array
		ints[i] = ints[len(ints)-1]
		ints = ints[:len(ints)-1]

		//converts []int to string to []byte
		data = []byte(strings.Trim(strings.Join(strings.Split(fmt.Sprint(ints), " "), ","), ""))
		//push values back
		err1 = b.Put(channelId, data)
		return err1
	})

	return err
}

// GetSubscribers is a helper function to find a channel's subscribers
// DEPRECATED : Subscribers are not stored in a database anymore
func GetSubscribers(channel []byte) ([]int, error) {
	db, err := OpenDB(SubscribeDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var data []int

	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketSubscribers))
		if bkt == nil {
			return nil
		}

		content := bkt.Get(channel)
		err = json.Unmarshal(content, &data)

		return err
	})

	return data, nil
}
