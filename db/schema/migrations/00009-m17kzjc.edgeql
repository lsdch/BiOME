CREATE MIGRATION m17kzjcbuahadpg2fgi3vkgmcszgpcck3lb4c3pt6uy77z57xelnoq
    ONTO m1eoedhqazvntxpkwxdxweink5jmtpdvnnahsc4sriftlfxkb35kyq
{
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          SET default := (ext::pgcrypto::gen_salt());
      };
  };
};
