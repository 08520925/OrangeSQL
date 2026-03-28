import { onMounted, onBeforeUnmount, type Ref } from "vue";
import { EditorView, keymap } from "@codemirror/view";
import { EditorState } from "@codemirror/state";
import { sql } from "@codemirror/lang-sql";
import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { oneDark } from "@codemirror/theme-one-dark";
import { syntaxHighlighting, defaultHighlightStyle } from "@codemirror/language";
import type { SqlEditorApi } from "./types";

/**
 * CodeMirror 6 による SQL エディタを管理する composable。
 */
export function useSqlEditor(
  containerRef: Ref<HTMLElement | null>,
  onExecute: () => void,
): SqlEditorApi {
  let view: EditorView | null = null;

  onMounted(() => {
    const el = containerRef.value;
    if (el == null) return;

    const executeKeymap = keymap.of([
      {
        key: "Ctrl-Enter",
        run: () => {
          onExecute();
          return true;
        },
      },
      {
        key: "Mod-Enter",
        run: () => {
          onExecute();
          return true;
        },
      },
    ]);

    const state = EditorState.create({
      doc: "",
      extensions: [
        executeKeymap,
        keymap.of([...defaultKeymap, ...historyKeymap]),
        history(),
        sql(),
        oneDark,
        syntaxHighlighting(defaultHighlightStyle),
        EditorView.lineWrapping,
        EditorView.theme({
          "&": { height: "100%", fontSize: "14px" },
          ".cm-scroller": { overflow: "auto" },
          ".cm-content": { fontFamily: "'Consolas', 'Courier New', monospace" },
        }),
      ],
    });

    view = new EditorView({ state, parent: el });
  });

  onBeforeUnmount(() => {
    view?.destroy();
    view = null;
  });

  function getValue(): string {
    return view?.state.doc.toString() ?? "";
  }

  function setValue(text: string): void {
    if (view == null) return;
    view.dispatch({
      changes: { from: 0, to: view.state.doc.length, insert: text },
    });
  }

  return { getValue, setValue };
}
