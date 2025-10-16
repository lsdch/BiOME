CREATE MIGRATION m145hzzzdqzcppbmqkwx7qvyfdwkbvbyw5pw5ayfm476cjyok57hmq
    ONTO m1a7sbqwojujqrehmkur3yrgekqc554maggjpnmu2qdupbmub4blzq
{
  ALTER TYPE occurrence::ExternalBioMat {
      ALTER LINK sequences {
          SET MULTI;
      };
      CREATE REQUIRED PROPERTY has_sequences := (EXISTS (.sequences));
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE REQUIRED PROPERTY has_sequences := (EXISTS (.specimens.sequences));
  };
  ALTER ALIAS occurrence::OccurrenceWithType USING (SELECT
      occurrence::Occurrence {
          required has_sequences := ([IS occurrence::InternalBioMat].has_sequences ?? ([IS occurrence::ExternalBioMat].has_sequences ?? false)),
          internal := [IS occurrence::InternalBioMat] {
              is_homogenous,
              is_congruent,
              seq_consensus
          },
          external := [IS occurrence::ExternalBioMat] {
              sequences,
              published_in,
              original_taxon,
              external_link,
              in_collection,
              item_vouchers,
              quantity,
              content_description
          }
      }
  );
  ALTER FUNCTION occurrence::insert_external_biomat(sampling: events::Sampling, data: std::json) USING (WITH
      identification := 
          (data)['identification']
      ,
      taxon := 
          taxonomy::taxonByName(<std::str>(identification)['taxon'])
  INSERT
      occurrence::ExternalBioMat
      {
          sampling := sampling,
          code := (<std::str>std::json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling)),
          sources := (SELECT
              references::DataSource
          FILTER
              (.code IN <std::str>std::json_array_unpack(std::json_get(data, 'sources')))
          ),
          external_link := <std::str>std::json_get(data, 'external_link'),
          original_taxon := <std::str>std::json_get(data, 'original_taxon'),
          quantity := <occurrence::QuantityType>std::json_get(data, 'quantity'),
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
          is_type := (<std::bool>std::json_get(data, 'is_type') ?? false)
      });
};
