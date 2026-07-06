<template>
  <div title="" :link="false" class="group-box d-flex">
    <div v-if="depth" class="box-indent-container">
      <div :class="['indent-lines connector', { terminal: terminal }]"></div>
      <div v-if="!terminal" class="indent-lines"></div>
    </div>
    <v-card
      variant="outlined"
      :class="['bg-main flex-grow-1 mb-3', { 'border-e-0 rounded-e-0': depth > 0 }]"
    >
      <div class="d-flex justify-space-between pa-3 cursor-pointer" v-ripple @click="toggleOpen()">
        <div>
          <span class="text-muted text-label-medium">
            {{ group.name }}
          </span>
          <v-badge
            inline
            :content="group.children.length"
            size="small"
            variant="outlined"
          ></v-badge>
        </div>
        <div class="d-flex align-center ga-2">
          <v-btn
            @click.stop="toggleOpen()"
            :icon="isOpen ? 'mdi-minus-box-outline' : 'mdi-plus-box-outline'"
            size="small"
            variant="plain"
            density="compact"
          ></v-btn>
          <v-menu>
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon="mdi-dots-vertical"
                color=""
                size="small"
                variant="plain"
                density="compact"
                rounded="100"
              ></v-btn>
            </template>
            <v-list>
              <v-list-item title="Edit" prepend-icon="mdi-pencil" @click.stop="emit('edit', group)">
              </v-list-item>
              <ConfirmDialog :title="`Delete group ${group.name}`" @confirm="emit('delete', group)">
                <v-card-text>
                  Are you sure that you want to delete this habitat group?
                  <ul>
                    <li>{{ group.children.length }} habitat(s) will be removed as well</li>
                    <li>
                      All downstream habitat groups depending on the removed habitats will also be
                      deleted
                    </li>
                    <li>
                      Survey data associated with the habitats will NOT be deleted, but will lose
                      the association.
                    </li>
                  </ul>
                  <v-alert variant="tonal" color="warning">
                    <strong>Warning:</strong> This action cannot be undone.
                  </v-alert>
                </v-card-text>
                <v-divider></v-divider>
                <template #activator="{ props }">
                  <v-list-item
                    title="Delete"
                    prepend-icon="mdi-delete"
                    v-bind="props"
                    @click.stop=""
                  >
                  </v-list-item>
                </template>
              </ConfirmDialog>
            </v-list>
          </v-menu>
        </div>
      </div>
      <v-expand-transition>
        <v-list v-if="isOpen">
          <template v-for="(child, index) in children" :key="child.id">
            <v-divider v-if="index > 0"></v-divider>
            <v-list-group v-if="child.children?.length" :title="child.name">
              <template #activator="{ props }">
                <v-list-item
                  :title="child.name"
                  :active="activeHabitat?.id === child.id"
                  color="primary"
                  @click="activeHabitat = child"
                  v-bind="props"
                >
                  <template #title>
                    {{ child.name }}
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
                v-model:active-habitat="activeHabitat"
                :group="getGroup(group.id)!"
                :depth="depth + 1"
                :terminal="index === child.children.length - 1"
                @edit="(g) => emit('edit', g)"
                @delete="(g) => emit('delete', g)"
                @add-group="(parent) => emit('add-group', parent)"
              />
            </v-list-group>
            <!-- <v-menu v-else content-class="bg-main">
              <v-list class="bg-main">
                <v-list-item v-if="child.upstream">
                  {{
                    child.upstream
                      ?.concat(child)
                      .map((u) => u.label)
                      .join(' 〉 ')
                  }}
                  <template #append>
                    <span class="text-muted text-caption"> Path </span>
                  </template>
                </v-list-item>
                <v-list-item
                  title="Add child group"
                  prepend-icon="mdi-plus"
                  @click.stop="emit('add-group', child)"
                >
                </v-list-item>
              </v-list>
              <template #activator="{ props, isActive }"> -->
            <v-list-item
              v-else
              :title="child.name"
              color="primary"
              :active="activeHabitat?.id === child.id"
              @click="activeHabitat = child"
            >
              <template #prepend>
                <v-icon icon="mdi-circle-medium"></v-icon>
              </template>
            </v-list-item>
            <!-- </template>
            </v-menu> -->
          </template>
        </v-list>
      </v-expand-transition>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/toolkit/ui/ConfirmDialog.vue'
import {
  HabitatGroupNode,
  HabitatNode,
  useHabitats
} from '@/features/registries/composables/habitats'
import { computed, ref } from 'vue'

const { group, depth = 0 } = defineProps<{
  group: HabitatGroupNode
  depth?: number
  terminal?: boolean
}>()

const { habitatGraph, getGroup, getHabitat } = useHabitats()
const children = computed(() => {
  return group.children.map((child) => getHabitat(child.id)!)
})

const isOpen = defineModel<boolean>('open', {
  default: false
})

const activeHabitat = defineModel<HabitatNode | undefined>('activeHabitat')

const emit = defineEmits<{
  edit: [group: HabitatGroupNode]
  delete: [group: HabitatGroupNode]
  'add-group': [group: HabitatNode]
}>()

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
  width: calc(var(--v-list-indent) + var(--list-indent-size) - 8px);
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
