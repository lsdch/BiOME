CREATE MIGRATION m1kufxx6w67wsocldhnobu3izg4jkvaxz74vgsoydcyosfmdbewdca
    ONTO m12odyfycilfqxqe73bcrkmxmfvdvnvtcz5zldssx6qucgiluioraq
{
  ALTER TYPE people::UserInvitation {
      EXTENDING people::Token LAST;
      ALTER PROPERTY expires {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY token {
          ALTER CONSTRAINT std::exclusive {
              DROP OWNED;
          };
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
};
