<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { RouterLink, RouterView } from 'vue-router';

import { useInboxStore } from '@/stores/inbox'
import { useThemeStore } from '@/stores/theme'

const inboxStore = useInboxStore();
const themeStore = useThemeStore();

onMounted(() => {
  inboxStore.initWebSocket();
  themeStore.initTheme();
});

onUnmounted(() => {
  inboxStore.closeWebSocket();
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
              <span v-if="inboxStore.inboxItemCount > 0" class="badge bg-danger ms-2">{{ inboxStore.inboxItemCount }}</span>
            </RouterLink>
          </li>
        </ul>
        <div class="ms-auto">
          <button class="btn btn-link nav-link" @click="themeStore.toggleTheme">
            <i :class="themeStore.theme === 'light' ? 'bi bi-moon-fill' : 'bi bi-sun-fill'"></i>
          </button>
        </div>
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

.btn-link {
  color: var(--bs-navbar-color);
  text-decoration: none;
}

.btn-link:hover {
  color: var(--bs-navbar-hover-color);
}
</style>
