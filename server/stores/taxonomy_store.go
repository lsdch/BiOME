package stores

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/models"
)

type TaxonomyStore struct {
}

func NewTaxonomyStore() *TaxonomyStore {
	return &TaxonomyStore{}
}

func (s *TaxonomyStore) SearchTaxa(ctx context.Context, q db.Querier, params models.ListTaxaParams) ([]models.Taxon, error) {
	stmt := SELECT(
		table.Taxa.AllColumns,
	)

	if params.SampledOnly {
		stmt = stmt.FROM(table.Taxa.INNER_JOIN(table.Occurrences, table.Taxa.ID.EQ(table.Occurrences.TaxonID)))
	} else {
		stmt = stmt.FROM(table.Taxa)
	}

	stmt = applyRanksFilter(params, stmt)
	stmt = applySearchFilter(stmt, params.SearchTerm)

	if params.Pagination.Limit > 0 {
		stmt = stmt.LIMIT(int64(params.Pagination.Limit))
	}
	stmt = stmt.OFFSET(int64(params.Pagination.Offset))

	sql, args := stmt.Sql()
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxa = []models.Taxon{}
	for rows.Next() {
		var t models.Taxon
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Authorship,
			&t.Rank,
			&t.Status,
			&t.GBIF_ID,
		)
		if err != nil {
			return nil, err
		}
		taxa = append(taxa, t)
	}

	return taxa, nil
}

func applyRanksFilter(params models.ListTaxaParams, stmt SelectStatement) SelectStatement {
	if len(params.Ranks) > 0 {
		ranks := make([]string, len(params.Ranks))
		for i, r := range params.Ranks {
			ranks[i] = string(r)
		}
		stmt = stmt.WHERE(table.Taxa.Rank.IN(StringArray(ranks...)))
	}
	return stmt
}

func applySearchFilter(stmt SelectStatement, searchTerm models.Optional[string]) SelectStatement {
	term, isSet := searchTerm.Get()
	if !isSet || len(term) < 3 {
		return stmt
	}

	if len(term) < 5 {
		return stmt.WHERE(table.Taxa.Name.LIKE(String(term + "%")))
	}

	return stmt.WHERE(TrgmMatch(table.Taxa.Name, String(term)))

}

func TrgmMatch(col Expression, value Expression) BoolExpression {
	return BoolExp(
		Raw("#arg1 % #arg2", map[string]interface{}{"#arg1": col, "#arg2": value}),
	)
}
