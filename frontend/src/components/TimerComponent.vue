<template>
  <div class="timer">
    <div class="timer-display" :class="{ 'time-low': timeLeft <= 30 }">
      {{ formatTime(timeLeft) }}
    </div>
    <div class="timer-controls">
      <button v-if="!isRunning" @click="startTimer" class="btn btn-primary">
        Start Timer
      </button>
      <button v-else @click="stopTimer" class="btn btn-danger">
        Stop Timer
      </button>
      <button @click="resetTimer" class="btn btn-secondary ms-2" :disabled="isRunning">
        Reset
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue';

const props = defineProps({
  duration: {
    type: Number,
    default: 10 // 2 minutes in seconds
  }
});

const emit = defineEmits(['timerComplete']);

const timeLeft = ref(props.duration);
const isRunning = ref(false);
let timerInterval: number | null = null;

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function startTimer() {
  if (!isRunning.value) {
    isRunning.value = true;
    timerInterval = window.setInterval(() => {
      if (timeLeft.value > 0) {
        timeLeft.value--;
      } else {
        stopTimer();
        showNotification();
        emit('timerComplete');
      }
    }, 1000);
  }
}

function stopTimer() {
  if (timerInterval) {
    clearInterval(timerInterval);
    timerInterval = null;
  }
  isRunning.value = false;
}

function resetTimer() {
  stopTimer();
  timeLeft.value = props.duration;
}

function showNotification() {
  // Request permission and show notification
  if ('Notification' in window) {
    Notification.requestPermission().then(permission => {
      if (permission === 'granted') {
        new Notification('Timer Complete!', {
          body: '2 minutes have passed. Time to move on!',
          icon: '/list-check.svg'
        });
      }
    });
  }
}

// Clean up on component unmount
onUnmounted(() => {
  if (timerInterval) {
    clearInterval(timerInterval);
  }
});
</script>

<style scoped>
.timer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  background-color: var(--bs-dark);
}

.timer-display {
  font-size: 2rem;
  font-weight: bold;
  font-family: monospace;
  color: var(--bs-dark-text);
}

.time-low {
  color: var(--bs-danger);
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

.timer-controls {
  display: flex;
  gap: 0.5rem;
}
</style>