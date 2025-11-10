<script setup lang="ts">
import { ref, computed } from "vue";
import { Icon } from "@iconify/vue";

interface Option {
  value: string;
  label: string;
}

defineProps<{
  modelValue: string;
  options: Option[];
  disabled?: boolean;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const isOpen = ref(false);
const selectElement = ref<HTMLElement | null>(null);

function selectOption(value: string) {
  emit("update:modelValue", value);
  isOpen.value = false;
}

const selectedLabel = (props: any) =>
  props.options.find((opt: Option) => opt.value === props.modelValue)?.label ||
  props.options[0]?.label;

const dropdownPosition = computed(() => {
  if (!selectElement.value || !isOpen.value) return "top";
  
  const rect = selectElement.value.getBoundingClientRect();
  const spaceBelow = window.innerHeight - rect.bottom;
  const spaceAbove = rect.top;
  
  // Open upward if there's more space above, otherwise open downward
  return spaceAbove > spaceBelow ? "top" : "bottom";
});
</script>

<template>
  <div class="custom-select" :class="{ disabled }" ref="selectElement">
    <button
      class="select-trigger"
      @click="isOpen = !isOpen"
      :disabled="disabled"
      type="button"
    >
      <span>{{ selectedLabel($props) }}</span>
      <Icon
        icon="tabler:chevron-down"
        width="16"
        height="16"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <transition name="dropdown">
      <div v-if="isOpen" class="select-dropdown" :class="{ 'open-up': dropdownPosition === 'top' }">
        <button
          v-for="option in options"
          :key="option.value"
          class="select-option"
          :class="{ active: modelValue === option.value }"
          @click="selectOption(option.value)"
          type="button"
        >
          {{ option.label }}
        </button>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.custom-select {
  position: relative;
  width: 100%;
}

.custom-select.disabled {
  opacity: 0.6;
  pointer-events: none;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-md);
  width: 100%;
  background: var(--bg);
  color: var(--fg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0.5rem 0.75rem;
  font-family: inherit;
  font-size: 0.875rem;
  cursor: pointer;
  transition: var(--transition);
  text-align: left;
}

.select-trigger:hover:not(:disabled) {
  border-color: transparent;
  opacity: 0.9;
}

.select-trigger:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.select-trigger svg {
  transition: transform 0.15s ease-out;
  flex-shrink: 0;
  color: var(--muted);
}

.select-trigger svg.rotate-180 {
  transform: rotate(180deg);
}

.select-dropdown {
  position: absolute;
  bottom: calc(100% + 0.25rem);
  left: 0;
  right: 0;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
  z-index: 10;
  box-shadow: none;
}

.select-dropdown.open-up {
  top: auto;
  bottom: calc(100% + 0.25rem);
}

.select-option {
  display: block;
  width: 100%;
  padding: 0.625rem 0.75rem;
  background: transparent;
  color: var(--fg);
  border: none;
  text-align: left;
  font-family: inherit;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background 0.1s ease-out;
}

.select-option:hover {
  background: var(--bg);
}

.select-option.active {
  background: var(--accent);
  color: white;
  font-weight: 500;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.12s ease-out;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .select-trigger svg,
  .select-option,
  .dropdown-enter-active,
  .dropdown-leave-active {
    transition: none;
  }
}
</style>
