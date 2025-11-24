CREATE MIGRATION m1pso6tvnywwbsqubu2q5nkesrpqypvn2p5av7pgmzxaqr7a7oscnq
    ONTO m1vvlq6acjo7mhvqkq3q5ageyjqwt27x7hshzxqmnpuresjtze6lxa
{
  ALTER FUNCTION occurrence::insert_external_occurrence(sampling: events::Sampling, data: std::json) USING (WITH
      identification := 
          (data)['identification']
      ,
      taxon := 
          taxonomy::taxonByName(<std::str>(identification)['taxon'])
  INSERT
      occurrence::ExternalOccurrence
      {
          sampling := sampling,
          code := (<std::str>std::json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling.code)),
          sources := (SELECT
              references::DataSource
          FILTER
              (.code IN <std::str>std::json_array_unpack(std::json_get(data, 'sources')))
          ),
          external_link := <std::str>std::json_get(data, 'external_link'),
          original_taxon := <std::str>std::json_get(data, 'original_taxon'),
          content_description := <std::str>std::json_get(data, 'content_description'),
          in_collection := <std::str>std::json_get(data, 'collection'),
          item_vouchers := <std::str>std::json_array_unpack(std::json_get(data, 'item_vouchers')),
          comments := <std::str>std::json_get(data, 'comments'),
          published_in := (SELECT
              DISTINCT (references::Article)
          FILTER
              (.code IN <std::str>std::json_array_unpack(std::json_get(data, 'published_in')))
          ),
          identification := (INSERT
              occurrence::Identification
              {
                  taxon := taxon,
                  identified_by := people::personByAlias(<std::str>std::json_get(identification, 'identified_by')),
                  identified_on := date::from_json_with_precision(std::json_get(identification, 'identified_on'))
              }),
          type_status := <occurrence::TypeStatus>std::json_get(data, 'type_status')
      });
  ALTER TYPE occurrence::ExternalOccurrence {
      ALTER PROPERTY quantity {
          CREATE ANNOTATION std::description := 'The number of specimens reported in this occurrence.';
      };
  };
  ALTER TYPE occurrence::ExternalOccurrence {
      ALTER PROPERTY quantity {
          SET TYPE tuple<lower: std::int32, upper: std::int32> USING (<tuple<lower: std::int32, upper: std::int32>>{});
      };
  };
  DROP SCALAR TYPE occurrence::QuantityType;
};
