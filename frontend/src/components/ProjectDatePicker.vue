<script setup lang="ts">
import { DatePicker } from 'v-calendar';
import 'v-calendar/style.css';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faCalendar } from '@fortawesome/free-regular-svg-icons';
import { ref } from 'vue';

const model = defineModel<string | null>();
const emit = defineEmits<{
  'clear': []
}>();

const datePickerRef = ref();

const formatDateForDisplay = (dateStr?: string) => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString();
};
</script>

<template>
  <DatePicker
    ref="datePickerRef"
    v-model.string="model"
    :min-date="new Date()"
    :first-day-of-week="1"
    title-position="left"
    :masks="{
      modelValue: 'YYYY-MM-DD',
    }"
  >
    <template #default="{ togglePopover }">
      <button class="btn btn-link" @click="togglePopover">
        <FontAwesomeIcon :icon="faCalendar" />
        {{ model ? formatDateForDisplay(model) : 'Add deadline' }}
      </button>
    </template>
    <template #footer>
      <button
        type="button"
        class="btn btn-primary w-75 mb-2 mx-auto d-block"
        :class="model ? 'btn-primary' : 'btn-primary disabled'"
        :disabled="!model"
        @click="() => {
          model = '';
          emit('clear');
          datePickerRef.hidePopover();
        }">
        Clear
      </button>
    </template>
  </DatePicker>
</template>