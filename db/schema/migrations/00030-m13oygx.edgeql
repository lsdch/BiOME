CREATE MIGRATION m13oygx3jdkp2tgvertlvzajmjhqatwdkotqyngkmnfdxrknfro6oa
    ONTO m12xvo3lo4db2q6c7tkxvwu5hhpfkhidiwwwzfy3cnx4r5kiynwveq
{
  ALTER TYPE default::Vocabulary {
      ALTER PROPERTY code {
          ALTER CONSTRAINT std::exclusive {
              SET DELEGATED;
          };
      };
  };
};
