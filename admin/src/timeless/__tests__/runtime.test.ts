// @vitest-environment jsdom
import { beforeEach, expect, test } from "vitest";
import { createSignal, dynamic, h, render, For, Show, onMount, onCleanup, lazy } from "../index";

beforeEach(() => document.body.replaceChildren());
const mount = (factory: () => any) => {
  const root = document.createElement("div"); document.body.append(root);
  return { root, dispose: render(factory, root) };
};

test("text and attributes update without duplicate DOM nodes", () => {
  const [value, set_value] = createSignal("first");
  const { root } = mount(() => h("button", { get class() { return value(); }, children: [dynamic(value)] }));
  set_value("second");
  expect(root.textContent).toBe("second");
  expect(root.querySelector("button")!.className).toBe("second");
});
test("For retains, reorders and removes records", async () => {
  const one = { name: "one" }, two = { name: "two" };
  const [records, set_records] = createSignal([one, two]);
  const { root } = mount(() => h(For, { get each() { return records(); }, children: (item: any) => h("span", { children: item.name }) }));
  await new Promise(resolve => setTimeout(resolve, 5));
  set_records([two]);
  expect(root.textContent).toBe("two");
  expect(root.querySelectorAll("span")).toHaveLength(1);
});
test("Show mounts lazily and removes hidden content", async () => {
  const [visible, set_visible] = createSignal(false);
  let creations = 0;
  const { root } = mount(() => h(Show, { get when() { return visible(); }, children: [lazy(() => { creations++; return h("span", { children: "shown" }); })] }));
  expect(creations).toBe(0);
  await new Promise(resolve => setTimeout(resolve, 5));
  set_visible(true); expect(root.textContent).toBe("shown");
  set_visible(false); expect(root.textContent).toBe("");
});
test("refs precede mount callbacks and cleanup executes", async () => {
  let element: HTMLElement | undefined;
  let mounted = false, cleaned = false;
  const { dispose } = mount(() => {
    onMount(() => { mounted = !!element; }); onCleanup(() => { cleaned = true; });
    return h("input", { ref: (node: HTMLElement) => { element = node; } });
  });
  await new Promise(resolve => setTimeout(resolve, 5));
  expect(mounted).toBe(true); dispose(); expect(cleaned).toBe(true);
});
test("boolean false attributes do not disable inputs and properties update", () => {
  const [disabled, set_disabled] = createSignal(false);
  const [value, set_value] = createSignal("initial");
  const { root } = mount(() => h("textarea", { get disabled() { return disabled(); }, get value() { return value(); } }));
  const textarea = root.querySelector("textarea")!;
  expect(textarea.disabled).toBe(false);
  expect(textarea.value).toBe("initial");
  set_value("updated"); expect(textarea.value).toBe("updated");
  set_disabled(true); expect(textarea.disabled).toBe(true);
  set_disabled(false); expect(textarea.disabled).toBe(false);
});
test("For removes rows whose root is conditional when models replace records", async () => {
  const [records, set_records] = createSignal([{ text: "one" }, { text: "two" }]);
  const { root } = mount(() => h(For, { get each() { return records(); }, children: (item: any) => h(Show, { when: true, children: [lazy(() => h("span", { children: item.text }))] }) }));
  await new Promise(resolve => setTimeout(resolve, 5));
  set_records([{ text: "three" }]);
  expect(root.textContent).toBe("three"); expect(root.querySelectorAll("span")).toHaveLength(1);
});
test("inline string styles preserve aspect-ratio layouts", () => {
  const { root } = mount(() => h("div", { style: "position: relative; padding-bottom: 50%;" }));
  const element = root.querySelector("div")!;
  expect(element.style.position).toBe("relative"); expect(element.style.paddingBottom).toBe("50%");
});
