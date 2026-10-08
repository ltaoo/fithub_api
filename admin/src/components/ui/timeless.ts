/** Bridge existing business stores to the local Timeless component models. */
import { vm } from "@timeless/timeless";
import { binding, onCleanup, onMount } from "@/timeless";
import { InputCore } from "@/domains/ui/form/input";
import { ButtonCore } from "@/domains/ui/button";
import { DialogCore } from "@/domains/ui/dialog";

export function view_props(props: any, excluded: string[] = []) {
  const attributes: Record<string, any> = {};
  const events: Record<string, any> = {};
  for (const key of Object.keys(props)) {
    if (["store", "children", "class", "style", "ref", ...excluded].includes(key)) continue;
    if (key.startsWith("on")) events[key] = props[key];
    else attributes[key] = binding(() => key.startsWith("aria-") && typeof props[key] === "boolean" ? String(props[key]) : props[key]);
  }
  return { ...events, attributes, class: props.class || "",
    style: binding(() => props.style || {}),
    onMounted(event: any) {
      const element = event.target.get$elm();
      props.ref?.(element);
      let previous = (props.class || "").split(/\s+/).filter(Boolean);
      const classname = binding(() => props.class || "");
      const unsubscribe = classname.subscribe({ onChange(value: string) {
        previous.forEach((name: string) => element.classList.remove(name));
        previous = value.split(/\s+/).filter(Boolean);
        previous.forEach((name: string) => element.classList.add(name));
      } });
      return () => { unsubscribe(); classname.destroy(); };
    },
  };
}

export function input_model(store: InputCore<any>) {
  let syncing = false;
  const model = new vm.InputCore({
    defaultValue: String(store.value ?? ""), placeholder: store.placeholder,
    disabled: store.disabled, allowClear: false,
    onChange(value) { if (!syncing) store.setValue(value); },
    onKeyDown(event) { if (event.key === "Enter") store.handleEnter(); },
    onBlur() { store.handleBlur(); },
  });
  const unsubscribe = store.onStateChange(state => {
    syncing = true;
    try {
      if (model.value !== String(state.value ?? "")) model.setValue(String(state.value ?? ""));
      model.disabled = state.disabled;
      model.setPlaceholder(state.placeholder); model.setLoading(state.loading);
    } finally { syncing = false; }
  });
  onCleanup(unsubscribe);
  const previous_focus = store.focus;
  store.focus = () => model.focus();
  onCleanup(() => { store.focus = previous_focus; });
  onMount(() => { store.setMounted(); if (store.autoFocus) model.focus(); });
  return model;
}

export function button_model(store: ButtonCore<any>, variant = "default", size = "default") {
  const model = new vm.ButtonCore({ variant, size, disabled: store.state.disabled, onClick: () => store.click() });
  const sync = () => {
    model.setLoading(store.state.loading);
    if (store.state.disabled) model.disable(); else model.enable();
  };
  sync(); onCleanup(store.onStateChange(sync));
  return model;
}

export function dialog_model(store: DialogCore) {
  const model = new vm.DialogCore({ ...store.state, footer: false });
  let syncing = false;
  const sync = () => {
    syncing = true;
    try {
      model.closeable = store.closeable;
      model.setTitle(store.title);
      if (store.open && !model.open) model.show(); else if (!store.open && model.open) model.hide();
    } finally { syncing = false; }
  };
  sync(); onCleanup(store.onStateChange(sync));
  onCleanup(model.onCancel(() => { if (!syncing && store.open) store.hide(); }));
  onCleanup(model.onUnmounted(() => store.present.unmount()));
  return model;
}
