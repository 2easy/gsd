<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <main>
    <h1 class="display-4 text-center mb-5">Getting Things Done</h1>
    
    <div class="row g-4">
      <div class="col-md-6">
        <RouterLink class="text-decoration-none" to="/projects">
          <div class="card h-100 shadow-sm hover-lift">
            <div class="card-body text-center">
              <h3 class="card-title text-primary">Projects</h3>
              <p class="card-text text-muted">Manage and organize your projects</p>
            </div>
          </div>
        </RouterLink>
      </div>
      
      <div class="col-md-6">
        <RouterLink class="text-decoration-none" to="/next-actions">
          <div class="card h-100 shadow-sm hover-lift">
            <div class="card-body text-center">
              <h3 class="card-title text-danger">Next Actions</h3>
              <p class="card-text text-muted">Track and complete your next actions</p>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </main>
</template>

<style scoped>
.hover-lift {
  transition: transform 0.2s ease;
}

.hover-lift:hover {
  transform: translateY(-5px);
}
</style>
