import { ErrorModel, SiteWithOccurrences } from "@/api";
import { MappingFilters } from "@/components/occurrence/OccurrenceDataFeedFilters.vue";
import { UseQueryReturnType } from "@tanstack/vue-query";
import { useLocalStorage } from "@vueuse/core";
import type { UUID } from "crypto";
import { v4 as uuidv4 } from "uuid";
import { computed, MaybeRef, ref, shallowReactive, unref } from "vue";

export type DataFeed = {
  id: UUID
  name?: string
  filters: MappingFilters
}

export type DataFeedContext = {
  datasets?: string[]
  taxa?: string[]
  whole_clade?: boolean
}

export type RegisteredDataFeed = {
  id: UUID,
} & DataFeed

const context = ref<DataFeedContext>({})
const contextEnabled = ref(false)

const dataFeeds = useLocalStorage<[DataFeed, ...Array<DataFeed>]>('maptool-data-feeds', [newDataFeed("Feed #1")]);

const data = shallowReactive(new Map<UUID, UseQueryReturnType<SiteWithOccurrences[], ErrorModel>>())

const anyLoading = computed(() => data.size > 0 && [...data.values()].some(remote => remote.isPending.value))

const allPending = computed(() => data.size > 0 && [...data.values()].every(remote => remote.isPending.value))

function register(dataFeed: MaybeRef<DataFeed>, remote: UseQueryReturnType<SiteWithOccurrences[], ErrorModel>) {
  data.set(unref(dataFeed).id, remote)
}

function newDataFeed(name?: string): DataFeed {
  return {
    id: uuidv4() as UUID,
    name: name || `Feed #${dataFeeds.value.length + 1}`,
    filters: contextEnabled.value ? { ...context.value } : {}
  }
}

function addDataFeed(): DataFeed {
  const feed = newDataFeed()
  dataFeeds.value.push(feed);
  return feed
}

function deleteFeed(id: UUID) {
  const index = dataFeeds.value.findIndex(df => df.id === id);
  if (index !== -1) {
    dataFeeds.value.splice(index, 1);
    data.delete(id);
  }
}

function duplicateFeed(feed: DataFeed): DataFeed {
  const newFeed: DataFeed = { name: `${feed.name} (copy)`, id: uuidv4() as UUID, filters: { ...feed.filters } };
  dataFeeds.value.push(newFeed);
  return newFeed;
}

function resetAll() {
  data.clear();
  dataFeeds.value = [newDataFeed("Feed #1")];
}

function applyContext() {
  dataFeeds.value.forEach(feed => {
    feed.filters.datasets = context.value.datasets?.concat(
      feed.filters.datasets?.filter(ds => !context.value.datasets?.includes(ds)) || []
    );
    feed.filters.taxa = context.value.taxa?.concat(
      feed.filters.taxa?.filter(t => !context.value.taxa?.includes(t)) || []
    );
    feed.filters.whole_clade = context.value.whole_clade
  });
}

export function useDataFeeds() {
  return {
    feeds: dataFeeds,
    data,
    anyLoading,
    allPending,
    context,
    contextEnabled,
    applyContext,
    register,
    deleteFeed,
    newDataFeed,
    addDataFeed,
    resetAll,
    duplicateFeed
  };
}