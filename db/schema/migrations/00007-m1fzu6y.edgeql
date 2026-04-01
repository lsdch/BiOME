CREATE MIGRATION m1fzu6yijdfmvftfujbpeumxznhayv25pvrxflw6oelllqoaqor6sq
    ONTO m1xwftyhymfxfkajzscdbxw5rgkpaq6nytumqqrgarzhl52ycdlyga
{
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY scientific_name {
          USING ((.name ++ ((' ' ++ .authorship) ?? '')));
      };
  };
};
