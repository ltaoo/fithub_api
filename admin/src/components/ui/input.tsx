import { View } from "@timeless/timeless";
import { Input as ShadcnInput } from "@timeless/shadcn";
import { JSX, nodes } from "@/timeless";
import { InputCore } from "@/domains/ui/form/input";
import { Input as FileInput } from "@/packages/ui/input";
import { input_model, view_props } from "./timeless";

export function Input(props: { store: InputCore<any>; prefix?: JSX.Element; class?: string; id?: string }): any {
  const store = props.store;
  if (store.type === "file") return FileInput({ store, class: props.class });
  const forwarded = view_props(props, ["prefix"]);
  const model = input_model(store);
  const input = ShadcnInput({ ...forwarded, store: model, onMounted(event) {
    model.focus = () => event.target.get$elm().focus();
    return forwarded.onMounted(event);
  }, class: props.prefix ? "pl-8" : "", attributes: {
    ...forwarded.attributes, type: store.type === "string" ? "text" : store.type,
    autocomplete: store.autoComplete ? "on" : "off",
  } });
  return View({ class: forwarded.class }, [
    View({ class: "relative w-full" }, [
      ...(props.prefix ? [View({ class: "pointer-events-none absolute left-3 top-1/2 -translate-y-1/2" }, nodes(props.prefix))] : []), input,
    ]),
  ]);
}
