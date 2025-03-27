CREATE MIGRATION m13345oyje6yqti7ta734iwitnnploqga5uzxzyeskoozibtypy7mq
    ONTO m1qhc5m73dsddd6uqldwbmkgfrdcyzkrwd6fvyjybvgjs2ndh5ieoq
{
  CREATE FUNCTION occurrence::biomat_code(taxon: taxonomy::Taxon, sampling: events::Sampling) ->  std::str USING ((((taxon.code ++ '[') ++ sampling.code) ++ ']'));
};
