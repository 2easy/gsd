<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
import { ref, onMounted } from 'vue';
import axios from 'axios';

const inboxItemCount = ref(0);

onMounted(async () => {
  try {
    const response = await axios.get('/api/inbox');
    inboxItemCount.value = response.data.length;
  } catch (error) {
    console.error('Failed to fetch inbox items:', error);
  }
});
</script>

<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-primary mb-4">
    <div class="container">
      <RouterLink class="navbar-brand" to="/"><i class="bi bi-list-check"></i> GSD</RouterLink>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav">
          <li class="nav-item">
            <RouterLink class="nav-link" active-class="active" to="/projects">Projects</RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink class="nav-link" active-class="active" to="/next-actions">Next Actions</RouterLink>
          </li>
          <li class="nav-item position-relative">
            <RouterLink class="nav-link d-flex align-items-center" active-class="active" to="/inbox">
              Inbox
              <span v-if="inboxItemCount > 0" class="badge bg-danger ms-2">{{ inboxItemCount }}</span>
            </RouterLink>
          </li>
        </ul>
      </div>
    </div>
  </nav>

  <div class="container">
    <RouterView />
  </div>
</template>

<style scoped>
.router-link-active {
  font-weight: bold;
}

.badge {
  font-size: 0.75rem;
  padding: 0.25em 0.5em;
}
</style>
