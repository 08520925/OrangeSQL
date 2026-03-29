<script setup lang="ts">
import { ref, nextTick } from "vue";

const props = defineProps<{
  value: string | null;
  editable: boolean;
  modified: boolean;
}>();

const emit = defineEmits<{
  update: [value: string | null];
  setNull: [];
}>();

const editing = ref<boolean>(false);
const inputRef = ref<HTMLInputElement | null>(null);
const editValue = ref<string>("");

function startEdit(): void {
  if (!props.editable) return;
  editing.value = true;
  editValue.value = props.value ?? "";
  void nextTick(() => {
    inputRef.value?.focus();
    inputRef.value?.select();
  });
}

function confirmEdit(): void {
  editing.value = false;
  emit("update", editValue.value);
}

function cancelEdit(): void {
  editing.value = false;
}

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === "Enter") {
    confirmEdit();
  } else if (e.key === "Escape") {
    cancelEdit();
  } else if (e.ctrlKey && e.shiftKey && e.key === "N") {
    e.preventDefault();
    editing.value = false;
    emit("setNull");
  }
}
</script>

<template>
  <td
    :class="{
      'null-cell': value === null && !editing,
      'modified-cell': modified,
      'pk-cell': !editable,
    }"
    @dblclick="startEdit"
  >
    <input
      v-if="editing"
      ref="inputRef"
      v-model="editValue"
      class="cell-input"
      @blur="confirmEdit"
      @keydown="handleKeydown"
    />
    <span v-else>{{ value === null ? 'NULL' : value }}</span>
  </td>
</template>

<style scoped>
td {
  padding: 4px 12px;
  border-bottom: 1px solid #333333;
  color: #cccccc;
  white-space: nowrap;
  cursor: default;
  position: relative;
}

.pk-cell {
  background: #2a2a2a;
  color: #888888;
}

.null-cell {
  color: #666666;
  font-style: italic;
}

.modified-cell {
  background: #3d3000;
}

.cell-input {
  width: 100%;
  padding: 2px 4px;
  background: #1a1a2e;
  border: 1px solid #007acc;
  color: #cccccc;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
  box-sizing: border-box;
  outline: none;
}
</style>
