<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { onBeforeRouteUpdate, type RouteLocationNormalized } from 'vue-router';
import draggable from 'vuedraggable';
import AgeBadge from '../components/AgeBadge.vue';
import { useNextActionsStore, type NextAction } from '../stores/nextActions';

const store = useNextActionsStore();
const newActionText = ref('');
const editingAction = ref<NextAction | null>(null);
const editForm = ref({ action: '', url: '' });

onMounted(store.fetchData);

// Add navigation hook to refresh data when route is updated
onBeforeRouteUpdate(async (to: RouteLocationNormalized) => {
  if (to.path === '/next-actions') {
    await store.fetchData();
  }
});

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && newActionText.value.trim() !== '') {
    addNextAction();
  }
};

const addNextAction = async () => {
  if (newActionText.value.trim() === '') return;
  await store.addNextAction(newActionText.value);
  newActionText.value = '';
};

const handleChange = async (event: any) => {
  if (store.sortBy !== 'position') return;
  
  if (event.moved) {
    const { newIndex } = event.moved;
    const movedAction = store.sortedActions[newIndex];
    await store.updateActionPosition(movedAction, newIndex);
  }
};

const handleInputChange = (event: Event, action: NextAction, field: 'project_id' | 'size' | 'energy') => {
  const target = event.target as HTMLInputElement | HTMLSelectElement;
  const value = target.value || undefined;
  store.updateActionField(action, field, value);
};

const openEditModal = (action: NextAction) => {
  editingAction.value = action;
  editForm.value = {
    action: action.action,
    url: action.url || ''
  };
};

const saveEdit = async () => {
  if (!editingAction.value) return;
  
  if (!editForm.value.action.trim()) {
    alert('Action text cannot be empty');
    return;
  }
  
  try {
    await store.updateActionField(editingAction.value, 'action', editForm.value.action.trim());
    if (editForm.value.url) {
      await store.updateActionField(editingAction.value, 'url', editForm.value.url.trim());
    }
    editingAction.value = null;
  } catch (error) {
    console.error('Failed to update action:', error);
    alert('Failed to save changes. Please try again.');
  }
};
</script>

