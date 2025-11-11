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

    // Save to My Shares
    const ttlMs = getTtlMs(p.ttl);
    const share = {
      id: j.id,
      syntax: p.syntax || "auto",
      createdAt: Date.now(),
      expiresAt: ttlMs ? Date.now() + ttlMs : null,
      allowEdit: p.allow_edit ?? false,
      password: !!p.password,
    };
    const stored = localStorage.getItem("binly-my-shares");
    const shares = stored ? JSON.parse(stored) : [];
    shares.push(share);
    localStorage.setItem("binly-my-shares", JSON.stringify(shares));

    // Direkt zum Paste navigieren
    router.push(url.value);
  } finally {
    isLoading.value = false;
  }
}

function getTtlMs(ttl: string): number | null {
  if (!ttl || ttl === "0") return null;
  const ttlMap: Record<string, number> = {
    "15m": 15 * 60 * 1000,
    "1h": 60 * 60 * 1000,
    "1d": 24 * 60 * 60 * 1000,
    "7d": 7 * 24 * 60 * 60 * 1000,
    "2h": 2 * 60 * 60 * 1000,
  };
  return ttlMap[ttl] || null;
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
