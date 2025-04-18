<template>
  <div class="container py-3">
    <div class="card">
      <div class="card-header">
        <h3 class="mb-0">Create a New Action</h3>
      </div>
      <div class="card-body">
        <form @submit.prevent="submitAction">
          <div class="mb-3">
            <label for="actionText" class="form-label">Action:</label>
            <input 
              type="text" 
              class="form-control" 
              id="actionText" 
              v-model="actionText" 
              required 
              :placeholder="inboxItem?.description || 'Enter action description'"
            />
          </div>
          <div class="mb-3">
            <label for="url" class="form-label">URL:</label>
            <input 
              type="url" 
              class="form-control" 
              id="url" 
              v-model="url" 
              :placeholder="inboxItem?.url || 'https://...'"
            />
          </div>
          <div class="mb-3">
            <label for="size" class="form-label">Size:</label>
            <select class="form-select" id="size" v-model="size">
              <option value="">Select size</option>
              <option value="small">Small</option>
              <option value="medium">Medium</option>
              <option value="big">Big</option>
            </select>
          </div>
          <div class="mb-3">
            <label for="energy" class="form-label">Energy Level:</label>
            <select class="form-select" id="energy" v-model="energy">
              <option value="">Select energy level</option>
              <option value="low">Low</option>
              <option value="high">High</option>
            </select>
          </div>
          <div class="d-grid gap-2">
            <button type="submit" class="btn btn-primary">Save Action</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits } from 'vue';
import axios from 'axios';

const props = defineProps<{
  inboxItem?: {
    description?: string;
    url?: string;
  };
}>();

const emit = defineEmits(['actionCreated']);

const actionText = ref(props.inboxItem?.description || '');
const url = ref(props.inboxItem?.url || '');
const size = ref('');
const energy = ref('');

const submitAction = async () => {
  try {
    const newAction = {
      action: actionText.value,
      url: url.value || null,
      size: size.value || null,
      energy: energy.value || null,
    };

    await axios.post('/api/next-actions', newAction);
    emit('actionCreated');

    // Reset the form fields
    actionText.value = '';
    url.value = '';
    size.value = '';
    energy.value = '';
  } catch (error) {
    console.error('Error saving action:', error);
    alert('Failed to save action. Please try again.');
  }
};
</script>

<style scoped>
.card {
  background-color: var(--bs-card-bg);
  border-color: var(--bs-border-color);
}

.card-header {
  background-color: var(--bs-card-cap-bg);
  border-bottom-color: var(--bs-border-color);
}

.form-control, .form-select {
  background-color: var(--bs-body-bg);
  border-color: var(--bs-border-color);
  color: var(--bs-body-color);
}

.form-control:focus, .form-select:focus {
  background-color: var(--bs-body-bg);
  border-color: var(--bs-primary);
  color: var(--bs-body-color);
  box-shadow: 0 0 0 0.25rem rgba(var(--bs-primary-rgb), 0.25);
}
</style>