<template>
  <v-form>
    <v-row>
      <v-col cols="12" sm="3" class="px-3 d-flex align-center justify-center">
        <v-avatar
          :class="iconHover ? 'avatar-hover' : ''"
          :size="120"
          variant="elevated"
          color="primary"
        >
          <img :src="img" :width="120" alt="alt" />
          <v-overlay
            @click="openIconDialog"
            v-model="iconHover"
            open-on-hover
            :close-delay="200"
            class="align-center justify-center cursor-pointer font-weight-black text-white"
            contained
            activator="parent"
            scrim="primary"
          >
            <!-- <v-icon>mdi-pencil</v-icon><br /> -->
            Change <br />icon
          </v-overlay>
        </v-avatar>
      </v-col>
      <v-col cols="12" sm="9" variant="text" class="pt-8 d-flex flex-column justify-center">
        <v-text-field
          label="Instance name"
          class="pb-4"
          hint="The name that is displayed in the navbar and front page"
          persistent-hint
        />
        <v-text-field
          label="Instance description"
          class="mb-5"
          hint="A short description of the database purpose to be displayed on the front page."
          persistent-hint
          clearable
        />
      </v-col>
    </v-row>
    <v-divider></v-divider>
    <v-switch
      label="Instance is public"
      class="mb-5"
      color="primary"
      hint="A private instance requires user authentication to get access to any data. A public instance allows read-only access to anonymous users on a subset of pages."
      persistent-hint
    />
    <v-divider></v-divider>
    <v-switch
      label="Allow contributor registration"
      color="primary"
      hint="If enabled, visitors may apply for an account with Contributor privileges."
      persistent-hint
    />
    <InstanceIconDialog v-model="dialog" />
  </v-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import InstanceIconDialog from './InstanceIconDialog.vue'

const img =
  'https://images.unsplash.com/photo-1600984575359-310ae7b6bdf2?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80'

const iconHover = ref(false)
const dialog = ref({
  open: false,
  iconSrc: img
})
function openIconDialog() {
  dialog.value.open = true
}
</script>

<style scoped>
.avatar-hover {
  border: 3px solid rgb(0, 112, 177);
}
</style>
