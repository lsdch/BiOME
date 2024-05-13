CREATE MIGRATION m1eegiinu6cy4hjigrrtjryokro6ntnzyef2yho4mojal6hti5fsnq
    ONTO m1jfssarnox2u433u7kkezp5kxw5ojnwk7sv6rdcpe4sfqr3brlbdq
{
  ALTER TYPE location::Habitat {
      ALTER LINK incompatible_from {
          ON TARGET DELETE ALLOW;
      };
  };
};
