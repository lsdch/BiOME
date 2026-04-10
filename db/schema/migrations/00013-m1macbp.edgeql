CREATE MIGRATION m1macbp2mbzhhbgblmehk4t2d6rnniwntgcproyqrpzfulawleg2eq
    ONTO m15jencp3tvkehhcuqtebe6j7i2zafv7ypggtpw6fo4obclfj6yvia
{
  ALTER TYPE taxonomy::Taxon {
      DROP PROPERTY anchor;
  };
};
