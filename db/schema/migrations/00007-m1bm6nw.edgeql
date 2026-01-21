CREATE MIGRATION m1bm6nwr4rz4s2byr7ybmzqljuc4vtfs6r4yom3kaarprpt3vrwnua
    ONTO m14dkru6mctwrtyf65fswdkahynnjtef2axylmyxoadhytuyzsihia
{
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY identified_by {
          SET MULTI;
      };
  };
};
