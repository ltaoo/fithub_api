import { vm } from "@timeless/timeless";
import { Popover as ShadcnPopover } from "@timeless/shadcn";
import { JSX, nodes, onCleanup } from "@/timeless";
import { PopoverCore } from "@/domains/ui/popover";
import { Align, Side } from "@/domains/ui/popper";
import { view_props } from "./timeless";

export function Popover(props: { store: PopoverCore; content: JSX.Element } & JSX.HTMLAttributes<HTMLElement>): any {
  const store = props.store;
  const model = new vm.PopoverCore({ side: store._side, align: store._align });
  let syncing = false;
  onCleanup(store.onShow(() => { syncing = true; model.show(); syncing = false; }));
  onCleanup(store.onHide(() => { syncing = true; model.hide(); syncing = false; }));
  onCleanup(model.onShow(() => { if (!syncing && !store.visible) { syncing = true; store.show(); syncing = false; } }));
  onCleanup(model.onHide(() => { if (!syncing && store.visible) { syncing = true; store.hide(); syncing = false; } }));
  if (store.visible) model.show();
  return ShadcnPopover({ ...view_props(props, ["content"]), store: model, content: nodes(props.content) }, nodes(props.children));
}
export function PurePopover(props: { content: JSX.Element; side?: Side; align?: Align } & JSX.HTMLAttributes<HTMLElement>): any {
  return Popover({ ...props, store: new PopoverCore({ side: props.side, align: props.align }) });
}
