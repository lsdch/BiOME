import { ErrorModel, SiteWithOccurrences } from "@/api";
import { MappingFilters } from "@/components/occurrence/OccurrenceDataFeedFilters.vue";
import { UseQueryReturnType } from "@tanstack/vue-query";
import type { UUID } from "crypto";
import { v4 as uuidv4 } from "uuid";
import { computed, MaybeRef, reactive, ref, shallowReactive, unref } from "vue";

export type DataFeed = {
  id: UUID
  name?: string
  filters: MappingFilters
}

export type RegisteredDataFeed = {
  id: UUID,
} & DataFeed

const dataFeeds = ref<[DataFeed, ...Array<DataFeed>]>([newDataFeed("Feed #1")]);

const remotes = shallowReactive(new Map<UUID, UseQueryReturnType<SiteWithOccurrences[], ErrorModel>>())

const anyLoading = computed(() => remotes.size > 0 && [...remotes.values()].some(remote => remote.isPending.value))

const allPending = computed(() => remotes.size > 0 && [...remotes.values()].every(remote => remote.isPending.value))

function register(dataFeed: MaybeRef<DataFeed>, remote: UseQueryReturnType<SiteWithOccurrences[], ErrorModel>) {
  remotes.set(unref(dataFeed).id, remote)
}

function newDataFeed(name?: string): DataFeed {
  return {
    id: uuidv4() as UUID,
    name: name || `Feed #${dataFeeds.value.length + 1}`,
    filters: {}
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
    remotes.delete(id);
  }
}

function duplicateFeed(feed: DataFeed): DataFeed {
  const newFeed: DataFeed = { name: `${feed.name} (copy)`, id: uuidv4() as UUID, filters: { ...feed.filters } };
  dataFeeds.value.push(newFeed);
  return newFeed;
}

function resetAll() {
  remotes.clear();
  dataFeeds.value = [newDataFeed("Feed #1")];
}

export function useDataFeeds() {
  return { registry: dataFeeds, remotes, anyLoading, allPending, register, deleteFeed, newDataFeed, addDataFeed, resetAll, duplicateFeed };
}