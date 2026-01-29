<template>
    <div v-if="result.columns.length">
        <table style="border-collapse: collapse; width: 100%;">
            <thead>
                <tr>
                    <th v-for="c in result.columns" :key="c" style="
              position: sticky;
              top: 0;
              background: #f6f6f6;
              border-bottom:1px solid #ddd;
              text-align:left;
              padding:8px;
            ">
                        {{ c }}
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(r, idx) in result.rows" :key="idx">
                    <td v-for="(cell, j) in r" :key="j" style="
              border-bottom:1px solid #eee;
              padding:8px;
              font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
            ">
                        {{ cell }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div v-else style="padding:12px; color:#666;">
        実行結果がここに表示されます
    </div>
</template>

<script setup lang="ts">
import type { QueryResult } from "@/entities/query";

defineProps<{
    result: QueryResult;
}>();
</script>
