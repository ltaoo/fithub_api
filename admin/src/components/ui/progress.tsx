import { Progress as ShadcnProgress } from "@timeless/shadcn";
import { ref } from "@timeless/timeless";
import { JSX, onCleanup } from "@/timeless";
import { ProgressCore } from "@/domains/ui/progress";
import { view_props } from "./timeless";

export function Progress(props: { store: ProgressCore } & JSX.HTMLAttributes<HTMLElement>): any {
  const value = ref(props.store.state.value || 0);
  onCleanup(props.store.onStateChange(state => value.as(state.value || 0)));
  return ShadcnProgress({ ...view_props(props), value, max: props.store.state.max });
}
