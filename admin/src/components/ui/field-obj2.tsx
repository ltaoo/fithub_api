import { JSX, Show } from "@/timeless";
import { useViewModelStore } from "@/hooks";
import { ObjectFieldCore } from "@/domains/ui/formv2";
import { FieldSet, FieldLegend, FieldGroup } from "./field-layout";

export function FieldObjV2(props: { store: ObjectFieldCore<any> } & JSX.HTMLAttributes<HTMLDivElement>) {
  const [state] = useViewModelStore(props.store);
  return <Show when={!state().hidden}>
    <FieldSet class={props.class}><Show when={state().label}><FieldLegend>{state().label}</FieldLegend></Show><FieldGroup>{props.children}</FieldGroup></FieldSet>
  </Show>;
}
