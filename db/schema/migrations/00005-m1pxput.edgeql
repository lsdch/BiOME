CREATE MIGRATION m1pxputeid7p7qtytkscw46mjflozaq3nej2xynq6fwgcl32ekdzca
    ONTO m1tldgeyc7mrdwl3rdqtigbvgstsgkfsthp4vt7jy5u3f5gmahlnja
{
  ALTER TYPE location::Country {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::max_len_value(3);
      };
  };
  ALTER TYPE location::Country {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::max_len_value(2);
      };
  };
};
