import { Dialog as ShadcnDialog } from "@timeless/shadcn";
import { View } from "@timeless/timeless";
import { JSX, h, nodes, createSignal, Show, onCleanup } from "@/timeless";
import { DialogCore } from "@/domains/ui/dialog";
import { Button } from "./button";
import { dialog_model, view_props } from "./timeless";

export function Dialog(props: { store: DialogCore; width?: string } & JSX.HTMLAttributes<HTMLElement>): any {
  const store = props.store;
  const [state, set_state] = createSignal(store.state);
  onCleanup(store.onStateChange(set_state));
  const forwarded = view_props({ ...props, get class() { return `${props.width || "w-fit max-w-[calc(100vw-2rem)]"} ${props.class || ""} admin-dialog ${state().closeable ? "" : "admin-dialog-locked"}`; } }, ["width"]);
  return ShadcnDialog({ ...forwarded, store: dialog_model(store) }, () => [
    View({ class: "p-4 overflow-y-auto min-h-0" }, nodes(props.children)),
    h(Show, { get when() { return state().footer; }, children: () => View({ class: "flex justify-end gap-2 border-t border-border bg-muted/50 p-4" }, [
      h(Button, { store: store.cancelBtn, variant: "outline", children: "取消" }),
      h(Button, { store: store.okBtn, children: "确认" }),
    ]) }),
  ]);
}
