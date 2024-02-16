CREATE MIGRATION m1plhlejp6dzzfr5hwkfv3yw4lznnkezriixvaf2svvvdtdehjxf7q
    ONTO m1qpv4cdoc6rnpda74fzszmsksnliokqwrnl5xnbpozpsm6fnk6asa
{
  ALTER TYPE people::Person {
      ALTER PROPERTY full_name {
          USING (std::array_join(std::array_agg({.first_name, .middle_names, .last_name}), ' '));
          SET REQUIRED USING (<std::str>{});
      };
  };
};
