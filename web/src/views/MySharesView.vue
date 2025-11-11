<template>
  <section class="my-shares-section">
    <h1>My Shares</h1>
    <p class="subtitle">Your recently shared pastes</p>

    <div v-if="shares.length === 0" class="empty-state">
      <p>
        No shares yet. <router-link to="/">Create your first paste</router-link>
      </p>
    </div>

    <div v-else class="shares-grid">
      <div v-for="share in shares" :key="share.id" class="share-card">
        <div class="share-header">
          <h3 class="share-id">
            <router-link :to="`/p/${share.id}`"
              >/p/<span class="id-value">{{ share.id }}</span></router-link
            >
          </h3>
          <span
            v-if="share.expiresAt"
            class="share-ttl"
            :class="getTtlClass(share.expiresAt)"
          >
            {{ formatTimeRemaining(share.expiresAt) }}
          </span>
          <span v-else class="share-ttl permanent">Never</span>
        </div>

        <div class="share-meta">
          <p class="share-syntax">
            <Icon icon="tabler:code" width="14" height="14" />
            {{ share.syntax || "auto" }}
          </p>
          <p class="share-created">
            <Icon icon="tabler:clock" width="14" height="14" />
            {{ formatDate(share.createdAt) }}
          </p>
          <p v-if="share.allowEdit" class="share-edit">
            <Icon icon="tabler:pencil" width="14" height="14" />
            Editable
          </p>
          <p v-if="share.password" class="share-password">
            <Icon icon="tabler:lock" width="14" height="14" />
            Protected
          </p>
        </div>

        <div class="share-actions">
          <button @click="copyLink(share.id)" class="action-btn">
            <Icon
              :icon="copiedId === share.id ? 'tabler:check' : 'tabler:copy'"
              width="14"
              height="14"
            />
            {{ copiedId === share.id ? "Copied" : "Copy" }}
          </button>
          <router-link :to="`/p/${share.id}`" class="action-btn view-link">
            <Icon icon="tabler:external-link" width="14" height="14" />
            View
          </router-link>
          <button @click="removeShare(share.id)" class="action-btn delete-btn">
            <Icon icon="tabler:trash" width="14" height="14" />
            Remove
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";

interface Share {
  id: string;
  syntax: string;
  createdAt: number;
  expiresAt: number | null;
  allowEdit: boolean;
  password: boolean;
}

const shares = ref<Share[]>([]);
const copiedId = ref<string | null>(null);
const tick = ref(0); // Force re-render for live countdown

// Load shares from localStorage
function loadShares() {
  const stored = localStorage.getItem("binly-my-shares");
  if (!stored) return;

  try {
    const allShares: Share[] = JSON.parse(stored);
    // Filter out expired shares
    shares.value = allShares.filter((s) => !isExpired(s.expiresAt));
    saveShares();
  } catch (error) {
    console.error("Failed to load shares:", error);
  }
}

function saveShares() {
  localStorage.setItem("binly-my-shares", JSON.stringify(shares.value));
}

function isExpired(expiresAt: number | null): boolean {
  if (!expiresAt) return false;
  return expiresAt < Date.now();
}

function formatTimeRemaining(expiresAt: number): string {
  const now = Date.now();
  const diff = expiresAt - now;

  if (diff <= 0) return "Expired";

  const hours = Math.floor(diff / (1000 * 60 * 60));
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

  if (hours > 24) {
    const days = Math.floor(hours / 24);
    return `${days}d left`;
  }
  if (hours > 0) {
    return `${hours}h ${minutes}m left`;
  }
  return `${minutes}m left`;
}

function formatDate(timestamp: number): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  // Less than a minute
  if (diff < 60000) return "Just now";

  // Less than an hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000);
    return `${minutes}m ago`;
  }

  // Less than a day
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000);
    return `${hours}h ago`;
  }

  // Format as date
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function getTtlClass(expiresAt: number | null): string {
  if (!expiresAt) return "";
  const diff = expiresAt - Date.now();
  if (diff < 0) return "expired";
  if (diff < 15 * 60 * 1000) return "critical"; // Less than 15 minutes
  if (diff < 60 * 60 * 1000) return "warning"; // Less than 1 hour
  return "ok";
}

