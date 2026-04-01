<script setup lang="ts">
import { ref, nextTick, computed } from "vue";

const props = defineProps<{
  value: string | null;
  editable: boolean;
  modified: boolean;
}>();

const emit = defineEmits<{
  update: [value: string | null];
  setNull: [];
  showDetail: [value: string];
}>();

const editing = ref<boolean>(false);
const inputRef = ref<HTMLInputElement | null>(null);
const editValue = ref<string>("");

const isLongText = computed(() => props.value != null && props.value.length > 100);

function openDetail(): void {
  if (props.value != null) {
    emit("showDetail", props.value);
  }
}

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

function setNull(): void {
  editing.value = false;
  emit("setNull");
}

function handleBlur(e: FocusEvent): void {
  // NULLボタンへのフォーカス移動時は確定しない
  const related = e.relatedTarget as HTMLElement | null;
  if (related?.classList.contains("null-btn")) return;
  confirmEdit();
}

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === "Enter") {
    confirmEdit();
  } else if (e.key === "Escape") {
    cancelEdit();
  } else if (e.ctrlKey && e.shiftKey && e.key === "N") {
    e.preventDefault();
    setNull();
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
    <div v-if="editing" class="edit-container">
      <input
        ref="inputRef"
        v-model="editValue"
        class="cell-input"
        @blur="handleBlur"
        @keydown="handleKeydown"
      />
      <button
        class="null-btn"
        title="NULLを設定 (Ctrl+Shift+N)"
        @mousedown.prevent="setNull"
      >NULL</button>
    </div>
    <template v-else>
      <span class="cell-text">{{ value === null ? 'NULL' : value }}</span>
      <button v-if="isLongText" class="expand-btn" @click.stop="openDetail">…</button>
    </template>
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
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cell-text {
  overflow: hidden;
  text-overflow: ellipsis;
}

.expand-btn {
  position: absolute;
  right: 2px;
  top: 50%;
  transform: translateY(-50%);
  background: #007acc;
  color: #ffffff;
  border: none;
  border-radius: 2px;
  cursor: pointer;
  font-size: 11px;
  padding: 1px 5px;
  line-height: 1;
}
.expand-btn:hover { background: #1a8ad4; }

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

.edit-container {
  display: flex;
  align-items: center;
  gap: 2px;
}

.cell-input {
  flex: 1;
  min-width: 0;
  padding: 2px 4px;
  background: #1a1a2e;
  border: 1px solid #007acc;
  color: #cccccc;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
  box-sizing: border-box;
  outline: none;
}

.null-btn {
  padding: 1px 4px;
  background: #5a4000;
  color: #e0a030;
  border: 1px solid #7a6020;
  border-radius: 2px;
  cursor: pointer;
  font-size: 10px;
  font-weight: 600;
  line-height: 1.2;
  white-space: nowrap;
  flex-shrink: 0;
}
.null-btn:hover { background: #7a5500; }
</style>
