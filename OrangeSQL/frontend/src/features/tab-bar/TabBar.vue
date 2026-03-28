<script setup lang="ts">
import type { Tab } from "./types";

defineProps<{
  tabs: Tab[];
  activeTabId: string;
}>();

const emit = defineEmits<{
  switchTab: [id: string];
  addTab: [];
  closeTab: [id: string];
}>();
</script>

<template>
  <div class="tab-bar">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab"
      :class="{ active: tab.id === activeTabId }"
      @click="emit('switchTab', tab.id)"
    >
      <span class="tab-title">{{ tab.title }}</span>
      <button
        v-if="tabs.length > 1"
        class="tab-close"
        @click.stop="emit('closeTab', tab.id)"
      >×</button>
    </div>
    <button class="tab-add" @click="emit('addTab')">+</button>
  </div>
</template>

<style scoped>
.tab-bar {
  display: flex;
  align-items: stretch;
  background: #252526;
  border-bottom: 1px solid #404040;
  height: 32px;
  flex-shrink: 0;
  overflow-x: auto;
}

.tab-bar::-webkit-scrollbar {
  height: 2px;
}

.tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 12px;
  min-width: 80px;
  max-width: 200px;
  cursor: pointer;
  color: #888888;
  font-size: 12px;
  border-right: 1px solid #333333;
  white-space: nowrap;
  user-select: none;
}

.tab.active {
  background: #1e1e1e;
  color: #cccccc;
}

.tab:hover:not(.active) {
  background: #2a2d2e;
}

.tab-title {
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-close {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  font-size: 14px;
  padding: 0 2px;
  border-radius: 3px;
  line-height: 1;
}

.tab-close:hover {
  background: #444444;
  color: #ffffff;
}

.tab-add {
  background: none;
  border: none;
  color: #888888;
  cursor: pointer;
  font-size: 18px;
  padding: 0 12px;
}

.tab-add:hover {
  color: #cccccc;
}
</style>
