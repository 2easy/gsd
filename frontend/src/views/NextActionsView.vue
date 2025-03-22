<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { onBeforeRouteUpdate, type RouteLocationNormalized } from 'vue-router';
import axios from 'axios';
import { v4 as uuidv4 } from 'uuid';
import 'bootstrap/dist/css/bootstrap.min.css';
import draggable from 'vuedraggable';

interface Project {
  id: string;
  name: string;
  position: number;
  deadline?: string;
}

interface NextAction {
  id: string;
  action: string;
  project_id?: string;
  url?: string;
  size?: 'small' | 'medium' | 'big';
  energy?: 'high' | 'low';
  created_at: string;
  completed_at?: string;
  position: number;
}

const nextActions = ref<NextAction[]>([]);
const projects = ref<Project[]>([]);
const newActionText = ref('');
const filterBy = ref<'all' | 'active' | 'completed'>('active');
const sortBy = ref<'position' | 'energy' | 'size'>('position');
const sortDirection = ref<'asc' | 'desc'>('asc');

const sortedActions = computed(() => {
  let filtered = nextActions.value.filter(action => {
    if (filterBy.value === 'active') return !action.completed_at;
    if (filterBy.value === 'completed') return !!action.completed_at;
    return true;
  });

  return filtered.sort((a, b) => {
    let comparison = 0;
    if (sortBy.value === 'position') {
      comparison = (a.position || 0) - (b.position || 0);
    } else if (sortBy.value === 'energy') {
      if (!a.energy && !b.energy) comparison = 0;
      else if (!a.energy) comparison = 1;
      else if (!b.energy) comparison = -1;
      else comparison = a.energy === 'high' ? -1 : 1;
    } else if (sortBy.value === 'size') {
      const sizeOrder = { small: 1, medium: 2, big: 3 };
      if (!a.size && !b.size) comparison = 0;
      else if (!a.size) comparison = 1;
      else if (!b.size) comparison = -1;
      else comparison = (sizeOrder[a.size] || 0) - (sizeOrder[b.size] || 0);
    }
    return sortDirection.value === 'asc' ? comparison : -comparison;
  });
});

const fetchData = async () => {
  const [actionsResponse, projectsResponse] = await Promise.all([
    axios.get('/api/next-actions'),
    axios.get('/api/projects')
  ]);
  nextActions.value = actionsResponse.data;
  projects.value = projectsResponse.data;
};

onMounted(fetchData);

// Add navigation hook to refresh data when route is updated
onBeforeRouteUpdate(async (to: RouteLocationNormalized) => {
  if (to.path === '/next-actions') {
    await fetchData();
  }
});

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && newActionText.value.trim() !== '') {
    addNextAction();
  }
};

const addNextAction = async () => {
  if (newActionText.value.trim() === '') return;

  const maxPosition = Math.max(0, ...nextActions.value.map(a => a.position || 0));
  const newAction: Partial<NextAction> = {
    id: uuidv4(),
    action: newActionText.value,
    position: maxPosition + 1,
    created_at: new Date().toISOString(),
    completed_at: undefined
  };

  const response = await axios.post('/api/next-actions', newAction);
  nextActions.value.push(response.data);
  newActionText.value = '';
};

const updateActionPosition = async (event: any) => {
  if (sortBy.value !== 'position') return;

  const movedAction = sortedActions.value[event.newIndex];
  const prevAction = event.newIndex > 0 ? sortedActions.value[event.newIndex - 1] : null;
  const nextAction = event.newIndex < sortedActions.value.length - 1 ? sortedActions.value[event.newIndex + 1] : null;

  let newPosition;
  if (!prevAction) {
    // If it's at the start
    newPosition = nextAction ? nextAction.position / 2 : 1;
  } else if (!nextAction) {
    // If it's at the end
    newPosition = prevAction.position + 1;
  } else {
    // If it's in between
    newPosition = (prevAction.position + nextAction.position) / 2;
  }

  // If sorting in descending order, we need to invert the position
  if (sortDirection.value === 'desc') {
    const maxPosition = Math.max(...nextActions.value.map(a => a.position || 0));
    newPosition = maxPosition - newPosition + 1;
  }

  try {
    await axios.patch(`/api/next-actions/${movedAction.id}`, {
      position: newPosition
    });
    const index = nextActions.value.findIndex(a => a.id === movedAction.id);
    if (index !== -1) {
      nextActions.value[index].position = newPosition;
    }
  } catch (error) {
    console.error('Failed to update action position:', error);
  }
};

const handleChange = async (event: any) => {
  if (sortBy.value !== 'position') return;
  
  if (event.moved) {
    const { newIndex } = event.moved;
    const movedAction = sortedActions.value[newIndex];
    const prevAction = newIndex > 0 ? sortedActions.value[newIndex - 1] : null;
    const nextAction = newIndex < sortedActions.value.length - 1 ? sortedActions.value[newIndex + 1] : null;

    let newPosition;
    if (!prevAction) {
      newPosition = nextAction ? nextAction.position / 2 : 1;
    } else if (!nextAction) {
      newPosition = prevAction.position + 1;
    } else {
      newPosition = (prevAction.position + nextAction.position) / 2;
    }

    try {
      const response = await axios.patch(`/api/next-actions/${movedAction.id}`, {
        position: newPosition
      });
      const index = nextActions.value.findIndex(a => a.id === movedAction.id);
      if (index !== -1) {
        nextActions.value[index] = response.data;
      }
    } catch (error) {
      console.error('Failed to update action position:', error);
    }
  }
};

