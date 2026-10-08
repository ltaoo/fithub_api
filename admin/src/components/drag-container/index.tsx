import { Component, createSignal, onMount, onCleanup, JSX } from "@/timeless";

import { ViewComponentProps } from "@/store/types";

import styles from "./styles.module.css";

export function DragContainer(props: { storage: ViewComponentProps["storage"] } & JSX.HTMLAttributes<HTMLDivElement>) {
  // const x = window.innerWidth / 2 - 100; // 假设容器宽度为200px
  // const y = window.innerHeight / 2 - 100;
  const { x, y } = props.storage.get("media_manager");
  const [position, setPosition] = createSignal({ x, y });
  const [isDragging, setIsDragging] = createSignal(false);
  const [isCollapsed, setIsCollapsed] = createSignal(false);
  const [dragOffset, setDragOffset] = createSignal({ x: 0, y: 0 });
  let containerRef: HTMLDivElement | undefined;

  const handleMouseDown = (e: MouseEvent) => {
    if (e.target instanceof HTMLElement && e.target.closest(`.${styles.header}`)) {
      setIsDragging(true);
      const rect = containerRef?.getBoundingClientRect();
      if (rect) {
        setDragOffset({
          x: e.clientX - rect.left,
          y: e.clientY - rect.top,
        });
      }
    }
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (isDragging()) {
      setPosition({
        x: e.clientX - dragOffset().x,
        y: e.clientY - dragOffset().y,
      });
    }
  };

  const handleMouseUp = () => {
    props.storage.set("media_manager", {
      x: position().x,
      y: position().y,
    });
    setIsDragging(false);
  };

  onMount(() => {
    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  });

  onCleanup(() => {
    document.removeEventListener("mousemove", handleMouseMove);
    document.removeEventListener("mouseup", handleMouseUp);
  });

  return (
    <div
      ref={containerRef}
      class={styles.container}
      style={{
        transform: `translate(${position().x}px, ${position().y}px)`,
        cursor: isDragging() ? "grabbing" : "default",
      }}
      onMouseDown={handleMouseDown}
    >
      <div class={styles.header}>
        <span class={styles.title}>{props.title}</span>
        <button class={styles.collapseButton} onClick={() => setIsCollapsed(!isCollapsed())}>
          {isCollapsed() ? "▼" : "▲"}
        </button>
      </div>
      <div class={`${styles.content} ${isCollapsed() ? styles.collapsed : ""}`}>{props.children}</div>
    </div>
  );
}
