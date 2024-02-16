package validations

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type BindingEdgeDB struct {
	ObjectName   string // e.g. 'people::Institution'
	PropertyName string // e.g. 'code'
	TypeCast     string // e.g. 'str'
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
