<template>
  <v-container>
    <v-row>
      <v-col
        v-if="items?.length"
        v-for="item in items"
        :key="item.ID"
        cols="12"
        sm="6"
        md="6"
        lg="4"
      >
        <AccountsPendingCard :item @invite="showInvitationForm" />
      </v-col>
      <v-col v-else> No pending account requests </v-col>
    </v-row>
    <InvitationFormDialog v-model="dialogOpen" v-bind="formContent" />
  </v-container>
</template>

<script setup lang="ts">
import { AccountService, PendingUserRequest } from '@/api'
import { handleErrors } from '@/api/responses'
import InvitationFormDialog, { InitialContent } from '@/components/account/InvitationFormDialog.vue'
import { ref } from 'vue'
import AccountsPendingCard from './AccountsPendingCard.vue'
import { useToggle } from '@vueuse/core'

const [dialogOpen, toggleInvitationDialog] = useToggle(false)

const formContent = ref<InitialContent>({
  role: 'Contributor'
})

function showInvitationForm(item: PendingUserRequest) {
  formContent.value.pending = item
  toggleInvitationDialog(true)
}

const items = ref(
  await AccountService.listPendingUserRequests().then(
    handleErrors((err) => {
      console.error('Failed to retrieve pending account requests:', err)
    })
  )
)
</script>

<style scoped></style>
