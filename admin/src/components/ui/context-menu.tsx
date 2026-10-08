import { vm } from "@timeless/timeless";
import { ContextMenu as ShadcnContextMenu } from "@timeless/shadcn";
import { JSX, nodes, onCleanup, createSignal, For, h } from "@/timeless";
import { ContextMenuCore } from "@/domains/ui/context-menu";
import { menu_items, menu_accessibility } from "./menu-model";
import { view_props } from "./timeless";

export function ContextMenu(props: { store: ContextMenuCore } & JSX.HTMLAttributes<HTMLElement>): any {
  const store = props.store;
  const [items, set_items] = createSignal(store.state.items);
  onCleanup(store.onStateChange(state => set_items(state.items)));
  const previous_show = store.show;
  const previous_hide = store.hide;
  let active_model: vm.ContextMenuCore | undefined;
  store.show = position => active_model?.show(position);
  store.hide = () => active_model?.hide({ reason: "business" });
  onCleanup(() => { store.show = previous_show; store.hide = previous_hide; });
  function MenuContent(): any {
    const model = new vm.ContextMenuCore({ items: menu_items(items()) });
    active_model = model;
    const forwarded = view_props(props);
    return ShadcnContextMenu({ ...forwarded, store: model, onMounted(event) {
      menu_accessibility(event); return forwarded.onMounted(event);
    } }, nodes(props.children));
  }
  return h(For, { get each() { return [items()]; }, children: () => h(MenuContent, {}) });
}
