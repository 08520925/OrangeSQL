<script setup lang="ts">
import { ref } from "vue";
import {
  downloadCSV,
  downloadJSON,
  copyCSV,
  copyMarkdown,
} from "./result-exporter";

const props = defineProps<{
  columns: string[];
  rows: (string | null)[][];
}>();

const copyMenuOpen = ref<boolean>(false);
const copyFeedback = ref<string>("");

function handleDownloadCSV(): void {
  downloadCSV(props.columns, props.rows);
}

function handleDownloadJSON(): void {
  downloadJSON(props.columns, props.rows);
}

async function handleCopyCSV(): Promise<void> {
  await copyCSV(props.columns, props.rows);
  showFeedback("CSV をコピーしました");
}

async function handleCopyMarkdown(): Promise<void> {
  await copyMarkdown(props.columns, props.rows);
  showFeedback("Markdown をコピーしました");
}

function showFeedback(msg: string): void {
  copyMenuOpen.value = false;
  copyFeedback.value = msg;
  setTimeout(() => {
    copyFeedback.value = "";
  }, 2000);
}
</script>

<template>
  <div class="export-toolbar">
    <button class="export-btn" title="CSV ダウンロード" @click="handleDownloadCSV">CSV</button>
    <button class="export-btn" title="JSON ダウンロード" @click="handleDownloadJSON">JSON</button>

    <div class="copy-wrapper">
      <button class="export-btn" @click="copyMenuOpen = !copyMenuOpen">コピー ▼</button>
      <div v-if="copyMenuOpen" class="copy-menu">
        <div class="copy-item" @click="handleCopyCSV">CSV としてコピー</div>
        <div class="copy-item" @click="handleCopyMarkdown">Markdown テーブルとしてコピー</div>
      </div>
      <div v-if="copyMenuOpen" class="copy-backdrop" @click="copyMenuOpen = false"></div>
    </div>

    <span v-if="copyFeedback !== ''" class="copy-feedback">{{ copyFeedback }}</span>
  </div>
</template>

<style scoped>
.export-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
}

.export-btn {
  padding: 2px 8px;
  background: #3a3a3a;
  color: #cccccc;
  border: 1px solid #555555;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
}
.export-btn:hover { background: #444444; }

.copy-wrapper { position: relative; }

.copy-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: #2d2d2d;
  border: 1px solid #555555;
  border-radius: 4px;
  min-width: 200px;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0,0,0,0.4);
}

.copy-item {
  padding: 6px 12px;
  cursor: pointer;
  font-size: 12px;
  color: #cccccc;
}
.copy-item:hover { background: #3a3a3a; }

.copy-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 99;
}

.copy-feedback {
  font-size: 11px;
  color: #4ec9b0;
  margin-left: 4px;
}
</style>
