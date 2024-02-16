CREATE MIGRATION m1qpv4cdoc6rnpda74fzszmsksnliokqwrnl5xnbpozpsm6fnk6asa
    ONTO m1vkutotun7xrnliqicc7tj7px2jcomxspwj3eahgdz6hfynd6p7xq
{
  ALTER TYPE people::Person {
      ALTER LINK institution {
          RENAME TO institutions;
      };
  };
};
