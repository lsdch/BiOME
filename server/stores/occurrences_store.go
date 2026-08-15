package stores

import (
	"context"
	"fmt"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
	"github.com/uber/h3-go/v4"
)

const MIN_FULL_TEXT_SEARCH_TERM_LENGTH = 3

type OccurrenceStore struct {
}

func NewOccurrenceStore() *OccurrenceStore {
	return &OccurrenceStore{}
}

func (s *OccurrenceStore) OccurrenceCoreTables() ReadableTable {
	return table.Occurrences.
		INNER_JOIN(table.Taxa,
			table.Occurrences.TaxonID.EQ(table.Taxa.ID)).
		INNER_JOIN(table.Samplings,
			table.Occurrences.SamplingID.EQ(table.Samplings.ID)).
		LEFT_JOIN(table.Countries,
			table.Samplings.SiteCountryCode.EQ(table.Countries.Code))

}

func (s *OccurrenceStore) H3CellFilter(resolution int64, cell h3.Cell) BoolExpression {
	resolutionExpr := RawInt(fmt.Sprintf("%d", resolution))
	h3CellEq := table.Samplings.H3Index.EQ(Int(int64(cell)))
	if resolution < 12 {
		h3CellEq = CAST(Func("h3_cell_to_parent", CAST(table.Samplings.H3Index).AS("h3index"), resolutionExpr)).AS_BIGINT().EQ(Int(int64(cell)))
	}
	return h3CellEq
}

func (s *OccurrenceStore) ListOccurringTaxaAtCell(ctx context.Context, q db.Querier, cell h3.Cell, resolution int64, params ListOccurrencesParams) ([]models.OccurrenceOverviewItem, error) {
	h3CellFilter := s.H3CellFilter(resolution, cell)

	countStmt := SELECT(
		table.Taxa.ID.AS("taxon_id"),
		COUNT(table.Occurrences.ID).AS("occurrences"),
	).
		FROM(s.OccurrenceCoreTables()).
		GROUP_BY(table.Taxa.ID)
	countStmt = params.ApplyFilters(countStmt, h3CellFilter)

	counts := countStmt.AsTable("counts")

	countOccurrences := IntegerColumn("occurrences").From(counts)
	countTaxonID := StringColumn("taxon_id").From(counts)

	ancestor := table.Taxa.AS("ancestor")

	occurrencesExpr := MAX(
		CASE().
			WHEN(
				ancestor.ID.EQ(countTaxonID),
			).THEN(countOccurrences).
			ELSE(
				Int32(0),
			),
	).AS("occurrences")

	stmt := SELECT(
		ancestor.ID,
		ancestor.Name,
		ancestor.Authorship,
		ancestor.Rank,
		ancestor.ParentID,
		occurrencesExpr,
	).
		FROM(
			counts.
				INNER_JOIN(
					table.TaxaClosure,
					table.TaxaClosure.DescendantID.EQ(countTaxonID),
				).
				INNER_JOIN(
					ancestor,
					ancestor.ID.EQ(table.TaxaClosure.AncestorID),
				),
		).
		GROUP_BY(
			ancestor.ID,
			ancestor.Name,
			ancestor.Authorship,
			ancestor.Rank,
			ancestor.ParentID,
		).ORDER_BY(ancestor.Rank.DESC())

	sql, args := stmt.Sql()
	// logrus.Infof("SQL Query: %s", sql)
	// logrus.Infof("args: %+v", args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.OccurrenceOverviewItem, 0)
	for rows.Next() {
		var (
			row biomedb.OccurrencesByTaxaOverviewRow
		)
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Authorship,
			&row.Rank,
			&row.ParentID,
			&row.OccurrencesCount,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, models.OccurrenceOverviewItemFromDB(row))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil

}

func (s *OccurrenceStore) ListSamplingsAtCell(ctx context.Context, q db.Querier, cell h3.Cell, resolution int64, params ListOccurrencesParams) ([]models.Sampling, error) {
	h3CellFilter := s.H3CellFilter(resolution, cell)
	return s.ListSamplings(ctx, q, params, h3CellFilter)
}

