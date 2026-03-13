CREATE MIGRATION m1kva3mqa6keqqabpj4tdc4cmcckhxgwy4u3hgtwoldth3o7qwkjra
    ONTO m1z5nazo2qcdec4ru2hv6356o73qxfgo27qmtcctzdepv6gnlqe6rq
{
  ALTER TYPE occurrence::Occurrence {
      DROP LINK collections;
  };
  ALTER TYPE occurrence::Occurrence {
      CREATE MULTI PROPERTY collections: tuple<name: std::str, vouchers: array<std::str>> {
          CREATE ANNOTATION std::description := 'Collections where the specimen(s) can be found, with optional voucher identifiers.';
      };
  };
  DROP TYPE references::Collection;
};
