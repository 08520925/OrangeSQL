import { ref, onMounted, onBeforeUnmount, type Ref } from "vue";
import * as monaco from "monaco-editor";
import type { SqlEditorApi } from "./types";

/**
 * Monaco Editor を管理する composable。
 * containerRef に渡した要素に Monaco を初期化する。
 */
export function useSqlEditor(
  containerRef: Ref<HTMLElement | null>,
  onExecute: () => void,
): SqlEditorApi {
  const editorInstance = ref<monaco.editor.IStandaloneCodeEditor | null>(null);

  onMounted(() => {
    const el = containerRef.value;
    if (el == null) return;

    const editor = monaco.editor.create(el, {
      value: "",
      language: "sql",
      theme: "vs-dark",
      automaticLayout: true,
      minimap: { enabled: false },
      fontSize: 14,
      lineNumbers: "on",
      scrollBeyondLastLine: false,
      wordWrap: "on",
    });

    // Ctrl+Enter で実行
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      onExecute();
    });

    editorInstance.value = editor;
  });

  onBeforeUnmount(() => {
    editorInstance.value?.dispose();
  });

  function getValue(): string {
    return editorInstance.value?.getValue() ?? "";
  }

  function setValue(sql: string): void {
    editorInstance.value?.setValue(sql);
  }

  return { getValue, setValue };
}
