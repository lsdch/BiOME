<template>
  <v-list>
    <v-list-item>
      <SiteAutocomplete :model-value="site" @update:model-value="addSite" class="pt-2" />
    </v-list-item>
    <v-divider />
    <v-list-item
      v-for="(site, index) in sites"
      :key="site.code"
      :title="site.name ?? site.code"
      :subtitle="site.name ? site.code : undefined"
      @click="emit('focusSite', site)"
    >
      <template #append>
        <v-icon icon="mdi-map-marker" :color="site.color"></v-icon>
        <v-btn
          variant="plain"
          icon="mdi-close-circle"
          size="small"
          color=""
          @click="sites.splice(index, 1)"
        ></v-btn>
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { SiteItem, SiteWithDistance, SiteWithScore } from '@/api'
import SiteAutocomplete from '@/features/site/components/SiteAutocomplete.vue'
import { Hsluv } from 'hsluv'
import { computed, reactive, ref } from 'vue'

const site = ref<SiteWithDistance | SiteWithScore>()

export type SiteItemWithColor = SiteItem & {
  color?: string
}

const sites = defineModel<SiteItemWithColor[]>({ default: reactive([]) })

const siteColors = computed(() => {
  return new Set(sites.value.map((s, i) => s.color))
})

const emit = defineEmits<{
  focusSite: [site: SiteItem]
}>()

function addSite(s?: SiteItem) {
  if (!s) return
  site.value = undefined
  sites.value.push({ ...s, color: generateColor() })
}

function generateColor(index: number = 0): string {
  const conv = new Hsluv()
  conv.hsluv_h = (index * 137.508) % 360 // use golden angle approximation
  conv.hsluv_s = 90
  conv.hsluv_l = 65
  conv.hsluvToHex()
  return siteColors.value.has(conv.hex) ? generateColor(index + 1) : conv.hex
}
</script>

<style scoped lang="scss"></style>
