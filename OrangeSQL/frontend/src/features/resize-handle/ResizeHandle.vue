<script setup lang="ts">
defineEmits<{
  mousedown: [e: MouseEvent];
}>();
</script>

<template>
  <div class="resize-handle" @mousedown="$emit('mousedown', $event)"></div>
</template>

<style scoped>
.resize-handle {
  height: 4px;
  cursor: row-resize;
  background: #404040;
  flex-shrink: 0;
}

.resize-handle:hover {
  background: #007acc;
}
</style>
