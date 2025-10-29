CREATE MIGRATION m1c3enlyoxfvhjbyemzecykb6h77ofsfihvmc7duqtxwqgtbplh6aa
    ONTO m146g2xi3yhebkb772izc3tckjon7z4btjv3ug6zmda6ouwcaocftq
{
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY confer {
          SET REQUIRED USING (<std::bool>false);
      };
  };
};
