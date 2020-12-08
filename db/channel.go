/* This file implements useful functions to work with channels in the database, such as Create/edit and GetInfos */
package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"student20_pop/define"
)

const bucketChannel = "channels"

/**
 * Function to create a new Object (LAO,Event...) and store it in the DB
 * @returns : error
 */
func writeChannel(obj interface{}, database string, secure bool) error {
	db, e := OpenDB(database)
	defer db.Close()
	if e != nil {
		return e
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b, err1 := tx.CreateBucketIfNotExists([]byte(bucketChannel))
		if err1 != nil {
			return err1
		}
		//generic adaptation
		var objID []byte
		switch obj.(type) {
		case define.LAO:
			// type assert
			objID = []byte(obj.(define.LAO).ID)
		case define.Meeting:
			objID = []byte(obj.(define.Meeting).ID)
		case define.Poll:
			objID = []byte(obj.(define.Poll).ID)
		case define.RollCall:
			objID = []byte(obj.(define.RollCall).ID)
		default:
			//TODO not sure for the error type
			return define.ErrRequestDataInvalid
		}
		//checks if there is already an entry with that ID if secure is true
		if secure {
			key := b.Get(objID)
			if key != nil {
				return define.ErrResourceAlreadyExists
			}
		} else {
			exists := b.Get(objID)
			if exists == nil {
				fmt.Printf("11")
				return define.ErrInvalidResource
			}
		}
		var dt []byte
		var err2 error
		switch obj.(type) {
		case define.LAO:
			// type assert
			dt, err2 = json.Marshal(obj.(define.LAO).ID)
		case define.Meeting:
			dt, err2 = json.Marshal(obj.(define.Meeting).ID)
		case define.Poll:
			dt, err2 = json.Marshal(obj.(define.Poll).ID)
		case define.RollCall:
			dt, err2 = json.Marshal(obj.(define.RollCall).ID)
		default:
			//TODO not sure for the error type
			return define.ErrRequestDataInvalid
		}
		// Marshal the Obj and store it
		if err2 != nil {
			return define.ErrRequestDataInvalid
		}
		err3 := b.Put(objID, dt)
		return err3
	})

	return err
}

/*writes a channel (LAO, meeting, rollCall, etc.) to the DB, returns an error if ID already is key in DB*/
func CreateChannel(obj interface{}, database string) error {
	return writeChannel(obj, database, true)
}

/*writes a channel (LAO, meeting, rollCall, etc.) to the DB, only if ID already exists, otherwise return an error*/
func UpdateChannel(obj interface{}, database string) error {
	return writeChannel(obj, database, false)
}

/*returns channel data from a given ID */
func GetChannel(id []byte, database string) []byte {
	db, e := OpenDB(database)
	defer db.Close()
	if e != nil {
		return nil
	}
	var data []byte
	e = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketChannel))
		data = b.Get(id)
		return nil
	})
	return data
}