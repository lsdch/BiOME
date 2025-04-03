<template>
  <EventFormDialog
    v-bind="{ ...props, ...$attrs }"
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :loading="activeMutation.isPending.value"
    :errors="schemaBindings.errors"
    @submit="submit"
  />
</template>

<script setup lang="ts">
import { $EventInput, $EventUpdate, Event, SiteItem } from '@/api'
import { SiteInput } from '@/api/adapters'
import { createEventMutation, updateEventMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { EventModel } from '@/models'
import { FormDialogProps } from '../toolkit/forms/FormDialog.vue'
import EventFormDialog from '../forms/EventFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Event>()
const props = defineProps<
  {
    site: SiteItem | SiteInput
  } & Omit<FormDialogProps, 'loading'>
>()

const create = defineFormCreate(createEventMutation(), {
  initial: EventModel.initialModel,
  schema: $EventInput,
  requestData({ performed_by, performed_on, programs }) {
    return {
      body: {
        performed_by: performed_by.map(({ alias }) => alias),
        performed_on,
        programs: programs.map(({ code }) => code)
      },
      path: { code: props.site.code }
    }
  }
})

const update = defineFormUpdate(updateEventMutation(), {
  schema: $EventUpdate,
  itemToModel: EventModel.fromEvent,
  requestData(item, { performed_by, performed_on, programs }) {
    return {
      body: {
        performed_by: performed_by.map(({ alias }) => alias),
        performed_on,
        programs: programs.map(({ code }) => code)
      },
      path: { id: item.id }
    }
  }
})

const { mode, model, schemaBindings, activeMutation, submit } = useMutationForm(item, {
  create,
  update,
  onSuccess: () => {
    dialog.value = false
  }
})
</script>

<style scoped lang="scss"></style>
