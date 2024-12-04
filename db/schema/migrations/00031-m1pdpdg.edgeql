CREATE MIGRATION m1pdpdgazingpqb75kiuk6znnpzfxmpq2tkjjxtkeonxbbffrrcgla
    ONTO m13oygx3jdkp2tgvertlvzajmjhqatwdkotqyngkmnfdxrknfro6oa
{
  ALTER TYPE datasets::Dataset {
      CREATE INDEX ON (.slug);
  };
  ALTER TYPE default::Vocabulary {
      CREATE INDEX ON ((.code, .label));
  };
  ALTER TYPE events::Program {
      CREATE INDEX ON ((.code, .label));
  };
  ALTER TYPE location::Site {
      CREATE INDEX ON (.code);
  };
  ALTER TYPE people::Institution {
      CREATE INDEX ON ((.code, .name));
  };
  ALTER TYPE people::Institution {
      DROP INDEX ON (.code);
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE INDEX ON (.email);
  };
  ALTER TYPE people::Person {
      CREATE INDEX ON ((.alias, .first_name, .last_name));
  };
  ALTER TYPE people::User {
      CREATE INDEX ON ((.email, .login));
  };
  ALTER TYPE samples::BioMaterial {
      CREATE INDEX ON (.code);
  };
  ALTER TYPE sampling::Habitat {
      CREATE INDEX ON (.label);
  };
  ALTER TYPE sampling::HabitatGroup {
      CREATE INDEX ON (.label);
  };
  ALTER TYPE taxonomy::Taxon {
      DROP INDEX ON (.name);
  };
  ALTER TYPE taxonomy::Taxon {
      DROP INDEX ON (.rank);
  };
  ALTER TYPE taxonomy::Taxon {
      DROP INDEX ON (.status);
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE INDEX ON ((.name, .code, .rank, .status));
  };
  ALTER TYPE tokens::UserInvitation {
      CREATE INDEX ON (.email);
  };
};
