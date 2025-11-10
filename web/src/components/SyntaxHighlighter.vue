<script setup lang="ts">
import { computed } from "vue";
import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark.css";

const props = defineProps<{
  code: string;
  language?: string;
  enabled?: boolean;
}>();

const highlightedCode = computed(() => {
  if (!props.code) return "";

  // If highlighting is disabled, return raw code
  if (props.enabled === false) {
    return props.code;
  }

  try {
    if (props.language && props.language !== "auto") {
      return hljs.highlight(props.code, {
        language: props.language,
        ignoreIllegals: true,
      }).value;
    }

    // Auto-detect language
    const result = hljs.highlightAuto(props.code);
    return result.value;
  } catch (error) {
    return props.code;
  }
});
</script>

<template>
  <pre class="syntax-highlighter"><code
    class="hljs"
    v-html="highlightedCode"
  /></pre>
</template>

<style scoped>
.syntax-highlighter {
  margin: 0;
  overflow-x: auto;
}

.syntax-highlighter code {
  padding: var(--spacing-md);
  display: block;
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 0.875rem;
  line-height: 1.6;
  color: inherit;
}

/* Override highlight.js dark theme for our design */
:deep(.hljs) {
  background: transparent;
  color: var(--fg);
}

:deep(.hljs-attr),
:deep(.hljs-attribute) {
  color: #79c0ff;
}

:deep(.hljs-literal),
:deep(.hljs-number),
:deep(.hljs-doctag) {
  color: #79c0ff;
}

:deep(.hljs-string) {
  color: #a5d6ff;
}

:deep(.hljs-title),
:deep(.hljs-section) {
  color: #79c0ff;
  font-weight: 600;
}

:deep(.hljs-built_in),
:deep(.hljs-selector-tag),
:deep(.hljs-type),
:deep(.hljs-class) {
  color: #f0883e;
}

:deep(.hljs-keyword),
:deep(.hljs-selector-class) {
  color: #ff7b72;
}

:deep(.hljs-comment),
:deep(.hljs-deletion),
:deep(.hljs-quote) {
  color: #8b949e;
}

:deep(.hljs-emphasis) {
  color: #d2a8ff;
  font-style: italic;
}

:deep(.hljs-strong) {
  color: #d2a8ff;
  font-weight: bold;
}

:deep(.hljs-bullet),
:deep(.hljs-link),
:deep(.hljs-symbol) {
  color: #3fb950;
}

:deep(.hljs-code) {
  background: var(--bg-secondary);
  padding: 0.125rem 0.375rem;
  border-radius: var(--radius-sm);
}
</style>
