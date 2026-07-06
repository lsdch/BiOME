package stores

import (
	"context"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

const MIN_FULL_TEXT_SEARCH_TERM_LENGTH = 3

type OccurrenceStore struct {
}

func NewOccurrenceStore() *OccurrenceStore {
	return &OccurrenceStore{}
}

func (s *OccurrenceStore) fromOccurrenceCoreTables(stmt SelectStatement) SelectStatement {
	return stmt.
		FROM(
			table.Occurrences.
				INNER_JOIN(table.Taxa,
					table.Occurrences.TaxonID.EQ(table.Taxa.ID)).
				INNER_JOIN(table.Samplings,
					table.Occurrences.SamplingID.EQ(table.Samplings.ID)).
				INNER_JOIN(table.Countries,
					table.Samplings.SiteCountryCode.EQ(table.Countries.Code)),
		)
}

func (s *OccurrenceStore) ListOccurrencesCount(ctx context.Context, q db.Querier, params ListOccurrencesParams) (int64, error) {
	stmt := s.fromOccurrenceCoreTables(SELECT(COUNT(Int(1))))
	searchParts := buildSearchParts(params.SearchTerm)
	stmt = params.ApplyFilters(stmt, searchParts)

	sql, args := stmt.Sql()
	row := q.QueryRow(ctx, sql, args...)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OccurrenceStore) ListOccurrences(ctx context.Context, q db.Querier, params ListOccurrencesParams) ([]models.Occurrence, error) {
	score := Float(0)
	searchParts := buildSearchParts(strings.TrimSpace(params.SearchTerm))
	score = searchParts.toScoreProjection()

	stmt := s.fromOccurrenceCoreTables(SELECT(
		table.Occurrences.AllColumns,
		table.Taxa.AllColumns,
		table.Samplings.AllColumns,
		table.Countries.AllColumns,
		score,
	))
	stmt = params.ApplyFilters(stmt, searchParts)
	stmt = params.ApplySorting(stmt, score)
	stmt = params.ApplyPagination(stmt)

	sql, args := stmt.Sql()
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
			&i.Sampling.ID,
			&i.Sampling.Notes,
			&i.Sampling.SiteCode,
			&i.Sampling.SiteName,
			&i.Sampling.SiteLocality,
			&i.Sampling.SiteCountryCode,
			&i.Sampling.CoordinatesPrecision,
			&i.Sampling.Coordinates,
			&i.Sampling.Latitude,
			&i.Sampling.Longitude,
			&i.Sampling.Altitude,
			&i.Sampling.EventDate,
			&i.Sampling.EventDatePrecision,
			&i.Sampling.PerformedBy,
			&i.Sampling.Duration,
			&i.Sampling.AccessPoints,
			&i.Sampling.H3Index,
			&i.Taxon.ID,
			&i.Taxon.GBIFID,
			&i.Taxon.Name,
			&i.Taxon.ScientificName,
			&i.Taxon.Rank,
			&i.Taxon.Status,
			&i.Taxon.Authorship,
			&i.Taxon.AcceptedTaxonID,
			&i.Taxon.ParentID,
			&i.Taxon.Comments,
			&i.Country.Code,
			&i.Country.Name,
			&i.Country.Continent,
			&i.Country.Subcontinent,
		)
		if err != nil {
			return nil, err
		}
		occurrences = append(occurrences, models.OccurrenceFromDB(i.Occurrence, i.Taxon, i.Sampling, i.Country))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return occurrences, nil
}

type ListOccurrencesParams struct {
	SearchTerm  string                                                   `json:"search_term,omitempty" query:"search_term"`
	Datasets    []ulid.ULID                                              `json:"datasets,omitempty" query:"datasets"`
	Taxa        []uuid.UUID                                              `json:"taxa,omitempty" query:"taxa"`
	WholeClade  bool                                                     `json:"whole_clade,omitempty" query:"whole_clade"`
	Confer      models.Optional[bool]                                    `json:"confer,omitempty" query:"confer"`
	TaxonRank   models.Optional[biomedb.TaxonRank]                       `json:"taxon_rank,omitempty" query:"taxon_rank"`
	TaxonStatus models.Optional[biomedb.TaxonStatus]                     `json:"taxon_status,omitempty" query:"taxon_status"`
	Pagination  models.Pagination                                        `json:"pagination"`
	OrderBy     models.Optional[models.SortBy[models.OccurrenceSortKey]] `json:"order_by"`
}

func (p *ListOccurrencesParams) ApplyFilters(stmt SelectStatement, parts searchParts) SelectStatement {
	stmt = p.applyDatasetFilter(stmt)
	stmt = parts.applySearchTermFilter(stmt)
	stmt = p.applyTaxonomyFilter(stmt)
	return stmt
}

func (p *ListOccurrencesParams) ApplySorting(stmt SelectStatement, score FloatExpression) SelectStatement {
	orderBy, ok := p.OrderBy.Get()
	if ok {
		stmt = stmt.ORDER_BY(score.DESC(), orderBy.ToOrderByClause())
	} else {
		// Default sorting by occurrence code if no order_by is provided
		stmt = stmt.ORDER_BY(score.DESC(), table.Occurrences.Code.ASC())
	}
	return stmt
}

