package users

import "github.com/edgedb/edgedb-go"

type PersonInput struct {
	FirstName string             `json:"first_name" edgedb:"first_name" validate:"required,alpha"`
	LastName  string             `json:"last_name" edgedb:"last_name" validate:"required,alpha"`
	Contact   edgedb.OptionalStr `json:"contact" edgedb:"contact"`
} // @name PersonInput

type Person struct {
	ID          edgedb.UUID `edgedb:"id" json:"id"`
	PersonInput `edgedb:"$inline"`
	FullName    string `json:"full_name" edgedb:"full_name"`
} // @name Person
