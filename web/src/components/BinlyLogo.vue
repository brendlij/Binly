<!-- BinlyLogo.vue -->
<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 96 96"
    fill="none"
    role="img"
    :aria-label="ariaLabel"
    v-bind="$attrs"
  >
    <title v-if="title">{{ title }}</title>

    <!-- optional circular background -->
    <circle v-if="showCircle" cx="48" cy="48" r="44" :fill="bgFill" />

    <!-- file-corner accent -->
    <path :fill="accent" d="M48 16v12h12z" />

    <!-- smooth lowercase b -->
    <path
      d="M38 20v48a12 12 0 1 0 12-12H38"
      fill="none"
      :stroke="stroke"
      :stroke-width="strokeWidth"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

type Mode = "light" | "dark" | "auto";

interface Props {
  size?: number | string; // px or CSS size
  mode?: Mode; // light | dark | auto (prefers-color-scheme)
  showCircle?: boolean; // show circular bg
  accent?: string; // triangle fill
  strokeWidth?: number; // B stroke width
  title?: string; // <title> for a11y
  label?: string; // aria-label override
}

const props = withDefaults(defineProps<Props>(), {
  size: 96,
  mode: "auto",
  showCircle: true,
  accent: "var(--accent, #6366F1)",
  strokeWidth: 6,
  title: "Binly logo",
  label: "Binly logo",
});

const prefersDark = ref(false);

onMounted(() => {
  if (props.mode === "auto" && window?.matchMedia) {
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    prefersDark.value = mq.matches;
    // keep in sync if system theme changes
    mq.addEventListener?.("change", (e) => (prefersDark.value = e.matches));
  }
});

const isDark = computed<boolean>(() => {
  if (props.mode === "dark") return true;
  if (props.mode === "light") return false;
  return prefersDark.value; // auto
});

const bgFill = computed(() => (isDark.value ? "#000" : "#fff"));
const stroke = computed(() => (isDark.value ? "#fff" : "#000"));

const ariaLabel = computed(() => props.label || "Binly logo");
</script>
