import { Button as ShadcnButton } from "@timeless/shadcn";
import { JSX, nodes } from "@/timeless";
import { ButtonCore } from "@/domains/ui/button";
import { button_model, view_props } from "./timeless";

export function Button<T = unknown>(props: {
  store: ButtonCore<T>; icon?: JSX.Element;
  variant?: "default" | "destructive" | "outline" | "subtle" | "ghost" | "link" | null;
  size?: "default" | "sm" | "lg" | null;
} & JSX.HTMLAttributes<HTMLButtonElement>): any {
  const model = button_model(props.store, props.variant === "subtle" ? "secondary" : props.variant || "default", props.size || "default");
  const forwarded = view_props(props, ["icon", "variant", "size"]);
  return ShadcnButton({ ...forwarded, store: model, prefix: nodes(props.icon), attributes: { ...forwarded.attributes, type: "button" } }, nodes(props.children));
}
