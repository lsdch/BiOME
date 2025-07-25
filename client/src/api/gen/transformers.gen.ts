// This file is auto-generated by @hey-api/openapi-ts

import type {
  ListAbioticParametersResponse,
  CreateAbioticParameterResponse,
  DeleteAbioticParameterResponse,
  UpdateAbioticParameterResponse,
  LoginResponse,
  ListPendingUserRequestsResponse,
  DeletePendingUserRequestResponse,
  GetPendingUserRequestResponse,
  RefreshSessionResponse,
  ClaimInvitationResponse,
  ListAnchorsResponse,
  ListBioMaterialResponse,
  UpdateExternalBioMatResponse,
  CreateExternalBioMatResponse,
  DeleteBioMaterialResponse,
  GetBioMaterialResponse,
  ListDataSourcesResponse,
  CreateDataSourceResponse,
  DeleteDataSourceResponse,
  UpdateDataSourceResponse,
  ListDatasetsResponse,
  UpdateDatasetResponse,
  ListOccurrenceDatasetsResponse,
  GetOccurrenceDatasetResponse,
  TogglePinDatasetResponse,
  ListSequenceDatasetsResponse,
  GetSequenceDatasetResponse,
  ListSiteDatasetsResponse,
  CreateSiteDatasetResponse,
  GetSiteDatasetResponse,
  GetDatasetResponse,
  DeleteEventResponse,
  UpdateEventResponse,
  EventAddExternalOccurrenceResponse,
  CreateSamplingAtEventResponse,
  UpdateSpottingResponse,
  ListFixativesResponse,
  CreateFixativeResponse,
  DeleteFixativeResponse,
  UpdateFixativeResponse,
  ListGenesResponse,
  CreateGeneResponse,
  DeleteGeneResponse,
  UpdateGeneResponse,
  ListHabitatGroupsResponse,
  CreateHabitatGroupResponse,
  DeleteHabitatGroupResponse,
  UpdateHabitatGroupResponse,
  SitesProximityResponse,
  SearchSitesResponse,
  OccurrencesBySiteResponse,
  ListOrganisationsResponse,
  CreateOrganisationResponse,
  DeleteOrganisationResponse,
  UpdateOrganisationResponse,
  ListPersonsResponse,
  CreatePersonResponse,
  DeletePersonResponse,
  UpdatePersonResponse,
  ListProgramsResponse,
  CreateProgramResponse,
  DeleteProgramResponse,
  UpdateProgramResponse,
  ListArticlesResponse,
  CreateArticleResponse,
  DeleteArticleResponse,
  UpdateArticleResponse,
  ListSamplingMethodsResponse,
  CreateSamplingMethodResponse,
  DeleteSamplingMethodResponse,
  UpdateSamplingMethodResponse,
  CreateSamplingResponse,
  DeleteSamplingResponse,
  UpdateSamplingResponse,
  SamplingAddExternalOccurrenceResponse,
  ListSequencesResponse,
  DeleteSequenceResponse,
  GetSequenceResponse,
  ListDataFeedsResponse,
  CreateDataFeedResponse,
  ListMapPresetsResponse,
  CreateUpdateMapPresetResponse,
  DeleteMapPresetResponse,
  ListSitesResponse,
  CreateSiteResponse,
  GetSiteResponse,
  UpdateSiteResponse,
  ListSiteEventsResponse,
  CreateEventResponse,
  SiteAddExternalOccurrenceResponse,
  GetTaxonomyResponse,
  ListTaxaResponse,
  CreateTaxonResponse,
  DeleteTaxonResponse,
  GetTaxonResponse,
  UpdateTaxonResponse
} from './types.gen'

const metaSchemaResponseTransformer = (data: any) => {
  data.created = new Date(data.created)
  data.last_updated = new Date(data.last_updated)
  if (data.modified) {
    data.modified = new Date(data.modified)
  }
  return data
}

const abioticParameterSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listAbioticParametersResponseTransformer = async (
  data: any
): Promise<ListAbioticParametersResponse> => {
  data = data.map((item: any) => {
    return abioticParameterSchemaResponseTransformer(item)
  })
  return data
}

export const createAbioticParameterResponseTransformer = async (
  data: any
): Promise<CreateAbioticParameterResponse> => {
  data = abioticParameterSchemaResponseTransformer(data)
  return data
}

