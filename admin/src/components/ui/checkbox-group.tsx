import { For, JSX } from "@/timeless";
import { CheckboxCore } from "@/domains/ui/checkbox";
import { CheckboxGroupCore } from "@/domains/ui/checkbox/group";
import { Checkbox } from "./checkbox";
import { createSignal, onCleanup } from "@/timeless";

export function CheckboxOption(props: { label: string; store: CheckboxCore } & JSX.HTMLAttributes<HTMLDivElement>) {
  return <div class={`inline-flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm ${props.class || ""}`}>
    <Checkbox store={props.store} aria-label={props.label} /><span>{props.label}</span>
  </div>;
}
export function CheckboxGroup<T>(props: { store: CheckboxGroupCore<T> } & JSX.HTMLAttributes<HTMLDivElement>) {
  const [state, set_state] = createSignal(props.store.state);
  onCleanup(props.store.onStateChange(set_state));
  return <div class={`flex flex-wrap items-center gap-2 ${props.class || ""}`}>
    <For each={state().options}>{option => <CheckboxOption store={option.core} label={option.label} />}</For>
  </div>;
}
