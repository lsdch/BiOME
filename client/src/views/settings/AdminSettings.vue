<template>
  <div :class="`d-flex ${xs ? 'flex-column' : 'flex-row'}`">
    <v-tabs v-model="tab" color="primary" slider-color="primary" :direction="direction">
      <v-tab
        value="Instance"
        text="Instance"
        prepend-icon="mdi-application-settings-outline"
        tile
        :rounded="false"
      />
      <v-tab
        value="Security"
        text="Security"
        :tile="true"
        prepend-icon="mdi-security"
        :rounded="false"
      />
      <v-tab value="Email" text="E-mail" prepend-icon="mdi-email" tile :rounded="false" />
    </v-tabs>
    <v-container>
      <v-row>
        <v-col>
          <v-window v-model="tab" class="pa-3" :direction="direction">
            <v-window-item value="Instance">
              <InstanceSettings />
            </v-window-item>
            <v-window-item value="Security">
              <SecuritySettings />
            </v-window-item>
            <v-window-item value="Email">
              <EmailSettings />
            </v-window-item>
          </v-window>
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

const { xs } = useDisplay()

type SettingsTab = 'Instance' | 'Security' | 'Email'

const direction = computed(() => {
  return xs.value ? 'horizontal' : 'vertical'
})

const tab = ref<SettingsTab>('Instance')
</script>

<style scoped></style>
