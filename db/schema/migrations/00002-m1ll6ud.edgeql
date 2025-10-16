CREATE MIGRATION m1ll6udrqcr4bfl6ywryw6aukliz6msicu3d2evmprky3p4it6qczq
    ONTO m1y7fygges5bpfvlvdjca3wuuwjrj6fqw5d5jijhrljefukykhcgxa
{
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
              DISTINCT (references::DataSource)
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
                  identified_by := people::personByAlias(<std::str>(identification)['identified_by']),
                  identified_on := date::from_json_with_precision((identification)['identified_on'])
              }),
          is_type := (<std::bool>std::json_get(data, 'is_type') ?? false)
      });
};