func (s *OccurrenceStore) ListSamplings(ctx context.Context, q db.Querier, params ListOccurrencesParams, filters ...BoolExpression) ([]models.Sampling, error) {

	logrus.Debugf("ListSamplingsAtCell called with params: %+v", params)

	stmt := SELECT(
		table.Samplings.AllColumns,
		table.Countries.AllColumns,
	).
		FROM(s.OccurrenceCoreTables()).
		GROUP_BY(table.Samplings.ID, table.Countries.Code)

	stmt = params.ApplyFilters(stmt, filters...)
	stmt = params.ApplySorting(stmt)
	// stmt = params.ApplySorting(stmt, RawFloat("0.0"))
	// stmt = params.ApplyPagination(stmt)

	sql, args := stmt.Sql()
	logrus.Infof("SQL Query: %s", sql)
	logrus.Infof("args: %+v", args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]models.Sampling, 0)
	for rows.Next() {
		var (
			sampling biomedb.SamplingsWithCountry
		)
		err := rows.Scan(
			&sampling.ID,
			&sampling.SourceSamplingHash,
			&sampling.Comments,
			&sampling.SiteCode,
			&sampling.SiteName,
			&sampling.SiteLocality,
			&sampling.SiteCountryCode,
			&sampling.CoordinatesPrecision,
			&sampling.Latitude,
			&sampling.Longitude,
			&sampling.Altitude,
			&sampling.EventDate,
			&sampling.EventDatePrecision,
			&sampling.PerformedBy,
			&sampling.Duration,
			&sampling.AccessPoints,
			&sampling.ImportBatchID,
			&sampling.H3Index,
			&sampling.SearchVector,
			&sampling.CountryCode,
			&sampling.CountryName,
			&sampling.CountryContinent,
			&sampling.CountrySubcontinent,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, models.NewSamplingFromDB(sampling))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *OccurrenceStore) ListOccurrencesAtCell(ctx context.Context, q db.Querier, cell h3.Cell, resolution int64, params ListOccurrencesParams) ([]models.BaseOccurrenceWithSamplingID, error) {
	h3CellFilter := s.H3CellFilter(resolution, cell)
	return s.ListBaseOccurrences(ctx, q, params, h3CellFilter)
}

