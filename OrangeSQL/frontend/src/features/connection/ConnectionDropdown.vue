<script setup lang="ts">
import { ref } from "vue";
import type { ConnectionProfile } from "./types";

defineProps<{
  profiles: ConnectionProfile[];
  activeId: string;
}>();

const emit = defineEmits<{
  connect: [id: string];
  openNew: [];
  openManager: [];
}>();

const open = ref<boolean>(false);

function toggle(): void {
  open.value = !open.value;
}

function handleConnect(id: string): void {
  open.value = false;
  emit("connect", id);
}

function handleNew(): void {
  open.value = false;
  emit("openNew");
}

function handleManager(): void {
  open.value = false;
  emit("openManager");
}
</script>

<template>
  <div class="dropdown-wrapper">
    <button class="dropdown-trigger" @click="toggle">
      {{ profiles.find(p => p.id === activeId)?.name ?? '未接続' }} ▼
    </button>
    <div v-if="open" class="dropdown-menu">
      <div
        v-for="p in profiles"
        :key="p.id"
        class="dropdown-item"
        :class="{ active: p.id === activeId }"
        @click="handleConnect(p.id)"
      >
        <span class="check">{{ p.id === activeId ? '●' : '' }}</span>
        <span>{{ p.name }}</span>
      </div>
      <div class="dropdown-divider"></div>
      <div class="dropdown-item" @click="handleNew">+ 新規接続</div>
      <div class="dropdown-item" @click="handleManager">⚙ 接続管理</div>
    </div>
    <div v-if="open" class="dropdown-backdrop" @click="open = false"></div>
  </div>
</template>

<style scoped>
.dropdown-wrapper { position: relative; }

.dropdown-trigger {
  background: #3a3a3a;
  color: #cccccc;
  border: 1px solid #555555;
  padding: 3px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}
.dropdown-trigger:hover { background: #444444; }

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  background: #2d2d2d;
  border: 1px solid #555555;
  border-radius: 4px;
  min-width: 200px;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0,0,0,0.4);
}

.dropdown-item {
  padding: 6px 12px;
  cursor: pointer;
  font-size: 12px;
  color: #cccccc;
  display: flex;
  align-items: center;
  gap: 8px;
}
.dropdown-item:hover { background: #3a3a3a; }
.dropdown-item.active { color: #ffffff; }

.check { width: 12px; font-size: 8px; }

.dropdown-divider {
  height: 1px;
  background: #444444;
  margin: 4px 0;
}

.dropdown-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 99;
}
</style>
