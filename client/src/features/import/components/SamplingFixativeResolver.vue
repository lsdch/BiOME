<template>
    <v-data-table :items="resolutions" :headers>
        <template #item.resolved_fixative_id="{ item }">
            <FixativePicker :prepend-icon="undefined" :model-value="item.resolved_fixative_id"
                @update:model-value="(resolved_fixative_id: UUID) => resolve(item, resolved_fixative_id)"
                item-value="id" hide-details density="compact" />
        </template>
    </v-data-table>
</template>

<script setup lang="ts">
import { SamplingMethodResolution } from '@/api';
import { getFixativesResolutionOptions, resolveFixativeMutation } from '@/api/gen/@tanstack/vue-query.gen';
import FixativePicker from '@/features/registries/components/FixativePicker.vue';
import { useMutation, useQuery } from '@tanstack/vue-query';

const { import_id } = defineProps<{
    import_id: UUID
}>()
const { data: resolutions, refetch } = useQuery(getFixativesResolutionOptions({ path: { id: import_id } }))
const { mutateAsync: resolveFixative, error } = useMutation(resolveFixativeMutation())

function resolve(item: SamplingMethodResolution, resolved_fixative_id: UUID) {
    resolveFixative({
        path: { id: item.import_id },
        body: { input_text: item.input_text, resolved_fixative_id: resolved_fixative_id, status: 'selected' }
    }, {
        onSuccess() {
            refetch()
        }
    })
}

const headers = [
    {
        key: 'input_text',
        title: 'Input Text'
    },
    {
        key: 'resolved_fixative_id',
        title: 'Resolved Fixative',
        sortable: false
    }
]
</script>

<style scoped lang="scss"></style>
