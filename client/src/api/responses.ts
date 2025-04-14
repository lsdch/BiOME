import { useFeedback } from "@/stores/feedback";

export type PaginatedList<Item> = {
  items: Item[];
  total_count: number;
}

export type ResponseBody<Data = unknown, Error = unknown> = {
  data: Data;
  error: undefined;
} | {
  data: undefined;
  error: Error;
}


export function mergeResponses<D = unknown, E = unknown>(responses: ResponseBody<D, E>[]): ResponseBody<D[], E[]> {
  return responses.reduce<ResponseBody<D[], E[]>>((acc, { data, error }) => {
    if (error !== undefined) {
      acc.error?.push(error) ?? (acc.error = [error])
    } else if (data !== undefined) {
      acc.data?.push(data)
    }
    return acc
  }, {
    data: [],
    error: undefined,
  })
}

export function useErrorHandler<D, E>(handler: (err: E) => any) {
  return ({ data, error }: ResponseBody<D, E>) => {
    return new Promise<D>((resolve, reject) => {
      if (error !== undefined) {
        handler(error)
        return reject()
      }
      return resolve(data!)
    })
  }
}

export function errorFeedback<Data, Error>(message: string) {
  return ({ data, error }: ResponseBody<Data, Error>) => {
    const { feedback } = useFeedback()
    return new Promise<Data>((resolve, reject) => {
      if (error != undefined) {
        feedback({ message, type: "error" })
        return reject()
      }
      return resolve(data!)
    })
  }
}