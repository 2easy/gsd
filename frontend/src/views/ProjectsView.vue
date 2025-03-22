<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { onBeforeRouteUpdate, type RouteLocationNormalized } from 'vue-router';
import axios from 'axios';
import { v4 as uuidv4 } from 'uuid';
import 'bootstrap/dist/css/bootstrap.min.css';
import draggable from 'vuedraggable';
import { faGripLines } from '@fortawesome/free-solid-svg-icons';
import ProjectDatePicker from '../components/ProjectDatePicker.vue';

interface Project {
  id: string;
  name: string;
  position: number;
  deadline?: string;
}

interface NextAction {
  id: string;
  action: string;
  project_id?: string;  // Changed from projectId
  url?: string;
  size?: 'small' | 'medium' | 'big';
  energy?: 'high' | 'low';
  created_at: string;   // Changed from createdAt
  completed_at?: string;  // Changed from completedAt
  position: number;
}

const projects = ref<Project[]>([]);
const newProjectName = ref('');
const sortMethod = ref<'position' | 'deadline'>('position');
const nextActions = ref<NextAction[]>([]);
const showProjectsWithoutActions = ref<boolean>(false);
const canBeSortedManually = computed(() => {
  return sortMethod.value === 'position' && !showProjectsWithoutActions.value;
});

const sortedProjects = computed(() => {
  var projectsList = [...projects.value];

  if (sortMethod.value === 'deadline') {
    return projectsList.sort((a, b) => {
      if (!a.deadline && !b.deadline) return 0;
      if (!a.deadline) return 1;
      if (!b.deadline) return -1;
      return new Date(a.deadline).getTime() - new Date(b.deadline).getTime();
    });
  }

  return projectsList.sort((a, b) => (a.position ?? 0) - (b.position ?? 0));
});

const filteredProjects = computed({
  get: () => {
    let projectsList = sortedProjects.value;

    if (showProjectsWithoutActions.value) {
      return projectsList.filter(project => {
        const projectActions = nextActions.value.filter(action => action.project_id === project.id && !action.completed_at);
        return projectActions.length === 0;
      });
    }

    return projectsList;
  },
  set: (value) => {
    // We need to set the new value to the projects array so that
    // handleDragChange can update the positions correctly
    if (sortMethod.value === 'position') {
      projects.value = value
    }
  }
});

const getNextActionsCount = (projectId: string) => {
  return nextActions.value.filter(action => action.project_id === projectId && !action.completed_at).length;
};

const projectsWithoutActions = computed(() => {
  return projects.value.filter(project => {
    const projectActions = nextActions.value.filter(action => action.project_id === project.id && !action.completed_at);
    return projectActions.length === 0;
  }).length;
});

const fetchData = async () => {
  const [projectsResponse, actionsResponse] = await Promise.all([
    axios.get('/api/projects'),
    axios.get('/api/next-actions')
  ]);
  projects.value = projectsResponse.data;
  nextActions.value = actionsResponse.data;
};

onMounted(fetchData);

// Add navigation hook to refresh data when route is updated
onBeforeRouteUpdate(async (to: RouteLocationNormalized) => {
  if (to.path === '/projects') {
    await fetchData();
  }
});

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && newProjectName.value.trim() !== '') {
    addProject();
  }
};

const addProject = async () => {
  if (newProjectName.value.trim() === '') return;

  const response = await axios.post('/api/projects', { 
    id: uuidv4(), 
    name: newProjectName.value,
  });
  projects.value.push(response.data);
  newProjectName.value = '';
};

const draggableHandleChange = async (event: any) => {
  if (sortMethod.value === 'deadline') {
    console.log('Sorting by deadline, ignoring drag');
    return;
  }

  if (event.moved) {
    const { newIndex, element } = event.moved;
    const prevProject = newIndex > 0 ? projects.value[newIndex - 1] : null;
    const nextProject = newIndex < projects.value.length - 1 ? projects.value[newIndex + 1] : null;

    let newPosition;
    if (!prevProject) {
      newPosition = nextProject ? nextProject.position / 2 : 1;
    } else if (!nextProject) {
      newPosition = prevProject.position + 1;
    } else {
      newPosition = (prevProject.position + nextProject.position) / 2;
    }

    try {
      const response = await axios.patch(`/api/projects/${element.id}`, {
        position: newPosition
      });
      element.position = response.data.position;
    } catch (error) {
      console.error('Failed to update project position:', error);
    }
  }
};

const deleteProject = async (projectId: string) => {
  try {
    console.log(projectId);
    await axios.delete(`/api/projects/${projectId}`);
    projects.value = projects.value.filter(p => p.id !== projectId);
  } catch (error) {
    console.error('Failed to delete project:', error);
  }
};

