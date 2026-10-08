import { vm } from "@timeless/timeless";
import { onCleanup } from "@/timeless";
import { MenuItemCore } from "@/domains/ui/menu/item";
import { MenuCore } from "@/domains/ui/menu";

export function menu_items(items: MenuItemCore[]): vm.MenuItemCore[] {
  return items.filter(item => !item._hidden).map(item => {
    const model = new vm.MenuItemCore({ label: item.label, icon: item.icon, tooltip: item.tooltip,
      shortcut: item.shortcut, disabled: item.state.disabled,
      menu: item.menu ? menu_model(item.menu) : undefined,
      onClick() { item.handleClick(); },
    });
    onCleanup(item.onStateChange(state => {
      model.label = state.label; model.icon = state.icon; model.shortcut = state.shortcut;
      if (state.disabled) model.disable(); else model.enable();
    }));
    return model;
  });
}
function menu_model(store: MenuCore): vm.MenuCore {
  const model = new vm.MenuCore({ items: menu_items(store.state.items) });
  onCleanup(store.onChange(state => model.setItems(menu_items(state.items))));
  return model;
}

export function menu_accessibility(event: any) {
  const element = event.target.get$elm();
  element.setAttribute("role", "menu");
  setTimeout(() => element.querySelectorAll('[class*="group/menubar-item"]').forEach((item: HTMLElement) => {
    item.setAttribute("role", "menuitem"); item.setAttribute("tabindex", "-1");
  }), 0);
}
