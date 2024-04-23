CREATE MIGRATION m1elskjwstzjcfust7yjko6jv53exahqsv5r5fuzoxrbyfvvu5ctga
    ONTO m1otljal5ypsejcqyp4axj4prfz5nw5criwr5e7jop7zxj7r7q5uya
{
  ALTER TYPE people::UserInvitation {
      ALTER PROPERTY dest {
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
                  SET REQUIRED USING (<std::str>{});
              };
          };
};
