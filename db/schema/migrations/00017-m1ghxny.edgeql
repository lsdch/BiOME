CREATE MIGRATION m1ghxnyikde74qaaa3ltk2d3wgqt6ozqki4g2tqe3mtiafc7ebg6sa
    ONTO m1sfktqdvf2xqrwx6af2sb5fs53vde2jcmih5t4ah32ica6epr5ztq
{
  ALTER TYPE occurrence::Identification {
      DROP PROPERTY qualifier;
  };
  DROP SCALAR TYPE occurrence::IdentificationQualifier;
};
