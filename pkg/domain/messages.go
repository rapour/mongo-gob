package domain

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gob.Register(Email{})
	gob.Register(Text{})
}

type MessageGob struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data []byte             `bson:"data"`
}

func (g MessageGob) GetMessage() (Message, error) {

	var msg Message
	dec := gob.NewDecoder(bytes.NewBuffer(g.Data))

	err := dec.Decode(&msg)

	return msg, err
}

type Message interface {
	GetGob() (MessageGob, error)
	GetType() string
	Marshal() ([]byte, error)
}

type Email struct {
	Title  string
	Conent string
}

func (e Email) GetType() string {
	return "Email"
}

func (e Email) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e Email) GetGob() (MessageGob, error) {

	var buf bytes.Buffer
	var container Message

	enc := gob.NewEncoder(&buf)

	container = e
	err := enc.Encode(&container)
	if err != nil {
		return MessageGob{}, err
	}

	return MessageGob{
		ID:   primitive.NewObjectID(),
		Data: buf.Bytes(),
	}, nil
}

type Text struct {
	Content string
}

func (t Text) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func (t Text) GetType() string {
	return "Text"
}

func (t Text) GetGob() (MessageGob, error) {

	var buf bytes.Buffer
	var container Message

	enc := gob.NewEncoder(&buf)

	container = t
	err := enc.Encode(&container)
	if err != nil {
		return MessageGob{}, err
	}

	return MessageGob{
		ID:   primitive.NewObjectID(),
		Data: buf.Bytes(),
	}, nil
}
