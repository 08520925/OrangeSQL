<template>
    <div ref="editorEl" style="height:100%;"></div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, defineExpose } from "vue";
import * as monaco from "monaco-editor";

const editorEl = ref<HTMLDivElement | null>(null);
let editor: monaco.editor.IStandaloneCodeEditor | null = null;

onMounted(() => {
    if (!editorEl.value) return;

    editor = monaco.editor.create(editorEl.value, {
        value: "SELECT id, name, age FROM users;",
        language: "sql",
        theme: "vs-dark",
        fontSize: 14,
        minimap: { enabled: false },
        automaticLayout: true,
    });

    // WebView2/Wails 初期レイアウト対策
    setTimeout(() => editor?.layout(), 0);
});

onBeforeUnmount(() => {
    editor?.dispose();
    editor = null;
});

function getValue(): string {
    return editor?.getValue() ?? "";
}

function setValue(v: string) {
    editor?.setValue(v);
}

defineExpose({ getValue, setValue });
</script>
