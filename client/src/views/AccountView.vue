<template>
  <div class="bg-surface fill-height">
    <TableToolbar title="Your account" icon="mdi-account"> </TableToolbar>
    <v-container>
      <v-row>
        <v-col>
          <v-card :title="model.identity.full_name">
            <template #append>
              <v-chip :text="model.role">
                <template #prepend>
                  <UserRoleIcon :role="model.role" class="mr-3" />
                </template>
              </v-chip>
            </template>
            <v-divider />
            <v-card-text>
              <v-text-field label="Login" v-model="model.login" prepend-inner-icon="mdi-account" />
              <v-text-field label="E-mail" v-model="model.email" prepend-inner-icon="mdi-at" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card> </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-expansion-panels>
            <v-expansion-panel>
              <v-expansion-panel-title>
                <v-icon>mdi-lock</v-icon>
                <h3>Password</h3>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <PasswordUpdate />
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import PasswordUpdate from '@/components/account/PasswordUpdate.vue'
import UserRoleIcon from '@/components/icons/UserRoleIcon.vue'
import { roleIcon } from '@/components/people/userRole'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useUserStore } from '@/stores/user'
import { ref } from 'vue'

const { user: currentUser } = useUserStore()

// Route is inaccessible to non-authenticated users
const model = ref(currentUser!)
</script>

<style scoped></style>
