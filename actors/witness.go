/* This file implements the functionalities a witness must provide. It uses the package db to interact with
a database.*/
package actors

import (
	b64 "encoding/base64"
	"encoding/json"
	"student20_pop/db"
	"student20_pop/define"
)

type Witness struct {
	PublicKey string
	database  string
}

func NewWitness(pkey string, db string) *Witness {
	return &Witness{
		PublicKey: pkey,
		database:  db,
	}
}

/*returns true if w is in the witness list of the event*/
func (w *Witness) IsWitness(id string) (bool, error) {
	data := db.GetChannel([]byte(id), w.database)
	lao := define.LAO{} //TODO currently is only for LAO. Need generic type for channel
	err := json.Unmarshal(data, &lao)
	if err != nil {
		return false, define.ErrEncodingFault
	}

	_, found := define.FindStr(lao.Witnesses, w.PublicKey)

	return found, nil
}

/** processes what is received from the WebSocket
 * Currently only supports updateProperties
 * msg : receivedMessage
 * returns, in order :
 * message to send on channel, or nil
 * channel for the message, or nil
 * response to the sender, or nil
 */
func (w *Witness) HandleWholeMessage(msg []byte, userId int) ([]byte, []byte, []byte) {
	generic, err := define.AnalyseGeneric(msg)
	if err != nil {
		return nil, nil, define.CreateResponse(define.ErrRequestDataInvalid, nil, generic)

	}

	var history []byte = nil
	var message []byte = nil
	var channel []byte = nil

	switch generic.Method {
	case "publish":
		message, channel, err = w.handlePublish(generic)
	default:
		message, channel, err = nil, nil, define.ErrRequestDataInvalid
	}

	return message, channel, define.CreateResponse(err, history, generic)
}

/** @returns, in order
 * message
 * channel
 * error
 */
func (w *Witness) handlePublish(generic define.Generic) ([]byte, []byte, error) {
	params, err := define.AnalyseParamsFull(generic.Params)
	if err != nil {
		return nil, nil, define.ErrRequestDataInvalid
	}

	message, err := define.AnalyseMessage(params.Message)
	if err != nil {
		return nil, nil, define.ErrRequestDataInvalid
	}

	data := define.Data{}
	base64Text := make([]byte, b64.StdEncoding.DecodedLen(len(message.Data)))
	l, _ := b64.StdEncoding.Decode(base64Text, message.Data)
	err = json.Unmarshal(base64Text[:l], &data)

	//data, err := define.AnalyseData(message.Data)
	if err != nil {
		return nil, nil, define.ErrRequestDataInvalid
	}

	switch data["object"] {
	case "lao":
		switch data["action"] {
		case "create":
			return w.handleCreateLAO(message, params.Channel, generic)
		case "update_properties":
			return w.handleUpdateProperties(message, params.Channel, generic)
		case "state":
			// just store in DB
		default:
			return nil, nil, define.ErrInvalidAction
		}

	case "message":
		switch data["action"] {
		case "witness":
			return w.handleWitnessMessage(message, params.Channel, generic)
		default:
			return nil, nil, define.ErrInvalidAction
		}
	case "roll call":
		switch data["action"] {
		case "create":
			//return w.handleCreateRollCall(message, params.Channel, generic)
		case "state":

		default:
			return nil, nil, define.ErrInvalidAction
		}
	case "meeting":
		switch data["action"] {
		case "create":
			//return w.handleCreateMeeting(message, params.Channel, generic)
		case "state":

		default:
			return nil, nil, define.ErrInvalidAction
		}
	case "poll":
		switch data["action"] {
		case "create":
			//return w.handleCreatePoll(message, params.Channel, generic)
		case "state":

		default:
			return nil, nil, define.ErrInvalidAction
		}
	default:
		return nil, nil, define.ErrRequestDataInvalid
	}

	return nil, nil, nil
}

func (w *Witness) handleCreateLAO(message define.Message, channel string, generic define.Generic) ([]byte, []byte, error) {
	if channel != "/root" {
		return nil, nil, define.ErrInvalidResource
	}

	data, err := define.AnalyseDataCreateLAO(message.Data)
	if err != nil {
		return nil, nil, define.ErrInvalidResource
	}

	err = define.LAOCreatedIsValid(data, message)
	if err != nil {
		return nil, nil, err
	}

	canalLAO := channel + data.ID

	err = db.CreateMessage(message, canalLAO, w.database)
	if err != nil {
		return nil, nil, err
	}

	lao := define.LAO{
		ID:            data.ID,
		Name:          data.Name,
		Creation:      data.Creation,
		LastModified:  data.Last_modified,
		OrganizerPKey: data.Organizer,
		Witnesses:     data.Witnesses,
	}
	err = db.CreateChannel(lao, w.database)

	return nil, nil, err
}

/*witness does not yet send stuff to channel*/
func (w *Witness) handleUpdateProperties(message define.Message, channel string, generic define.Generic) ([]byte, []byte, error) {
	data, err := define.AnalyseDataCreateLAO(message.Data)
	if err != nil {
		return nil, nil, define.ErrInvalidResource
	}
	err = define.LAOCreatedIsValid(data, message)
	if err != nil {
		return nil, nil, err
	}

	//stores received message in DB
	canalLAO := channel + data.ID
	err = db.CreateMessage(message, canalLAO, w.database)
	if err != nil {
		return nil, nil, err
	}

	//toSign := message.Sender + string(message.Data)

	//TODO create a response signing the message

	return nil, nil, err
}

func (w *Witness) handleWitnessMessage(message define.Message, channel string, generic define.Generic) ([]byte, []byte, error) {

	//shall a witness increment count on base message as well ?

	data, err := define.AnalyseDataCreateLAO(message.Data)
	if err != nil {
		return nil, nil, define.ErrInvalidResource
	}
	err = define.LAOCreatedIsValid(data, message)
	if err != nil {
		return nil, nil, err
	}

	//stores received message in DB
	canalLAO := channel + data.ID
	err = db.CreateMessage(message, canalLAO, w.database)
	return nil, nil, err
}