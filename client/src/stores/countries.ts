import { getCountriesSummaryOptions, listCountriesOptions } from "@/api/gen/@tanstack/vue-query.gen";
import { useQuery } from "@tanstack/vue-query";
import { defineStore } from "pinia";

export const useCountries = defineStore("countries", () => {

  const { data: countries, error, isPending, refetch } = useQuery({
    ...getCountriesSummaryOptions(),
    gcTime: Infinity,
    initialData: []
  })

  function findCountry(code: string) {
    return countries.value.find(({ code: c }) => c === code)
  }

  return { countries, isPending, error, refetch, findCountry }
})