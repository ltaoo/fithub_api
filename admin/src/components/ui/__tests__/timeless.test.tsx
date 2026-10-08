// @vitest-environment jsdom
import { beforeEach, expect, test, vi } from "vitest";
import { h, render } from "@/timeless";
import { Input } from "../input";
import { Textarea } from "../textarea";
import { Checkbox } from "../checkbox";
import { Button } from "../button";
import { Select } from "../select";
import { Dialog } from "../dialog";
import { DropdownMenu } from "../dropdown-menu";
import { InputCore, ButtonCore, CheckboxCore, SelectCore, DialogCore, DropdownMenuCore, MenuItemCore } from "@/domains/ui";

beforeEach(() => {
  document.body.replaceChildren();
  vi.spyOn(HTMLElement.prototype, "getBoundingClientRect").mockReturnValue({ x: 0, y: 0, left: 0, top: 0, right: 800, bottom: 600, width: 800, height: 600, toJSON() {} });
});
const settle = () => new Promise(resolve => setTimeout(resolve, 150));
function mount(component: any, props: any) {
  const root = document.createElement("div");
  document.body.append(root);
  const dispose = render(() => h(component, props), root);
  return { root, dispose };
}

test("input and textarea synchronize both directions and preserve password type", async () => {
  const store = new InputCore({ defaultValue: "one", type: "password" });
  const { root, dispose } = mount(Input, { store, id: "password" });
  await settle();
  const input = root.querySelector("input")!;
  expect(input.type).toBe("password");
  expect(input.id).toBe("password");
  expect(input.value).toBe("one");
  store.setValue("two");
  expect(input.value).toBe("two");
  input.value = "three";
  input.dispatchEvent(new Event("input", { bubbles: true }));
  expect(store.value).toBe("three");
  dispose();
  const textarea_store = new InputCore({ defaultValue: "text" });
  const textarea_view = mount(Textarea, { store: textarea_store });
  await settle();
  const textarea = textarea_view.root.querySelector("textarea")!;
  textarea.value = "edited";
  textarea.dispatchEvent(new Event("input", { bubbles: true }));
  expect(textarea_store.value).toBe("edited");
  textarea_store.setValue("updated");
  expect(textarea.value).toBe("updated");
  textarea_view.dispose();
});

test("button reflects loading and invokes the original business callback once", async () => {
  const click = vi.fn();
  const store = new ButtonCore({ onClick: click });
  const { root, dispose } = mount(Button, { store, children: "保存" });
  await settle();
  const button = root.querySelector("button")!;
  button.click();
  expect(click).toHaveBeenCalledTimes(1);
  store.setLoading(true);
  expect(button.className).toContain("opacity-70");
  button.click();
  expect(click).toHaveBeenCalledTimes(1);
  dispose();
});

test("checkbox toggles once, reflects programmatic changes and respects disabled", async () => {
  const change = vi.fn();
  const store = new CheckboxCore({ onChange: change });
  const { root, dispose } = mount(Checkbox, { store });
  await settle();
  const button = root.querySelector("button")!;
  button.click();
  expect(store.checked).toBe(true);
  expect(change).toHaveBeenCalledTimes(1);
  expect(button.getAttribute("aria-checked")).toBe("true");
  store.uncheck();
  expect(button.getAttribute("aria-checked")).toBe("false");
  dispose();
  const disabled_store = new CheckboxCore({ disabled: true });
  const disabled_view = mount(Checkbox, { store: disabled_store });
  await settle();
  disabled_view.root.querySelector("button")!.click();
  expect(disabled_store.checked).toBe(false);
  disabled_view.dispose();
});

test("select renders complete dynamic options and preserves numeric values", async () => {
  const store = new SelectCore({ defaultValue: 1, options: [{ value: 1, label: "一" }, { value: 2, label: "二" }] });
  const { root, dispose } = mount(Select, { store });
  await settle();
  root.querySelector("button")!.dispatchEvent(new MouseEvent("pointerdown", { bubbles: true }));
  await settle();
  const option = [...document.querySelectorAll<HTMLElement>('[role="option"]')].find(item => item.textContent === "二");
  expect(option).toBeDefined();
  option!.click();
  expect(store.value).toBe(2);
  store.setOptions([{ value: 3, label: "三" }]);
  store.setValue(3);
  await settle();
  expect(root.textContent).toContain("三");
  dispose();
});