func (s *OccurrenceStore) ListBaseOccurrences(ctx context.Context, q db.Querier, params ListOccurrencesParams, filters ...BoolExpression) ([]models.BaseOccurrenceWithSamplingID, error) {
	logrus.Debugf("ListOccurrencesAtCell called with params: %+v", params)

	stmt := SELECT(
		table.Occurrences.AllColumns,
		table.Taxa.AllColumns,
	).FROM(s.OccurrenceCoreTables())

	stmt = params.ApplyFilters(stmt, filters...)
	stmt = params.ApplySorting(stmt)
	// stmt = params.ApplySorting(stmt, RawFloat("0.0"))
	// stmt = params.ApplyPagination(stmt)

	sql, args := stmt.Sql()
	// logrus.Infof("SQL Query: %s", sql)
	// logrus.Infof("args: %+v", args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]models.BaseOccurrenceWithSamplingID, 0)
	for rows.Next() {
		var (
			o biomedb.Occurrence
			t biomedb.Taxon
		)
		err := rows.Scan(
			&o.ID,
			&o.Code,
			&o.SamplingID,
			&o.TypeStatus,
			&o.Comments,
			&o.TaxonID,
			&o.VerbatimIdentification,
			&o.IdentifiedBy,
			&o.IdentificationDate,
			&o.IdentificationDatePrecision,
			&o.IdentificationConfer,
			&o.IdentificationAddendum,
			&o.ContentDescription,
			&o.QuantityExact,
			&o.QuantityLower,
			&o.QuantityUpper,
			&o.Sources,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.ImportBatchID,
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
		results = append(results, models.BaseOccurrenceFromDB(o, t).WithSamplingID(o.SamplingID))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *OccurrenceStore) listH3CellsWithSummaryAtResolution(ctx context.Context, q db.Querier, fromTables ReadableTable, resolution int64, params FilterParams) ([]models.H3CellWithRichness, error) {
	logrus.Debugf("listH3CellsAtResolution called with params: %+v", params)

	resolutionExpr := RawInt(fmt.Sprintf("%d", resolution))
	h3Cell := Func("h3_cell_to_parent", CAST(table.Samplings.H3Index).AS("h3index"), resolutionExpr)

	stmt := SELECT(
		h3Cell.AS("h3_index"),
		COUNT(DISTINCT(table.Occurrences.ID)).AS("occurrences_count"),
		COUNT(DISTINCT(table.Occurrences.SamplingID)).AS("samplings_count"),
		Raw(`
			COUNT(DISTINCT ancestor.id)
			FILTER (WHERE ancestor.rank = 'species')
		`).AS("species_richness"),
		Raw(`
			COUNT(DISTINCT ancestor.id)
			FILTER (WHERE ancestor.rank = 'genus')
		`).AS("genus_richness"),
		Raw(`
			COUNT(DISTINCT ancestor.id)
			FILTER (WHERE ancestor.rank = 'family')
		`).AS("family_richness"),
	).FROM(fromTables).
		GROUP_BY(h3Cell)

	stmt = params.ApplyFilters(stmt)

	sql, args := stmt.Sql()
	logrus.Infof("SQL Query: %s", sql)
	logrus.Infof("args: %+v", args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]models.H3CellWithRichness, 0)
	for rows.Next() {
		var (
			cellIndex string
			cell      models.H3CellWithRichness
		)
		err := rows.Scan(
			&cellIndex,
			&cell.OccurrencesCount,
			&cell.SamplingsCount,
			&cell.SpeciesRichness,
			&cell.GenusRichness,
			&cell.FamilyRichness,
		)
		cell.H3Index = h3.CellFromString(cellIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, cell)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *OccurrenceStore) ListSamplingsH3(ctx context.Context, q db.Querier, resolution int64, params ListSamplingsParams) ([]models.H3CellWithRichness, error) {
	tables := table.Samplings.
		LEFT_JOIN(table.Countries, table.Samplings.SiteCountryCode.EQ(table.Countries.Code)).
		LEFT_JOIN(table.Occurrences, table.Occurrences.SamplingID.EQ(table.Samplings.ID)).
		LEFT_JOIN(table.Taxa, table.Occurrences.TaxonID.EQ(table.Taxa.ID)).
		LEFT_JOIN(table.TaxaClosure, table.TaxaClosure.DescendantID.EQ(table.Taxa.ID)).
		LEFT_JOIN(table.Taxa.AS("ancestor"), table.TaxaClosure.AncestorID.EQ(table.Taxa.AS("ancestor").ID))
	return s.listH3CellsWithSummaryAtResolution(ctx, q, tables, resolution, params)
}

func (s *OccurrenceStore) ListOccurrencesH3(ctx context.Context, q db.Querier, resolution int64, params ListOccurrencesParams) ([]models.H3CellWithRichness, error) {

	logrus.Debugf("ListOccurrencesH3 called with params: %+v", params)

	tables := s.OccurrenceCoreTables().
		INNER_JOIN(table.TaxaClosure, table.TaxaClosure.DescendantID.EQ(table.Taxa.ID)).
		INNER_JOIN(table.Taxa.AS("ancestor"), table.TaxaClosure.AncestorID.EQ(table.Taxa.AS("ancestor").ID))
	return s.listH3CellsWithSummaryAtResolution(ctx, q, tables, resolution, params)

	// resolutionExpr := RawInt(fmt.Sprintf("%d", resolution))
	// h3Cell := Func("h3_cell_to_parent", CAST(table.Samplings.H3Index).AS("h3index"), resolutionExpr)
	// // h3CellBigInt := CAST(h3Cell).AS_BIGINT()
	// stmt := SELECT(
	// 	h3Cell.AS("h3_index"),
	// 	COUNT(DISTINCT(table.Occurrences.ID)).AS("occurrences_count"),
	// 	COUNT(DISTINCT(table.Occurrences.SamplingID)).AS("samplings_count"),
	// 	Raw(`
	// 		COUNT(DISTINCT ancestor.id)
	// 		FILTER (WHERE ancestor.rank = 'species')
	// 	`).AS("species_richness"),
	// 	Raw(`
	// 		COUNT(DISTINCT ancestor.id)
	// 		FILTER (WHERE ancestor.rank = 'genus')
	// 	`).AS("genus_richness"),
	// 	Raw(`
	// 		COUNT(DISTINCT ancestor.id)
	// 		FILTER (WHERE ancestor.rank = 'family')
	// 	`).AS("family_richness"),
	// ).FROM(
	// 	s.OccurrenceCoreTables().
	// 		INNER_JOIN(table.TaxaClosure, table.TaxaClosure.DescendantID.EQ(table.Taxa.ID)).
	// 		INNER_JOIN(table.Taxa.AS("ancestor"), table.TaxaClosure.AncestorID.EQ(table.Taxa.AS("ancestor").ID)),
	// ).GROUP_BY(h3Cell)

	// searchParts := buildSearchParts(params.SearchTerm)
	// stmt = params.ApplyFilters(stmt, searchParts)

	// sql, args := stmt.Sql()
	// // logrus.Infof("SQL Query: %s", sql)
	// // logrus.Infof("args: %+v", args)
	// rows, err := q.Query(ctx, sql, args...)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// result := make([]models.H3CellWithRichness, 0)
	// for rows.Next() {
	// 	var (
	// 		cellIndex string
	// 		cell      models.H3CellWithRichness
	// 	)
	// 	err := rows.Scan(
	// 		&cellIndex,
	// 		&cell.OccurrencesCount,
	// 		&cell.SamplingsCount,
	// 		&cell.SpeciesRichness,
	// 		&cell.GenusRichness,
	// 		&cell.FamilyRichness,
	// 	)
	// 	cell.H3Index = h3.CellFromString(cellIndex)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	result = append(result, cell)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	// return result, nil
}

func (s *OccurrenceStore) ListOccurrencesCount(ctx context.Context, q db.Querier, params ListOccurrencesParams) (int64, error) {

	logrus.Debugf("ListOccurrencesCount called with params: %+v", params)
	stmt := SELECT(COUNT(STAR)).FROM(s.OccurrenceCoreTables())
	stmt = params.ApplyFilters(stmt)

	sql, args := stmt.Sql()
	// logrus.Infof("args: %+v", args)
	row := q.QueryRow(ctx, sql, args...)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OccurrenceStore) ListOccurrences(ctx context.Context, q db.Querier, params ListOccurrencesParams) ([]models.Occurrence, error) {

	logrus.Debugf("ListOccurrences called with params: %+v", params)

	score := Float(0.0)
	searchParts := buildSearchParts(strings.TrimSpace(params.SearchTerm))
	score = searchParts.toScoreProjection()

	stmt := SELECT(
		table.Occurrences.AllColumns,
		table.Samplings.AllColumns,
		table.Taxa.AllColumns,
		table.Countries.AllColumns,
		score.AS("score"),
	).FROM(s.OccurrenceCoreTables())
	stmt = params.ApplyFilters(stmt)
	stmt = params.ApplySortingWithScore(stmt, FloatColumn("score"))
	stmt = params.ApplyPagination(stmt)

	sql, args := stmt.Sql()
	// logrus.Infof("SQL Query: %s", sql)
	// logrus.Infof("args: %+v", args)
	occurrencesRows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer occurrencesRows.Close()

	occurrences, err := s.ScanOccurrenceRow(occurrencesRows)
	if err != nil {
		return nil, err
	}

	return occurrences, nil
}

func (s *OccurrenceStore) ScanOccurrenceRow(rows pgx.Rows) ([]models.Occurrence, error) {
	occurrences := make([]models.Occurrence, 0, 64)

	for rows.Next() {
		var score float32
		var i biomedb.GetOccurrenceByIDRow
		err := rows.Scan(
			&i.Occurrence.ID,
			&i.Occurrence.Code,
			&i.Occurrence.SamplingID,
			&i.Occurrence.TypeStatus,
			&i.Occurrence.Comments,
			&i.Occurrence.TaxonID,
			&i.Occurrence.VerbatimIdentification,
			&i.Occurrence.IdentifiedBy,
			&i.Occurrence.IdentificationDate,
			&i.Occurrence.IdentificationDatePrecision,
			&i.Occurrence.IdentificationConfer,
			&i.Occurrence.IdentificationAddendum,
			&i.Occurrence.ContentDescription,
			&i.Occurrence.QuantityExact,
			&i.Occurrence.QuantityLower,
			&i.Occurrence.QuantityUpper,
			&i.Occurrence.Sources,
			&i.Occurrence.CreatedAt,
			&i.Occurrence.UpdatedAt,
			&i.Occurrence.ImportBatchID,
			&i.SamplingsWithCountry.ID,
			&i.SamplingsWithCountry.SourceSamplingHash,
			&i.SamplingsWithCountry.Comments,
			&i.SamplingsWithCountry.SiteCode,
			&i.SamplingsWithCountry.SiteName,
			&i.SamplingsWithCountry.SiteLocality,
			&i.SamplingsWithCountry.SiteCountryCode,
			&i.SamplingsWithCountry.CoordinatesPrecision,
			&i.SamplingsWithCountry.Latitude,
			&i.SamplingsWithCountry.Longitude,
			&i.SamplingsWithCountry.Altitude,
			&i.SamplingsWithCountry.EventDate,
			&i.SamplingsWithCountry.EventDatePrecision,
			&i.SamplingsWithCountry.PerformedBy,
			&i.SamplingsWithCountry.Duration,
			&i.SamplingsWithCountry.AccessPoints,
			&i.SamplingsWithCountry.ImportBatchID,
			&i.SamplingsWithCountry.H3Index,
			&i.SamplingsWithCountry.SearchVector,
			&i.Taxon.ID,
			&i.Taxon.GBIFID,
			&i.Taxon.Name,
			&i.Taxon.ScientificName,
			&i.Taxon.Rank,
			&i.Taxon.Status,
			&i.Taxon.Authorship,
			&i.Taxon.AcceptedTaxonID,
			&i.Taxon.ParentID,
			&i.Taxon.SearchVector,
			&i.Taxon.Comments,
			&i.SamplingsWithCountry.CountryCode,
			&i.SamplingsWithCountry.CountryName,
			&i.SamplingsWithCountry.CountryContinent,
			&i.SamplingsWithCountry.CountrySubcontinent,
			&score,
		)
		if err != nil {
			return nil, err
		}
		occurrences = append(occurrences, models.OccurrenceFromDB(i.Occurrence, i.Taxon, i.SamplingsWithCountry))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return occurrences, nil
}