const toggleComplete = async (action: NextAction) => {
  try {
    const response = await axios.patch(`/api/next-actions/${action.id}`, {
      completed_at: action.completed_at ? null : new Date().toISOString()
    });
    
    // Update with the returned action data
    const index = nextActions.value.findIndex(a => a.id === action.id);
    if (index !== -1) {
      nextActions.value[index] = response.data;
    }
  } catch (error) {
    console.error('Failed to toggle completion:', error);
  }
};

const updateActionField = async (action: NextAction, field: keyof NextAction, value: any) => {
  try {
    const response = await axios.patch(`/api/next-actions/${action.id}`, {
      [field]: value
    });
    const index = nextActions.value.findIndex(a => a.id === action.id);
    if (index !== -1) {
      nextActions.value[index] = response.data;
    }
  } catch (error) {
    console.error(`Failed to update ${field}:`, error);
  }
};

const deleteAction = async (actionId: string) => {
  try {
    await axios.delete(`/api/next-actions/${actionId}`);
    nextActions.value = nextActions.value.filter(a => a.id !== actionId);
  } catch (error) {
    console.error('Failed to delete action:', error);
  }
};

const getProjectName = (projectId: string) => {
  const project = projects.value.find(p => p.id === projectId);
  return project?.name || '';
};

const handleInputChange = (event: Event, action: NextAction, field: keyof NextAction) => {
  const target = event.target as HTMLInputElement | HTMLSelectElement;
  const value = target.value || undefined;
  updateActionField(action, field, value);
};

const editingAction = ref<NextAction | null>(null);
const editForm = ref({ action: '', url: '' });

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
    const response = await axios.patch(`/api/next-actions/${editingAction.value.id}`, {
      action: editForm.value.action.trim(),
      url: editForm.value.url.trim() || undefined
    });
    
    const index = nextActions.value.findIndex(a => a.id === editingAction.value?.id);
    if (index !== -1) {
      nextActions.value[index] = response.data;
    }
    editingAction.value = null; // Close modal only after successful save
  } catch (error) {
    console.error('Failed to update action:', error);
    alert('Failed to save changes. Please try again.');
  }
};

const calculateAgeDays = (created_at: string): string => {
  const created = new Date(created_at);
  const now = new Date();
  const diffMs = now.getTime() - created.getTime();
  const diffDays = diffMs / (1000 * 60 * 60 * 24);
  const roundedDays = Math.round(diffDays * 2) / 2; // Round to nearest 0.5
  return roundedDays === 1 ? '1 day' : `${roundedDays} days`;
};

const toggleSort = (sort: typeof sortBy.value) => {
  if (sortBy.value === sort) {
    // Toggle direction if clicking the same sort
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    // Set default direction for new sort type
    sortBy.value = sort;
    if (sort === 'position') sortDirection.value = 'asc';
    else if (sort === 'energy') sortDirection.value = 'desc';
    else if (sort === 'size') sortDirection.value = 'asc';
  }
};
</script>

<template>
  <div class="container py-4">
    <header class="pb-3 mb-4 border-bottom">
      <h1 class="display-4 fw-bold text-center text-danger">Next Actions</h1>
    </header>

    <main>
      <div class="row">
        <div class="col-md-10 mx-auto">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <div class="input-group" style="max-width: 70%;">
              <input
                v-model="newActionText"
                type="text"
                class="form-control"
                placeholder="Add a next action"
                aria-label="New action"
                @keypress="handleKeyPress"
              />
            </div>

            <div class="btn-group me-2">
              <button
                v-for="f in ['all', 'active', 'completed'] as const"
                :key="f"
                class="btn"
                :class="filterBy === f ? 'btn-primary' : 'btn-outline-primary'"
                @click="filterBy = f"
              >
                {{ f.charAt(0).toUpperCase() + f.slice(1) }}
              </button>
            </div>

            <div class="btn-group">
              <button
                v-for="s in ['position', 'energy', 'size'] as const"
                :key="s"
                class="btn"
                :class="sortBy === s ? 'btn-danger' : 'btn-outline-danger'"
                @click="toggleSort(s)"
              >
                {{ s.charAt(0).toUpperCase() + s.slice(1) }}
                <span v-if="sortBy === s">
                  {{ sortDirection === 'asc' ? '↑' : '↓' }}
                </span>
              </button>
            </div>
          </div>

          <draggable 
            :list="sortedActions"
            class="list-group mb-4" 
            item-key="id"
            @change="handleChange"
            :animation="200"
            :disabled="sortBy !== 'position'"
          >
            <template #item="{ element: action }">
              <div class="list-group-item list-group-item-action">
                <div class="d-flex w-100 justify-content-between align-items-center">
                  <div class="d-flex align-items-center flex-grow-1">
                    <div class="drag-handle me-3" v-if="sortBy === 'position'">⋮⋮</div>
                    <div class="form-check me-3">
                      <input
                        class="form-check-input"
                        type="checkbox"
                        :checked="!!action.completed_at"
                        @change="toggleComplete(action)"
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
                          <span class="badge bg-secondary ms-2" title="Age">{{ calculateAgeDays(action.created_at) }}</span>
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
                            <option v-for="project in projects" :key="project.id" :value="project.id">
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
                            @click="deleteAction(action.id)"
                          ></button>
                        </div>
                      </div>
                      <div v-if="action.project_id" class="text-muted small">
                        Project: {{ getProjectName(action.project_id) }}
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
    <div v-if="editingAction" class="modal d-block" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Edit Action</h5>
            <button type="button" class="btn-close" @click="editingAction = null"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">Action</label>
              <input type="text" class="form-control" v-model="editForm.action">
            </div>
            <div class="mb-3">
              <label class="form-label">URL (optional)</label>
              <input type="url" class="form-control" v-model="editForm.url" placeholder="https://...">
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="editingAction = null">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveEdit">Save changes</button>
          </div>
        </div>
      </div>
      <div class="modal-backdrop fade show"></div>
    </div>
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
</style>