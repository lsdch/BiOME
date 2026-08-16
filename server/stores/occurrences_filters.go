package stores

import (
	"slices"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
	"github.com/sirupsen/logrus"

	. "github.com/go-jet/jet/v2/postgres"
)

type FilterParams interface {
	ApplyFilters(stmt SelectStatement, constraints ...BoolExpression) SelectStatement
}

type BufferUnit = models.EventDatePrecision

type DateBuffer struct {
	Value int32      `json:"value" query:"buffer"`
	Unit  BufferUnit `json:"unit" query:"buffer_unit"`
}

type DateFilterParams struct {
	From           string                                    `json:"from,omitempty" query:"from" pattern:"^\\d{4}(-\\d{2})?(-\\d{2})?$"`
	fromParsed     models.Optional[models.DateWithPrecision] `json:"-"`
	To             string                                    `json:"to,omitempty" query:"to" pattern:"^\\d{4}(-\\d{2})?(-\\d{2})?$"`
	toParsed       models.Optional[models.DateWithPrecision] `json:"-"`
	Buffer         models.Optional[string]                   `json:"buffer,omitempty" query:"buffer" format:"duration"`
	bufferValue    time.Duration                             `json:"-"`
	IncludeUnknown bool                                      `json:"include_unknown,omitempty" query:"include_unknown"`
}

func (i *DateFilterParams) Resolve(ctx huma.Context) []error {
	if fromDate := i.From; fromDate != "" {
		from, err := models.ParseDateWithPrecision(fromDate)
		if err != nil {
			return []error{err}
		}
		i.fromParsed = models.NewOptionalFromPtr(from)
	}
	if toDate := i.To; toDate != "" {
		to, err := models.ParseDateWithPrecision(toDate)
		if err != nil {
			return []error{err}
		}
		if to != nil {
			*to = to.UpperBound()
			logrus.Infof("Upper bound for toDate %s is %s", toDate, to.Date.Format("2006-01-02"))
		}
		i.toParsed = models.NewOptionalFromPtr(to)
	}
	if i.Buffer.IsSet {
		duration, err := time.ParseDuration(i.Buffer.Value)
		if err != nil {
			return []error{err}
		}
		i.bufferValue = duration
	}
	return nil
}

func DateExpFromDateWithPrecision(date models.DateWithPrecision) DateExpression {
	return Date(
		int(date.Date.Year()),
		time.Month(date.Date.Month()),
		int(date.Date.Day()),
	)
}

type ListSamplingsParams struct {
	// Datasets             []types.ULID `json:"datasets,omitempty" query:"datasets"`
	Batches              []uuid.UUID      `json:"batches,omitempty" query:"batches"`
	TargetTaxa           []uuid.UUID      `json:"target_taxa,omitempty" query:"target_taxa"`
	TargetTaxaWholeClade bool             `json:"target_taxa_whole_clade,omitempty" query:"target_taxa_whole_clade"`
	Countries            []string         `json:"countries,omitempty" query:"countries"`
	Date                 DateFilterParams `json:",inline" query:"date,deepObject"`
	models.Pagination    `json:"pagination" query:"pagination"`
}

func (p ListSamplingsParams) ApplyFilters(stmt SelectStatement, constraints ...BoolExpression) SelectStatement {
	filters := append(
		slices.Concat(
			p.dateFilter(),
			p.batchFilter(),
			p.targetTaxaFilter(),
			p.countryFilter(),
		),
		constraints...,
	)
	if len(filters) == 0 {
		return stmt
	}
	return stmt.WHERE(AND(filters...))
}

func (p *ListSamplingsParams) dateFilter() []BoolExpression {
	if !p.Date.fromParsed.IsSet && !p.Date.toParsed.IsSet {
		return nil
	}

	var filters []BoolExpression

	if p.Date.fromParsed.IsSet {
		dateFrom := p.Date.fromParsed.Value
		if p.Date.Buffer.IsSet {
			dateFrom.Date.Add(-p.Date.bufferValue)
		}
		lowerBound := DateExpFromDateWithPrecision(dateFrom)
		lowerBoundFilter := table.Samplings.EventDate.GT_EQ(lowerBound)
		filters = append(filters, lowerBoundFilter)
	}

	if p.Date.toParsed.IsSet {
		dateTo := p.Date.toParsed.Value
		if p.Date.Buffer.IsSet {
			dateTo.Date.Add(p.Date.bufferValue)
		}
		upperBound := DateExpFromDateWithPrecision(dateTo)
		upperBoundFilter := table.Samplings.EventDate.LT_EQ(upperBound)
		filters = append(filters, upperBoundFilter)
	}

	if p.Date.IncludeUnknown {
		unknownFilter := table.Samplings.EventDate.IS_NULL()
		return []BoolExpression{OR(unknownFilter, AND(filters...))}
	}
	return filters
}

