CREATE MIGRATION m1xeampwczkkzmxspppuhrczxk6jybxqcg72zyn7tekvoppmxi6tva
    ONTO m1eqgei3f7r2v4eku3ka6i3bzzgiky2fqbx7ldw45nrzz3vrykrjoa
{
  ALTER SCALAR TYPE location::CoordinatesPrecision EXTENDING enum<`<100m`, `<1km`, `<10km`, `10-100km`, Unknown>;
};
