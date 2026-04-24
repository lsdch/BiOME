<template>
  <div id="habitats-view-container" class="fill-height" @click="activeHabitat = undefined">
    <div class="fill-height overflow-y-auto">
      <v-container id="habitats-list-container">
        <HabitatGroupBox
          v-for="group in habitatGraph.rootGroups"
          :key="group.id"
          :group="group"
          open
          v-model:active-habitat="activeHabitat"
          @click.stop
          @add-group="(parent) => openForm(undefined, parent)"
          @edit="(group) => openForm(group)"
          @delete="({ label }) => deleteHabitatGroup({ path: { label } })"
        />

        <HabitatGroupFormDialogMutation
          v-if="isGranted('Maintainer')"
          @success="refetch()"
          v-model:dialog="formDialog"
          :item="selectedGroup"
          :depends="parentHabitat?.label"
        >
        </HabitatGroupFormDialogMutation>
      </v-container>
    </div>
    <!-- attach="#router-view-suspense-container" -->
    <!-- :model-value="activeHabitat !== undefined" -->
    <v-bottom-sheet
      attach="#habitats-view-container"
      :model-value="true"
      contained
      :scrim="false"
      persistent
      :close-on-back="false"
      :capture-focus="false"
      no-click-animation
      stick-to-target
      scroll-strategy="block"
      @click.stop
    >
      <v-card>
        <v-list>
          <v-list-item
            v-if="activeHabitat?.upstream"
            :title="
              activeHabitat.upstream
                ?.concat(activeHabitat)
                .map((u) => u.label)
                .join(' 〉 ')
            "
          >
            <template #title>
              {{ activeHabitat.upstream?.map((u) => u.label).join(' 〉 ') }}
              <span> 〉 <v-chip :text="activeHabitat.label"></v-chip></span>
            </template>
            <template #append>
              <span class="text-muted text-caption"> {{ activeHabitat.group.label }} </span>
            </template>
          </v-list-item>
          <template v-if="isGranted('Maintainer')">
            <v-list-item
              v-if="activeHabitat"
              title="Add child group"
              :active="formDialog"
              color="primary"
              prepend-icon="mdi-plus"
              @click.stop="openForm(undefined, activeHabitat)"
            >
            </v-list-item>
            <v-list-item
              v-else
              title="Add root habitat group"
              color="primary"
              :active="formDialog"
              prepend-icon="mdi-plus"
              @click="openForm()"
            />
          </template>
        </v-list>
      </v-card>
    </v-bottom-sheet>
  </div>
</template>

<script setup lang="ts">
import { deleteHabitatGroupMutation } from '@/api/gen/@tanstack/vue-query.gen'
import HabitatGroupFormDialogMutation from '@/components/forms/HabitatGroupFormDialogMutation.vue'
import HabitatGroupBox from '@/features/registries/components/HabitatGroupBox.vue'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { useMutation } from '@tanstack/vue-query'
import { ref } from 'vue'
import { HabitatGroupNode, HabitatNode, useHabitats } from '../composables/habitats'

const { habitatGraph, refetch } = useHabitats()
const { isGranted } = useUserStore()

const formDialog = ref(false)
const selectedGroup = ref<HabitatGroupNode>()
const parentHabitat = ref<HabitatNode>()

const activeHabitat = ref<HabitatNode>()

const { feedback } = useFeedback()

const { mutate: deleteHabitatGroup } = useMutation({
  ...deleteHabitatGroupMutation(),
  onSuccess: () => {
    feedback({ type: 'success', message: 'Habitat group deleted' })
    refetch()
  },
  onError(err) {
    console.error('Failed to delete habitat group', err)
    feedback({ type: 'error', message: 'Failed to delete habitat group' })
  }
})
function openForm(group?: HabitatGroupNode, parent?: HabitatNode) {
  selectedGroup.value = group
  parentHabitat.value = parent
  formDialog.value = true
}
</script>

<style scoped>
#habitats-view-container {
  position: relative;
}
#habitats-list-container {
  padding-bottom: 50vh;
}
</style>
