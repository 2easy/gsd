import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'
import { v4 as uuidv4 } from 'uuid'

export interface Project {
  id: string;
  name: string;
  position: number;
  deadline?: string;
}

export interface NextAction {
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

export const useNextActionsStore = defineStore('nextActions', () => {
  const nextActions = ref<NextAction[]>([])
  const projects = ref<Project[]>([])
  const filterBy = ref<'all' | 'active' | 'completed'>('active')
  const sortBy = ref<'position' | 'energy' | 'size'>('position')
  const sortDirection = ref<'asc' | 'desc'>('asc')

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

  async function fetchData() {
    const [actionsResponse, projectsResponse] = await Promise.all([
      axios.get('/api/next-actions'),
      axios.get('/api/projects')
    ]);
    
    nextActions.value = actionsResponse.data.map((action: Partial<NextAction>) => {
      if (!action.created_at) {
        console.warn('Missing created_at for action:', action);
        action.created_at = new Date().toISOString();
      }
      return action as NextAction;
    });
    projects.value = projectsResponse.data;
  }

  async function addNextAction(actionText: string) {
    if (actionText.trim() === '') return;

    const maxPosition = Math.max(0, ...nextActions.value.map(a => a.position || 0));
    const newAction: Partial<NextAction> = {
      id: uuidv4(),
      action: actionText,
      position: maxPosition + 1,
      created_at: new Date().toISOString(),
      completed_at: undefined
    };

    const response = await axios.post('/api/next-actions', newAction);
    nextActions.value.push(response.data);
  }

  async function updateActionPosition(movedAction: NextAction, newIndex: number) {
    if (sortBy.value !== 'position') return;

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
  }

  async function toggleComplete(action: NextAction) {
    try {
      const response = await axios.patch(`/api/next-actions/${action.id}`, {
        completed_at: action.completed_at ? null : new Date().toISOString()
      });
      
      const index = nextActions.value.findIndex(a => a.id === action.id);
      if (index !== -1) {
        nextActions.value[index] = response.data;
      }
    } catch (error) {
      console.error('Failed to toggle completion:', error);
    }
  }

  async function updateActionField(action: NextAction, field: keyof NextAction, value: any) {
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
  }

  async function deleteAction(actionId: string) {
    try {
      await axios.delete(`/api/next-actions/${actionId}`);
      nextActions.value = nextActions.value.filter(a => a.id !== actionId);
    } catch (error) {
      console.error('Failed to delete action:', error);
    }
  }

  function getProjectName(projectId: string) {
    const project = projects.value.find(p => p.id === projectId);
    return project?.name || '';
  }

  function toggleSort(sort: typeof sortBy.value) {
    if (sortBy.value === sort) {
      sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
    } else {
      sortBy.value = sort;
      if (sort === 'position') sortDirection.value = 'asc';
      else if (sort === 'energy') sortDirection.value = 'desc';
      else if (sort === 'size') sortDirection.value = 'asc';
    }
  }

  return {
    nextActions,
    projects,
    filterBy,
    sortBy,
    sortDirection,
    sortedActions,
    fetchData,
    addNextAction,
    updateActionPosition,
    toggleComplete,
    updateActionField,
    deleteAction,
    getProjectName,
    toggleSort
  }
}) 