func (p *ListSamplingsParams) batchFilter() []BoolExpression {
	if len(p.Batches) == 0 {
		return nil
	}
	var batchesExpr = make([]Expression, 0, len(p.Batches))
	for _, id := range p.Batches {
		batchesExpr = append(batchesExpr, UUID(id))
	}

	return []BoolExpression{table.Samplings.ImportBatchID.IN(batchesExpr...)}
}

func (p ListSamplingsParams) targetTaxaFilter() []BoolExpression {
	if len(p.TargetTaxa) == 0 {
		return nil
	}

	var targetTaxaExpr = make([]Expression, 0, len(p.TargetTaxa))
	for _, id := range p.TargetTaxa {
		targetTaxaExpr = append(targetTaxaExpr, UUID(id))
	}

	if p.TargetTaxaWholeClade {
		return []BoolExpression{EXISTS(
			SELECT(RawInt("1")).
				FROM(table.TaxaClosure).
				WHERE(
					table.TaxaClosure.DescendantID.EQ(table.Taxa.ID).
						AND(table.TaxaClosure.AncestorID.IN(targetTaxaExpr...)),
				),
		),
		}
	} else {
		return []BoolExpression{table.Taxa.ID.IN(targetTaxaExpr...)}
	}
}

func (p ListSamplingsParams) countryFilter() []BoolExpression {
	if len(p.Countries) == 0 {
		return nil
	}
	var countriesExpr = make([]Expression, 0, len(p.Countries))
	for _, code := range p.Countries {
		countriesExpr = append(countriesExpr, String(code))
	}
	return []BoolExpression{table.Samplings.SiteCountryCode.IN(countriesExpr...)}
}

type ListOccurrencesParams struct {
	ListSamplingsParams `json:",inline"`
	models.SortBy[models.OccurrenceSortKey]
	SearchTerm  string                                       `json:"search_term,omitempty" query:"search_term"`
	Datasets    []types.ULID                                 `json:"datasets,omitempty" query:"datasets"`
	Batches     []uuid.UUID                                  `json:"batches,omitempty" query:"batches"`
	Taxa        []uuid.UUID                                  `json:"taxa,omitempty" query:"taxa"`
	WholeClade  bool                                         `json:"whole_clade,omitempty" query:"whole_clade"`
	Confer      models.Optional[bool]                        `json:"confer,omitempty" query:"confer"`
	TypeStatus  models.Optional[models.OccurrenceTypeStatus] `json:"type_status,omitempty" query:"type_status"`
	TaxonRank   models.Optional[models.TaxonRank]            `json:"taxon_rank,omitempty" query:"taxon_rank"`
	TaxonStatus models.Optional[models.TaxonStatus]          `json:"taxon_status,omitempty" query:"taxon_status"`
	// TargetTaxa           []uuid.UUID                                  `json:"target_taxa,omitempty" query:"target_taxa"`
	// TargetTaxaWholeClade bool                                         `json:"target_taxa_whole_clade,omitempty" query:"target_taxa_whole_clade"`
	// Countries            []string                                     `json:"countries,omitempty" query:"countries"`
	// models.Pagination    `json:"pagination" query:"pagination"`
}

func (p ListOccurrencesParams) ApplyFilters(stmt SelectStatement, constraints ...BoolExpression) SelectStatement {

	parts := buildSearchParts(p.SearchTerm)

	filters := append(
		slices.Concat(
			p.taxonomyFilters(),
			p.targetTaxaFilter(),
			p.datasetFilter(),
			p.typeStatusFilter(),
			p.countryFilter(),
			p.batchFilter(),
			p.dateFilter(),
			constraints,
		),
		parts.toExpression(),
	)
	return stmt.WHERE(AND(filters...))
}

