CREATE MIGRATION m1wytcfnhahhtbe72lczcaaueqjdyazr7b65x2vecflwvpvb5qqcoq
    ONTO m1rg77poeapfh3u77ttpitjsmh7ciw3d253a3hjkqdjdvqiblzv77a
{
  ALTER TYPE people::EmailConfirmation {
      ALTER PROPERTY email {
          CREATE CONSTRAINT std::exclusive;
      };
  };
};
