import { vm } from "@timeless/timeless";
import { DropdownMenu as ShadcnDropdownMenu } from "@timeless/shadcn";
import { JSX, nodes, onCleanup, createSignal, For, h } from "@/timeless";
import { DropdownMenuCore } from "@/domains/ui/dropdown-menu";
import { menu_items, menu_accessibility } from "./menu-model";
import { view_props } from "./timeless";

export function DropdownMenu(props: { store: DropdownMenuCore } & JSX.HTMLAttributes<HTMLElement>): any {
  const store = props.store;
  const [items, set_items] = createSignal(store.items);
  onCleanup(store.onStateChange(state => set_items(state.items)));
  const previous_toggle = store.toggle;
  const previous_hide = store.hide;
  let active_model: vm.DropdownMenuCore | undefined;
  store.toggle = position => active_model?.toggle(position);
  store.hide = () => active_model?.hide();
  onCleanup(() => { store.toggle = previous_toggle; store.hide = previous_hide; });
  function MenuContent(): any {
    const model = new vm.DropdownMenuCore({ trigger: "click", items: menu_items(items()),
      side: store.menu.popper.state.placedSide, align: store.menu.popper.state.placedAlign,
      onHidden() { store.menu.hide(); },
    });
    active_model = model;
    onCleanup(model.menu.onShow(() => { store.menu.state.open = true; }));
    onCleanup(store.menu.onShow(() => model.menu.show()));
    onCleanup(store.menu.onHide(() => model.hide()));
    const forwarded = view_props(props);
    return ShadcnDropdownMenu({ ...forwarded, store: model, onMounted(event) {
      menu_accessibility(event); return forwarded.onMounted(event);
    } }, nodes(props.children));
  }
  return h(For, { get each() { return [items()]; }, children: () => h(MenuContent, {}) });
}
