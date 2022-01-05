package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type URL struct {
	URI string `bson:"uri"`
	Age int    `bson:"age"`
}

func (*URL) Name() string {
	return "url"
}

type Student struct {
	Class string `bson:"class"`
}

func (*Student) Name() string {
	return "student"
}

type Value interface {
	Name() string
}

type MyStruct struct {
	Image Value  `bson:"image" json:"image"`
	Type  string `bson:"type" json:"type"`
}

type GetType struct {
	Type string `bson:"type" json:"type"`
}

type H MyStruct

func (m *MyStruct) UnmarshalBSON(data []byte) error {
	var getType GetType
	if err := bson.Unmarshal(data, &getType); err != nil {
		return err
	}
	log.Printf("%+v\n", getType)
	// switch type
	switch getType.Type {
	case "url":
		aux := &struct {
			Image *URL             `bson:"image"`
			*H    `bson:",inline"` // inline 是重点
		}{
			H: (*H)(m),
		}
		if err := bson.Unmarshal(data, &aux); err != nil {
			log.Printf("err: %+v\n", err)
			return err
		}
		log.Printf("res: %+v\n", aux.Image)
		log.Printf("res: %+v\n", aux.H)
		m.Image = aux.Image
	default:
		return fmt.Errorf("no implement")
	}

	return nil
}

func main() {
	m := &MyStruct{
		Image: &URL{
			URI: "foobar",
			Age: 18,
		},
		Type: "url",
	}
	res, err := bson.Marshal(m)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v\n",bson.Raw(res))
	fmt.Printf("%+v\n", bson.Raw(res))

	var unmarshalled MyStruct
	if err := bson.Unmarshal(res, &unmarshalled); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", unmarshalled)
	fmt.Printf("%+v\n", unmarshalled.Image.(*URL).URI)
}
