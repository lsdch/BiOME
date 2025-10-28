CREATE MIGRATION m1p2hjulxab3ghtdfrafijxshetqmefrbxdvnchs4p5x2cqbalbuhq
    ONTO m16ybjkvw3prnvznuep5aw36f6gffcn24kbuf53jltu5t4je7lyvna
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          RESET readonly;
      };
  };
};
