CREATE MIGRATION m1tizvfhsjc74f3stjynmssynymh5ssikvzczohlootiiyngxhnkwq
    ONTO m1453uwr6ge3m45sl7mlotihq276wmselda4twzs42a3ynovzys5dq
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
          external := (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN {
              original_link := occurrence::ExternalBioMat.original_link,
              in_collection := occurrence::ExternalBioMat.in_collection,
              item_vouchers := occurrence::ExternalBioMat.item_vouchers,
              quantity := occurrence::ExternalBioMat.quantity,
              content_description := occurrence::ExternalBioMat.content_description
          } ELSE {})
      }
  );
};