export const deleteAbioticParameterResponseTransformer = async (
  data: any
): Promise<DeleteAbioticParameterResponse> => {
  data = abioticParameterSchemaResponseTransformer(data)
  return data
}

export const updateAbioticParameterResponseTransformer = async (
  data: any
): Promise<UpdateAbioticParameterResponse> => {
  data = abioticParameterSchemaResponseTransformer(data)
  return data
}

const authenticationResponseSchemaResponseTransformer = (data: any) => {
  data.auth_token_expiration = new Date(data.auth_token_expiration)
  return data
}

export const loginResponseTransformer = async (data: any): Promise<LoginResponse> => {
  data = authenticationResponseSchemaResponseTransformer(data)
  return data
}

const pendingUserRequestSchemaResponseTransformer = (data: any) => {
  data.created_on = new Date(data.created_on)
  return data
}

export const listPendingUserRequestsResponseTransformer = async (
  data: any
): Promise<ListPendingUserRequestsResponse> => {
  data = data.map((item: any) => {
    return pendingUserRequestSchemaResponseTransformer(item)
  })
  return data
}

export const deletePendingUserRequestResponseTransformer = async (
  data: any
): Promise<DeletePendingUserRequestResponse> => {
  data = pendingUserRequestSchemaResponseTransformer(data)
  return data
}

export const getPendingUserRequestResponseTransformer = async (
  data: any
): Promise<GetPendingUserRequestResponse> => {
  data = pendingUserRequestSchemaResponseTransformer(data)
  return data
}

export const refreshSessionResponseTransformer = async (
  data: any
): Promise<RefreshSessionResponse> => {
  data = authenticationResponseSchemaResponseTransformer(data)
  return data
}

export const claimInvitationResponseTransformer = async (
  data: any
): Promise<ClaimInvitationResponse> => {
  data = authenticationResponseSchemaResponseTransformer(data)
  return data
}

const taxonWithParentRefSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listAnchorsResponseTransformer = async (data: any): Promise<ListAnchorsResponse> => {
  data = data.map((item: any) => {
    return taxonWithParentRefSchemaResponseTransformer(item)
  })
  return data
}

const codeHistorySchemaResponseTransformer = (data: any) => {
  data.time = new Date(data.time)
  return data
}

const dateWithPrecisionSchemaResponseTransformer = (data: any) => {
  if (data.date) {
    data.date = new Date(data.date)
  }
  return data
}

const optionalDateWithPrecisionSchemaResponseTransformer = (data: any) => {
  if (data.date) {
    data.date = new Date(data.date)
  }
  return data
}

const siteItemSchemaResponseTransformer = (data: any) => {
  if (data.last_visited) {
    data.last_visited = optionalDateWithPrecisionSchemaResponseTransformer(data.last_visited)
  }
  return data
}

const eventWithParticipantsSchemaResponseTransformer = (data: any) => {
  data.performed_on = dateWithPrecisionSchemaResponseTransformer(data.performed_on)
  data.site = siteItemSchemaResponseTransformer(data.site)
  return data
}

const geneSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const taxonSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const identificationSchemaResponseTransformer = (data: any) => {
  data.identified_on = dateWithPrecisionSchemaResponseTransformer(data.identified_on)
  data.meta = metaSchemaResponseTransformer(data.meta)
  data.taxon = taxonSchemaResponseTransformer(data.taxon)
  return data
}

const articleSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const dataSourceSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const seqReferenceSchemaResponseTransformer = (data: any) => {
  data.db = dataSourceSchemaResponseTransformer(data.db)
  return data
}

const externalBioMatSequenceSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  data.gene = geneSchemaResponseTransformer(data.gene)
  data.identification = identificationSchemaResponseTransformer(data.identification)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return articleSchemaResponseTransformer(item)
    })
  }
  if (data.referenced_in) {
    data.referenced_in = data.referenced_in.map((item: any) => {
      return seqReferenceSchemaResponseTransformer(item)
    })
  }
  return data
}

const externalBioMatContentSchemaResponseTransformer = (data: any) => {
  data.sequences = data.sequences.map((item: any) => {
    return externalBioMatSequenceSchemaResponseTransformer(item)
  })
  return data
}

const optionalDataSourceSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const optionalExternalBioMatSpecificSchemaResponseTransformer = (data: any) => {
  if (data.content) {
    data.content = data.content.map((item: any) => {
      return externalBioMatContentSchemaResponseTransformer(item)
    })
  }
  if (data.original_source) {
    data.original_source = optionalDataSourceSchemaResponseTransformer(data.original_source)
  }
  return data
}

const occurrenceReferenceSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const fixativeSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const habitatSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const samplingMethodSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const samplingTargetSchemaResponseTransformer = (data: any) => {
  if (data.taxa) {
    data.taxa = data.taxa.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  return data
}

const samplingInnerSchemaResponseTransformer = (data: any) => {
  if (data.fixatives) {
    data.fixatives = data.fixatives.map((item: any) => {
      return fixativeSchemaResponseTransformer(item)
    })
  }
  if (data.habitats) {
    data.habitats = data.habitats.map((item: any) => {
      return habitatSchemaResponseTransformer(item)
    })
  }
  if (data.methods) {
    data.methods = data.methods.map((item: any) => {
      return samplingMethodSchemaResponseTransformer(item)
    })
  }
  data.target = samplingTargetSchemaResponseTransformer(data.target)
  return data
}

const optionalTaxonSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

const bioMaterialSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  if (data.external) {
    data.external = optionalExternalBioMatSpecificSchemaResponseTransformer(data.external)
  }
  data.identification = identificationSchemaResponseTransformer(data.identification)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  data.sampling = samplingInnerSchemaResponseTransformer(data.sampling)
  if (data.seq_consensus) {
    data.seq_consensus = optionalTaxonSchemaResponseTransformer(data.seq_consensus)
  }
  return data
}

const samplingSchemaResponseTransformer = (data: any) => {
  if (data.fixatives) {
    data.fixatives = data.fixatives.map((item: any) => {
      return fixativeSchemaResponseTransformer(item)
    })
  }
  if (data.habitats) {
    data.habitats = data.habitats.map((item: any) => {
      return habitatSchemaResponseTransformer(item)
    })
  }
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.methods) {
    data.methods = data.methods.map((item: any) => {
      return samplingMethodSchemaResponseTransformer(item)
    })
  }
  if (data.occurring_taxa) {
    data.occurring_taxa = data.occurring_taxa.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  if (data.samples) {
    data.samples = data.samples.map((item: any) => {
      return bioMaterialSchemaResponseTransformer(item)
    })
  }
  data.target = samplingTargetSchemaResponseTransformer(data.target)
  return data
}

const bioMaterialWithDetailsSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  data.event = eventWithParticipantsSchemaResponseTransformer(data.event)
  if (data.external) {
    data.external = optionalExternalBioMatSpecificSchemaResponseTransformer(data.external)
  }
  data.identification = identificationSchemaResponseTransformer(data.identification)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  data.sampling = samplingSchemaResponseTransformer(data.sampling)
  if (data.seq_consensus) {
    data.seq_consensus = optionalTaxonSchemaResponseTransformer(data.seq_consensus)
  }
  return data
}

const paginatedListBioMaterialWithDetailsSchemaResponseTransformer = (data: any) => {
  data.items = data.items.map((item: any) => {
    return bioMaterialWithDetailsSchemaResponseTransformer(item)
  })
  return data
}

