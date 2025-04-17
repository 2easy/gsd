<template>
  <div>
    <div class="wizard-header">
      <!-- Progress bar -->
      <div v-if="inboxItems.length > 0" class="progress mb-3" style="height: 20px;">
        <div 
          class="progress-bar bg-success" 
          role="progressbar" 
          :style="{ width: `${(currentIndex / inboxItems.length) * 100}%` }"
          :aria-valuenow="(currentIndex / inboxItems.length) * 100" 
          aria-valuemin="0" 
          aria-valuemax="100">
          {{ currentIndex }} of {{ inboxItems.length }} items processed
        </div>
      </div>
    </div>
    <div v-if="currentIndex < inboxItems.length">
    <div v-for="(stage, index) in stages" :key="index" class="stage">
      <div class="stage-content">
        <!-- Display inbox item details in the first stage -->
        <div v-if="index === 0" class="inbox-item-details">
          <p><strong>Description:</strong> {{ currentItem.description }}</p>
          <p><strong>URL:</strong> <a :href="currentItem.url" target="_blank">{{ currentItem.url }}</a></p>
        </div>
        <p>{{ stage.question }}</p>
        <div>
          <button
            v-for="(option, idx) in stage.options"
            :key="idx"
            :class="{ selected: stage.answer === option.label }"
            @click="handleAnswer(index, option)"
          >
            {{ option.label }}
          </button>
        </div>
        <!-- Render the component if defined -->
        <component
          v-if="stage.component"
          :is="stage.component"
          :inbox-item="currentItem"
          @action-created="handleActionCreated"
          @timer-complete="handleTimerComplete"
          @move-to-next-item="handleMoveToNextItem"
        />
      </div>
    </div>
    </div>
    <div v-else class="empty-state">
      <div class="alert alert-success text-center">
        <h4 class="alert-heading mb-3">🎉 All Done!</h4>
        <p class="mb-4">All inbox items have been successfully processed.</p>
        <div class="d-flex justify-content-center gap-3">
          <button class="btn btn-outline-primary" @click="goToInbox">
            <i class="bi bi-inbox"></i> Return to Inbox
          </button>
          <button class="btn btn-primary" @click="goToNextActions">
            <i class="bi bi-check2-square"></i> View Next Actions
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, type PropType } from 'vue';
import CreateActionForm from './CreateActionForm.vue';
import TimerComponent from './TimerComponent.vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

type InboxItem = {
  id: string;
  description: string;
  url?: string;
};

const router = useRouter();

const props = defineProps({
  inboxItems: {
    type: Array as PropType<InboxItem[]>,
    required: true,
  },
});

const currentIndex = ref(0)
const currentItem = computed(() => props.inboxItems[currentIndex.value])

interface Option {
  label: string;
  nextStage?: (answer: string) => Stage;
}

interface Stage {
  question: string;
  options: Option[];
  completed: boolean;
  answer?: string;
  component?: any; // Use 'any' for dynamic component rendering
  props?: Record<string, any>; // Add props field to Stage interface
}

const initializeStages = () => [
  {
    question: "Is this item actionable?",
    options: [
      {
        label: "Yes",
        nextStage: (answer: string) => ({
          question: "Takes less than 2 minutes?",
          options: [
            {
              label: "Yes",
              nextStage: (answer: string) => ({
                question: "Do it now.",
                // Only timer controls; remove stage-level option buttons
                options: [],
                component: TimerComponent,
                completed: false,
              })
            },
            { 
              label: "No",
                nextStage: (answer: string) => ({
                    question: "Who should do it?",
                    options: [
                    {
                        label: "Me",
                        nextStage: (answer: string) => ({
                            question: "Create Action",
                            component: CreateActionForm,
                            props: {
                                inboxItem: currentItem.value,
                                onActionCreated: () => handleActionCreated(),
                            },
                            completed: false,
                        })
                    },
                    { label: "Someone Else" }
                    ],
                    completed: false,
                })
                
            },

          ],
          completed: false
        })
      },
      {
        label: "No",
        nextStage: (answer: string) => ({
          question: "What should be done with it?",
          options: [
            { label: "Trash it" },
            { label: "Reference it" }
          ],
          completed: false
        })
      }
    ],
    completed: false
  }
];

const stages = ref<Stage[]>(initializeStages());

