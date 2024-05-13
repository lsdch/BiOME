CREATE MIGRATION m14tp3p5svkgciveqdodoi5j2rcul7hsykjrzkriv63u2dwl2532ba
    ONTO m1xm6ztndgvghtb7vh6i3fwxpfla3a3afgtmdgicimp423zp7wxmeq
{
  ALTER TYPE location::Habitat {
      ALTER LINK incompatibleFrom {
          RENAME TO incompatible_from;
      };
  };
};
