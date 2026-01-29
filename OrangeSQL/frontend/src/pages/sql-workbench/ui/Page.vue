<template>
    <main style="
      padding: 12px;
      height: 100vh;
      box-sizing: border-box;
      display: grid;
      grid-template-columns: 260px 1fr;
      grid-template-rows: auto 1fr auto 1fr;
      gap: 10px;
    ">
        <!-- Header -->
        <header style="grid-column: 1 / -1; display:flex; align-items:center; gap:10px;">
            <h2 style="margin:0;">OrangeSQL</h2>
            <button @click="run" :disabled="running">実行</button>
            <span v-if="error" style="color:#c00;">{{ error }}</span>
        </header>

        <!-- Left: Schema -->
        <section style="grid-row: 2 / 5; border:1px solid #ccc; border-radius:8px; overflow:auto;">
            <SchemaSidebar />
        </section>

        <!-- Editor -->
        <section style="border:1px solid #ccc; border-radius:8px; overflow:hidden; height:100%;">
            <EditorPanel ref="editorRef" />
        </section>

        <!-- Result summary -->
        <section style="display:flex; align-items:center; gap:10px;">
            <div>Rows: {{ result.rows.length }}</div>
            <div v-if="result.columns.length">Columns: {{ result.columns.join(", ") }}</div>
        </section>

        <!-- Results -->
        <section style="border:1px solid #ccc; border-radius:8px; overflow:auto;">
            <ResultsPanel :result="result" />
        </section>
    </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { SchemaSidebar } from "@/widgets/schema-sidebar";
import { EditorPanel } from "@/widgets/editor-panel";
import { ResultsPanel } from "@/widgets/results-panel";
import type { QueryResult } from "@/entities/query";
import { useRunSql } from "@/features/run-sql";

const editorRef = ref<InstanceType<typeof EditorPanel> | null>(null);

const result = ref<QueryResult>({ columns: [], rows: [] });

const { runSql, running, error } = useRunSql();

async function run() {
    const sql = editorRef.value?.getValue() ?? "";
    const res = await runSql(sql);
    if (res) result.value = res;
}
</script>
