<script setup lang="ts">
import { ref } from "vue";
import type { ConnectionProfile } from "./types";

const props = defineProps<{
  editProfile?: ConnectionProfile;
}>();

const emit = defineEmits<{
  save: [name: string, driver: string, path: string];
  close: [];
}>();

const name = ref<string>(props.editProfile?.name ?? "");
const path = ref<string>(props.editProfile?.path ?? "");
const error = ref<string>("");

function handleSave(): void {
  if (name.value.trim() === "" || path.value.trim() === "") {
    error.value = "接続名とファイルパスは必須です";
    return;
  }
  error.value = "";
  emit("save", name.value.trim(), "sqlite", path.value.trim());
}
</script>

<template>
  <div class="dialog-backdrop" @click.self="emit('close')">
    <div class="dialog">
      <div class="dialog-title">{{ editProfile != null ? '接続を編集' : '新規接続' }}</div>
      <label class="field">
        <span class="field-label">接続名</span>
        <input v-model="name" class="field-input" placeholder="例: 開発DB" />
      </label>
      <label class="field">
        <span class="field-label">ドライバ</span>
        <select class="field-input" disabled>
          <option>SQLite</option>
        </select>
      </label>
      <label class="field">
        <span class="field-label">ファイルパス</span>
        <input v-model="path" class="field-input" placeholder="例: C:/work/dev.db" />
      </label>
      <div v-if="error !== ''" class="error">{{ error }}</div>
      <div class="dialog-actions">
        <button class="btn btn-cancel" @click="emit('close')">キャンセル</button>
        <button class="btn btn-primary" @click="handleSave">{{ editProfile != null ? '保存' : '接続' }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-backdrop {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 200;
}
.dialog {
  background: #2d2d2d; border: 1px solid #555555; border-radius: 6px;
  padding: 20px; width: 400px; max-width: 90vw;
}
.dialog-title { font-size: 14px; font-weight: 600; color: #e0e0e0; margin-bottom: 16px; }
.field { display: block; margin-bottom: 12px; }
.field-label { display: block; font-size: 12px; color: #888888; margin-bottom: 4px; }
.field-input {
  width: 100%; padding: 6px 8px; background: #1e1e1e; border: 1px solid #555555;
  border-radius: 3px; color: #cccccc; font-size: 13px; box-sizing: border-box;
}
.field-input:focus { border-color: #007acc; outline: none; }
.error { color: #f44747; font-size: 12px; margin-bottom: 12px; }
.dialog-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
.btn { padding: 6px 16px; border: none; border-radius: 3px; cursor: pointer; font-size: 12px; }
.btn-cancel { background: #3a3a3a; color: #cccccc; }
.btn-cancel:hover { background: #444444; }
.btn-primary { background: #007acc; color: #ffffff; }
.btn-primary:hover { background: #1a8ad4; }
</style>
