<template>
  <div>
    <h1>Process Inbox Items</h1>
    <ProcessInboxWizard 
      :inbox-items="inboxItems" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import ProcessInboxWizard from '@/components/ProcessInboxWizard.vue';

const inboxItems = ref<{id: string, description: string, url?: string}[]>([]);
const currentItem = ref<{id: string, description: string, url?: string} | null>(null);

const fetchInboxItems = async () => {
  try {
    const response = await axios.get('/api/inbox');
    inboxItems.value = response.data;
    currentItem.value = inboxItems.value[0] || null;
  } catch (error) {
    console.error('Failed to fetch inbox items:', error);
  }
};

onMounted(fetchInboxItems);

</script>

<style scoped>
/* Add styles here if needed */
</style>