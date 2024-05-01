<template>
  <div :class="`d-flex ${xs ? 'flex-column' : 'flex-row'}`">
    <v-tabs v-model="tab" color="primary" slider-color="primary" :direction="direction">
      <v-tab
        value="instance"
        text="Instance"
        prepend-icon="mdi-application-settings-outline"
        tile
        :rounded="false"
        href="#instance"
      />
      <v-tab
        value="security"
        text="Security"
        :tile="true"
        prepend-icon="mdi-security"
        :rounded="false"
        href="#security"
      />
      <v-tab
        value="email"
        text="E-mail"
        prepend-icon="mdi-email"
        tile
        :rounded="false"
        href="#email"
      />
    </v-tabs>
    <v-container>
      <v-row>
        <v-col>
          <Suspense>
            <template #default>
              <v-window v-model="tab" class="pa-3" :direction="direction">
                <v-window-item value="instance">
                  <InstanceSettings />
                </v-window-item>
                <v-window-item value="security">
                  <SecuritySettings />
                </v-window-item>
                <v-window-item value="email">
                  <EmailSettings />
                </v-window-item>
              </v-window>
            </template>
            <template #fallback>
              <v-skeleton-loader type="article, article, article"> </v-skeleton-loader>
            </template>
          </Suspense>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import InstanceSettings from '@/components/settings/InstanceSettings.vue'
import SecuritySettings from '@/components/settings/SecuritySettings.vue'
import EmailSettings from '@/components/settings/EmailSettings.vue'
import { useDisplay } from 'vuetify'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { watch } from 'vue'

const { xs } = useDisplay()

type SettingsTab = 'instance' | 'security' | 'email'

const direction = computed(() => {
  return xs.value ? 'horizontal' : 'vertical'
})

const route = useRoute()
function hashToTab(hash: string) {
  return hash.startsWith('#') ? (hash.slice(1) as SettingsTab) : 'instance'
}
const tab = ref(hashToTab(route.hash))
watch(
  () => route.hash,
  (hash: string) => {
    tab.value = hashToTab(hash)
  }
)
</script>

<style scoped></style>
