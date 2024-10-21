<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Account request"
    :toolbar="{
      title: 'Account requests',
      icon: 'mdi-account-plus'
    }"
    :fetch-items="AccountService.listPendingUserRequests"
    :delete="
      (item: PendingUserRequest) =>
        AccountService.deletePendingUserRequest({ path: { email: item.email } })
    "
    show-actions="delete"
  >
    <template #[`item.email_verified`]="{ value }">
      <v-icon
        :icon="value ? 'mdi-check-circle' : 'mdi-clock'"
        :color="value ? 'success' : 'warning'"
      />
    </template>
    <template #[`expanded-row-inject`]="{ item }">
      <div class="mx-5 my-3 d-flex align-start">
        <div class="font-italic">
          {{ item.motive ?? 'No motive provided.' }}
        </div>
        <v-btn
          class="ml-auto"
          color="primary"
          text="Send invitation"
          prepend-icon="mdi-account-plus"
          variant="plain"
          @click="showInvitationForm(item)"
        />
      </div>
    </template>
  </CRUDTable>
  <InvitationFormDialog v-model="dialogOpen" v-bind="formContent" />
</template>

<script setup lang="ts">
import { AccountService, PendingUserRequest, PeopleService } from '@/api'
import { handleErrors } from '@/api/responses'
import InvitationFormDialog, { InitialContent } from '@/components/account/InvitationFormDialog.vue'
import { ref } from 'vue'
import AccountsPendingCard from './AccountsPendingCard.vue'
import { useToggle } from '@vueuse/core'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { title: 'Name', key: 'full_name' },
  { title: 'E-mail', key: 'email' },
  { title: 'E-mail verified', key: 'email_verified', align: 'center' },
  { title: 'Institution', key: 'institution' }
]

const [dialogOpen, toggleInvitationDialog] = useToggle(false)

const selection = ref<string[]>([])

const formContent = ref<InitialContent>({
  role: 'Contributor'
})

async function deleteSelectedItems() {
  if (!selection.value.length) return
  const responses = await Promise.all(
    selection.value.map((email) => AccountService.deletePendingUserRequest({ path: { email } }))
  )
}

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
