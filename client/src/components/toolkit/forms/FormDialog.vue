<template>
  <CardDialog v-model="dialog" v-bind="props">
    <template #append>
      <v-btn
        color="primary"
        type="submit"
        @click="emit('submit')"
        :loading="loading"
        :text="btnText"
      />
    </template>
  </CardDialog>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import CardDialog, { CardDialogProps } from './CardDialog.vue'

const dialog = defineModel<boolean>({ default: false })

const emit = defineEmits<{ submit: [] }>()

const props = withDefaults(defineProps<CardDialogProps & { btnText?: string }>(), {
  btnText: 'Submit',
  closeText: 'Cancel'
})
</script>

<style scoped></style>
