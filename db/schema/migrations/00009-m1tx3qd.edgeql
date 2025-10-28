CREATE MIGRATION m1tx3qdxao3zkirxehdavlzn6dux3ajqteyjmsqavnrmiftrkgf6iq
    ONTO m1p2hjulxab3ghtdfrafijxshetqmefrbxdvnchs4p5x2cqbalbuhq
{
  ALTER TYPE occurrence::Identification {
      CREATE PROPERTY addendum: std::str {
          CREATE ANNOTATION std::description := "Additional information about the identification, e.g. 'group venustus' or 'form A'.";
      };
  };
};
