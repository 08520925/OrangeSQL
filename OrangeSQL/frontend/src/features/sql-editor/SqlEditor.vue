<script setup lang="ts">
import { ref } from "vue";
import { useSqlEditor } from "./use-sql-editor";
import type { SqlEditorApi } from "./types";

const emit = defineEmits<{
  execute: [];
}>();

const containerRef = ref<HTMLElement | null>(null);

const editorApi: SqlEditorApi = useSqlEditor(containerRef, () => {
  emit("execute");
});

defineExpose({
  getValue: editorApi.getValue,
  setValue: editorApi.setValue,
  getStatementAtCursor: editorApi.getStatementAtCursor,
  refreshCompletions: editorApi.refreshCompletions,
});
</script>

<template>
  <div ref="containerRef" class="sql-editor-container"></div>
</template>

<style scoped>
.sql-editor-container {
  width: 100%;
  height: 100%;
}
</style>
