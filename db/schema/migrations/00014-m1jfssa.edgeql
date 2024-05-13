CREATE MIGRATION m1jfssarnox2u433u7kkezp5kxw5ojnwk7sv6rdcpe4sfqr3brlbdq
    ONTO m14tp3p5svkgciveqdodoi5j2rcul7hsykjrzkriv63u2dwl2532ba
{
  ALTER TYPE location::Habitat {
      ALTER LINK in_group {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
};
