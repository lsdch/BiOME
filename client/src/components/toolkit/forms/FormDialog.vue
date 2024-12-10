<template>
  <v-form @submit.prevent>
    <template #="{ isValid, isDisabled }">
      <CardDialog v-model="dialog" v-bind="{ ...$props, ...$attrs }">
        <template #append>
          <v-btn
            color="primary"
            variant="flat"
            :loading="loading"
            :text="btnText"
            :disabled="!isValid.value || isDisabled.value"
            @click="emit('submit')"
          />
        </template>

        <!-- Default slot -->
        <slot />
      </CardDialog>
    </template>
  </v-form>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import CardDialog, { CardDialogProps } from './CardDialog.vue'

const dialog = defineModel<boolean>({ default: false })

const emit = defineEmits<{ submit: [] }>()

withDefaults(defineProps<CardDialogProps & { btnText?: string }>(), {
  btnText: 'Submit',
  closeText: 'Cancel'
})
</script>

<style scoped></style>
