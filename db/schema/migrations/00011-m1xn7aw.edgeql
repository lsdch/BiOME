CREATE MIGRATION m1xn7awaoifsolgq55fpssj4556k4gkcd5476qhxxxysr65sy6fj7a
    ONTO m1muhdccmtusoyegmsrid2qgzjsbx4yvcnwyiw3mwt4dw37uxb7pxa
{
  ALTER FUNCTION taxonomy::is_in_clade(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (SELECT
      ((ancestor IN taxon.lineage) OR (taxon = ancestor))
  );
};
