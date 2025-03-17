package geoapify

import (
	"context"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
)

const CREDIT_LIMIT = 3_000

type GeoapifyUsage struct {
	ID       geltypes.UUID      `gel:"id" json:"id"`
	Date     geltypes.LocalDate `gel:"date" json:"date"`
	Requests int32              `gel:"requests" json:"requests"`
}

func GeoapifyUsageList(e geltypes.Executor) (usages []GeoapifyUsage, err error) {
	err = e.Query(context.Background(),
		`#edgeql
			select admin::GeoapifyUsage { * } order by .date desc
		`, &usages)
	return
}

func TodayGeoapifyUsage(e geltypes.Executor) (current GeoapifyUsage, err error) {
	err = e.QuerySingle(context.Background(),
		`#edgeql
			select admin::GeoapifyUsage { * }
			filter .date = cal::to_local_date(datetime_of_transaction(), 'UTC')
		`, &current)
	if db.IsNoData(err) {
		year, month, day := time.Now().UTC().Date()
		return GeoapifyUsage{Date: geltypes.NewLocalDate(year, month, day)}, nil
	}
	return
}

func TrackGeoapifyUsage(e geltypes.Executor, requests int32) (current GeoapifyUsage, err error) {
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