test("dialog retains business submit/loading/cancel and honors closeable", async () => {
  const submit = vi.fn();
  const cancel = vi.fn();
  const store = new DialogCore({ title: "编辑", onOk: submit, onCancel: cancel });
  const { dispose } = mount(Dialog, { store, children: "内容" });
  await settle();
  store.show();
  await settle();
  const buttons = [...document.querySelectorAll<HTMLButtonElement>("button")];
  const confirm = buttons.find(button => button.textContent?.includes("确认"));
  expect(confirm).toBeDefined();
  confirm!.click();
  expect(submit).toHaveBeenCalledTimes(1);
  store.okBtn.setLoading(true);
  expect(confirm!.className).toContain("opacity-70");
  buttons.find(button => button.textContent === "取消")!.click();
  await settle();
  expect(store.open).toBe(false);
  expect(cancel).toHaveBeenCalledTimes(1);
  dispose();
});

test("dropdown menu invokes business items", async () => {
  const click = vi.fn();
  const store = new DropdownMenuCore({ items: [new MenuItemCore({ label: "编辑", onClick: click })] });
  const { root, dispose } = mount(DropdownMenu, { store, children: h("button", { children: "菜单" }) });
  await settle();
  root.querySelector("button")!.dispatchEvent(new MouseEvent("pointerdown", { bubbles: true }));
  await settle();
  const item = [...document.querySelectorAll<HTMLElement>('[role="menuitem"]')].find(item => item.textContent?.includes("编辑"));
  expect(item).toBeDefined();
  item!.click();
  expect(click).toHaveBeenCalledTimes(1);
  dispose();
});

test("button retains reactive navigation classes and attributes", async () => {
  const { createSignal } = await import("@/timeless");
  const [active, set_active] = createSignal(false);
  const { root, dispose } = mount(Button, { store: new ButtonCore(), get class() { return active() ? "bg-accent" : "text-muted-foreground"; }, get "aria-current"() { return active() ? "page" : undefined; } });
  await settle();
  const button = root.querySelector("button")!;
  expect(button.className).not.toContain("[object Object]");
  set_active(true);
  expect(button.className).toContain("bg-accent");
  expect(button.className).not.toContain("text-muted-foreground");
  expect(button.getAttribute("aria-current")).toBe("page");
  dispose();
});

test("menu options can be replaced and disabled items cannot run business actions", async () => {
  const click = vi.fn();
  const store = new DropdownMenuCore({ items: [new MenuItemCore({ label: "旧操作" })] });
  const { root, dispose } = mount(DropdownMenu, { store, children: h("button", { children: "菜单" }) });
  await settle();
  store.setItems([new MenuItemCore({ label: "新操作", onClick: click, disabled: true })]);
  await settle();
  root.querySelector("button")!.dispatchEvent(new MouseEvent("pointerdown", { bubbles: true }));
  await settle();
  const item = [...document.querySelectorAll<HTMLElement>('[role="menuitem"]')].find(item => item.textContent === "新操作");
  expect(item).toBeDefined();
  item!.click();
  expect(click).not.toHaveBeenCalled();
  dispose();
});

test("date picker reflects business values and later programmatic edits", async () => {
  const { DatePicker } = await import("../date-picker");
  const { DatePickerCore } = await import("@/domains/ui/date-picker");
  const store = DatePickerCore({ today: new Date(2026, 9, 8) });
  store.setValue(new Date(2026, 9, 8));
  const { root, dispose } = mount(DatePicker, { store });
  await settle();
  expect(root.textContent).toContain("2026/10/08");
  store.setValue(new Date(2026, 9, 9));
  expect(root.textContent).toContain("2026/10/09");
  dispose();
});

test("non-closeable dialogs hide the close control and can reopen after cancellation", async () => {
  const store = new DialogCore({ title: "保存", closeable: false });
  const { dispose } = mount(Dialog, { store });
  await settle(); store.show(); await settle();
  expect(document.querySelector(".admin-dialog-locked")).not.toBeNull();
  store.hide(); await settle(); await settle();
  store.show(); await settle();
  expect(store.open).toBe(true);
  expect(document.querySelectorAll(".admin-dialog-locked")).toHaveLength(1);
  dispose();
});
