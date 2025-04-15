CREATE MIGRATION m1wl3ufuahzlvzeebode3yyqmpqgs4vzk3dxf3mntqvduifeugocga
    ONTO m1hmseivjwsdw3q4tqbktdgtpokjprry2y36k3zis777zgwbsjhdla
{
  ALTER FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) {
      RENAME TO taxonomy::is_in_clade;
  };
};
