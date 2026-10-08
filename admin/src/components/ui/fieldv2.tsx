import { JSX, Show } from "@/timeless";
import { useViewModelStore } from "@/hooks";
import { SingleFieldCore } from "@/domains/ui/formv2";
import { Label } from "./label";
import { FieldDescription } from "./field-layout";

export function FieldV2(props: { store: SingleFieldCore<any> } & JSX.HTMLAttributes<HTMLDivElement>) {
  const [state] = useViewModelStore(props.store);
  return <Show when={!state().hidden}>
    <div class={`space-y-2 py-2 ${props.class || ""}`}>
      <Label>{state().label}</Label>
      <div>{props.children}</div>
      <Show when={state().error}><FieldDescription class="text-destructive">{state().error?.message}</FieldDescription></Show>
    </div>
  </Show>;
}
