CREATE MIGRATION m16mzxn3hd5q3fcxsihhiamrcp3npff3bxz2ks6nbpaxy6hbnm3tua
    ONTO m1fgjgxn27luncdge6xbyn4d3i3mpf3ita2ue4kpls2nr6fupf53ja
{
  ALTER TYPE occurrence::BioMaterial {
      CREATE REQUIRED PROPERTY is_type: std::bool {
          SET default := false;
      };
  };
  ALTER TYPE occurrence::Identification {
      DROP PROPERTY is_type;
  };
};
