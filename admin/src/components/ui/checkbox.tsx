import { vm, refobj, computed } from "@timeless/timeless";
import { Checkbox as ShadcnCheckbox } from "@timeless/shadcn";
import { JSX, onCleanup } from "@/timeless";
import { CheckboxCore } from "@/domains/ui/checkbox";
import { view_props } from "./timeless";

export function Checkbox(props: { store: CheckboxCore } & JSX.HTMLAttributes<HTMLDivElement>): any {
  const store = props.store;
  let syncing = false;
  const model = new vm.CheckboxCore({ ...store.state,
    onChange(value) { if (!syncing && !store.disabled && value !== store.checked) store.toggle(); },
  });
  onCleanup(store.onStateChange(state => {
    syncing = true;
    try { model.disabled = state.disabled; model.setValue(!!state.checked); model.setStatus(model.status); }
    finally { syncing = false; }
  }));
  const state = refobj(model.state);
  onCleanup(model.onStateChange(next => state.as(next)));
  const forwarded = view_props(props);
  return ShadcnCheckbox({ ...forwarded, id: props.id, store: model, attributes: {
    ...forwarded.attributes, type: "button", role: "checkbox",
    "aria-label": store.label || props["aria-label"],
    "aria-checked": computed(state, s => String(!!s.checked)),
    disabled: computed(state, s => s.disabled ? "" : undefined),
  } });
}
