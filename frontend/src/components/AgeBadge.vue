<template>
<span class="badge text-bg-secondary" :title="`Created ${createdAt}`">
    {{ age }}
</span>
</template>

<script setup lang="ts">
import { computed, withDefaults, defineProps } from 'vue';

const props = withDefaults(defineProps<{
  createdAt: string
}>(), {});

const age = computed(() => {
  const createdDate = new Date(props.createdAt);
  const now = new Date();
  const diffMs = now.getTime() - createdDate.getTime();
  const diffDays = diffMs / (1000 * 60 * 60 * 24);
  const roundedDays = Math.round(diffDays * 2) / 2; // Round to nearest 0.5
  if (roundedDays === 0) return 'NEW';
  return roundedDays === 1 ? '1 day' : `${roundedDays} days`;
});
</script>

<style scoped>
.badge {
  margin: 0 0.5em; /* Added margin on both sides */
}
</style>