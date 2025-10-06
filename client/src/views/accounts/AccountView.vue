<template>
  <div class="bg-main">
    <TableToolbar :title="model.identity.full_name" icon="mdi-account">
    </TableToolbar>
    <v-container>
      <v-row>
        <v-col>
          <v-card :title="model.role">
            <template #prepend>
              <UserRole.Icon :role="model.role" class="mr-3" />
            </template>
            <v-divider/>
            <template v-if="model.role === 'Admin'">
              <v-card-text>
              <p>
                As an <b class="text-red">administrator</b> you have
                <b>full access to the system</b>.
              </p>
              <p>
                You are responsible for :
                <ul>
                  <li>
                    managing user account invitations and validate pending user
                    account requests.
                  </li>
                  <li>
                    maintaining database referentials, such as taxonomy and
                    metadata registries (e.g. genes, sampling methods, habitats, etc).
                  </li>
                  <li>
                    maintaining the system configuration and settings
                  </li>
                </ul>
              </p>
            </v-card-text>
              <v-divider/>
            <v-card-text >

              <v-list>
                <v-list-item title="Occurrence datasets" subtitle="Full read/write access">
                </v-list-item>
                <v-list-item title="Metadata registries" subtitle="Update and maintain">
                </v-list-item>
                <v-list-item title="User accounts" subtitle="Manage invitations and requests">
                </v-list-item>
                <v-list-item title="System configuration" subtitle="Instance branding, configure connections to external services"/>
                <v-list-item title="Dashboard management" subtitle="Organize dashboard widgets and pinned datasets"/>
              </v-list>
            </v-card-text>
          </template>
            <v-card-text>
              <p>
                As a <b class="text-green">maintainer</b> you have
                <b>full read/write access to the platform data and metadata</b>.
              </p>
              <p>You are not in charge of managing user accounts, and global platform configuration.</p>
              <v-list>
                <v-list-item title="Occurrence datasets" subtitle="Full read/write access">
                </v-list-item>
                <v-list-item title="Metadata registries" subtitle="Update and maintain">
                </v-list-item>
              </v-list>
            </v-card-text>
            <v-card-text>
              <p>
                As a <b class="text-blue">contributor</b>
                you have <b>full read access to the system,
                and may submit and manage your own datasets</b>
              </p>
              <v-list>
                <v-list-item lines="two" title="Occurrence datasets">
                  <template #subtitle>
                    Full read access <br />
                    Submit and modify your own datasets
                  </template>
                </v-list-item>
                <v-list-item lines="two" title="Metadata registries">
                  <template #subtitle>
                    Read-only access <br />
                    Contact a maintainer or administrator for updates if existing registries do not
                    fit your data.
                  </template>
                </v-list-item>
              </v-list>
            </v-card-text>
            <v-card-text>
              <p>
                As a <b class="text-orange">visitor</b> you have <b>read-only access to the system</b>.
              </p>
              <p>
                If you wish to contribute data to the platform, you may apply for a <b>contributor</b> account.
              </p>
              <v-list>
                <v-list-item title="Occurrence datasets" subtitle="Read-only access"> </v-list-item>
                <v-list-item title="Metadata registries" subtitle="Read-only access"> </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card>
            <v-card-text>
              <v-text-field label="Login" v-model="model.login" prepend-inner-icon="mdi-account" />
              <v-text-field label="E-mail" v-model="model.email" prepend-inner-icon="mdi-at" />
            </v-card-text>
          </v-card>
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
import { UserRole } from '@/api'
import PasswordUpdate from '@/components/users/PasswordUpdate.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useUserStore } from '@/stores/user'
import { ref } from 'vue'

const { user: currentUser } = useUserStore()

// Route is inaccessible to non-authenticated users
const model = ref(currentUser!)
</script>

<style scoped></style>
