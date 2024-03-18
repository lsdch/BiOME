module taxonomy {

  scalar type Rank extending enum<Kingdom, Phylum, Class, Order, Family, Genus, Species, Subspecies>;
  # Ignore FORM, VARIETY, UNRANKED

  scalar type TaxonStatus extending enum<Accepted, Synonym, Unclassified>;

  type Taxon extending default::Auditable {
    GBIF_ID: int32 {
      constraint exclusive;
    };
    required name: str {
      constraint min_len_value(4);
    };

    constraint expression on (not contains(.name, " "))
      except (.rank = Rank.Species or .rank = Rank.Subspecies);
    constraint expression on (len(str_split(.name, " ")) = 2)
      except (.rank != Rank.Species);
    constraint expression on (len(str_split(.name, " ")) = 3)
      except (.rank != Rank.Subspecies);

    required rank: Rank;

    required status: TaxonStatus;
    required code: str {
      constraint exclusive;
      rewrite insert, update using (
        with chopped := str_split(.name, " "),
        suffix := "[syn]" if .status = TaxonStatus.Synonym else ""
        select (
          if (__specified__.code and len(.code) > 0) then .code
          else (
            if (not .rank in {Rank.Species, Rank.Subspecies}) then (.name ++ suffix)
            else (str_upper(chopped[0][:3]) ++ array_join(chopped[1:], "_") ++ suffix)
          )
        )
      )
    };
    required anchor: bool {
      annotation description := "Signals whether this taxon was manually imported";
      default := false;
    }
    authorship: str;

    kingdom: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Kingdom
          then .parent
          else .parent.kingdom
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    phylum: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Phylum
          then .parent
          else .parent.phylum
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    class: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Class
          then .parent
          else .parent.class
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    order: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Order
          then .parent
          else .parent.order
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    family: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Family
          then .parent
          else .parent.family
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    genus: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Genus
          then .parent
          else .parent.genus
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    species: Taxon {
      rewrite insert, update using (
        if exists .parent then (
          if .parent.rank = Rank.Species
          then .parent
          else .parent.species
        ) else <Taxon>{}
      );
      on target delete allow;
    };
    parent: Taxon {
      on target delete delete source;
    };
    constraint expression on (exists .parent) except (.rank = Rank.Kingdom);

    constraint exclusive on ((.name, .status));

    multi link children := .<parent[is Taxon];

    comment: str;

    index on (.name);
    index on (.rank);
    index on (.status);
  }
}
