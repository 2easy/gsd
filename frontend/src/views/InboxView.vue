<template>
  <div class="inbox-container">
    <h1 class="d-flex justify-content-between align-items-center">
      Inbox
      <button 
        @click="goToProcessItems" 
        class="btn btn-primary" 
        :disabled="inboxItems.length === 0"
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
      <div v-for="item in inboxItems" :key="item.id" class="list-group-item">
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

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import DecoratedLink from '@/components/DecoratedLink.vue';
import { useRouter } from 'vue-router';
import AgeBadge from '@/components/AgeBadge.vue';

type InboxItem = {
  id: string;
  description: string;
  url?: string;
  created_at: string;
};

const inboxItems = ref<InboxItem[]>([]); // Ensure inboxItems is initialized as an array
const newItemText = ref('');

const router = useRouter();

onMounted(async () => {
  try {
    const response = await axios.get('/api/inbox');
    inboxItems.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error('Failed to fetch inbox items:', error);
  }
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

  try {
    const response = await axios.post('/api/inbox', { description, url });
    console.log('Created inbox item:', response.data);
    inboxItems.value.push(response.data);
    newItemText.value = '';
  } catch (error) {
    console.error('Failed to create inbox item:', error);
  }
}

function goToProcessItems() {
  router.push('/process-items');
}
</script>