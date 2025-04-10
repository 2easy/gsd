import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

type InboxItem = {
  id: string;
  description: string;
  url?: string;
  created_at: string;
};

export const useInboxStore = defineStore('inbox', () => {
  const inboxItems = ref<InboxItem[]>([])

  const inboxItemCount = computed(() => inboxItems.value.length)
  let ws: WebSocket | null = null

  function initWebSocket() {
    ws = new WebSocket('ws://localhost:8081/api/ws')
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.type === 'inbox_item_created') {
        inboxItems.value.push(data.data)
      }
    }
    ws.onclose = () => {
      setTimeout(() => { initWebSocket() }, 1000)
    }
  }

  function closeWebSocket() {
    ws?.close()
  }

  async function fetchInboxItems() {
    try {
      const response = await axios.get('/api/inbox');
      inboxItems.value = Array.isArray(response.data) ? response.data : [];
    } catch (error) {
      console.error('Failed to fetch inbox items:', error);
    }
  }

  async function createInboxItem(description: string, url?: string) {
    try {
      const response = await axios.post('/api/inbox', { description, url });
      console.log('Created inbox item:', response.data);
      inboxItems.value = [...inboxItems.value, response.data];
    } catch (error) {
      console.error('Failed to create inbox item:', error);
    }
  }

  return { inboxItems, inboxItemCount, fetchInboxItems, createInboxItem, initWebSocket, closeWebSocket }
})