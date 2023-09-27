CREATE MIGRATION m14ytq2iiu6hdrp5vuhphjzouphbytocdfely5fjxozsffv6hkhyaa
    ONTO m17dalqln6hw3ab2ljhv5n2rpabbpki3wle7azbigal47cpb7etk5q
{
  ALTER TYPE taxonomy::Taxon {
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 3)) EXCEPT ((.rank != taxonomy::Rank.Subspecies));
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 2)) EXCEPT ((.rank != taxonomy::Rank.Species));
  };
};
