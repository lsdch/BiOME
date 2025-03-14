CREATE MIGRATION m12lk2vyqr7ey6ddbd4tbzpqthemqjtc6kcmu53ojumrwsdkxqplpq
    ONTO m1pxputeid7p7qtytkscw46mjflozaq3nej2xynq6fwgcl32ekdzca
{
  ALTER TYPE location::Country {
      CREATE REQUIRED PROPERTY boundaries: ext::postgis::geometry {
          SET REQUIRED USING (<ext::postgis::geometry>{});
      };
  };
};