func (p *ListOccurrencesParams) ApplyPagination(stmt SelectStatement) SelectStatement {
	if p.Pagination.Limit > 0 {
		stmt = stmt.LIMIT(int64(p.Pagination.Limit))
	}
	stmt = stmt.OFFSET(int64(p.Pagination.Offset))
	return stmt
}

func (p *ListOccurrencesParams) applyDatasetFilter(stmt SelectStatement) SelectStatement {
	if len(p.Datasets) == 0 {
		return stmt
	}

	var datasetsExpr = make([]Expression, 0, len(p.Datasets))
	for _, id := range p.Datasets {
		datasetsExpr = append(datasetsExpr, String(id.String()))
	}

	od := table.OccurrencesDatasets
	return stmt.WHERE(
		EXISTS(
			SELECT(Int(1)).
				FROM(od).
				WHERE(
					od.OccurrenceID.EQ(table.Occurrences.ID).
						AND(od.DatasetID.IN(datasetsExpr...)),
				),
		),
	)
}

func (p *ListOccurrencesParams) applyTaxonomyFilter(stmt SelectStatement) SelectStatement {
	hasTaxaFilter := len(p.Taxa) > 0
	hasRankFilter := p.TaxonRank.IsSet
	hasStatusFilter := p.TaxonStatus.IsSet
	hasConferFilter := p.Confer.IsSet

	if !hasTaxaFilter && !hasRankFilter && !hasStatusFilter && !hasConferFilter {
		return stmt
	}

	if hasTaxaFilter {
		var taxaExpr = make([]Expression, 0, len(p.Taxa))
		for _, id := range p.Taxa {
			taxaExpr = append(taxaExpr, String(id.String()))
		}
		if p.WholeClade {
			stmt = stmt.WHERE(
				EXISTS(
					SELECT(Int(1)).
						FROM(table.TaxaClosure).
						WHERE(
							table.TaxaClosure.DescendantID.EQ(table.Taxa.ID).
								AND(table.TaxaClosure.AncestorID.IN(taxaExpr...)),
						),
				),
			)
		} else {
			stmt = stmt.WHERE(
				table.Taxa.ID.IN(taxaExpr...),
			)
		}
	}

	// 2. rank
	if hasRankFilter {
		stmt = stmt.WHERE(
			table.Taxa.Rank.EQ(String(string(p.TaxonRank.Value))),
		)
	}

	// 3. status
	if hasStatusFilter {
		stmt = stmt.WHERE(
			table.Taxa.Status.EQ(String(string(p.TaxonStatus.Value))),
		)
	}

	// 4. confer (occurrences)
	if hasConferFilter {
		stmt = stmt.WHERE(
			table.Occurrences.IdentificationConfer.EQ(Bool(p.Confer.Value)),
		)
	}

	return stmt
}

type searchParts struct {
	term     string
	exact    BoolExpression
	prefix   BoolExpression
	fts      BoolExpression
	fts_rank FloatExpression
}

func buildSearchParts(term string) searchParts {

	if len(term) == 0 {
		return searchParts{}
	}

	t := String(term)
	like := String(term + "%")

	return searchParts{
		term: term,
		exact: OR(
			table.Occurrences.Code.EQ(t),
			table.Taxa.Name.EQ(t),
			table.Samplings.SiteCode.EQ(t),
		),
		prefix: OR(
			table.Occurrences.Code.LIKE(like),
			table.Taxa.Name.LIKE(like),
			table.Samplings.SiteCode.LIKE(like),
			table.Samplings.SiteName.LIKE(like),
		),
		fts: RawBool(`
			taxa.search_vector @@ plainto_tsquery('simple', #term)
			OR samplings.search_vector @@ plainto_tsquery('simple', #term)
		`, RawArgs{"#term": term}),
		fts_rank: RawFloat(`
			(ts_rank_cd(taxa.search_vector, plainto_tsquery('simple', #term)) * 0.6)
			+ (ts_rank_cd(samplings.search_vector, plainto_tsquery('simple', #term)) * 0.4)
		`, RawArgs{"#term": term}),
	}
}

func (s *searchParts) hasTerm() bool {
	return s.term != ""
}

func (s *searchParts) hasFullTextSearch() bool {
	return len(s.term) >= MIN_FULL_TEXT_SEARCH_TERM_LENGTH
}

func (s searchParts) toExpression() BoolExpression {
	if !s.hasTerm() {
		return Bool(true)
	}
	components := []BoolExpression{s.exact, s.prefix}
	if s.hasFullTextSearch() {
		components = append(components, s.fts)
	}
	return OR(components...)
}

func (s searchParts) toScoreProjection() FloatExpression {
	if !s.hasTerm() {
		return Float(0)
	}

	exact := Float(100)
	prefixScore := Float(60)
	ftsScore := Float(30)

	score := FloatExp(
		CASE().
			WHEN(s.exact).THEN(exact).
			WHEN(s.prefix).THEN(prefixScore).
			ELSE(Float(0)),
	)

	if s.hasFullTextSearch() {
		score = score.ADD(FloatExp(LEAST(s.fts_rank)).MUL(ftsScore))
	}
	return score
}

func (s searchParts) applySearchTermFilter(stmt SelectStatement) SelectStatement {
	if s.hasTerm() {
		return stmt.WHERE(s.toExpression())
	}
	return stmt
}
