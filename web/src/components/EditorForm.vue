<script setup lang="ts">
import { computed, ref, watchEffect, nextTick } from "vue";
import { Icon } from "@iconify/vue";
import CustomSelect from "./CustomSelect.vue";

type Ttl = "0" | "15m" | "1h" | "1d" | "7d";

const props = defineProps<{
  mode?: "create" | "edit";
  initialContent?: string;
  initialSyntax?: string;
  initialTtl?: Ttl;
  initialAllowEdit?: boolean;
  showSyntax?: boolean;
  showTtl?: boolean;
  showAllowEdit?: boolean;
  showPassword?: boolean;
  submitLabel?: string;
  maxLen?: number;
  isLoading?: boolean;
}>();

const emit = defineEmits<{
  (
    e: "submit",
    payload: {
      content: string;
      syntax?: string;
      ttl?: Ttl;
      allow_edit?: boolean;
      password?: string;
    }
  ): void;
}>();

const content = ref(props.initialContent ?? "");
const syntax = ref(props.initialSyntax ?? "auto");
const ttl = ref<Ttl>(props.initialTtl ?? "1d");
const allowEdit = ref(props.initialAllowEdit ?? false);
const password = ref("");
const error = ref<string | null>(null);

const remain = computed(() =>
  props.maxLen ? props.maxLen - content.value.length : undefined
);

const remainPercentage = computed(() => {
  if (!remain.value || !props.maxLen) return 100;
  return Math.max(0, (remain.value / props.maxLen) * 100);
});

watchEffect(() => {
  error.value = null;
});

function onSubmit() {
  if (!content.value.trim()) {
    error.value = "Content is required";
    return;
  }
  if (props.maxLen && content.value.length > props.maxLen) {
    error.value = `Content exceeds limit (>${props.maxLen} characters)`;
    return;
  }
  emit("submit", {
    content: content.value,
    ...(props.showSyntax ? { syntax: syntax.value } : {}),
    ...(props.showTtl ? { ttl: ttl.value } : {}),
    ...(props.showAllowEdit ? { allow_edit: allowEdit.value } : {}),
    ...(props.showPassword && password.value
      ? { password: password.value }
      : {}),
  });
}

function handleEditorTab(e: KeyboardEvent) {
  if (e.key === "Tab") {
    e.preventDefault();
    const textarea = e.target as HTMLTextAreaElement;
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;

    // Insert 2 spaces for indentation
    const indent = "  ";
    content.value =
      content.value.substring(0, start) + indent + content.value.substring(end);

    // Move cursor after the inserted indent
    nextTick(() => {
      textarea.selectionStart = textarea.selectionEnd = start + indent.length;
    });
  }
}
</script>

<template>
  <div class="editor-form">
    <div class="editor-wrapper">
      <textarea
        v-model="content"
        class="editor"
        :placeholder="
          mode === 'edit' ? 'Update content…' : 'Paste code or text here…'
        "
        :disabled="isLoading"
        @keydown="handleEditorTab"
      />
      <div v-if="remain !== undefined" class="capacity-bar">
        <div
          class="capacity-fill"
          :style="{ width: remainPercentage + '%' }"
          :class="{
            'capacity-full': remainPercentage < 10,
            'capacity-warning': remainPercentage < 25,
          }"
        />
      </div>
    </div>

    <div class="controls">
      <div class="controls-row">
        <CustomSelect
          v-if="showSyntax"
          v-model="syntax"
          :options="[
            { value: 'auto', label: 'auto' },
            { value: 'js', label: 'js' },
            { value: 'ts', label: 'ts' },
            { value: 'py', label: 'py' },
            { value: 'go', label: 'go' },
            { value: 'rust', label: 'rust' },
            { value: 'json', label: 'json' },
            { value: 'xml', label: 'xml' },
            { value: 'sql', label: 'sql' },
          ]"
          :disabled="isLoading"
        />

        <CustomSelect
          v-if="showTtl"
          v-model="ttl"
          :options="[
            { value: '15m', label: '15 minutes' },
            { value: '1h', label: '1 hour' },
            { value: '1d', label: '1 day' },
            { value: '7d', label: '7 days' },
            { value: '0', label: 'Never' },
          ]"
          :disabled="isLoading"
        />

        <label v-if="showAllowEdit" class="control-checkbox">
          <input type="checkbox" v-model="allowEdit" :disabled="isLoading" />
          <span>Allow editing</span>
        </label>

        <input
          v-if="showPassword"
          type="password"
          v-model="password"
          placeholder="Optional password"
          class="control-input"
          :disabled="isLoading"
        />
      </div>

      <div class="controls-row controls-bottom">
        <span v-if="remain !== undefined" class="character-count">
          {{ remain.toLocaleString() }} characters remaining
        </span>

        <button
          @click="onSubmit"
          class="submit-button"
          :disabled="isLoading || !content.trim()"
        >
          <Icon
            :icon="mode === 'edit' ? 'tabler:check' : 'tabler:plus'"
            width="18"
            height="18"
          />
          <span>{{
            isLoading
              ? "Processing..."
              : submitLabel ?? (mode === "edit" ? "Save" : "Create")
          }}</span>
        </button>
      </div>
    </div>

    <transition name="fade">
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </transition>
  </div>
