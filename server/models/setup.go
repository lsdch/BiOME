package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go"
	log "github.com/sirupsen/logrus"
)

func ConnectDB() (db *edgedb.Client) {
	ctx := context.Background()
	db, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %+v", err)
	}

	return
}

var DB *edgedb.Client = ConnectDB()

type Expr interface {
	String() string
	Args() map[string]interface{}
}

type QueryFilter struct {
	Field string
	Op    string
	Arg   string
	Value interface{}
}

func (filter *QueryFilter) String() string {
	return fmt.Sprintf("%s %s %s", filter.Field, filter.Op, filter.Arg)
}

func (filter *QueryFilter) Args() map[string]interface{} {
	key := strings.Split(filter.Arg, "$")[1]
	return map[string]interface{}{key: filter.Value}
}

type FilterGroup struct {
	Operator   string
	Components []Expr
}

func (group *FilterGroup) ComponentStrings() []string {
	var res []string
	for _, comp := range group.Components {
		if comp != nil {
			res = append(res, comp.String())
		}
	}
	return res
}

func (group *FilterGroup) String() string {
	components := group.ComponentStrings()
	return strings.Join(components, fmt.Sprintf(" %s ", group.Operator))
}

func (group *FilterGroup) Args() map[string]interface{} {
	args := make(map[string]interface{})
	for _, component := range group.Components {
		if component != nil {
			for k, v := range component.Args() {
				args[k] = v
			}
		}
	}
	return args
}

func (group *FilterGroup) Add(expr Expr) *FilterGroup {
	group.Components = append(group.Components, expr)
	return group
}

type QueryBuilder struct {
	Query string
	Expr  Expr
}

func (qb *QueryBuilder) String() string {
	filter_str := qb.Expr.String()
	if filter_str != "" {
		return fmt.Sprintf("%s filter %s", qb.Query, qb.Expr.String())
	} else {
		return qb.Query
	}
}

func (qb *QueryBuilder) Args() map[string]interface{} {
	return qb.Expr.Args()
}
