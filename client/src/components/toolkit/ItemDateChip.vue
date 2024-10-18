<template>
  <v-chip label color="#777" variant="text" v-bind="$attrs">
    <template v-slot:prepend>
      <v-icon color="#777" class="mr-3" size="small" :icon="icon" />
    </template>
    {{ displayDate(date) }}
  </v-chip>
</template>

<script setup lang="ts">
import moment, { type Moment, type MomentInput } from 'moment'
import { ComputedRef, computed } from 'vue'

type UpdateIcon = 'mdi-update'
type CreatedIcon = 'mdi-content-save'
type Icon = UpdateIcon | CreatedIcon | string

type Props = {
  date: MomentInput
  icon: Icon
  iconColor?: string
}
const props = withDefaults(defineProps<Props>(), {
  iconColor: '#777'
})

const date: ComputedRef<Moment> = computed(() => {
  return moment(props.date)
})
const icon: ComputedRef<string> = computed(() => {
  switch (props.icon) {
    case 'updated':
      return 'mdi-update'
    case 'created':
      return 'mdi-content-save'
    default:
      return props.icon
  }
})

function displayDate(date: Moment) {
  return date.calendar()
}
</script>

<style scoped></style>
