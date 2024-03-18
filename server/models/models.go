package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/edgedb/edgedb-go"
	"golang.org/x/exp/maps"
)

type UserShortIdentity struct {
	edgedb.Optional
	Name  string `json:"name" edgedb:"name"`
	Alias string `json:"alias" edgedb:"alias"`
}

type Meta struct {
	ID          edgedb.UUID             `edgedb:"id" json:"-" swaggerignore:"true"`
	Created     time.Time               `edgedb:"created" json:"created" example:"2023-09-01T16:41:10.921097+00:00" binding:"required"`
	Modified    edgedb.OptionalDateTime `edgedb:"modified" json:"modified" example:"2023-09-02T20:39:10.218057+00:00"`
	LastUpdated time.Time               `edgedb:"lastUpdated" json:"last_updated" example:"2023-09-02T20:39:10.218057+00:00"`
	CreatedBy   UserShortIdentity       `json:"created_by" edgedb:"created_by"`
	UpdatedBy   UserShortIdentity       `json:"updated_by" edgedb:"updated_by"`
} // @name Meta

func StructToMap(val interface{}) map[string]interface{} {
	//The name of the tag you will use for fields of struct
	const tagTitle = "edgedb"

	var data map[string]interface{} = make(map[string]interface{})
	varType := reflect.TypeOf(val)
	if varType.Kind() != reflect.Struct {
		// Provided value is not an interface, do what you will with that here
		fmt.Println("Not a struct")
		return nil
	}

	value := reflect.ValueOf(val)
	for i := 0; i < varType.NumField(); i++ {
		if !value.Field(i).CanInterface() {
			//Skip unexported fields
			continue
		}
		tag, ok := varType.Field(i).Tag.Lookup(tagTitle)
		var fieldName string
		if ok && len(tag) > 0 {
			fieldName = tag
		} else {
			fieldName = varType.Field(i).Name
		}
		if varType.Field(i).Type.Kind() != reflect.Struct {
			data[fieldName] = value.Field(i).Interface()
		} else if fieldName == "$inline" {
			maps.Copy(data, StructToMap(value.Field(i).Interface()))
		} else {
			data[fieldName] = StructToMap(value.Field(i).Interface())
		}
	}
	return data
}
