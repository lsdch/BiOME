package validations

import (
	"context"
	"darco/proto/db"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type BindingEdgeDB struct {
	ObjectName   string // e.g. 'people::Institution'
	PropertyName string // e.g. 'code'
	TypeCast     string // e.g. 'str'
}

func (b *BindingEdgeDB) UniqueQuery(args any) bool {
	query := fmt.Sprintf("select exists (select %s filter .%s = <%s>$0)", b.ObjectName, b.PropertyName, b.TypeCast)
	var exists bool
	err := db.Client().QuerySingle(context.Background(), query, &exists, args)
	if err != nil {
		logrus.Errorf("Unique validation query failed: %v with query %s", err, query)
	}
	return !exists
}

func ParseEdgeDBBindings(s string, typeCast string) BindingEdgeDB {
	var edgeDBBinding = strings.SplitN(s, ".", 2)
	if len(edgeDBBinding) != 2 {
		logrus.Errorf("Invalid field specification for 'exist_all' validator. Expected format like '{ObjectName}.{PropertyName}', got %s", s)
	}
	return BindingEdgeDB{
		ObjectName:   edgeDBBinding[0],
		PropertyName: edgeDBBinding[1],
		TypeCast:     typeCast,
	}
}
