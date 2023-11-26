package mongo

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type Model interface {
	GetID() primitive.ObjectID
}

type Translator interface {
	Translate(isZHLang bool)
}

type JSONTime time.Time

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return ([]byte)(strconv.FormatInt(time.Time(*t).Unix(), 10)), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	num, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*t = JSONTime(time.Unix(int64(num), 0))
	return
}

func (t *JSONTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	retByte := make([]byte, 0)
	retByte = bsoncore.AppendTime(retByte, time.Time(*t))
	return bsontype.DateTime, retByte, nil
}

func (t *JSONTime) UnmarshalBSONValue(ty bsontype.Type, data []byte) error {
	if ty == bsontype.DateTime {
		if d, _, ok := bsoncore.ReadTime(data); ok {
			*t = JSONTime(d)
		}
	}
	return nil
}

type BaseModel struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreateTime JSONTime           `json:"create_time" bson:"create_time"`
	UpdateTime JSONTime           `json:"update_time" bson:"update_time"`
	DeleteTime JSONTime           `json:"-" bson:"delete_time"`

	// soft delete
	IsDeleted bool `json:"-" bson:"is_deleted"`

	Creator primitive.ObjectID `json:"creator" bson:"creator"`
}

func (m *BaseModel) GetID() primitive.ObjectID {
	return m.ID
}

func NewBaseModel(creator string) (BaseModel, error) {
	now := JSONTime(time.Now())
	m := BaseModel{
		CreateTime: now,
		UpdateTime: now,
		DeleteTime: JSONTime{},
		IsDeleted:  false,
	}
	if creator != "" {
		creatorObj, err := primitive.ObjectIDFromHex(creator)
		if err != nil {
			return BaseModel{}, err
		}
		m.Creator = creatorObj
	}

	return m, nil
}

func IDsToObjIDs(ids []string) ([]primitive.ObjectID, error) {
	var res []primitive.ObjectID
	for _, id := range ids {
		obj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func ObjIDsToStrings(ids []primitive.ObjectID) []string {
	var res []string
	for _, id := range ids {
		res = append(res, id.Hex())
	}
	return res
}
