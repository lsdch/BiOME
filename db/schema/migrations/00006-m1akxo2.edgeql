CREATE MIGRATION m1akxo2vcp7okt344m7cj2h2gc6egtacn3msh4zkvkw4pfxgrnequq
    ONTO m1bd75yvmwip7cx7aac3h3rjn3bkjn5brf3gefmxblo26pif6y73ca
{
  ALTER TYPE seq::Sequence {
      CREATE REQUIRED LINK identification: occurrence::Identification {
          SET REQUIRED USING (<occurrence::Identification>{});
      };
  };
};
