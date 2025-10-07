CREATE MIGRATION m1f7bean47xanlluwcpovqiry3k3x6a6ajdoxjpfefk3maj7n3jcdq
    ONTO m1rfx6iijta632dhgiolfcmtjkh7wkptwjc3nbrw4qdy3vkslmz5ca
{
  DROP FUNCTION occurrence::insert_external_biomat(sampling: events::Sampling, data: std::json);
  DROP FUNCTION occurrence::insert_internal_biomat(sampling: events::Sampling, data: std::json);
  DROP FUNCTION seq::insert_external_seq(sampling: events::Sampling, data: std::json);
  DROP FUNCTION date::from_json_with_precision(value: std::json);
  CREATE FUNCTION date::from_json_with_precision(value: OPTIONAL std::json) -> OPTIONAL tuple<date: std::datetime, precision: date::DatePrecision> USING ((IF NOT (EXISTS (value)) THEN <tuple<date: std::datetime, precision: date::DatePrecision>>{} ELSE std::assert_exists((
      date := std::to_datetime(<std::int64>(value)['date']['year'], (<std::int64>std::json_get(value, 'date', 'month') ?? 1), (<std::int64>std::json_get(value, 'date', 'day') ?? 1), 0, 0, 0, 'UTC'),
      precision := <date::DatePrecision>(value)['precision']
  ), message := ('Failed to parse date with precision from JSON: ' ++ std::to_str(value)))));
  CREATE FUNCTION occurrence::insert_external_biomat(sampling: events::Sampling, data: std::json) ->  occurrence::ExternalBioMat USING (WITH
      identification := 
          (data)['identification']
      ,
      taxon := 
          taxonomy::taxonByName(<std::str>(identification)['taxon'])
      ,
      publications := 
          std::json_array_unpack(std::json_get(data, 'published_in'))
  INSERT
      occurrence::ExternalBioMat
      {
          sampling := sampling,
          code := (<std::str>std::json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling)),
          original_source := (WITH
              src := 
                  <std::str>std::json_get(data, 'original_source')
          SELECT
              (IF EXISTS (src) THEN (default::get_vocabulary(src))[IS references::DataSource] ELSE <references::DataSource>{})
          ),
          original_link := <std::str>std::json_get(data, 'original_link'),
          quantity := <occurrence::QuantityType>std::json_get(data, 'quantity'),
          content_description := <std::str>std::json_get(data, 'content_description'),
          in_collection := <std::str>std::json_get(data, 'collection'),
          item_vouchers := <std::str>std::json_array_unpack(std::json_get(data, 'item_vouchers')),
          comments := <std::str>std::json_get(data, 'comments'),
          published_in := (SELECT
              DISTINCT ((FOR p IN publications
              UNION 
                  (SELECT
                      references::Article {
                          @original_source := <std::bool>std::json_get(p, 'original')
                      }
                  FILTER
                      (.code = <std::str>(p)['code'])
                  )))
          ),
          identification := (INSERT
              occurrence::Identification
              {
                  taxon := taxon,
                  identified_by := people::personByAlias(<std::str>(identification)['identified_by']),
                  identified_on := date::from_json_with_precision((identification)['identified_on'])
              }),
          is_type := (<std::bool>std::json_get(data, 'is_type') ?? false)
      });
  CREATE FUNCTION occurrence::insert_internal_biomat(sampling: events::Sampling, data: std::json) ->  occurrence::InternalBioMat USING (WITH
      identification := 
          (data)['identification']
      ,
      taxon := 
          taxonomy::taxonByName(<std::str>(identification)['taxon'])
      ,
      publications := 
          std::json_array_unpack(std::json_get(data, 'published_in'))
  INSERT
      occurrence::InternalBioMat
      {
          sampling := sampling,
          code := (<std::str>std::json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling)),
          identification := (INSERT
              occurrence::Identification
              {
                  taxon := taxon,
                  identified_by := people::personByAlias(<std::str>(identification)['identified_by']),
                  identified_on := date::from_json_with_precision((identification)['identified_on'])
              }),
          is_type := (<std::bool>std::json_get(data, 'is_type') ?? false),
          comments := <std::str>std::json_get(data, 'comments'),
          published_in := (SELECT
              DISTINCT ((FOR p IN publications
              UNION 
                  (SELECT
                      references::Article {
                          @original_source := <std::bool>std::json_get(p, 'original')
                      }
                  FILTER
                      (.code = <std::str>(p)['code'])
                  )))
          )
      });
  CREATE FUNCTION seq::insert_external_seq(sampling: events::Sampling, data: std::json) ->  seq::ExternalSequence USING (INSERT
      seq::ExternalSequence
      {
          sampling := sampling,
          code := <std::str>(data)['code'],
          label := <std::str>std::json_get(data, 'label'),
          sequence := <std::str>std::json_get(data, 'sequence'),
          gene := seq::geneByCode(<std::str>(data)['gene']),
          legacy := <tuple<id: std::int32, code: std::str, alignment_code: std::str>>std::json_get(data, 'legacy'),
          origin := <seq::ExtSeqOrigin>std::json_get(data, 'origin'),
          published_in := (WITH
              pubs := 
                  std::json_array_unpack(std::json_get(data, 'published_in'))
          SELECT
              DISTINCT ((FOR p IN pubs
              UNION 
                  (SELECT
                      references::Article {
                          @original_source := <std::bool>std::json_get(p, 'original')
                      }
                  FILTER
                      (.code = <std::str>(p)['code'])
                  )))
          ),
          identification := (WITH
              identification := 
                  (data)['identification']
          INSERT
              occurrence::Identification
              {
                  identified_by := people::personByAlias(<std::str>(identification)['identified_by']),
                  identified_on := date::from_json_with_precision((identification)['identified_on']),
                  taxon := taxonomy::taxonByName(<std::str>(identification)['taxon'])
              }),
          referenced_in := (FOR ref IN std::json_array_unpack(std::json_get(data, 'referenced_in'))
          INSERT
              references::SeqReference
              {
                  db := references::dataSourceByCode(<std::str>(ref)['db']),
                  accession := <std::str>(ref)['accession'],
                  is_origin := <std::bool>std::json_get(ref, 'is_origin')
              }),
          specimen_identifier := <std::str>std::json_get(data, 'specimen_identifier'),
          original_taxon := <std::str>std::json_get(data, 'original_taxon'),
          source_sample := (WITH
              source_sample := 
                  <std::str>std::json_get(data, 'source_sample')
          SELECT
              (IF EXISTS (source_sample) THEN occurrence::externalBiomatByCode(source_sample) ELSE <occurrence::ExternalBioMat>{})
          )
      });
};