const handleAnswer = async (stageIndex: number, option: Option) => {
  const stage = stages.value[stageIndex];

  // Reset all stages after the current one if the answer changes
  if (stage.answer !== option.label) {
    stages.value = stages.value.slice(0, stageIndex + 1);
  }

  stage.completed = true;
  stage.answer = option.label;

  // After timer completion, ask if the task was completed
  if (stage.question === "Did you complete the task?") {
    if (option.label === "Yes") {
      // User completed the task, remove it and proceed to next item
      try {
        await axios.delete(`/api/inbox/${currentItem.value.id}`);
        stages.value = initializeStages();
        currentIndex.value++;
      } catch (error) {
        console.error('Error deleting inbox item:', error);
        alert('Failed to delete inbox item. Please try again.');
      }
    } else {
      // User did not complete the task, ask to create a next action
      stages.value.push({
        question: "Create Action",
        component: CreateActionForm,
        props: {
          inboxItem: currentItem.value,
          onActionCreated: () => handleActionCreated(),
        },
        completed: false,
        options: [],
      });
    }
    return;
  }

  if (option.label === "Trash it") {
    try {
      // Send delete request to the backend server
      await fetch(`/api/inbox/${currentItem.value.id}`, {
        method: 'DELETE',
      });

      // Reset wizard stages and move to the next item
      stages.value = initializeStages();

      currentIndex.value++;
    } catch (error) {
      console.error("Failed to delete item:", error);
    }
    return;
  }

  if (option.nextStage) {
    const nextStage = option.nextStage(option.label);
    stages.value.push(nextStage);
  }
};

const handleActionCreated = async () => {
  if (currentItem.value?.id) {
    try {
      await axios.delete(`/api/inbox/${currentItem.value.id}`);

      // Reset wizard stages and move to the next item
      stages.value = initializeStages();

      currentIndex.value++;

    } catch (error) {
      console.error('Error deleting inbox item:', error);
      alert('Failed to delete inbox item. Please try again.');
    }
  }
};

const goToInbox = () => {
  router.push('/inbox');
};

const goToNextActions = () => {
  router.push('/next-actions');
};

const handleTimerComplete = () => {
  // After timer finishes, ask user if the task was completed
  stages.value.push({
    question: "Did you complete the task?",
    options: [
      { label: "Yes" },
      { label: "No" },
    ],
    completed: false,
  });
};

const handleMoveToNextItem = async () => {
  if (currentItem.value?.id) {
    try {
      await axios.delete(`/api/inbox/${currentItem.value.id}`);
      stages.value = initializeStages();
      currentIndex.value++;
    } catch (error) {
      console.error('Error deleting inbox item:', error);
      alert('Failed to delete inbox item. Please try again.');
    }
  }
};
</script>

<style scoped>
.stage {
  margin-bottom: 1rem;
  padding: 1rem;
  border: 1px solid var(--bs-border-color);
  border-radius: 0.5rem;
  background-color: var(--bs-body-bg);
}

.stage-content {
  margin-bottom: 0.5rem;
}

button {
  margin-right: 0.5rem;
  padding: 0.5rem 1rem;
  border: 2px solid transparent;
  border-radius: 0.25rem;
  background-color: var(--bs-primary);
  color: var(--bs-white);
  cursor: pointer;
  transition: all 0.3s ease;
}

button.selected {
  background-color: var(--bs-primary);
  border-color: var(--bs-primary);
}

button:not(.selected) {
  background-color: var(--bs-body-bg);
  color: var(--bs-primary);
  border-color: var(--bs-primary);
}

button:hover {
  background-color: var(--bs-primary);
  color: var(--bs-white);
  border-color: var(--bs-primary);
}

.inbox-item-details {
  margin-bottom: 1rem;
  padding: 0.5rem;
  border: 1px solid var(--bs-border-color);
  border-radius: 0.25rem;
  background-color: var(--bs-body-tertiary-bg);
}

.inbox-item-details p {
  margin: 0.5rem 0;
}

.inbox-item-details a {
  color: var(--bs-link-color);
  text-decoration: none;
}

.inbox-item-details a:hover {
  text-decoration: underline;
}

.page-nav {
  font-size: 2rem; /* Make navigation bigger */
  cursor: pointer; /* Change cursor to pointer to indicate clickability */
  margin: 0 1rem; /* Add spacing around the chevrons */
  color: #007bff; /* Add a color to make it more visually appealing */
  transition: color 0.3s ease; /* Smooth transition for hover effect */
}

.page-nav:hover {
  color: #0056b3; /* Darker color on hover */
}

.page-nav.disabled {
  color: #ccc; /* Gray out disabled chevrons */
  cursor: not-allowed; /* Indicate that the chevron is not clickable */
}

.paging {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 1rem 0; /* Add spacing around the navigation */
}

.empty-state {
  padding: 2rem;
  max-width: 600px;
  margin: 0 auto;
}

.empty-state .alert {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.empty-state .btn {
  min-width: 160px;
  transition: transform 0.2s ease;
}

.empty-state .btn:hover {
  transform: translateY(-1px);
}
</style>