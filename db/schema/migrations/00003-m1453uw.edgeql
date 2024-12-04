CREATE MIGRATION m1453uwr6ge3m45sl7mlotihq276wmselda4twzs42a3ynovzys5dq
    ONTO m1ezxsuag3lmhzroybwww3qbledrfl4aeskzw7vvjnhymaff6hly6a
{
  ALTER ALIAS occurrence::AnyBioMaterial USING (SELECT
      occurrence::BioMaterial {
          *,
          published_in: {
              *
          },
          sampling: {
              *
          },
          identification: {
              *
          },
          meta: {
              *
          },
          type := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown')),
          external := (SELECT
              occurrence::BioMaterial[IS occurrence::ExternalBioMat] {
                  original_link,
                  in_collection,
                  item_vouchers,
                  quantity,
                  content_description
              }
          )
      }
  );
};
