package models

import "github.com/edgedb/edgedb-go"

type Creatable[Item Updatable[Item]] interface {
	Create(db *edgedb.Client) (Item, error)
}

type Updatable[T any] interface {
	Update(db *edgedb.Client) (T, error)
}

type ItemFinder[ID any, Item Updatable[Item]] func(db *edgedb.Client, id ID) (Item, error)