export const listBioMaterialResponseTransformer = async (
  data: any
): Promise<ListBioMaterialResponse> => {
  data = paginatedListBioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const updateExternalBioMatResponseTransformer = async (
  data: any
): Promise<UpdateExternalBioMatResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const createExternalBioMatResponseTransformer = async (
  data: any
): Promise<CreateExternalBioMatResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const deleteBioMaterialResponseTransformer = async (
  data: any
): Promise<DeleteBioMaterialResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const getBioMaterialResponseTransformer = async (
  data: any
): Promise<GetBioMaterialResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const listDataSourcesResponseTransformer = async (
  data: any
): Promise<ListDataSourcesResponse> => {
  data = data.map((item: any) => {
    return dataSourceSchemaResponseTransformer(item)
  })
  return data
}

export const createDataSourceResponseTransformer = async (
  data: any
): Promise<CreateDataSourceResponse> => {
  data = dataSourceSchemaResponseTransformer(data)
  return data
}

export const deleteDataSourceResponseTransformer = async (
  data: any
): Promise<DeleteDataSourceResponse> => {
  data = dataSourceSchemaResponseTransformer(data)
  return data
}

export const updateDataSourceResponseTransformer = async (
  data: any
): Promise<UpdateDataSourceResponse> => {
  data = dataSourceSchemaResponseTransformer(data)
  return data
}

const datasetSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listDatasetsResponseTransformer = async (data: any): Promise<ListDatasetsResponse> => {
  data = data.map((item: any) => {
    return datasetSchemaResponseTransformer(item)
  })
  return data
}

export const updateDatasetResponseTransformer = async (
  data: any
): Promise<UpdateDatasetResponse> => {
  data = datasetSchemaResponseTransformer(data)
  return data
}

const occurrenceDatasetListItemSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listOccurrenceDatasetsResponseTransformer = async (
  data: any
): Promise<ListOccurrenceDatasetsResponse> => {
  data = data.map((item: any) => {
    return occurrenceDatasetListItemSchemaResponseTransformer(item)
  })
  return data
}

const samplingEventWithOccurrencesSchemaResponseTransformer = (data: any) => {
  data.date = dateWithPrecisionSchemaResponseTransformer(data.date)
  if (data.occurring_taxa) {
    data.occurring_taxa = data.occurring_taxa.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  data.target = samplingTargetSchemaResponseTransformer(data.target)
  return data
}

const siteWithOccurrencesSchemaResponseTransformer = (data: any) => {
  if (data.last_visited) {
    data.last_visited = optionalDateWithPrecisionSchemaResponseTransformer(data.last_visited)
  }
  data.samplings = data.samplings.map((item: any) => {
    return samplingEventWithOccurrencesSchemaResponseTransformer(item)
  })
  return data
}

const occurrenceDatasetSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  data.sites = data.sites.map((item: any) => {
    return siteWithOccurrencesSchemaResponseTransformer(item)
  })
  return data
}

export const getOccurrenceDatasetResponseTransformer = async (
  data: any
): Promise<GetOccurrenceDatasetResponse> => {
  data = occurrenceDatasetSchemaResponseTransformer(data)
  return data
}

export const togglePinDatasetResponseTransformer = async (
  data: any
): Promise<TogglePinDatasetResponse> => {
  data = datasetSchemaResponseTransformer(data)
  return data
}

const eventInnerSchemaResponseTransformer = (data: any) => {
  data.performed_on = dateWithPrecisionSchemaResponseTransformer(data.performed_on)
  data.site = siteItemSchemaResponseTransformer(data.site)
  return data
}

const optionalBioMaterialSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  if (data.external) {
    data.external = optionalExternalBioMatSpecificSchemaResponseTransformer(data.external)
  }
  data.identification = identificationSchemaResponseTransformer(data.identification)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  data.sampling = samplingInnerSchemaResponseTransformer(data.sampling)
  if (data.seq_consensus) {
    data.seq_consensus = optionalTaxonSchemaResponseTransformer(data.seq_consensus)
  }
  return data
}

const optionalExtSeqSpecificsBioMaterialSchemaResponseTransformer = (data: any) => {
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  if (data.referenced_in) {
    data.referenced_in = data.referenced_in.map((item: any) => {
      return seqReferenceSchemaResponseTransformer(item)
    })
  }
  if (data.source_sample) {
    data.source_sample = optionalBioMaterialSchemaResponseTransformer(data.source_sample)
  }
  return data
}

const sequenceSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  data.event = eventInnerSchemaResponseTransformer(data.event)
  if (data.external) {
    data.external = optionalExtSeqSpecificsBioMaterialSchemaResponseTransformer(data.external)
  }
  data.gene = geneSchemaResponseTransformer(data.gene)
  data.identification = identificationSchemaResponseTransformer(data.identification)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  data.sampling = samplingInnerSchemaResponseTransformer(data.sampling)
  return data
}

const sequenceDatasetSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  data.sequences = data.sequences.map((item: any) => {
    return sequenceSchemaResponseTransformer(item)
  })
  data.sites = data.sites.map((item: any) => {
    return siteItemSchemaResponseTransformer(item)
  })
  return data
}

export const listSequenceDatasetsResponseTransformer = async (
  data: any
): Promise<ListSequenceDatasetsResponse> => {
  data = data.map((item: any) => {
    return sequenceDatasetSchemaResponseTransformer(item)
  })
  return data
}

export const getSequenceDatasetResponseTransformer = async (
  data: any
): Promise<GetSequenceDatasetResponse> => {
  data = sequenceDatasetSchemaResponseTransformer(data)
  return data
}

const siteDatasetSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  data.sites = data.sites.map((item: any) => {
    return siteItemSchemaResponseTransformer(item)
  })
  return data
}

