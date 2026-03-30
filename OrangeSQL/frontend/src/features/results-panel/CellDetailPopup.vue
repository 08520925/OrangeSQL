<script setup lang="ts">
defineProps<{
  value: string;
}>();

const emit = defineEmits<{
  close: [];
}>();

async function copyToClipboard(text: string): Promise<void> {
  await navigator.clipboard.writeText(text);
}
</script>

<template>
  <div class="popup-backdrop" @click.self="emit('close')">
    <div class="popup">
      <div class="popup-header">
        <span class="popup-title">セル内容 ({{ value.length }}文字)</span>
        <div class="popup-actions">
          <button class="btn-copy" @click="copyToClipboard(value)">コピー</button>
          <button class="btn-close" @click="emit('close')">✕</button>
        </div>
      </div>
      <pre class="popup-body">{{ value }}</pre>
    </div>
  </div>
</template>

<style scoped>
.popup-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.popup {
  background: #2d2d2d;
  border: 1px solid #555555;
  border-radius: 6px;
  width: 600px;
  max-width: 90vw;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #404040;
}

.popup-title {
  font-size: 13px;
  font-weight: 600;
  color: #e0e0e0;
}

.popup-actions {
  display: flex;
  gap: 8px;
}

.btn-copy {
  padding: 3px 10px;
  background: #007acc;
  color: #ffffff;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
}
.btn-copy:hover { background: #1a8ad4; }

.btn-close {
  padding: 3px 8px;
  background: #3a3a3a;
  color: #cccccc;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}
.btn-close:hover { background: #444444; }

.popup-body {
  padding: 16px;
  overflow: auto;
  flex: 1;
  color: #cccccc;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>
