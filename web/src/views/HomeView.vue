<script setup lang="ts">
import { ref } from "vue";
import EditorForm from "../components/EditorForm.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const url = ref<string | null>(null),
  err = ref<string | null>(null),
  isLoading = ref(false);

async function createPaste(p: any) {
  err.value = null;
  url.value = null;
  isLoading.value = true;
  try {
    const body = new URLSearchParams({
      content: p.content,
      ...(p.syntax ? { syntax: p.syntax } : {}),
      ...(p.ttl ? { ttl: p.ttl } : {}),
      ...(p.allow_edit ? { allow_edit: "on" } : {}),
      ...(p.password ? { password: p.password } : {}),
    });
    const res = await fetch("/api/p", { method: "POST", body });
    if (!res.ok) {
      err.value = "Failed to create paste";
      return;
    }
    const j = await res.json();
    url.value = `/p/${j.id}`;
    // Direkt zum Paste navigieren
    router.push(url.value);
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <section class="home-section">
    <h1>New Paste</h1>
    <p class="subtitle">Share code and text instantly</p>

    <EditorForm
      mode="create"
      :show-syntax="true"
      :show-ttl="true"
      :show-allow-edit="true"
      :show-password="true"
      :max-len="524288"
      submit-label="Create"
      :is-loading="isLoading"
      @submit="createPaste"
    />

    <transition name="fade">
      <div v-if="err" class="error-box">
        {{ err }}
      </div>
    </transition>
  </section>
</template>

<style scoped>
.home-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  width: 100%;
  max-width: 100%;
}

.home-section h1 {
  margin: 0 0 0.25rem 0;
  font-size: 1.5rem;
}

.subtitle {
  color: var(--fg-secondary);
  font-size: 0.9375rem;
  margin: 0;
}

.error-box {
  background: transparent;
  border: 1px solid var(--error);
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
