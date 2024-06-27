import { Country, LocationService } from "@/api";
import { handleErrors } from "@/api/responses";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useCountries = defineStore("countries", () => {
  const countries = ref<Country[]>([])

  async function fetch() {
    countries.value = await LocationService.listCountries().then(handleErrors(err => {
      console.error("Failed to fetch countries: ", err)
    }))
    return countries.value
  }

  function findCountry(code: string) {
    return countries.value.find(({ code: c }) => c === code)
  }

  return { countries, fetch, findCountry }
})