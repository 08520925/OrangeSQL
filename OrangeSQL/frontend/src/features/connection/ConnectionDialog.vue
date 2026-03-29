<script setup lang="ts">
import { ref, computed } from "vue";
import type { ConnectionProfile, CreateProfileRequest, DriverType } from "./types";
import { DRIVERS, DEFAULT_PORTS } from "./types";

const props = defineProps<{
  editProfile?: ConnectionProfile;
}>();

const emit = defineEmits<{
  save: [req: CreateProfileRequest];
  close: [];
}>();

const driver = ref<DriverType>(props.editProfile?.driver ?? "sqlite");
const name = ref<string>(props.editProfile?.name ?? "");
const path = ref<string>(props.editProfile?.path ?? "");
const host = ref<string>(props.editProfile?.host ?? "localhost");
const port = ref<number>(props.editProfile?.port ?? DEFAULT_PORTS[driver.value]);
const user = ref<string>(props.editProfile?.user ?? "");
const password = ref<string>(props.editProfile?.password ?? "");
const dbName = ref<string>(props.editProfile?.dbName ?? "");
const sslMode = ref<string>(props.editProfile?.sslMode ?? "disable");
const error = ref<string>("");

const isTCP = computed<boolean>(() => driver.value !== "sqlite");

function onDriverChange(): void {
  port.value = DEFAULT_PORTS[driver.value];
}

function handleSave(): void {
  if (name.value.trim() === "") {
    error.value = "接続名は必須です";
    return;
  }

  if (driver.value === "sqlite") {
    if (path.value.trim() === "") {
      error.value = "ファイルパスは必須です";
      return;
    }
    error.value = "";
    emit("save", { name: name.value.trim(), driver: driver.value, path: path.value.trim() });
    return;
  }

  // TCP 系バリデーション
  if (host.value.trim() === "") {
    error.value = "ホストは必須です";
    return;
  }
  if (dbName.value.trim() === "") {
    error.value = "データベース名は必須です";
    return;
  }

  error.value = "";
  emit("save", {
    name: name.value.trim(),
    driver: driver.value,
    host: host.value.trim(),
    port: port.value,
    user: user.value.trim(),
    password: password.value,
    dbName: dbName.value.trim(),
    sslMode: sslMode.value,
  });
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
        <select v-model="driver" class="field-input" :disabled="editProfile != null" @change="onDriverChange">
          <option v-for="d in DRIVERS" :key="d.value" :value="d.value">{{ d.label }}</option>
        </select>
      </label>

      <!-- SQLite -->
      <template v-if="!isTCP">
        <label class="field">
          <span class="field-label">ファイルパス</span>
          <input v-model="path" class="field-input" placeholder="例: C:/work/dev.db" />
        </label>
      </template>

      <!-- TCP 系 -->
      <template v-if="isTCP">
        <div class="field-row">
          <label class="field field-grow">
            <span class="field-label">ホスト</span>
            <input v-model="host" class="field-input" placeholder="localhost" />
          </label>
          <label class="field field-port">
            <span class="field-label">ポート</span>
            <input v-model.number="port" class="field-input" type="number" />
          </label>
        </div>
        <label class="field">
          <span class="field-label">ユーザー</span>
          <input v-model="user" class="field-input" placeholder="例: postgres" />
        </label>
        <label class="field">
          <span class="field-label">パスワード</span>
          <input v-model="password" class="field-input" type="password" />
        </label>
        <label class="field">
          <span class="field-label">データベース</span>
          <input v-model="dbName" class="field-input" placeholder="例: mydb" />
        </label>
        <label v-if="driver === 'postgres'" class="field">
          <span class="field-label">SSL モード</span>
          <select v-model="sslMode" class="field-input">
            <option value="disable">disable</option>
            <option value="require">require</option>
            <option value="verify-ca">verify-ca</option>
            <option value="verify-full">verify-full</option>
          </select>
        </label>
      </template>

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
.field-input:disabled { opacity: 0.5; }
.field-row { display: flex; gap: 8px; }
.field-grow { flex: 1; }
.field-port { width: 80px; }
.error { color: #f44747; font-size: 12px; margin-bottom: 12px; }
.dialog-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
.btn { padding: 6px 16px; border: none; border-radius: 3px; cursor: pointer; font-size: 12px; }
.btn-cancel { background: #3a3a3a; color: #cccccc; }
.btn-cancel:hover { background: #444444; }
.btn-primary { background: #007acc; color: #ffffff; }
.btn-primary:hover { background: #1a8ad4; }
</style>
