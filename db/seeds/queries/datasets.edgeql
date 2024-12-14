with data := <json>$0
for item in json_array_unpack(data) union (
  with dataset := (insert datasets::Dataset {
    label := <str>item['label'],
    slug := <str>item['slug'],
    description := <str>json_get(item, 'description'),
    maintainers := (
      select people::Person filter .alias in <str>json_array_unpack(item['maintainers'])
    )
  }),
  sites := (for site in json_array_unpack(json_get(item, 'sites')) union (
    with created_site := (insert location::Site {
      datasets := dataset,
      name := <str>site['name'],
      code := <str>site['code'],
      locality := <str>json_get(site, 'locality'),
      country := (
        select location::Country filter .code = <str>json_get(site, 'country_code')
      ),
      description := <str>json_get(site, 'description'),
      coordinates := ((
        precision := <location::CoordinatesPrecision>site['coordinates']['precision'],
        latitude := <float32>site['coordinates']['latitude'],
        longitude := <float32>site['coordinates']['longitude'],
      ))
    }),
    events := (for event in json_array_unpack(json_get(site,'events')) union (
      with e := (insert events::Event {
        site := created_site,
        dataset := dataset,
        performed_on := ((
          date := <datetime>event['performed_on']['date'],
          precision := <date::DatePrecision>event['performed_on']['precision'],
        )),
        programs := (
          select events::Program
          filter .code in <str>json_array_unpack(json_get(event, 'programs'))
        ),
        performed_by := (
          select people::Person
          filter .alias in <str>json_array_unpack(json_get(event, 'performed_by'))
        )
      }),

      spotting := (
        if exists json_get(event,'spotting')
        then (
          insert events::Spotting {
            event := e,
            target_taxa := (
              select taxonomy::Taxon
              filter .code in <str>json_array_unpack(json_get(event, 'spotting', "target_taxa"))
            )
          }
        ) else {}
      ),

      abiotics := (
        for am in json_array_unpack(json_get(event, 'abiotic_measurements')) union (
          insert events::AbioticMeasurement {
            event := e,
            param := (
              select events::AbioticParameter
              filter .code = <str>am['param']
            ),
            value := <float32>am['value']
          }
        )
      ),

      samplings := (
        for sampling in json_array_unpack(json_get(event, 'samplings')) union (
        with s := (
          insert events::Sampling {
            event := e,
            fixatives := (
              select samples::Fixative
              filter .code in <str>json_array_unpack(json_get(sampling, 'fixatives'))
            ),
            methods := (
              select events::SamplingMethod
              filter .code in <str>json_array_unpack(json_get(sampling, 'methods'))
            ),
            sampling_target := <events::SamplingTarget>sampling['target']['kind'],
            target_taxa := (
              select taxonomy::Taxon
              filter .code in <str>json_array_unpack(json_get(sampling, 'target', 'target_taxa'))
            ),
            habitats := (
              select sampling::Habitat
              filter .label in <str>json_array_unpack(json_get(sampling, 'habitats'))
            ),
            access_points := <str>json_array_unpack(json_get(sampling, 'access_points')),
            sampling_duration := (
              <int32>json_get(sampling, 'duration')
            ),
            comments := <str>json_get(sampling, 'comments')
          }
        ),
        intbiomats := (for intbm in json_array_unpack(json_get(sampling, 'internal_biomat')) union (
          with identification := (
            insert occurrence::Identification {
              taxon := (select taxonomy::Taxon filter .name = <str>intbm['identification']['taxon']),
              identified_by := (select people::Person filter .alias = <str>json_get(intbm, 'identification', 'identified_by')),
              identified_on := ((
                date := <datetime>intbm['identification']['identified_on']['date'],
                precision := <date::DatePrecision>intbm['identification']['identified_on']['precision'],
              ))
            }
          )
          insert occurrence::InternalBioMat {
            sampling := s,
            code := ((identification.taxon.code) ++ "[" ++ (s.code) ++ "]"),
            comments := <str>json_get(intbm, "comments"),
            is_type := <bool>json_get(intbm, "is_type")?? false,
            published_in := (
              with pubs := json_array_unpack(json_get(intbm, 'published_in'))
              select assert_distinct((for p in pubs union (
                select references::Article {
                  @original_source := <bool>json_get(p, 'original')
                } filter .code = <str>p['code']
              )))
            ),
            identification := identification
          }
        )),
        extbiomats := (for extbm in json_array_unpack(json_get(sampling, 'external_biomat')) union (
          with identification := (
            insert occurrence::Identification {
              taxon := (select taxonomy::Taxon filter .name = <str>extbm['identification']['taxon']),
              identified_by := (select people::Person filter .alias = <str>json_get(extbm, 'identification', 'identified_by')),
              identified_on := ((
                date := <datetime>extbm['identification']['identified_on']['date'],
                precision := <date::DatePrecision>extbm['identification']['identified_on']['precision'],
              ))
            }
          ),
          ebm := (insert occurrence::ExternalBioMat {
            sampling := s,
            code := ((identification.taxon.code) ++ "[" ++ (s.code) ++ "]"),
            quantity := <occurrence::QuantityType>extbm['quantity'],
            content_description := <str>json_get(extbm, "content_description"),
            in_collection := <str>json_get(extbm, "in_collection"),
            item_vouchers := <str>json_array_unpack(json_get(extbm, "item_vouchers")),
            original_link := <str>json_get(extbm, "original_link"),
            comments := <str>json_get(extbm, "comments"),
            is_type := <bool>json_get(extbm, "is_type") ?? false,
            published_in := (
              with pubs := json_array_unpack(json_get(extbm, 'published_in'))
              select assert_distinct((for p in pubs union (
                select references::Article {
                  @original_source := <bool>json_get(p, 'original')
                } filter .code = <str>p['code']
              )))
            ),
            identification := identification
          }),
          ebm_seqs := (for exseq in json_array_unpack(json_get(extbm, 'sequences')) union (
            insert seq::ExternalSequence {
              sampling := s,
              source_sample := ebm,
              published_in := (
                select references::Article filter .code = <str>json_get(exseq, 'published_in')
              ),
              specimen_identifier := <str>exseq['specimen_identifier'],
              original_taxon := <str>json_get(exseq, 'original_taxon'),
              referenced_in := (
                for ref in json_array_unpack(json_get(exseq, 'referenced_in')) union (
                  insert seq::SeqReference {
                    accession := <str>ref['accession'],
                    is_origin := <bool>json_get(ref, 'is_origin') ?? false,
                    db := (select seq::SeqDB filter .code = <str>ref['db'])
                  }
                )
              ),
              code := <str>exseq['code'],
              origin := <seq::ExtSeqOrigin>exseq['origin'],
              gene := (
                select seq::Gene filter .code = <str>exseq['gene']
              ),
              label := <str>json_get(exseq, 'label'),
              comments := <str>json_get(exseq, 'comments'),
              sequence := <str>json_get(exseq, 'sequence'),
              identification := (
                insert occurrence::Identification {
                  taxon := (select taxonomy::Taxon filter .name = <str>exseq['identification']['taxon']),
                  identified_by := (select people::Person filter .alias = <str>json_get(exseq, 'identification', 'identified_by')),
                  identified_on := ((
                    date := <datetime>exseq['identification']['identified_on']['date'],
                    precision := <date::DatePrecision>exseq['identification']['identified_on']['precision'],
                  ))
                }
              )
            }
          )),
          select ebm
        )),
        ext_sequences := (for extseq in json_array_unpack(json_get(sampling, 'sequences')) union (
          insert seq::ExternalSequence {
            sampling:= s,
            code := <str>extseq['code'],
            origin := <seq::ExtSeqOrigin>extseq['origin'],
            label := <str>json_get(extseq, 'label'),
            gene := (
              select seq::Gene filter .code = <str>extseq['gene']
            ),
            comments := <str>json_get(extseq, 'comments'),
            sequence := <str>json_get(extseq, 'sequence'),
            referenced_in := (
              for ref in json_array_unpack(json_get(extseq, 'referenced_in')) union (
                insert seq::SeqReference {
                  accession := <str>ref['accession'],
                  is_origin := <bool>json_get(ref, 'is_origin') ?? false,
                  db := (select seq::SeqDB filter .code = <str>ref['db'])
                }
              )
            ),
            published_in := (
              select assert_single(references::Article filter .code in <str>json_get(extseq, 'published_in'))
            ),
            specimen_identifier := <str>json_get(extseq, 'specimen_identifier'),
            original_taxon := <str>json_get(extseq, 'original_taxon'),
            identification := (
              insert occurrence::Identification {
                taxon := (select taxonomy::Taxon filter .name = <str>extseq['identification']['taxon']),
                identified_by := (select people::Person filter .alias = <str>json_get(extseq, 'identification', 'identified_by')),
                identified_on := ((
                  date := <datetime>extseq['identification']['identified_on']['date'],
                  precision := <date::DatePrecision>extseq['identification']['identified_on']['precision'],
                ))
              }
            )
          }
        )
        )
        select s
      )),
      select e
    )),
    select created_site
  )),
  select dataset
);