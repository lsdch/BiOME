CREATE MIGRATION m1zrjyp3fro6qmrtbyf6ugee5tsijjjuhvmtlnrwi7736lpudq4dea
    ONTO m15czjri44kvqhifef3m5p53s3roabgaozwtpocw7hloxswcmnia6q
{
  ALTER SCALAR TYPE taxonomy::TaxonStatus EXTENDING enum<Accepted, Synonym, Doubtful, Unreferenced, Unclassified>;
};
