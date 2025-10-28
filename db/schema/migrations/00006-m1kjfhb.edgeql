CREATE MIGRATION m1kjfhbmuwa65j4byyzeato7bzon6ydqxjpl2kttrz45izfmwqddba
    ONTO m1nkgduvf2df6eqyrsimwsqmfpsbeu25obmxcsy5klvv2eedzjxvsa
{
  CREATE FUNCTION occurrence::insert_external_occurrence(sampling: events::Sampling, data: std::json) ->  occurrence::ExternalOccurrence USING (WITH
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
  CREATE FUNCTION occurrence::insert_internal_biomat(sampling: events::Sampling, data: std::json) ->  occurrence::InternalBioMat USING (WITH
      identification := 
          (data)['identification']
      ,
      taxon := 
          taxonomy::taxonByName(<std::str>(identification)['taxon'])
  INSERT
      occurrence::InternalBioMat
      {
          sampling := sampling,
          code := (<std::str>std::json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling.code)),
          identification := (INSERT
              occurrence::Identification
              {
                  taxon := taxon,
                  identified_by := people::personByAlias(<std::str>(identification)['identified_by']),
                  identified_on := date::from_json_with_precision((identification)['identified_on'])
              }),
          is_type := <std::bool>std::json_get(data, 'is_type'),
          comments := <std::str>std::json_get(data, 'comments')
      });
};
