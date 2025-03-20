import { ErrorModel, GeoapifyResult } from "@/api";
import { getGeoapifyStatusOptions, reverseGeocodeOptions } from "@/api/gen/@tanstack/vue-query.gen";
import { Coordinates, MaybeCoordinates } from "@/components/maps";
import { QueryObserverOptions, useQuery, useQueryClient, UseQueryOptions } from "@tanstack/vue-query";
import { defineStore } from "pinia";
import { computed, MaybeRef, unref, watch } from "vue";

export const useGeoapify = defineStore("geoapify", () => {

  const client = useQueryClient()

  const { data: status, error } = useQuery({
    ...getGeoapifyStatusOptions(),
    gcTime: Infinity,
  })

  const isAvailable = computed(() => !!status.value?.available)

  function reverseGeocodeQuery(
    coords: MaybeRef<MaybeCoordinates>,
    options?: Omit<
      Partial<QueryObserverOptions<GeoapifyResult, ErrorModel, GeoapifyResult, any, any>>,
      'enabled' | 'queryKey' | 'queryFn'
    > & {
      enabled?: MaybeRef<boolean>
    }
  ) {

    const q = useQuery(computed(() => ({
      enabled: (
        isAvailable.value &&
        Coordinates.isValidCoordinates(unref(coords)) &&
        (unref(options?.enabled) ?? true)
      ),

      ...reverseGeocodeOptions({
        body: unref(coords) as Coordinates
      }),
    })))

    // Update API usage
    watch(q.isFetching, (fetching, wasFetching) => {
      if (!fetching && wasFetching) {
        client.invalidateQueries({
          queryKey: getGeoapifyStatusOptions().queryKey
        })
      }
    })
    return q
  }

  watch(error, (err) => {
    console.error("Failed to get Geoapify status", err)
  })


  return { status, reverseGeocodeQuery }
})