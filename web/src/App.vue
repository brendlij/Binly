<template>
  <div id="app" class="app" :class="{ 'light-mode': !isDark }">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="title">
          <BinlyLogo
            :size="64"
            :show-circle="false"
            :mode="isDark ? 'dark' : 'light'"
          />
          <span class="title-text">inly</span>
        </router-link>
        <nav class="navbar">
          <router-link to="/my-shares" class="nav-link">My Shares</router-link>
          <a href="/docs" class="nav-link">Docs</a>
          <a
            href="https://github.com/brendlij/Binly"
            target="_blank"
            rel="noopener noreferrer"
            class="nav-link"
          >
            <Icon icon="tabler:brand-github" width="18" height="18" />
          </a>
        </nav>
        <button
          class="theme-toggle"
          @click="toggleTheme"
          :title="`Switch to ${isDark ? 'light' : 'dark'} mode`"
        >
          <Icon
            :icon="isDark ? 'mdi:white-balance-sunny' : 'mdi:moon-new'"
            width="20"
            height="20"
          />
        </button>
      </div>
    </header>

    <main class="main">
      <div class="router-view-wrapper">
        <router-view />
      </div>
    </main>

    <footer class="footer">
      <a
        href="https://github.com/brendlij/Binly"
        target="_blank"
        rel="noopener noreferrer"
        class="footer-link"
      >
        <Icon icon="tabler:brand-github" width="18" height="18" />
        <span>Binly</span>
      </a>
      <span class="footer-version">v1</span>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";
import BinlyLogo from "./components/BinlyLogo.vue";

const isDark = ref(true);

onMounted(() => {
  const saved = localStorage.getItem("binly-theme");
  if (saved) {
    isDark.value = saved === "dark";
  } else {
    isDark.value = window.matchMedia("(prefers-color-scheme: dark)").matches;
  }
  applyTheme();
});

function toggleTheme() {
  isDark.value = !isDark.value;
  localStorage.setItem("binly-theme", isDark.value ? "dark" : "light");
  applyTheme();
}

function applyTheme() {
  const root = document.documentElement;
  if (isDark.value) {
    root.style.colorScheme = "dark";
  } else {
    root.style.colorScheme = "light";
  }
}
</script>

<style>
:root {
  /* Dark mode (default) */
  --bg: #0b0b0c;
  --bg-secondary: #141416;
  --fg: #f5f5f6;
  --fg-secondary: #b9b9ba;
  --border: #2a2a2b;
  --accent: #6366f1; /* matches logo fold */
  --accent-hover: #4f46e5; /* slightly deeper violet */
  --muted: #707070;
  --error: #ef4444;
  --success: #22c55e;
  --scrollbar: var(--border); /* matches logo fold */

  /* Spacing & sizing */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;

  /* Border radius */
  --radius-sm: 0.25rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;

  /* Transitions */
  --transition: all 0.2s ease-out;

  color-scheme: dark;
}

.light-mode {
  --bg: #fafafa;
  --bg-secondary: #f2f2f3;
  --fg: #0c0c0d;
  --fg-secondary: #4c4c4d;
  --border: #e6e6e7;
  --accent: #6366f1; /* same accent */
  --accent-hover: #4f46e5;
  --scrollbar: var(--border); /* same as accent */
  --muted: #888;
  --error: #ef4444;
  --success: #22c55e;

  color-scheme: light;
}

html,
body,
#app {
  margin: 0;
  padding: 0;
  height: 100%;
  background: var(--bg);
  color: var(--fg);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto Mono",
    "JetBrains Mono", monospace;
  font-weight: 400;
  line-height: 1.5;
  letter-spacing: -0.01em;
  transition: var(--transition);
  -webkit-text-size-adjust: 100%;
  /* Custom scrollbar */
  scrollbar-color: var(--scrollbar) transparent;
  scrollbar-width: thin;
}

* {
  box-sizing: border-box;
}

/* Custom scrollbar for webkit browsers */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
  margin-top: 10px;
  margin-bottom: 10px;
}

