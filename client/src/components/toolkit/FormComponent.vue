<script setup lang="ts" generic="ItemType extends { id: string }">
defineProps<{
  onSuccess: any
}>()

// defineEmits<{
//   (event: 'created', item: ItemType): void
//   (event: 'edited', item: ItemType): void
// }>()
</script>