export const listSiteDatasetsResponseTransformer = async (
  data: any
): Promise<ListSiteDatasetsResponse> => {
  data = data.map((item: any) => {
    return siteDatasetSchemaResponseTransformer(item)
  })
  return data
}

export const createSiteDatasetResponseTransformer = async (
  data: any
): Promise<CreateSiteDatasetResponse> => {
  data = siteDatasetSchemaResponseTransformer(data)
  return data
}

export const getSiteDatasetResponseTransformer = async (
  data: any
): Promise<GetSiteDatasetResponse> => {
  data = siteDatasetSchemaResponseTransformer(data)
  return data
}

export const getDatasetResponseTransformer = async (data: any): Promise<GetDatasetResponse> => {
  data = datasetSchemaResponseTransformer(data)
  return data
}

const abioticMeasurementSchemaResponseTransformer = (data: any) => {
  data.param = abioticParameterSchemaResponseTransformer(data.param)
  return data
}

const eventSchemaResponseTransformer = (data: any) => {
  if (data.abiotic_measurements) {
    data.abiotic_measurements = data.abiotic_measurements.map((item: any) => {
      return abioticMeasurementSchemaResponseTransformer(item)
    })
  }
  data.meta = metaSchemaResponseTransformer(data.meta)
  data.performed_on = dateWithPrecisionSchemaResponseTransformer(data.performed_on)
  if (data.samplings) {
    data.samplings = data.samplings.map((item: any) => {
      return samplingSchemaResponseTransformer(item)
    })
  }
  data.site = siteItemSchemaResponseTransformer(data.site)
  if (data.spottings) {
    data.spottings = data.spottings.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  return data
}

export const deleteEventResponseTransformer = async (data: any): Promise<DeleteEventResponse> => {
  data = eventSchemaResponseTransformer(data)
  return data
}

export const updateEventResponseTransformer = async (data: any): Promise<UpdateEventResponse> => {
  data = eventSchemaResponseTransformer(data)
  return data
}

export const eventAddExternalOccurrenceResponseTransformer = async (
  data: any
): Promise<EventAddExternalOccurrenceResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const createSamplingAtEventResponseTransformer = async (
  data: any
): Promise<CreateSamplingAtEventResponse> => {
  data = samplingSchemaResponseTransformer(data)
  return data
}

export const updateSpottingResponseTransformer = async (
  data: any
): Promise<UpdateSpottingResponse> => {
  data = data.map((item: any) => {
    return taxonSchemaResponseTransformer(item)
  })
  return data
}

export const listFixativesResponseTransformer = async (
  data: any
): Promise<ListFixativesResponse> => {
  data = data.map((item: any) => {
    return fixativeSchemaResponseTransformer(item)
  })
  return data
}

export const createFixativeResponseTransformer = async (
  data: any
): Promise<CreateFixativeResponse> => {
  data = fixativeSchemaResponseTransformer(data)
  return data
}

export const deleteFixativeResponseTransformer = async (
  data: any
): Promise<DeleteFixativeResponse> => {
  data = fixativeSchemaResponseTransformer(data)
  return data
}

export const updateFixativeResponseTransformer = async (
  data: any
): Promise<UpdateFixativeResponse> => {
  data = fixativeSchemaResponseTransformer(data)
  return data
}

export const listGenesResponseTransformer = async (data: any): Promise<ListGenesResponse> => {
  data = data.map((item: any) => {
    return geneSchemaResponseTransformer(item)
  })
  return data
}

export const createGeneResponseTransformer = async (data: any): Promise<CreateGeneResponse> => {
  data = geneSchemaResponseTransformer(data)
  return data
}

export const deleteGeneResponseTransformer = async (data: any): Promise<DeleteGeneResponse> => {
  data = geneSchemaResponseTransformer(data)
  return data
}

export const updateGeneResponseTransformer = async (data: any): Promise<UpdateGeneResponse> => {
  data = geneSchemaResponseTransformer(data)
  return data
}

const habitatGroupSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listHabitatGroupsResponseTransformer = async (
  data: any
): Promise<ListHabitatGroupsResponse> => {
  data = data.map((item: any) => {
    return habitatGroupSchemaResponseTransformer(item)
  })
  return data
}

export const createHabitatGroupResponseTransformer = async (
  data: any
): Promise<CreateHabitatGroupResponse> => {
  data = habitatGroupSchemaResponseTransformer(data)
  return data
}

export const deleteHabitatGroupResponseTransformer = async (
  data: any
): Promise<DeleteHabitatGroupResponse> => {
  data = habitatGroupSchemaResponseTransformer(data)
  return data
}

export const updateHabitatGroupResponseTransformer = async (
  data: any
): Promise<UpdateHabitatGroupResponse> => {
  data = habitatGroupSchemaResponseTransformer(data)
  return data
}

const siteWithDistanceSchemaResponseTransformer = (data: any) => {
  if (data.last_visited) {
    data.last_visited = optionalDateWithPrecisionSchemaResponseTransformer(data.last_visited)
  }
  return data
}

export const sitesProximityResponseTransformer = async (
  data: any
): Promise<SitesProximityResponse> => {
  data = data.map((item: any) => {
    return siteWithDistanceSchemaResponseTransformer(item)
  })
  return data
}

const siteWithScoreSchemaResponseTransformer = (data: any) => {
  if (data.last_visited) {
    data.last_visited = optionalDateWithPrecisionSchemaResponseTransformer(data.last_visited)
  }
  return data
}

export const searchSitesResponseTransformer = async (data: any): Promise<SearchSitesResponse> => {
  data = data.map((item: any) => {
    return siteWithScoreSchemaResponseTransformer(item)
  })
  return data
}

export const occurrencesBySiteResponseTransformer = async (
  data: any
): Promise<OccurrencesBySiteResponse> => {
  data = data.map((item: any) => {
    return siteWithOccurrencesSchemaResponseTransformer(item)
  })
  return data
}

const organisationSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listOrganisationsResponseTransformer = async (
  data: any
): Promise<ListOrganisationsResponse> => {
  data = data.map((item: any) => {
    return organisationSchemaResponseTransformer(item)
  })
  return data
}

export const createOrganisationResponseTransformer = async (
  data: any
): Promise<CreateOrganisationResponse> => {
  data = organisationSchemaResponseTransformer(data)
  return data
}

export const deleteOrganisationResponseTransformer = async (
  data: any
): Promise<DeleteOrganisationResponse> => {
  data = organisationSchemaResponseTransformer(data)
  return data
}

export const updateOrganisationResponseTransformer = async (
  data: any
): Promise<UpdateOrganisationResponse> => {
  data = organisationSchemaResponseTransformer(data)
  return data
}

const personSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listPersonsResponseTransformer = async (data: any): Promise<ListPersonsResponse> => {
  data = data.map((item: any) => {
    return personSchemaResponseTransformer(item)
  })
  return data
}

export const createPersonResponseTransformer = async (data: any): Promise<CreatePersonResponse> => {
  data = personSchemaResponseTransformer(data)
  return data
}

export const deletePersonResponseTransformer = async (data: any): Promise<DeletePersonResponse> => {
  data = personSchemaResponseTransformer(data)
  return data
}

export const updatePersonResponseTransformer = async (data: any): Promise<UpdatePersonResponse> => {
  data = personSchemaResponseTransformer(data)
  return data
}

const programSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listProgramsResponseTransformer = async (data: any): Promise<ListProgramsResponse> => {
  data = data.map((item: any) => {
    return programSchemaResponseTransformer(item)
  })
  return data
}

export const createProgramResponseTransformer = async (
  data: any
): Promise<CreateProgramResponse> => {
  data = programSchemaResponseTransformer(data)
  return data
}

export const deleteProgramResponseTransformer = async (
  data: any
): Promise<DeleteProgramResponse> => {
  data = programSchemaResponseTransformer(data)
  return data
}

export const updateProgramResponseTransformer = async (
  data: any
): Promise<UpdateProgramResponse> => {
  data = programSchemaResponseTransformer(data)
  return data
}

export const listArticlesResponseTransformer = async (data: any): Promise<ListArticlesResponse> => {
  data = data.map((item: any) => {
    return articleSchemaResponseTransformer(item)
  })
  return data
}

export const createArticleResponseTransformer = async (
  data: any
): Promise<CreateArticleResponse> => {
  data = articleSchemaResponseTransformer(data)
  return data
}

export const deleteArticleResponseTransformer = async (
  data: any
): Promise<DeleteArticleResponse> => {
  data = articleSchemaResponseTransformer(data)
  return data
}

export const updateArticleResponseTransformer = async (
  data: any
): Promise<UpdateArticleResponse> => {
  data = articleSchemaResponseTransformer(data)
  return data
}

export const listSamplingMethodsResponseTransformer = async (
  data: any
): Promise<ListSamplingMethodsResponse> => {
  data = data.map((item: any) => {
    return samplingMethodSchemaResponseTransformer(item)
  })
  return data
}

export const createSamplingMethodResponseTransformer = async (
  data: any
): Promise<CreateSamplingMethodResponse> => {
  data = samplingMethodSchemaResponseTransformer(data)
  return data
}

export const deleteSamplingMethodResponseTransformer = async (
  data: any
): Promise<DeleteSamplingMethodResponse> => {
  data = samplingMethodSchemaResponseTransformer(data)
  return data
}

export const updateSamplingMethodResponseTransformer = async (
  data: any
): Promise<UpdateSamplingMethodResponse> => {
  data = samplingMethodSchemaResponseTransformer(data)
  return data
}

export const createSamplingResponseTransformer = async (
  data: any
): Promise<CreateSamplingResponse> => {
  data = samplingSchemaResponseTransformer(data)
  return data
}

export const deleteSamplingResponseTransformer = async (
  data: any
): Promise<DeleteSamplingResponse> => {
  data = samplingSchemaResponseTransformer(data)
  return data
}

export const updateSamplingResponseTransformer = async (
  data: any
): Promise<UpdateSamplingResponse> => {
  data = samplingSchemaResponseTransformer(data)
  return data
}

export const samplingAddExternalOccurrenceResponseTransformer = async (
  data: any
): Promise<SamplingAddExternalOccurrenceResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

export const listSequencesResponseTransformer = async (
  data: any
): Promise<ListSequencesResponse> => {
  data = data.map((item: any) => {
    return sequenceSchemaResponseTransformer(item)
  })
  return data
}

export const deleteSequenceResponseTransformer = async (
  data: any
): Promise<DeleteSequenceResponse> => {
  data = sequenceSchemaResponseTransformer(data)
  return data
}

const sequenceWithDetailsSchemaResponseTransformer = (data: any) => {
  if (data.code_history) {
    data.code_history = data.code_history.map((item: any) => {
      return codeHistorySchemaResponseTransformer(item)
    })
  }
  data.event = eventInnerSchemaResponseTransformer(data.event)
  if (data.external) {
    data.external = optionalExtSeqSpecificsBioMaterialSchemaResponseTransformer(data.external)
  }
  data.gene = geneSchemaResponseTransformer(data.gene)
  data.identification = identificationSchemaResponseTransformer(data.identification)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.published_in) {
    data.published_in = data.published_in.map((item: any) => {
      return occurrenceReferenceSchemaResponseTransformer(item)
    })
  }
  data.sampling = samplingSchemaResponseTransformer(data.sampling)
  return data
}