::-webkit-scrollbar-thumb {
  background: var(--scrollbar);
  border-radius: 4px;
  margin-right: 2px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--accent);
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  padding: var(--spacing-md) var(--spacing-lg);
  background: var(--bg);
  border-bottom: none;
  position: sticky;
  top: 0;
  z-index: 100;
  transition: var(--transition);
  border-bottom: 1px solid var(--border);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  gap: var(--spacing-lg);
}

.navbar {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin: 0 auto;
}

.nav-link {
  color: var(--fg-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: var(--transition);
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  border-radius: var(--radius-md);
}

.nav-link:hover {
  color: var(--accent);
  background: transparent;
}

.title {
  text-decoration: none;
  font-weight: 600;
  font-size: 1.125rem;
  letter-spacing: -0.02em;
  color: var(--fg);
  transition: var(--transition);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-bottom: none !important;
  align-items: flex-start;
}

.title-text {
  margin-top: +1.8rem;
  margin-left: -1.5rem;
}

.title:hover {
  color: var(--accent);
  border-bottom: none;
}

.title svg {
  transition: transform 0.3s ease-out;
  flex-shrink: 0;
}

.title:hover svg {
  transform: none;
}

.title:hover .title-text {
  transform: none;
  transition: transform 0.3s ease-out;
}

.theme-toggle {
  background: transparent;
  border: none;
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-family: inherit;
  color: var(--fg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  transition: var(--transition);
  min-width: 36px;
  height: 36px;
}

.theme-toggle:hover {
  background: transparent;
  border-color: transparent;
  color: var(--accent);
}

.theme-toggle:active {
  transform: scale(0.95);
}

.main {
  flex: 1;
  padding: var(--spacing-lg);
  max-width: 1225px;
  margin: 0 auto;
  width: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  min-height: 0; /* Important for flex children to respect height */
}

/* Ensure router-view fills available space */
:deep(.router-view-wrapper) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.footer {
  text-align: center;
  font-size: 0.8125rem;
  color: var(--muted);
  padding: var(--spacing-lg);
  border-top: none;
  background: var(--bg);
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  border-top: 1px solid var(--border);
}

.footer-link {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  color: var(--muted);
  text-decoration: none;
  transition: var(--transition);
  font-weight: 500;
}

.footer-link:hover {
  color: var(--accent);
}

.footer-version {
  color: var(--muted);
  opacity: 0.65;
  font-size: 0.75rem;
}

/* Global element styles */
a {
  color: var(--accent);
  text-decoration: none;
  transition: var(--transition);
  border-bottom: none;
}

a:hover {
  border-bottom: none;
  opacity: 0.8;
}

button {
  background: var(--accent);
  border: 1px solid var(--accent);
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-family: inherit;
  color: white;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: var(--transition);
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

button:hover {
  background: var(--accent-hover);
}

button:active {
  transform: scale(0.98);
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

button:focus {
  outline: none;
}

textarea,
select,
input {
  background: var(--bg);
  color: var(--fg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0.625rem 0.75rem;
  font-family: inherit;
  font-size: 1rem;
  transition: var(--transition);
  -webkit-user-select: text;
  user-select: text;
}

textarea::placeholder,
input::placeholder {
  color: var(--muted);
}

textarea:focus,
select:focus,
input:focus {
  outline: none;
  border-color: transparent;
  opacity: 0.9;
}

textarea {
  resize: vertical;
  min-height: 200px;
  line-height: 1.6;
}

h1,
h2,
h3,
h4,
h5,
h6 {
  margin: 0;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--fg);
}

h1 {
  font-size: 1.875rem;
  margin-bottom: var(--spacing-lg);
}

h2 {
  font-size: 1.5rem;
  margin-bottom: var(--spacing-md);
}

code {
  background: var(--bg);
  padding: 0.125rem 0.375rem;
  border-radius: var(--radius-sm);
  font-size: 0.875em;
  color: var(--accent);
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
}

pre {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  overflow-x: auto;
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 0.875rem;
  line-height: 1.6;
}

pre code {
  background: transparent;
  padding: 0;
  color: var(--fg);
}

@media (max-width: 640px) {
  .main {
    padding: var(--spacing-md);
  }

  .header {
    padding: var(--spacing-md);
  }

  .footer {
    padding: var(--spacing-md);
  }

  h1 {
    font-size: 1.5rem;
  }
}
</style>
