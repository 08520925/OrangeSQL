<script setup lang="ts">
// レイアウト骨格のみ。各機能コンポーネントは Issue #7〜#10 で実装する。
</script>

<template>
  <div class="app-layout">
    <header class="header-bar">
      <span class="app-name">OrangeSQL</span>
    </header>
    <div class="main-area">
      <aside class="sidebar">
        <p class="placeholder">サイドバー（#9 で実装）</p>
      </aside>
      <div class="content-area">
        <div class="editor-area">
          <p class="placeholder">SQLエディタ（#7 で実装）</p>
        </div>
        <div class="results-area">
          <p class="placeholder">結果パネル（#8 で実装）</p>
        </div>
        <div class="status-bar">
          <span>Ready</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
  color: #cccccc;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 13px;
}

.header-bar {
  height: 40px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
  flex-shrink: 0;
}

.app-name {
  font-weight: 600;
  font-size: 14px;
  color: #e0e0e0;
}

.main-area {
  display: flex;
  flex: 1;
  min-height: 0;
}

.sidebar {
  width: 250px;
  background: #252526;
  border-right: 1px solid #404040;
  overflow-y: auto;
  flex-shrink: 0;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.editor-area {
  flex: 1;
  min-height: 0;
  border-bottom: 1px solid #404040;
}

.results-area {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.status-bar {
  height: 28px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  background: #007acc;
  color: #ffffff;
  font-size: 12px;
  flex-shrink: 0;
}

.placeholder {
  padding: 20px;
  color: #666666;
  font-style: italic;
}
</style>
