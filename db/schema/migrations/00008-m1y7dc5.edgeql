CREATE MIGRATION m1y7dc5hxzyu74a62utzt3hl6evgnlfxmv3qyrxwbll4recrvjcsqq
    ONTO m1fzu6yijdfmvftfujbpeumxznhayv25pvrxflw6oelllqoaqor6sq
{
  ALTER TYPE taxonomy::Taxon {
      DROP PROPERTY scientific_name;
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE REQUIRED PROPERTY scientific_name: std::str {
          SET default := ((.name ++ ((' ' ++ .authorship) ?? '')));
          CREATE CONSTRAINT std::exclusive;
          CREATE REWRITE
              INSERT 
              USING ((.name ++ ((' ' ++ .authorship) ?? '')));
          CREATE REWRITE
              UPDATE 
              USING ((.name ++ ((' ' ++ .authorship) ?? '')));
      };
  };
};
