import { SimpleSelectCore } from "@/domains/ui/simple-select";
import { SelectCore } from "@/domains/ui/select";
import { onCleanup, h } from "@/timeless";
import { Select } from "./select";

export function SimpleSelect(props: { store: SimpleSelectCore }): any {
  const store = props.store;
  const model = new SelectCore({ defaultValue: store.value, options: store.options, onChange(value) { if (value !== null) store.select(value); } });
  onCleanup(store.onStateChange(state => { model.setOptions(state.options); model.setValue(state.value); }));
  return h(Select, { store: model });
}
