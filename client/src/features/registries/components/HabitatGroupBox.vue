<template>
  <div title="" :link="false" class="group-box d-flex">
    <div v-if="depth" class="box-indent-container">
      <div :class="['indent-lines connector', { terminal: terminal }]"></div>
      <div v-if="!terminal" class="indent-lines"></div>
    </div>
    <v-card variant="outlined" class="bg-main flex-grow-1 mb-3 me-3">
      <div class="d-flex justify-space-between pa-3 cursor-pointer" v-ripple @click="toggleOpen()">
        <div>
          <span class="text-muted text-label-medium">
            {{ group.label }}
          </span>
          <v-badge
            inline
            :content="group.children.length"
            size="small"
            variant="outlined"
          ></v-badge>
        </div>
        <v-btn
          @click.stop="toggleOpen()"
          icon="mdi-plus-box-outline"
          size="small"
          variant="plain"
          density="compact"
        ></v-btn>
      </div>
      <v-expand-transition>
        <v-list v-if="isOpen">
          <template v-for="(child, index) in group.children" :key="child.id">
            <v-divider v-if="index > 0"></v-divider>
            <v-list-group v-if="child.children?.length" :title="child.label">
              <template #activator="{ props }">
                <v-list-item :title="child.label" v-bind="props">
                  <template #title>
                    {{ child.label }}
                    <v-badge
                      inline
                      :content="child.children.length"
                      size="small"
                      variant="outlined"
                    ></v-badge>
                  </template>
                  <template #prepend>
                    <v-icon icon="mdi-circle-medium"></v-icon>
                  </template>
                </v-list-item>
              </template>
              <HabitatGroupBox
                v-if="child.children?.length"
                v-for="(group, index) in child.children"
                :key="group.id"
                v-model:open="childrenState[group.id]"
                :group="group"
                :depth="depth + 1"
                :terminal="index === child.children.length - 1"
              />
            </v-list-group>
            <v-list-item v-else :title="child.label">
              <template #prepend>
                <v-icon icon="mdi-circle-medium"></v-icon>
              </template>
            </v-list-item>
          </template>
        </v-list>
      </v-expand-transition>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { HabitatGroup, HabitatRecord } from '@/api'

export type HabitatNode = HabitatRecord & {
  children?: HabitatGroupNode[]
  type: 'habitat'
}

export type HabitatGroupNode = HabitatGroup & {
  children: HabitatNode[]
  type: 'group'
}

const { group, depth = 0 } = defineProps<{
  group: HabitatGroupNode
  depth?: number
  terminal?: boolean
}>()

const isOpen = defineModel<boolean>('open', {
  default: false
})

function toggleOpen() {
  isOpen.value = !isOpen.value
}

const childrenState = ref<Record<string, boolean>>({})
</script>

<style lang="scss">
@use 'vuetify';

.v-list-group div.box-indent-container {
  display: flex;
  flex-direction: column;
  width: calc(var(--v-list-indent) + var(--list-indent-size) + 0px);
  .indent-lines {
    margin-left: 50%;
    width: 50%;
    flex-grow: 1;
    border-left: 1px solid rgba(var(--v-theme-on-surface), 0.5);
    &.connector {
      flex-grow: 0;
      height: 16px;
      border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.5);
    }
    &.terminal {
      border-bottom-left-radius: 4px;
    }
  }
}
</style>
