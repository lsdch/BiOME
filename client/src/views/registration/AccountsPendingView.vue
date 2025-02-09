<template>
  <CRUDTable
    class="fill-height"
    expand-on-click
    :headers
    entity-name="Account request"
    :toolbar="{
      title: 'Account requests',
      icon: 'mdi-account-plus'
    }"
    :fetch-items="listPendingUserRequestsOptions"
    :delete="{
      mutation: deletePendingUserRequestMutation,
      params: (item: PendingUserRequest) => ({ path: { email: item.email } })
    }"
    appendActions="delete"
  >
    <!-- Email column -->
    <template #item.email="{ value }">
      <a :href="`mailto:${value}`" @click.stop> {{ value }}</a>
    </template>

    <!-- Email verified column -->
    <template #header.email_verified="slotProps">
      <IconTableHeader icon="mdi-email-check" class="justify-center" v-bind="slotProps" />
    </template>
    <template #[`item.email_verified`]="{ value }">
      <v-icon
        :icon="value ? 'mdi-check-circle' : 'mdi-clock'"
        :color="value ? 'success' : 'warning'"
      />
    </template>

    <!-- Organisation column -->
    <template #item.organisation="{ value }">
      <span :class="{ 'font-italic': !value }">
        {{ value ?? 'N/A' }}
      </span>
    </template>

    <template #[`expanded-row-inject`]="{ item }">
      <div class="mx-5 my-2">
        <div v-if="mobile" class="my-1 d-flex align-center">
          <v-icon icon="mdi-domain" size="small" />
          <span class="font-weight-bold mx-1">Organisation: </span>
          <span>
            {{ item.organisation ?? 'N/A' }}
          </span>
        </div>
        <div class="font-weight-bold">Motive:</div>
        <div class="my-1 font-italic">
          {{ item.motive ? `"${item.motive}"` : 'No motive provided.' }}
        </div>
      </div>
    </template>
    <template #expanded-row-footer="{ item }">
      <div class="d-flex align-center">
        <ItemDateChip icon="created" :date="item.created_on" />
        <v-spacer></v-spacer>
        <v-btn
          class="ml-auto"
          color="primary"
          text="Send invitation"
          prepend-icon="mdi-account-plus"
          variant="plain"
          density="compact"
          @click="showInvitationForm(item)"
        />
      </div>
    </template>

    <!-- Disable export btn -->
    <template #footer.prepend />
  </CRUDTable>
  <InvitationFormDialog v-model="dialogOpen" v-bind="formContent" />
</template>

<script setup lang="ts">
import { AccountService, PendingUserRequest } from '@/api'
import {
  deletePendingUserRequestMutation,
  listPendingUserRequestsOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import InvitationFormDialog, { InitialContent } from '@/components/account/InvitationFormDialog.vue'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import IconTableHeader from '@/components/toolkit/tables/IconTableHeader.vue'
import { useToggle } from '@vueuse/core'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()

const headers: CRUDTableHeader<PendingUserRequest>[] = [
  { title: 'Name', key: 'full_name' },
  { title: 'E-mail', key: 'email' },
  { title: 'E-mail verified', key: 'email_verified', align: 'center' },
  { title: 'Organisation', key: 'organisation', hide: mobile }
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
</script>

<style scoped></style>
