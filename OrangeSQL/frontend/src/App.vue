<template>
  <main
    style="
      padding: 12px;
      height: 100vh;
      box-sizing: border-box;
      display: grid;
      grid-template-rows: auto 1fr auto 1fr;
      gap: 10px;
    "
  >
    <!-- Header -->
    <header style="display:flex; align-items:center; gap:10px;">
      <h2 style="margin:0;">OrangeSQL</h2>
      <button @click="run" :disabled="running">実行</button>
    </header>

    <!-- Error Display -->
    <div
      v-if="error"
      style="
        background-color: #fff5f5;
        color: #c53030;
        border: 1px solid #fed7d7;
        border-radius: 8px;
        padding: 12px;
        white-space: pre-wrap;
        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      "
    >
      <strong>エラー:</strong><br />{{ error }}
    </div>

    <!-- Monaco Editor -->
    <div
      style="
        border:1px solid #ccc;
        border-radius:8px;
        overflow:hidden;
        height:100%;
      "
    >
      <div ref="editorEl" style="height:100%;"></div>
    </div>

    <!-- Result summary -->
    <div style="display:flex; align-items:center; gap:10px;">
      <div>Rows: {{ result.rows.length }}</div>
      <div v-if="result.columns.length">
        Columns: {{ result.columns.join(", ") }}
      </div>
    </div>

    <!-- Result grid -->
    <div
      style="
        border:1px solid #ccc;
        border-radius:8px;
        overflow:auto;
      "
    >
      <table
        v-if="result.columns.length"
        style="border-collapse: collapse; width: 100%;"
      >
        <thead>
          <tr>
            <th
              v-for="c in result.columns"
              :key="c"
              style="
                position: sticky;
                top: 0;
                background: #f6f6f6;
                border-bottom:1px solid #ddd;
                text-align:left;
                padding:8px;
              "
            >
              {{ c }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(r, idx) in result.rows" :key="idx">
            <td
              v-for="(cell, j) in r"
              :key="j"
              style="
                border-bottom:1px solid #eee;
                padding:8px;
                font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
              "
            >
              {{ cell }}
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else style="padding:12px; color:#666;">
        実行結果がここに表示されます
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import * as monaco from "monaco-editor";
import { ExecuteSQL } from "./wailsjs/go/main/App";

type QueryResult = {
  columns: string[];
  rows: string[][];
};

const editorEl = ref<HTMLDivElement | null>(null);
let editor: monaco.editor.IStandaloneCodeEditor;

const running = ref(false);
const error = ref("");

const result = ref<QueryResult>({
  columns: [],
  rows: [],
});

onMounted(() => {
  if (!editorEl.value) return;

  editor = monaco.editor.create(editorEl.value, {
    value: "SELECT * FROM users;",
    language: "sql",
    theme: "vs-dark",
    fontSize: 14,
    minimap: { enabled: false },
    automaticLayout: true,
  });

  // Wails + WebView2 対策（必須）
  setTimeout(() => {
    editor.layout();
  }, 0);
});

async function run() {
  error.value = "";
  running.value = true;

  try {
    const sql = editor.getValue();
    const res = await ExecuteSQL(sql);
    result.value = res;
  } catch (e: any) {
    error.value = e?.message ?? String(e);
  } finally {
    running.value = false;
  }
}
</script>
