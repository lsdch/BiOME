import { DatasetInner, OrganisationInner, Person, PersonInner, PersonInput, PersonUser, Program, ProgramInput, ProgramUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type ProgramModel = Omit<ProgramInput, "funding_agencies" | "managers" | "datasets"> & {
  funding_agencies: OrganisationInner[];
  managers: PersonInner[];
  datasets: DatasetInner[];
}

export function initialModel(): Reactive<ProgramModel> {
  return reactive({
    code: '',
    label: '',
    funding_agencies: [],
    managers: [],
    datasets: []
  })
}

export function fromProgram({
  id, meta, $schema, ...rest
}: Program): ProgramModel {
  return rest
}

export function toRequestBody({ managers, funding_agencies, datasets, ...model }: ProgramModel): ProgramInput {
  return {
    ...model,
    managers: managers.map(({ alias }) => alias),
    funding_agencies: funding_agencies.map(({ code }) => code),
    datasets: datasets.map(({ slug }) => slug)
  } satisfies ProgramUpdate
}
