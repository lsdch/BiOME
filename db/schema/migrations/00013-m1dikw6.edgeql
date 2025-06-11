CREATE MIGRATION m1dikw6komtq5cnna6zqflk54353gkmxkseo5iiix3xoj4dg47uzyq
    ONTO m1maa34zvuro4hoybz3pfnev3mgj7pbp7tgoa7oicqmez3aod6jxka
{
  ALTER TYPE location::Site {
      ALTER PROPERTY last_visited {
          SET SINGLE;
      };
  };
};
