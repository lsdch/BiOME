<template>
  <v-form @submit.prevent>
    <template #="{ isValid, isDisabled }">
      <CardDialog v-model="model" v-bind="props">
        <template #subtitle v-if="$slots['subtitle']">
          <slot name="subtitle" />
        </template>
        <template #append>
          <v-btn
            color="primary"
            variant="flat"
            :loading="loading"
            v-bind="
              $vuetify.display.smAndUp
                ? {
                    text: btnText
                  }
                : {
                    icon: 'mdi-floppy',
                    size: 'small'
                  }
            "
            :disabled="!isValid.value || isDisabled.value"
            @click="emit('submit')"
            rounded="sm"
          />
        </template>

        <!-- Default slot -->
        <slot />

        <!-- Expose activator slot -->
        <template #activator="slotData" v-if="slots.activator">
          <slot name="activator" v-bind="slotData" />
        </template>
      </CardDialog>
    </template>
  </v-form>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { VCard, VDialog } from 'vuetify/components'
import CardDialog, { CardDialogProps } from '../ui/CardDialog.vue'
export type FormDialogProps = CardDialogProps & { btnText?: string }

// dialog state exposed from CardDialog
const model = defineModel<boolean>()

const emit = defineEmits<{ submit: [] }>()

const props = withDefaults(defineProps<FormDialogProps>(), {
  btnText: 'Submit',
  closeText: 'Cancel'
})

const slots = defineSlots<VCard['$slots'] & Pick<VDialog['$slots'], 'activator'>>()
</script>

<style scoped></style>