export const getSequenceResponseTransformer = async (data: any): Promise<GetSequenceResponse> => {
  data = sequenceWithDetailsSchemaResponseTransformer(data)
  return data
}

const dataFeedSpecSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listDataFeedsResponseTransformer = async (
  data: any
): Promise<ListDataFeedsResponse> => {
  data = data.map((item: any) => {
    return dataFeedSpecSchemaResponseTransformer(item)
  })
  return data
}

export const createDataFeedResponseTransformer = async (
  data: any
): Promise<CreateDataFeedResponse> => {
  data = dataFeedSpecSchemaResponseTransformer(data)
  return data
}

const mapToolPresetSchemaResponseTransformer = (data: any) => {
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listMapPresetsResponseTransformer = async (
  data: any
): Promise<ListMapPresetsResponse> => {
  data = data.map((item: any) => {
    return mapToolPresetSchemaResponseTransformer(item)
  })
  return data
}

export const createUpdateMapPresetResponseTransformer = async (
  data: any
): Promise<CreateUpdateMapPresetResponse> => {
  data = mapToolPresetSchemaResponseTransformer(data)
  return data
}

export const deleteMapPresetResponseTransformer = async (
  data: any
): Promise<DeleteMapPresetResponse> => {
  data = mapToolPresetSchemaResponseTransformer(data)
  return data
}

const siteSchemaResponseTransformer = (data: any) => {
  if (data.events) {
    data.events = data.events.map((item: any) => {
      return eventSchemaResponseTransformer(item)
    })
  }
  if (data.last_visited) {
    data.last_visited = optionalDateWithPrecisionSchemaResponseTransformer(data.last_visited)
  }
  data.meta = metaSchemaResponseTransformer(data.meta)
  return data
}

export const listSitesResponseTransformer = async (data: any): Promise<ListSitesResponse> => {
  data = data.map((item: any) => {
    return siteSchemaResponseTransformer(item)
  })
  return data
}

export const createSiteResponseTransformer = async (data: any): Promise<CreateSiteResponse> => {
  data = siteSchemaResponseTransformer(data)
  return data
}

export const getSiteResponseTransformer = async (data: any): Promise<GetSiteResponse> => {
  data = siteSchemaResponseTransformer(data)
  return data
}

export const updateSiteResponseTransformer = async (data: any): Promise<UpdateSiteResponse> => {
  data = siteSchemaResponseTransformer(data)
  return data
}

export const listSiteEventsResponseTransformer = async (
  data: any
): Promise<ListSiteEventsResponse> => {
  data = data.map((item: any) => {
    return eventSchemaResponseTransformer(item)
  })
  return data
}

export const createEventResponseTransformer = async (data: any): Promise<CreateEventResponse> => {
  data = eventSchemaResponseTransformer(data)
  return data
}

export const siteAddExternalOccurrenceResponseTransformer = async (
  data: any
): Promise<SiteAddExternalOccurrenceResponse> => {
  data = bioMaterialWithDetailsSchemaResponseTransformer(data)
  return data
}

const taxonomySchemaResponseTransformer = (data: any) => {
  if (data.children) {
    data.children = data.children.map((item: any) => {
      return taxonomySchemaResponseTransformer(item)
    })
  }
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.parent) {
    data.parent = optionalTaxonSchemaResponseTransformer(data.parent)
  }
  return data
}

