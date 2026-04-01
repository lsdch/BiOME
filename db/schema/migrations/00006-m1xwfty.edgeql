CREATE MIGRATION m1xwftyhymfxfkajzscdbxw5rgkpaq6nytumqqrgarzhl52ycdlyga
    ONTO m1zrjyp3fro6qmrtbyf6ugee5tsijjjuhvmtlnrwi7736lpudq4dea
{
  ALTER TYPE taxonomy::Taxon {
      CREATE REQUIRED PROPERTY scientific_name := ((.name ++ (std::str_pad_start(.authorship, 1) ?? '')));
  };
};
