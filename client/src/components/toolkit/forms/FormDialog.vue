<template>
  <v-form>
    <template #="{ isValid }">
      <CardDialog v-model="dialog" v-bind="{ ...$props, ...$attrs }">
        <template #append>
          <v-btn
            color="primary"
            type="submit"
            @click="emit('submit')"
            :loading="loading"
            :text="btnText"
            :disabled="!isValid.value"
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
