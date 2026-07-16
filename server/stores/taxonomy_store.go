package stores

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
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
	logrus.Debugf("SearchTaxa SQL: %s, args: %v", sql, args)
	for i, arg := range args {
		logrus.Debugf("arg[%d]: %T %#v", i, arg, arg)
	}
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxa = []models.Taxon{}
	for rows.Next() {
		var t biomedb.Taxon
		err := rows.Scan(
			&t.ID,
			&t.GBIFID,
			&t.Name,
			&t.ScientificName,
			&t.Rank,
			&t.Status,
			&t.Authorship,
			&t.AcceptedTaxonID,
			&t.ParentID,
			&t.SearchVector,
			&t.Comments,
		)
		if err != nil {
			return nil, err
		}
		taxa = append(taxa, *models.TaxonFromDB(&t))
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
	if !isSet || len(term) < 2 {
		return stmt
	}

	return stmt.WHERE(OR(
		table.Taxa.Name.LIKE(String("%"+term+"%")),
		TrgmMatch(table.Taxa.Name, String(term)),
	))

}

func TrgmMatch(col Expression, value Expression) BoolExpression {
	return BoolExp(
		BinaryOperator(col, value, "%"),
	)
}