async function copyLink(id: string) {
  try {
    const link = `${window.location.origin}/p/${id}`;
    await navigator.clipboard.writeText(link);
    copiedId.value = id;
    setTimeout(() => {
      copiedId.value = null;
    }, 2000);
  } catch (error) {
    console.error("Failed to copy link:", error);
  }
}

function removeShare(id: string) {
  shares.value = shares.value.filter((s) => s.id !== id);
  saveShares();

  // Delete from server
  fetch(`/api/p/${id}`, { method: "DELETE" }).catch((error) => {
    console.error(`Failed to delete paste ${id} from server:`, error);
  });
}

// Cleanup expired shares periodically
onMounted(() => {
  loadShares();

  // Live countdown update every second
  const countdownInterval = setInterval(() => {
    tick.value++; // Trigger re-render
  }, 1000);

  // Cleanup expired shares every minute
  const cleanupInterval = setInterval(() => {
    const before = shares.value.length;
    shares.value = shares.value.filter((s) => !isExpired(s.expiresAt));
    if (shares.value.length !== before) {
      saveShares();
    }
  }, 60000); // Check every minute

  return () => {
    clearInterval(countdownInterval);
    clearInterval(cleanupInterval);
  };
});
</script>

<style scoped>
.my-shares-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  width: 100%;
}

.my-shares-section h1 {
  margin: 0 0 0.25rem 0;
  font-size: 1.5rem;
}

.subtitle {
  color: var(--fg-secondary);
  font-size: 0.9375rem;
  margin: 0;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--fg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: transparent;
}

.empty-state a {
  color: var(--accent);
  font-weight: 500;
}

.shares-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--spacing-lg);
}

.share-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: transparent;
  transition: var(--transition);
}

.share-card:hover {
  border-color: var(--accent);
  opacity: 0.95;
}

.share-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-sm);
}

.share-id {
  margin: 0;
  font-size: 0.9375rem;
  font-weight: 500;
  word-break: break-all;
  flex: 1;
}

.share-id a {
  color: var(--fg);
  text-decoration: none;
  transition: var(--transition);
}

.share-id a:hover {
  color: var(--accent);
}

.id-value {
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  color: var(--accent);
}

.share-ttl {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
  background: var(--bg-secondary);
  color: var(--fg-secondary);
}

.share-ttl.ok {
  background: transparent;
  color: var(--success);
  border: 1px solid var(--success);
}

.share-ttl.warning {
  background: transparent;
  color: #f59e0b;
  border: 1px solid #f59e0b;
}

.share-ttl.critical {
  background: transparent;
  color: var(--error);
  border: 1px solid var(--error);
}

.share-ttl.expired {
  background: transparent;
  color: var(--muted);
  border: 1px solid var(--muted);
  opacity: 0.5;
}

.share-ttl.permanent {
  background: transparent;
  color: var(--accent);
  border: 1px solid var(--accent);
}

.share-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  border-top: 1px solid var(--border);
  border-bottom: 1px solid var(--border);
}

.share-meta p {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  margin: 0;
  font-size: 0.8125rem;
  color: var(--fg-secondary);
}

.share-edit {
  color: var(--accent);
}

.share-password {
  color: var(--accent);
}

.share-created {
  white-space: nowrap;
}

.share-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.action-btn {
  flex: 1;
  padding: 0.5rem 0.75rem;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--fg-secondary);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  text-decoration: none;
  font-family: inherit;
}

.action-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: transparent;
}

.action-btn.delete-btn:hover {
  border-color: var(--error);
  color: var(--error);
}

.action-btn.view-link {
  color: var(--accent);
  border-color: var(--accent);
}

@media (max-width: 768px) {
  .shares-grid {
    grid-template-columns: 1fr;
  }

  .share-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .share-ttl {
    align-self: flex-start;
  }

  .share-actions {
    gap: var(--spacing-xs);
  }

  .action-btn {
    font-size: 0.75rem;
    padding: 0.375rem 0.5rem;
  }

  .action-btn span {
    display: none;
  }
}
</style>
