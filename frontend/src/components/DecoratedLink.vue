<template>
  <a :href="url" target="_blank" class="decorated-link">
    <FontAwesomeIcon :icon="iconClass" class="link-icon" />
    <slot></slot>
  </a>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faSlack, faGithub } from '@fortawesome/free-brands-svg-icons';
import { faLink } from '@fortawesome/free-solid-svg-icons';

interface Props {
  url: string;
}

const props = defineProps<Props>();

const iconClass = computed(() => {
  if (props.url.includes('slack.com')) {
    return faSlack;
  } else if (props.url.includes('github.com')) {
    return faGithub;
  } else {
    return faLink;
  }
});
</script>

<style scoped>
.decorated-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: inherit;
}

.link-icon {
  font-size: 16px;
  margin-right: 8px;
}
</style>