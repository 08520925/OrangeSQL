import { ref, onBeforeUnmount } from "vue";

/**
 * エディタ/結果パネルの上下分割リサイズを管理する composable。
 * editorRatio は 0〜1 の値で、エディタが占める割合を表す。
 */
export function useResize(containerSelector: string) {
  const editorRatio = ref<number>(0.5);
  let dragging = false;

  function onMouseDown(_e: MouseEvent): void {
    dragging = true;
    document.addEventListener("mousemove", onMouseMove);
    document.addEventListener("mouseup", onMouseUp);
    document.body.style.cursor = "row-resize";
    document.body.style.userSelect = "none";
  }

  function onMouseMove(e: MouseEvent): void {
    if (!dragging) return;
    const container = document.querySelector(containerSelector);
    if (container == null) return;

    const rect = container.getBoundingClientRect();
    const statusBarHeight = 28;
    const availableHeight = rect.height - statusBarHeight;
    const offsetY = e.clientY - rect.top;
    let ratio = offsetY / availableHeight;

    // 最小高さ 100px を確保
    const minRatio = 100 / availableHeight;
    const maxRatio = 1 - minRatio;
    ratio = Math.max(minRatio, Math.min(maxRatio, ratio));

    editorRatio.value = ratio;
  }

  function onMouseUp(): void {
    dragging = false;
    document.removeEventListener("mousemove", onMouseMove);
    document.removeEventListener("mouseup", onMouseUp);
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
  }

  onBeforeUnmount(() => {
    document.removeEventListener("mousemove", onMouseMove);
    document.removeEventListener("mouseup", onMouseUp);
  });

  return { editorRatio, onMouseDown };
}