func (p *ListOccurrencesParams) ApplySorting(stmt SelectStatement) SelectStatement {
	if p.SortBy.IsSet() {
		stmt = stmt.ORDER_BY(p.SortBy.ToOrderByClause())

	}
	return stmt
}

func (p *ListOccurrencesParams) ApplySortingWithScore(stmt SelectStatement, score ColumnFloat) SelectStatement {
	if p.SortBy.IsSet() {
		stmt = stmt.ORDER_BY(score.DESC(), p.SortBy.ToOrderByClause())
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

func (p *ListOccurrencesParams) batchFilter() []BoolExpression {
	if len(p.Batches) == 0 {
		return nil
	}
	var batchesExpr = make([]Expression, 0, len(p.Batches))
	for _, id := range p.Batches {
		batchesExpr = append(batchesExpr, UUID(id))
	}

	return []BoolExpression{table.Occurrences.ImportBatchID.IN(batchesExpr...)}
}

func (p *ListOccurrencesParams) datasetFilter() []BoolExpression {

	if len(p.Datasets) == 0 {
		return nil
	}
	var datasetsExpr = make([]Expression, 0, len(p.Datasets))
	for _, id := range p.Datasets {
		datasetsExpr = append(datasetsExpr, String(id.String()))
	}

	od := table.OccurrencesDatasets
	return []BoolExpression{EXISTS(
		SELECT(Bool(true)).
			FROM(od).
			WHERE(
				od.OccurrenceID.EQ(table.Occurrences.ID).
					AND(od.DatasetID.IN(datasetsExpr...)),
			),
	)}
}

func (p *ListOccurrencesParams) typeStatusFilter() []BoolExpression {
	if !p.TypeStatus.IsSet {
		return nil
	}
	return []BoolExpression{table.Occurrences.TypeStatus.EQ(StringExp(CAST(String(string(p.TypeStatus.Value))).AS("public.occurrence_type_status")))}
}

func (p *ListOccurrencesParams) taxonomyFilters() []BoolExpression {
	hasTaxaFilter := len(p.Taxa) > 0
	hasRankFilter := p.TaxonRank.IsSet
	hasStatusFilter := p.TaxonStatus.IsSet
	hasConferFilter := p.Confer.IsSet

	if !hasTaxaFilter && !hasRankFilter && !hasStatusFilter && !hasConferFilter {
		return nil
	}

	expressions := make([]BoolExpression, 0, 4)

	if hasTaxaFilter {
		var taxaExpr = make([]Expression, 0, len(p.Taxa))
		for _, id := range p.Taxa {
			taxaExpr = append(taxaExpr, UUID(id))
		}
		if p.WholeClade {
			expressions = append(expressions, EXISTS(
				SELECT(RawInt("1")).
					FROM(table.TaxaClosure).
					WHERE(
						table.TaxaClosure.DescendantID.EQ(table.Taxa.ID).
							AND(table.TaxaClosure.AncestorID.IN(taxaExpr...)),
					),
			),
			)
		} else {
			expressions = append(expressions, table.Taxa.ID.IN(taxaExpr...))
		}
	}

	// 2. rank
	if hasRankFilter {
		expressions = append(expressions, table.Taxa.Rank.EQ(StringExp(CAST(String(string(p.TaxonRank.Value))).AS("public.taxon_rank"))))
	}

	// 3. status
	if hasStatusFilter {
		expressions = append(expressions, table.Taxa.Status.EQ(StringExp(CAST(String(string(p.TaxonStatus.Value))).AS("public.taxon_status"))))
	}

	// 4. confer (occurrences)
	if hasConferFilter {
		expressions = append(expressions, table.Occurrences.IdentificationConfer.EQ(Bool(p.Confer.Value)))
	}

	return expressions
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

	t := Text(term)
	like := Text("%" + term + "%")

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
		return RawFloat("0.0")
	}

	exact := RawFloat("100.0")
	prefixScore := RawFloat("60.0")
	ftsScore := RawFloat("30.0")

	score := FloatExp(
		CAST(
			CASE().
				WHEN(s.exact).THEN(exact).
				WHEN(s.prefix).THEN(prefixScore).
				ELSE(RawFloat("0.0")),
		).AS("float"),
	)

	if s.hasFullTextSearch() {
		score = score.ADD(FloatExp(LEAST(s.fts_rank)).MUL(ftsScore))
	}
	return score
}
