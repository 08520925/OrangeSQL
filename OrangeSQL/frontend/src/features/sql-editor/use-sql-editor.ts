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

  /**
   * カーソル位置にある SQL 文を返す。
   * セミコロンで区切られた複数文がある場合、カーソルが置かれている文だけを返す。
   */
  function getStatementAtCursor(): string {
    if (view == null) return "";
    const doc = view.state.doc.toString();
    const cursorPos = view.state.selection.main.head;

    // セミコロンで分割し、各文の開始・終了位置を記録
    let pos = 0;
    const statements: { start: number; end: number; text: string }[] = [];
    for (const part of doc.split(";")) {
      const end = pos + part.length;
      const trimmed = part.trim();
      if (trimmed.length > 0) {
        statements.push({ start: pos, end, text: trimmed });
      }
      pos = end + 1; // +1 for semicolon
    }

    // カーソル位置を含む文を探す
    for (const stmt of statements) {
      if (cursorPos >= stmt.start && cursorPos <= stmt.end + 1) {
        return stmt.text;
      }
    }

    // 見つからなければ最後の文、それもなければ全文
    const last = statements[statements.length - 1];
    return last?.text ?? doc.trim();
  }

  return { getValue, setValue, getStatementAtCursor };
}