export const getTaxonomyResponseTransformer = async (data: any): Promise<GetTaxonomyResponse> => {
  data = taxonomySchemaResponseTransformer(data)
  return data
}

export const listTaxaResponseTransformer = async (data: any): Promise<ListTaxaResponse> => {
  data = data.map((item: any) => {
    return taxonWithParentRefSchemaResponseTransformer(item)
  })
  return data
}

const taxonWithRelativesSchemaResponseTransformer = (data: any) => {
  if (data.children) {
    data.children = data.children.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.parent) {
    data.parent = optionalTaxonSchemaResponseTransformer(data.parent)
  }
  return data
}

export const createTaxonResponseTransformer = async (data: any): Promise<CreateTaxonResponse> => {
  data = taxonWithRelativesSchemaResponseTransformer(data)
  return data
}

export const deleteTaxonResponseTransformer = async (data: any): Promise<DeleteTaxonResponse> => {
  data = taxonWithRelativesSchemaResponseTransformer(data)
  return data
}

const lineageSchemaResponseTransformer = (data: any) => {
  if (data.class) {
    data.class = optionalTaxonSchemaResponseTransformer(data.class)
  }
  if (data.family) {
    data.family = optionalTaxonSchemaResponseTransformer(data.family)
  }
  if (data.genus) {
    data.genus = optionalTaxonSchemaResponseTransformer(data.genus)
  }
  if (data.kingdom) {
    data.kingdom = optionalTaxonSchemaResponseTransformer(data.kingdom)
  }
  if (data.order) {
    data.order = optionalTaxonSchemaResponseTransformer(data.order)
  }
  if (data.phylum) {
    data.phylum = optionalTaxonSchemaResponseTransformer(data.phylum)
  }
  if (data.species) {
    data.species = optionalTaxonSchemaResponseTransformer(data.species)
  }
  if (data.subspecies) {
    data.subspecies = optionalTaxonSchemaResponseTransformer(data.subspecies)
  }
  return data
}

const taxonWithLineageSchemaResponseTransformer = (data: any) => {
  if (data.children) {
    data.children = data.children.map((item: any) => {
      return taxonSchemaResponseTransformer(item)
    })
  }
  data.lineage = lineageSchemaResponseTransformer(data.lineage)
  data.meta = metaSchemaResponseTransformer(data.meta)
  if (data.parent) {
    data.parent = optionalTaxonSchemaResponseTransformer(data.parent)
  }
  return data
}

export const getTaxonResponseTransformer = async (data: any): Promise<GetTaxonResponse> => {
  data = taxonWithLineageSchemaResponseTransformer(data)
  return data
}

export const updateTaxonResponseTransformer = async (data: any): Promise<UpdateTaxonResponse> => {
  data = taxonSchemaResponseTransformer(data)
  return data
}
