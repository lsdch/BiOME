import { Article, ArticleInput, ArticleUpdate } from "@/api";
import { Optional } from "ts-toolbelt/out/Object/Optional";
import { reactive, Reactive } from "vue";

export type ArticleFormModel = Optional<ArticleInput, 'year'>

export function initialModel(): Reactive<ArticleFormModel> {
  return reactive({
    authors: []
  })
}

export function fromArticle({ id, meta, $schema, ...rest }: Article): ArticleFormModel {
  return rest
}

export function toRequestBody({ ...model }: ArticleFormModel): ArticleInput {
  return {
    ...model,
    year: model.year!
  } satisfies ArticleUpdate
}