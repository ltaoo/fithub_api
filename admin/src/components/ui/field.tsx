import { JSX, createSignal, Show, onCleanup } from "@/timeless";
import { FormFieldCore } from "@/domains/ui/form/field";
import { Label } from "./label";

export function Field(props: { store: FormFieldCore<any>; extra?: JSX.Element } & JSX.HTMLAttributes<HTMLDivElement>) {
  const [state, set_state] = createSignal(props.store.state);
  onCleanup(props.store.onStateChange(set_state));
  return <Show when={!state().hidden}>
    <div class={`space-y-2 py-2 ${props.class || ""}`}>
      <div class="flex items-center justify-between gap-2"><Label>{props.store.label}</Label>{props.extra}</div>
      <div>{props.children}</div>
    </div>
  </Show>;
}
