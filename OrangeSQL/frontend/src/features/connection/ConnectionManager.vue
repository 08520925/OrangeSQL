<script setup lang="ts">
import type { ConnectionProfile } from "./types";

defineProps<{
  profiles: ConnectionProfile[];
  activeId: string;
}>();

const emit = defineEmits<{
  edit: [profile: ConnectionProfile];
  delete: [id: string];
  close: [];
}>();
</script>

<template>
  <div class="dialog-backdrop" @click.self="emit('close')">
    <div class="dialog">
      <div class="dialog-title">接続管理</div>
      <div v-if="profiles.length === 0" class="empty">プロファイルなし</div>
      <div v-for="p in profiles" :key="p.id" class="profile-row">
        <span class="profile-name">{{ p.name }}</span>
        <span class="profile-driver">{{ p.driver }}</span>
        <div class="profile-actions">
          <button class="btn-sm" @click="emit('edit', p)">編集</button>
          <button
            class="btn-sm btn-danger"
            :disabled="p.id === activeId"
            @click="emit('delete', p.id)"
          >削除</button>
        </div>
      </div>
      <div class="dialog-actions">
        <button class="btn btn-cancel" @click="emit('close')">閉じる</button>
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
  padding: 20px; width: 450px; max-width: 90vw;
}
.dialog-title { font-size: 14px; font-weight: 600; color: #e0e0e0; margin-bottom: 16px; }
.empty { color: #666666; font-size: 12px; padding: 8px 0; }
.profile-row {
  display: flex; align-items: center; gap: 8px; padding: 8px 0;
  border-bottom: 1px solid #3a3a3a; font-size: 13px;
}
.profile-name { flex: 1; color: #cccccc; }
.profile-driver { color: #666666; font-size: 11px; }
.profile-actions { display: flex; gap: 4px; }
.btn-sm {
  padding: 3px 8px; border: none; border-radius: 3px; cursor: pointer;
  font-size: 11px; background: #3a3a3a; color: #cccccc;
}
.btn-sm:hover { background: #444444; }
.btn-sm:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-danger:not(:disabled):hover { background: #c72e2e; color: #ffffff; }
.dialog-actions { display: flex; justify-content: flex-end; margin-top: 16px; }
.btn { padding: 6px 16px; border: none; border-radius: 3px; cursor: pointer; font-size: 12px; }
.btn-cancel { background: #3a3a3a; color: #cccccc; }
.btn-cancel:hover { background: #444444; }
</style>
