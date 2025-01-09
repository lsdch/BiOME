<template>
  <FormDialog v-model="dialogOpen" title="New invitation">
    <v-row>
      <v-col>
        <v-text-field v-model="model.email" label="E-mail" v-bind="field('email')"></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <PersonPicker label="Identity" :searchValue="pending?.full_name" restrict="unregistered" />
      </v-col>
      <v-col>
        <UserRolePicker v-model="model.role" />
      </v-col>
    </v-row>
  </FormDialog>
</template>

<script setup lang="ts">
import { $InvitationInput, InvitationInput, PendingUserRequest, UserRole } from '@/api'
import { ref, watch } from 'vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { useSchema } from '../toolkit/forms/schema'
import UserRolePicker from './UserRolePicker.vue'
import PersonPicker from '../people/PersonPicker.vue'

export type InitialContent = { pending?: PendingUserRequest; role?: UserRole }

const props = defineProps<InitialContent>()

const dialogOpen = defineModel<boolean>({ required: true })

watch(
  dialogOpen,
  (isOpen) => {
    if (!isOpen) return
    model.value.email = props.pending?.email ?? ''
    model.value.role = props.role ?? 'Visitor'
  },
  { immediate: true }
)

const model = ref<InvitationInput>({
  email: props.pending?.email ?? '',
  role: 'Contributor'
})

const { field } = useSchema($InvitationInput)
</script>

<style scoped></style>