</template>

<style scoped>
.editor-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.editor-wrapper {
  position: relative;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.editor {
  width: 100%;
  min-height: 300px;
  padding: var(--spacing-md);
  background: var(--bg);
  color: var(--fg);
  border: none;
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 0.9375rem;
  line-height: 1.6;
  letter-spacing: 0;
  resize: vertical;
  transition: var(--transition);
}

.editor:focus {
  outline: none;
}

.capacity-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: rgba(37, 99, 235, 0.1);
}

.capacity-fill {
  height: 100%;
  background: var(--accent);
  transition: all 0.3s ease-out;
}

.capacity-fill.capacity-warning {
  background: #f59e0b;
}

.capacity-fill.capacity-full {
  background: var(--error);
}

.controls {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.controls-row {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
  align-items: center;
}

.controls-bottom {
  justify-content: space-between;
  align-items: center;
}

.control-input {
  flex: 1;
  min-width: 120px;
  color: var(--fg);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0.5rem 0.75rem;
  font-family: inherit;
  font-size: 0.875rem;
  transition: var(--transition);
}

.control-input:hover:not(:disabled) {
  border-color: transparent;
  opacity: 0.9;
}

.control-input:focus {
  outline: none;
  border-color: transparent;
  opacity: 0.9;
}

.control-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.control-checkbox {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  font-size: 0.9375rem;
  user-select: none;
}

.control-checkbox input {
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  cursor: pointer;
  width: 16px;
  height: 16px;
  min-width: 16px;
  min-height: 16px;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: var(--bg);
  transition: var(--transition);
  flex-shrink: 0;
  position: relative;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
}

.control-checkbox input:checked {
  background: var(--accent);
  border-color: var(--accent);
}

.control-checkbox input:checked::after {
  content: "✓";
  position: absolute;
  color: white;
  font-size: 11px;
  font-weight: bold;
  line-height: 1;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.control-checkbox input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.character-count {
  font-size: 0.8125rem;
  color: var(--muted);
}

.submit-button {
  padding: 0.625rem 1.5rem;
  background: var(--accent);
  border: 1px solid var(--accent);
  border-radius: var(--radius-md);
  color: white;
  font-weight: 600;
  font-size: 0.9375rem;
  cursor: pointer;
  transition: background 0.15s ease-out, border 0.15s ease-out;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.submit-button:hover:not(:disabled) {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  background: transparent;
  border: 1px solid var(--error);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  color: var(--error);
  font-size: 0.9375rem;
  animation: slideIn 0.2s ease-out;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.12s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .editor {
    min-height: 250px;
    font-size: 16px; /* Prevents iOS zoom */
  }

  .controls-row {
    flex-direction: column;
  }

  .control-select,
  .control-input {
    width: 100%;
    min-width: unset;
  }

  .controls-bottom {
    flex-direction: column-reverse;
  }

  .character-count {
    text-align: center;
    width: 100%;
  }

  .submit-button {
    width: 100%;
  }
}
</style>
