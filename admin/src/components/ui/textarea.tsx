import { Textarea as ShadcnTextarea } from "@timeless/shadcn";
import { View } from "@timeless/timeless";
import { JSX } from "@/timeless";
import { InputCore } from "@/domains/ui/form/input";
import { input_model, view_props } from "./timeless";

export function Textarea(props: { store: InputCore<string> } & JSX.HTMLAttributes<HTMLTextAreaElement>): any {
  const forwarded = view_props(props);
  const model = input_model(props.store);
  return View({ class: forwarded.class }, [ShadcnTextarea({ ...forwarded, store: model, onMounted(event) {
    model.focus = () => event.target.get$elm().focus();
    return forwarded.onMounted(event);
  } })]);
}
