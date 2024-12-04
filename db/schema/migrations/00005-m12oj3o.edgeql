CREATE MIGRATION m12oj3opdaq4wi74zsi7yheyuunvxukyeecspztj3wpcawgydbvcuq
    ONTO m1tizvfhsjc74f3stjynmssynymh5ssikvzczohlootiiyngxhnkwq
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
          external := std::assert_single((IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN {
              original_link := occurrence::ExternalBioMat.original_link,
              in_collection := occurrence::ExternalBioMat.in_collection,
              item_vouchers := occurrence::ExternalBioMat.item_vouchers,
              quantity := occurrence::ExternalBioMat.quantity,
              content_description := occurrence::ExternalBioMat.content_description
          } ELSE {}))
      }
  );
};
