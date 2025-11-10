<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { Icon } from "@iconify/vue";
import EditorForm from "../components/EditorForm.vue";
import SyntaxHighlighter from "../components/SyntaxHighlighter.vue";
const route = useRoute();
const id = route.params.id as string;

const content = ref(""),
  allowEdit = ref(false),
  syntax = ref("auto");
const needPw = ref(false),
  pw = ref(""),
  err = ref<string | null>(null);
const editing = ref(false);
const isLoading = ref(false);
const copySuccess = ref(false);
const copyLinkSuccess = ref(false);
const expandContent = ref(false);
const syntaxHighlightingEnabled = ref(true);

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(content.value);
    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2000);
  } catch (err) {
    console.error("Failed to copy:", err);
  }
}

async function copyLinkToClipboard() {
  try {
    const link = `${window.location.origin}/p/${id}`;
    await navigator.clipboard.writeText(link);
    copyLinkSuccess.value = true;
    setTimeout(() => {
      copyLinkSuccess.value = false;
    }, 2000);
  } catch (err) {
    console.error("Failed to copy link:", err);
  }
}

async function load() {
  isLoading.value = true;
  try {
    const res = await fetch(`/api/p/${id}`);
    if (res.status === 401) {
      needPw.value = true;
      return;
    }
    if (!res.ok) {
      err.value = "Paste not found or has expired";
      return;
    }
    const j = await res.json();
    content.value = j.content;
    syntax.value = j.syntax || "auto";
    allowEdit.value = j.allow_edit;
  } finally {
    isLoading.value = false;
  }
}

async function auth() {
  isLoading.value = true;
  try {
    const body = new URLSearchParams({ password: pw.value });
    const res = await fetch(`/api/p/${id}/auth`, { method: "POST", body });
    if (res.ok) {
      needPw.value = false;
      await load();
    } else {
      err.value = "Invalid password";
    }
  } finally {
    isLoading.value = false;
  }
}

async function save(p: any) {
  isLoading.value = true;
  try {
    const body = new URLSearchParams({ content: p.content });
    const res = await fetch(`/api/p/${id}`, { method: "POST", body });
    if (res.ok || res.status === 204) {
      editing.value = false;
      await load();
    } else {
      err.value = "Save failed";
    }
  } finally {
    isLoading.value = false;
  }
}

onMounted(load);
</script>

