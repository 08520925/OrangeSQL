<script setup lang="ts">
import type { TableEntry, ColumnEntry } from "./types";

defineProps<{
  tables: TableEntry[];
  loading: boolean;
  isExpanded: (name: string) => boolean;
  getColumns: (name: string) => ColumnEntry[];
}>();

const emit = defineEmits<{
  selectTable: [tableName: string];
  toggleExpand: [tableName: string];
  refresh: [];
}>();
</script>

<template>
  <div class="schema-sidebar">
    <div class="sidebar-header">
      <span class="sidebar-title">テーブル</span>
      <button class="refresh-btn" :disabled="loading" @click="emit('refresh')">↻</button>
    </div>
    <div v-if="loading" class="loading">読み込み中...</div>
    <ul v-else class="table-list">
      <li v-for="table in tables" :key="table.name" class="table-group">
        <div class="table-item">
          <button class="expand-btn" @click="emit('toggleExpand', table.name)">
            {{ isExpanded(table.name) ? '▼' : '▶' }}
          </button>
          <span
            class="table-label"
            @click="emit('selectTable', table.name)"
          >
            <span class="table-icon">{{ table.type === 'view' ? '👁' : '▦' }}</span>
            <span class="table-name">{{ table.name }}</span>
          </span>
        </div>
        <ul v-if="isExpanded(table.name)" class="column-list">
          <li v-for="col in getColumns(table.name)" :key="col.name" class="column-item">
            <span class="col-pk">{{ col.pk ? '🔑' : '  ' }}</span>
            <span class="col-name">{{ col.name }}</span>
            <span class="col-type">{{ col.type }}</span>
            <span v-if="col.notNull" class="col-notnull">NOT NULL</span>
          </li>
        </ul>
      </li>
    </ul>
    <div v-if="!loading && tables.length === 0" class="empty">テーブルなし</div>
  </div>
</template>

<style scoped>
.schema-sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid #404040;
}

.sidebar-title {
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  color: #888888;
}

.refresh-btn {
  background: none;
  border: none;
  color: #cccccc;
  cursor: pointer;
  font-size: 16px;
  padding: 2px 6px;
  border-radius: 3px;
}

.refresh-btn:hover { background: #3a3a3a; }
.refresh-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.loading, .empty {
  padding: 12px;
  color: #666666;
  font-size: 12px;
}

.table-list {
  list-style: none;
  margin: 0;
  padding: 0;
  overflow-y: auto;
  flex: 1;
}

.table-group {
  border-bottom: 1px solid #2a2a2a;
}

.table-item {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.expand-btn {
  background: none;
  border: none;
  color: #888888;
  cursor: pointer;
  font-size: 10px;
  padding: 6px 4px 6px 8px;
  width: 24px;
  text-align: center;
}

.expand-btn:hover { color: #cccccc; }

.table-label {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  padding: 6px 8px 6px 0;
  cursor: pointer;
}

.table-label:hover { background: #2a2d2e; }

.table-icon { font-size: 12px; width: 16px; text-align: center; }
.table-name { color: #cccccc; }

.column-list {
  list-style: none;
  margin: 0;
  padding: 0 0 4px 0;
}

.column-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 2px 8px 2px 32px;
  font-size: 12px;
  color: #999999;
}

.col-pk { width: 16px; font-size: 10px; }
.col-name { color: #bbbbbb; min-width: 60px; }
.col-type { color: #6a9955; font-size: 11px; }
.col-notnull { color: #666666; font-size: 10px; }
</style>
