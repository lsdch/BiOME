CREATE MIGRATION m1eb5xvjiif2cw45hdn7aynjcojpsdzydjebeckbaika3347ikkscq
    ONTO m1yv6b46o6qyge2mdxsarylysutuzy3wngbnzwaissokz5dwj7gxjq
{
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK kingdom {
          CREATE REWRITE
              INSERT 
              USING ((.parent.kingdom IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.kingdom IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
  };
};
