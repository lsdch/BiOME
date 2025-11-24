CREATE MIGRATION m1vvlq6acjo7mhvqkq3q5ageyjqwt27x7hshzxqmnpuresjtze6lxa
    ONTO m1iuvfvffsf5ttu3y2ak5462llvvji3r66xfot26h5uopwzau6iarq
{
  ALTER TYPE occurrence::ExternalOccurrence {
      CREATE PROPERTY quantity: range<std::int32>;
  };
  ALTER ALIAS occurrence::OccurrenceWithType USING (SELECT
      occurrence::Occurrence {
          required has_sequences := ([IS occurrence::InternalBioMat].has_sequences ?? ([IS occurrence::ExternalOccurrence].has_sequences ?? false)),
          internal := [IS occurrence::InternalBioMat] {
              is_homogenous,
              is_congruent,
              seq_consensus
          },
          external := [IS occurrence::ExternalOccurrence] {
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
          quantity := <range<std::int32>>std::json_get(data, 'quantity'),
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
};
