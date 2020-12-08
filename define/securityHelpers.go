/* functions to check signatures and verify data correctness */
package define

import (
	//b64 "encoding/base64"
	"bytes"
	ed "crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/rogpeppe/godef/go/ast"
	"strconv"
	"time"
)

const MaxTimeBetweenLAOCreationAndPublish = 600

// TODO if we use the json Schema, don't need to check structure correctness
func LAOCreatedIsValid(data DataCreateLAO, message Message) error {
	//the last modified timestamp is equal to the creation timestamp,
	if data.Creation != data.Last_modified {
		fmt.Printf("%v, %v", data, data.Last_modified)
		fmt.Printf("sec1")
		return ErrInvalidResource
	}
	//the timestamp is reasonably recent with respect to the server’s clock,
	if data.Creation > time.Now().Unix() || data.Creation-time.Now().Unix() > MaxTimeBetweenLAOCreationAndPublish {
		fmt.Printf("sec2")
		return ErrInvalidResource
	}
	//the attestation is valid,
	str := []byte(data.Organizer)
	str = append(str, []byte(strconv.FormatInt(data.Creation, 10))...)
	str = append(str, []byte(data.Name)...)
	hash := sha256.Sum256(str)
	//hash64 := b64.StdEncoding.EncodeToString(hash[:])

	if !bytes.Equal([]byte(data.ID), hash[:]) {
		//if(hash64 != data.ID) {
		fmt.Printf("sec3 \n")
		fmt.Printf("%v, %v", hash, data.ID)
		return ErrInvalidResource
	}

	return nil
}

func MeetingCreatedIsValid(data DataCreateMeeting, message Message) error {
	//the last modified timestamp is equal to the creation timestamp,
	if data.Creation != data.Last_modified {
		return ErrInvalidResource
	}
	//the timestamp is reasonably recent with respect to the server’s clock,
	if data.Creation > time.Now().Unix() || data.Creation-time.Now().Unix() > MaxTimeBetweenLAOCreationAndPublish {
		return ErrInvalidResource
	}

	//we start after the creation and we end after the start
	if data.Start < data.Creation || data.End < data.Start {
		return ErrInvalidResource
	}
	//need to meet some	where
	if data.Location == "" {
		return ErrInvalidResource
	}
	return nil
}

func PollCreatedIsValid(data DataCreatePoll, message Message) error {
	return nil
}

func RollCallCreatedIsValid(data DataCreateRollCall, message Message) error {
	return nil
}

func MessageIsValid(msg Message) error {
	// the message_id is valid
	str := []byte(msg.Data)
	str = append(str, []byte(msg.Signature)...)
	hash := sha256.Sum256(str)

	if !bytes.Equal([]byte(msg.Message_id), hash[:]) {
		return ErrInvalidResource
	}

	// the signature is valid
	err := VerifySignature(msg.Sender, msg.Data, msg.Signature)
	if err != nil {
		return err
	}

	// the witness signatures are valid (check on every message??)
	return VerifyWitnessSignatures()
}

/*
	we check that Sign(sender||data) is the given signature
*/
func VerifySignature(publicKey string, data []byte, signature string) error {
	//check the size of the key as it will panic if we plug it in Verify
	if len(publicKey) != ed.PublicKeySize {
		return ErrRequestDataInvalid
	}
	//check the validity of the signature
	//TODO method is defined supposing args are encrypted
	//the key is in base64 so we need to decrypt it before using it
	keyDecoded, err := Decode(publicKey)
	if err != nil {
		return ErrEncodingFault
	}
	//data is also in base64 so we need to decrypt it before using it
	dataDecoded, err := Decode(string(data))
	if err != nil {
		return ErrEncodingFault
	}
	if ed.Verify(keyDecoded, dataDecoded, []byte(signature)) {
		return nil
	}
	//invalid signature
	return ErrRequestDataInvalid
}

//TODO be careful about the size and the order !
/*Maybe have a fixed size byte ?
To handle checks while the slice is in construction, the slice must have full space
from the beginning. We should check how to create fixed length arrays in go. And
instead of appending in witness_message, put them in the slot which matches the slot
of the witness id in witness[]

	Witness[1,2,3...]
	witnessSignature[_,_,_./.]
	WitnessSignatures[3,6,2,1]
*/
func VerifyWitnessSignatures(publicKeys []byte, signatures []byte, data string, sender string) error {
	senderDecoded, err := Decode(sender)
	if err != nil {
		return ErrEncodingFault
	}
	dataDecoded, err := Decode(data)
	if err != nil {
		return ErrEncodingFault
	}
	toCheck := append(senderDecoded, dataDecoded...)
	for i := 0; i < len(signatures); i++ {
		err := VerifySignature(string(publicKeys[i]), toCheck, string(signatures[i]))
		if err != nil {
			return err
		}
	}
	return nil
}