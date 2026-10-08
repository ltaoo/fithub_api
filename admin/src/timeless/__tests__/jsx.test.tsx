// @vitest-environment jsdom
import { expect, test } from "vitest";
import { createSignal, For, Show, render } from "../index";
import { Content } from "../../packages/ui/popper";
import { vi } from "vitest";

test("floating content recalculates placement after conditional mounting", async () => {
  const setFloating = vi.fn();
  const place = vi.fn();
  const store = {
    state: { x: 0, y: 0, isPlaced: false, strategy: "fixed", placedSide: "bottom", placedAlign: "start" },
    onStateChange: vi.fn(), setFloating, place,
  };
  const [visible, set_visible] = createSignal(false);
  const root = document.createElement("div"); document.body.append(root);
  render(() => <Show when={visible()}><Content store={store as any}>Menu</Content></Show>, root);
  set_visible(true);
  await new Promise(resolve => setTimeout(resolve, 5));
  expect(setFloating).toHaveBeenCalledOnce();
  expect(place).toHaveBeenCalledOnce();
  expect(setFloating.mock.invocationCallOrder[0]).toBeLessThan(place.mock.invocationCallOrder[0]);
});

test("compiled TSX preserves lazy props, control-flow children and ref assignments", async () => {
  const [records, set_records] = createSignal(["one", "two"]);
  let input: HTMLInputElement | undefined;
  const root = document.createElement("div"); document.body.append(root);
  render(() => <div>
    <input ref={input} disabled={false} />
    <For each={records()}>
      {(item, index) => <Show when={item !== "hidden"}><span>{index()}:{item}</span></Show>}
    </For>
  </div>, root);
  await new Promise(resolve => setTimeout(resolve, 5));
  expect(input!.disabled).toBe(false); expect(root.textContent?.trim()).toBe("0:one1:two");
  set_records(["three"]); expect(root.textContent?.trim()).toBe("0:three");
});
