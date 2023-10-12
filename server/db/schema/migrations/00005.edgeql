CREATE MIGRATION m1kghg74bjihfngeirongjeooiyy5vvcfb5qa6osb66gs7kjb6lejq
    ONTO m16dj7xaefgahqazq34efjxkudsj3mpmhk5rdmolnrwlk3hgbov2gq
{
  ALTER TYPE taxonomy::Taxon {
      CREATE INDEX ON (.name);
      CREATE INDEX ON (.rank);
      CREATE INDEX ON (.status);
  };
};
