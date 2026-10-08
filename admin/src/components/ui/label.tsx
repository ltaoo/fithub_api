import { Label as ShadcnLabel } from "@timeless/shadcn";
import { JSX, nodes } from "@/timeless";
import { view_props } from "./timeless";

export function Label(props: JSX.HTMLAttributes<HTMLDivElement> & JSX.AriaAttributes & { for?: string }): any {
  return ShadcnLabel(view_props(props), nodes(props.children));
}