<template>
  <div class="container py-4">
    <header class="pb-3 mb-4 border-bottom">
      <h1 class="display-4 fw-bold text-center text-primary">Next Actions</h1>
    </header>

    <main>
      <div class="row">
        <div class="col-md-10 mx-auto">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <div class="input-group" style="max-width: 50%;">
              <input 
                v-model="newActionText" 
                type="text" 
                class="form-control" 
                placeholder="New action" 
                aria-label="New action"
                @keypress="handleKeyPress"
              />
            </div>

            <div class="d-flex align-items-center gap-3">
              <div class="btn-group me-2">
                <button 
                  class="btn position-relative"
                  :class="[
                    store.filterBy === 'completed' ? 'btn-success' : 'btn-outline-success',
                  ]"
                  @click="store.filterBy = store.filterBy === 'completed' ? 'all' : 'completed'"
                >
                  {{ store.filterBy === 'completed' ? 'Show All' : 'Show Completed' }}
                </button>
              </div>

              <div class="btn-group">
                <button 
                  class="btn"
                  :class="store.sortBy === 'position' ? 'btn-primary' : 'btn-outline-primary'"
                  @click="store.toggleSort('position')"
                >
                  Position
                </button>
                <button 
                  class="btn"
                  :class="store.sortBy === 'energy' ? 'btn-primary' : 'btn-outline-primary'"
                  @click="store.toggleSort('energy')"
                >
                  Energy
                </button>
                <button 
                  class="btn"
                  :class="store.sortBy === 'size' ? 'btn-primary' : 'btn-outline-primary'"
                  @click="store.toggleSort('size')"
                >
                  Size
                </button>
              </div>
            </div>
          </div>

          <draggable 
            :list="store.sortedActions"
            class="list-group mb-4" 
            item-key="id"
            @change="handleChange"
            :animation="200"
            :disabled="store.sortBy !== 'position'"
          >
            <template #item="{ element: action }">
              <div class="list-group-item list-group-item-action">
                <div class="d-flex w-100 justify-content-between align-items-center">
                  <div class="d-flex align-items-center flex-grow-1">
                    <div class="drag-handle me-3" v-if="store.sortBy === 'position'">⋮⋮</div>
                    <div class="form-check me-3">
                      <input
                        class="form-check-input"
                        type="checkbox"
                        :checked="!!action.completed_at"
                        @change="store.toggleComplete(action)"
                      >
                    </div>
                    <div class="d-flex flex-column flex-grow-1">
                      <div class="d-flex justify-content-between align-items-center">
                        <div class="d-flex align-items-center">
                          <h5 class="mb-1" :class="{ 'completed-action': action.completed_at }">
                            <i v-if="action.url" class="bi bi-link-45deg me-1"></i>
                            <a v-if="action.url" 
                               :href="action.url" 
                               target="_blank" 
                               rel="noopener noreferrer"
                               class="text-reset">
                                {{ action.action.charAt(0).toUpperCase() + action.action.slice(1) }}
                            </a>
                            <span v-else>{{ action.action.charAt(0).toUpperCase() + action.action.slice(1) }}</span>
                          </h5>
                          <AgeBadge :createdAt="action.created_at" />
                          <button 
                            class="btn btn-link btn-sm p-0 ms-2" 
                            @click="openEditModal(action)"
                            title="Edit action"
                          >
                            <i class="bi bi-pencil"></i>
                          </button>
                        </div>
                        <div class="d-flex align-items-center">
                          <select
                            class="form-select form-select-sm me-2"
                            style="width: auto;"
                            :value="action.project_id || ''"
                            @change="(e) => handleInputChange(e, action, 'project_id')"
                          >
                            <option value="">No Project</option>
                            <option v-for="project in store.projects" :key="project.id" :value="project.id">
                              {{ project.name }}
                            </option>
                          </select>

                          <select
                            class="form-select form-select-sm me-2"
                            style="width: auto;"
                            :value="action.size || ''"
                            @change="(e) => handleInputChange(e, action, 'size')"
                          >
                            <option value="">Size</option>
                            <option value="small">Small</option>
                            <option value="medium">Medium</option>
                            <option value="big">Big</option>
                          </select>

                          <select
                            class="form-select form-select-sm me-2"
                            style="width: auto;"
                            :value="action.energy || ''"
                            @change="(e) => handleInputChange(e, action, 'energy')"
                          >
                            <option value="">Energy</option>
                            <option value="high">High</option>
                            <option value="low">Low</option>
                          </select>

                          <button
                            class="btn btn-close"
                            aria-label="Delete action"
                            @click="store.deleteAction(action.id)"
                          ></button>
                        </div>
                      </div>
                      <div v-if="action.project_id" class="text-muted small">
                        Project: {{ store.getProjectName(action.project_id) }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </div>
    </main>

    <!-- Edit Modal -->
    <div v-if="editingAction" class="modal fade show" style="display: block;" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Edit Action</h5>
            <button type="button" class="btn-close" @click="editingAction = null"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveEdit">
              <div class="mb-3">
                <label class="form-label">Action</label>
                <input type="text" class="form-control" v-model="editForm.action">
              </div>
              <div class="mb-3">
                <label class="form-label">URL (optional)</label>
                <input type="url" class="form-control" v-model="editForm.url" placeholder="https://...">
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="editingAction = null">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveEdit">Save changes</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="editingAction" class="modal-backdrop fade show"></div>
  </div>
</template>

<style scoped>
.drag-handle {
  cursor: grab;
  display: flex;
  align-items: center;
  font-size: 24px;
  letter-spacing: -0.2em;
  line-height: 1;
  color: #6c757d;
  font-weight: bold;
  opacity: 0.7;
}

.drag-handle:active {
  cursor: grabbing;
}

.list-group-item {
  cursor: default;
  transition: background-color 0.2s;
}

.list-group-item:hover {
  background-color: rgba(0, 0, 0, 0.02);
}

.form-select {
  cursor: pointer;
}

.btn-close {
  font-size: 0.8rem;
  opacity: 0.5;
  transition: opacity 0.2s;
}

.btn-close:hover {
  opacity: 1;
}

.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: -1;
}

.btn span {
  margin-left: 4px;
}

.completed-action {
  text-decoration: line-through;
  opacity: 0.6;
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

.modal-content {
  background-color: var(--bs-body-bg);
  border-color: var(--bs-border-color);
}

.modal-header {
  border-bottom-color: var(--bs-border-color);
}

.modal-footer {
  border-top-color: var(--bs-border-color);
}

.btn-close {
  filter: var(--bs-btn-close-white-filter);
}
</style>