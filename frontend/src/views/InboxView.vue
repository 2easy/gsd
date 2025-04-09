<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { useInboxStore } from '@/stores/inbox'
import AgeBadge from '@/components/AgeBadge.vue';
import DecoratedLink from '@/components/DecoratedLink.vue';

const inboxStore = useInboxStore();
const newItemText = ref('');

onMounted(async () => {
  await inboxStore.fetchInboxItems();
});

async function createInboxItem() {
  if (!newItemText.value.trim()) return;

  const text = newItemText.value.trim();
  const urlRegex = /(https?:\/\/[^\s]+)/g;
  const urlMatch = text.match(urlRegex);
  const url = urlMatch ? urlMatch[0] : undefined;
  let description = url ? text.replace(url, '').trim() : text;
  
  // If description is empty and we have a URL, use truncated URL as description
  if (!description && url) {
    description = url.length > 50 ? url.substring(0, 47) + '...' : url;
  }

  await inboxStore.createInboxItem(description, url);
  newItemText.value = '';
}

const router = useRouter();
function goToProcessInbox() {
  router.push('/process-inbox');
}
</script>

<template>
  <div class="inbox-container">
    <h1 class="d-flex justify-content-between align-items-center">
      Inbox
      <button 
        @click="goToProcessInbox" 
        class="btn btn-primary" 
        :disabled="inboxStore.inboxItemCount === 0"
      >
        Process Items
      </button>
    </h1>
    <input
      v-model="newItemText"
      @keyup.enter="createInboxItem"
      placeholder="Add a new inbox item"
      class="form-control mb-3"
    />
    <div class="list-group">
      <div v-for="item in inboxStore.inboxItems" :key="item.id" class="list-group-item">
        <template v-if="item.url">
          <DecoratedLink :url="item.url">{{ item.description }}</DecoratedLink>
        </template>
        <template v-else>
          <span>{{ item.description }}</span>
        </template>
        <AgeBadge :createdAt="item.created_at" />
      </div>
    </div>
  </div>
</template>