const updateProjectDeadline = async (date: string | null | undefined, project: Project) => {
  try {
    console.log('Updating deadline for project:', project.id, 'to', date);
    const response = await axios.patch(`/api/projects/${project.id}`, {
      deadline: date
    });

    // Update the project in the main list with the returned data
    const index = projects.value.findIndex(p => p.id === project.id);
    if (index !== -1) {
      projects.value[index] = response.data;
    }
    console.log('Project updated:', response.data);
  } catch (error) {
    console.error('Failed to update project deadline:', error);
  }
};

</script>

<template>
  <div class="container py-4">
    <header class="pb-3 mb-4 border-bottom">
      <h1 class="display-4 fw-bold text-center text-primary">Projects</h1>
    </header>

    <main>
      <div class="row">
        <div class="col-md-10 mx-auto">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <div class="input-group" style="max-width: 50%;">
              <input 
                v-model="newProjectName" 
                type="text" 
                class="form-control" 
                placeholder="New project name" 
                aria-label="New project name"
                @keypress="handleKeyPress"
              />
            </div>

            <div class="d-flex align-items-center gap-3">
              <div class="btn-group me-2">
                <button 
                  class="btn position-relative"
                  :class="[
                    showProjectsWithoutActions ? 'btn-warning' : 'btn-outline-warning',
                    { 'disabled': projectsWithoutActions === 0 }
                  ]"
                  @click="showProjectsWithoutActions = !showProjectsWithoutActions"
                  :title="projectsWithoutActions === 0 ? 'No projects without actions' : undefined"
                >
                  {{ showProjectsWithoutActions ? 'Show All' : 'Show w/o Actions' }}
                  <span 
                    v-if="projectsWithoutActions > 0 && !showProjectsWithoutActions" 
                    class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger"
                  >
                    {{ projectsWithoutActions }}
                  </span>
                </button>
              </div>

              <div class="btn-group">
                <button 
                  class="btn"
                  :class="sortMethod === 'position' ? 'btn-primary' : 'btn-outline-primary'"
                  @click="sortMethod = 'position'"
                >
                  Position
                </button>
                <button 
                  class="btn"
                  :class="sortMethod === 'deadline' ? 'btn-primary' : 'btn-outline-primary'"
                  @click="sortMethod = 'deadline'"
                >
                  Deadline
                </button>
              </div>
            </div>
          </div>

          <draggable
            v-model="filteredProjects"
            :disabled="sortMethod === 'deadline'"
            :animation="200"
            item-key="id"
            handle=".drag-handle"
            class="list-group"
            ghost-class="ghost"
            @change="draggableHandleChange"
          >
            <template #item="{ element }">
              <div class="list-group-item" :class="{ 'not-draggable': sortMethod === 'deadline' }">
                <FontAwesomeIcon :icon="faGripLines" v-if="canBeSortedManually" class="drag-handle"/>
                {{ (element.name.split(' ').map((word: string): string => 
                  word.charAt(0).toUpperCase() + word.slice(1)) as string[]).join(' ') }}
                <!-- {{ element.position }} -->

                <span class="badge me-2" :class="getNextActionsCount(element.id) === 0 ? 'bg-danger' : 'bg-secondary'">
                  {{ getNextActionsCount(element.id) }} Actions
                </span>

                <ProjectDatePicker
                  v-model="element.deadline"
                  :project-id="element.id"
                  @update:model-value="(date) => updateProjectDeadline(date, element)"
                  @clear="() => updateProjectDeadline('', element)"
                />

                <button 
                  class="btn-close float-end" 
                  aria-label="Close" 
                  @click="deleteProject(element.id)"
                ></button>

              </div>
            </template>
          </draggable>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.project-title {
  font-size: 1.3rem;
  color: #0d6efd;
  font-weight: 600;
  margin: 0;
}

.cursor-move {
  cursor: grab;
}

.cursor-move:active {
  cursor: grabbing;
}

.btn-close {
  font-size: 0.8rem;
  opacity: 0.5;
  transition: opacity 0.2s;
}

.btn-close:hover {
  opacity: 1;
}

.drag-handle {
  cursor: grab;
  margin-right: 10px;
}

.drag-handle:active {
  cursor: grabbing;
}

.list-group-item {
  cursor: default;
}

/* Remove the modal overlay styles since we don't need them anymore */
.modal-overlay,
.date-picker-modal::before {
  display: none;
}

.project-title {
  cursor: pointer;
  transition: color 0.2s;
}

.project-title:hover {
  color: #0b5ed7;
}
</style>
