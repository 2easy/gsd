import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useThemeStore = defineStore('theme', () => {
  const theme = ref(localStorage.getItem('theme') || 'light');

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light';
    localStorage.setItem('theme', theme.value);
    document.documentElement.setAttribute('data-bs-theme', theme.value);
  };

  const initTheme = () => {
    document.documentElement.setAttribute('data-bs-theme', theme.value);
  };

  return {
    theme,
    toggleTheme,
    initTheme,
  };
}); 