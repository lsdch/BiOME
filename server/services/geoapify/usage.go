package geoapify

import (
	"context"
	"darco/proto/db"
	"time"

	"github.com/edgedb/edgedb-go"
)

type GeoapifyUsage struct {
	ID       edgedb.UUID      `edgedb:"id" json:"id"`
	Date     edgedb.LocalDate `edgedb:"date" json:"date"`
	Requests int32            `edgedb:"requests" json:"requests"`
}

func GeoapifyUsageList(e edgedb.Executor) (usages []GeoapifyUsage, err error) {
	err = e.Query(context.Background(),
		`#edgeql
			select admin::GeoapifyUsage { * } order by .date desc
		`, &usages)
	return
}

func TodayGeoapifyUsage(e edgedb.Executor) (current GeoapifyUsage, err error) {
	err = e.QuerySingle(context.Background(),
		`#edgeql
			select admin::GeoapifyUsage { * }
			filter .date = cal::to_local_date(datetime_of_transaction(), 'UTC')
		`, &current)
	if db.IsNoData(err) {
		year, month, day := time.Now().UTC().Date()
		return GeoapifyUsage{Date: edgedb.NewLocalDate(year, month, day)}, nil
	}
	return
}

func TrackGeoapifyUsage(e edgedb.Executor, requests int32) (current GeoapifyUsage, err error) {
	err = e.QuerySingle(context.Background(),
		`#edgeql
			select (
				insert admin::GeoapifyUsage {
					date := cal::to_local_date(datetime_of_transaction(), 'UTC'),
					requests := <int32>$0
				} unless conflict on .date else (
					update admin::GeoapifyUsage set {
						requests := .requests + <int32>$0
					}
				)
			) { * }
		`, &current, requests)
	return
}