<template>
  <section class="paste-section">
    <div v-if="needPw" class="auth-form">
      <div class="auth-header">
        <h2>Password Required</h2>
        <p class="subtitle">This paste is protected</p>
      </div>
      <input
        type="password"
        v-model="pw"
        placeholder="Enter password"
        @keyup.enter="auth"
      />
      <button @click="auth" :disabled="isLoading">
        {{ isLoading ? "Unlocking..." : "Unlock" }}
      </button>
      <transition name="fade">
        <div v-if="err" class="error-box">
          {{ err }}
        </div>
      </transition>
    </div>

    <div v-else>
      <div class="paste-header">
        <div class="paste-id-section">
          <h2 class="paste-id">
            /p/<span class="id-value">{{ id }}</span>
          </h2>
          <button
            @click="copyLinkToClipboard"
            class="copy-link-btn"
            :title="copyLinkSuccess ? 'Link copied!' : 'Copy link'"
          >
            <Icon
              :icon="copyLinkSuccess ? 'tabler:check' : 'tabler:link'"
              width="14"
              height="14"
            />
            <p style="margin-left:0.2rem">{{ copyLinkSuccess ? "Link copied" : "Copy link" }}</p>
          </button>
        </div>
        <div class="header-actions">
          <button
            @click="expandContent = !expandContent"
            class="action-btn"
            :title="expandContent ? 'Collapse' : 'Expand'"
          >
            <Icon
              :icon="
                expandContent ? 'tabler:chevron-up' : 'tabler:chevron-down'
              "
              width="16"
              height="16"
            />
          </button>
          <label class="syntax-toggle">
            <input
              type="checkbox"
              v-model="syntaxHighlightingEnabled"
            />
            <span :title="syntaxHighlightingEnabled ? 'Syntax highlighting enabled' : 'Syntax highlighting disabled'">
              <Icon
                :icon="syntaxHighlightingEnabled ? 'tabler:highlight' : 'tabler:highlight-off'"
                width="16"
                height="16"
              />
            </span>
          </label>
          <button
            @click="copyToClipboard"
            class="action-btn"
            :title="copySuccess ? 'Copied!' : 'Copy to clipboard'"
          >
            <Icon
              :icon="copySuccess ? 'tabler:check' : 'tabler:copy'"
              width="16"
              height="16"
            />
            <span>{{ copySuccess ? "Copied" : "Copy" }}</span>
          </button>
          <a :href="`/api/raw/${id}`" target="_blank" class="raw-link">
            <Icon icon="tabler:download" width="16" height="16" />
            <span>Raw</span>
          </a>
        </div>
      </div>

      <div
        v-if="!editing"
        class="content-display"
        :class="{ expanded: expandContent }"
      >
        <SyntaxHighlighter 
          :code="content" 
          :language="syntax" 
          :enabled="syntaxHighlightingEnabled" 
        />
      </div>

      <div v-if="!editing && allowEdit" class="edit-section">
        <button @click="editing = true">
          <Icon icon="tabler:pencil" width="18" height="18" />
          <span>Edit</span>
        </button>
      </div>

      <div v-if="editing" class="edit-form">
        <EditorForm
          mode="edit"
          :initial-content="content"
          :show-syntax="false"
          :show-ttl="false"
          :show-allow-edit="false"
          :show-password="false"
          :max-len="524288"
          submit-label="Save"
          :is-loading="isLoading"
          @submit="save"
        />
        <button
          class="cancel-btn"
          @click="editing = false"
          :disabled="isLoading"
        >
          Cancel
        </button>
      </div>

      <transition name="fade">
        <div v-if="err" class="error-box">
          {{ err }}
        </div>
      </transition>
    </div>
  </section>
</template>

<style scoped>
.paste-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
}

.auth-header h2 {
  margin-bottom: var(--spacing-sm);
}

.subtitle {
  color: var(--fg-secondary);
  font-size: 0.9375rem;
  margin: 0;
}

.paste-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
}

.paste-id-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.copy-link-btn {
  flex-shrink: 0;
  padding: 0.4rem 0.6rem;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--fg);
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-family: inherit;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

.copy-link-btn:hover {
  border-color: transparent;
  color: var(--accent);
}

.copy-link-btn:active {
  transform: scale(0.95);
}

.copy-link-btn:focus {
  outline: none;
}

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

.action-btn,
.raw-link {
  white-space: nowrap;
  padding: 0.5rem 1rem;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-family: inherit;
  color: var(--fg);
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

.action-btn:hover:not(:disabled),
.raw-link:hover {
  border-color: transparent;
  background: transparent;
  color: var(--accent);
}

.action-btn:active {
  transform: scale(0.98);
}

.action-btn:focus {
  outline: none;
}

.syntax-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  transition: var(--transition);
}

.syntax-toggle:hover {
  background: transparent;
  color: var(--accent);
}

.syntax-toggle input {
  display: none;
}

.syntax-toggle span {
  display: flex;
  align-items: center;
  justify-content: center;
}

.paste-id {
  margin: 0;
  word-break: break-all;
}

.id-value {
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  color: var(--accent);
}

.content-display {
  margin-bottom: var(--spacing-lg);
  user-select: text;
  max-height: 600px;
  overflow-y: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  transition: max-height 0.12s ease-out;
}

.content-display.expanded {
  max-height: calc(100vh - 280px);
}

.content-display :deep(pre) {
  margin: 0;
  user-select: text;
}

.content-display :deep(code) {
  user-select: text;
}

.content-display pre {
  margin: 0;
}

.edit-section {
  display: flex;
  gap: var(--spacing-md);
}

.edit-section button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.cancel-btn {
  background: transparent;
  border: none;
  color: var(--fg);
  padding: 0.5rem 1rem;
}

.cancel-btn:hover {
  border-color: transparent;
  color: var(--accent);
}

.error-box {
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  color: var(--error);
  font-size: 0.9375rem;
  animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.12s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
