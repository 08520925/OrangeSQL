import { ref, computed } from "vue";
import type { CellEdit, NewRow } from "./sql-generator";
import {
  generateUpdates,
  generateInserts,
  generateDeletes,
} from "./sql-generator";

/**
 * テーブル編集状態を管理する composable。
 */
export function useTableEdit(
  tableName: string,
  columns: string[],
  pkColumns: string[],
  originalRows: (string | null)[][],
) {
  const edits = ref<CellEdit[]>([]);
  const newRows = ref<NewRow[]>([]);
  const deletedRows = ref<Set<number>>(new Set());

  const hasChanges = computed<boolean>(
    () =>
      edits.value.length > 0 ||
      newRows.value.some((r) => r.values.some((v) => v !== null && v !== "")) ||
      deletedRows.value.size > 0,
  );

  function updateCell(rowIndex: number, colIndex: number, newValue: string | null): void {
    const original = originalRows[rowIndex]?.[colIndex] ?? null;
    // 変更なしなら edit を削除
    if (newValue === original) {
      edits.value = edits.value.filter(
        (e) => !(e.rowIndex === rowIndex && e.colIndex === colIndex),
      );
      return;
    }
    // 既存の edit を更新 or 追加
    const existing = edits.value.find(
      (e) => e.rowIndex === rowIndex && e.colIndex === colIndex,
    );
    if (existing != null) {
      existing.newValue = newValue;
    } else {
      edits.value.push({ rowIndex, colIndex, newValue });
    }
  }

  function addNewRow(): void {
    newRows.value.push({
      values: columns.map(() => null),
    });
  }

  function updateNewRowCell(newRowIndex: number, colIndex: number, value: string | null): void {
    const row = newRows.value[newRowIndex];
    if (row != null) {
      row.values[colIndex] = value;
    }
  }

  function toggleDeleteRow(rowIndex: number): void {
    const s = new Set(deletedRows.value);
    if (s.has(rowIndex)) {
      s.delete(rowIndex);
    } else {
      s.add(rowIndex);
    }
    deletedRows.value = s;
  }

  function discardAll(): void {
    edits.value = [];
    newRows.value = [];
    deletedRows.value = new Set();
  }

  function generateAllSQL(): string[] {
    const updates = generateUpdates(tableName, columns, pkColumns, originalRows, edits.value);
    const inserts = generateInserts(tableName, columns, pkColumns, newRows.value);
    const deletes = generateDeletes(tableName, columns, pkColumns, originalRows, deletedRows.value);
    return [...updates, ...inserts, ...deletes];
  }

  function changeSummary(): { updates: number; inserts: number; deletes: number } {
    const byRow = new Set(edits.value.map((e) => e.rowIndex));
    return {
      updates: byRow.size,
      inserts: newRows.value.filter((r) => r.values.some((v) => v !== null && v !== "")).length,
      deletes: deletedRows.value.size,
    };
  }

  return {
    edits,
    newRows,
    deletedRows,
    hasChanges,
    updateCell,
    addNewRow,
    updateNewRowCell,
    toggleDeleteRow,
    discardAll,
    generateAllSQL,
    changeSummary,
  };
}
