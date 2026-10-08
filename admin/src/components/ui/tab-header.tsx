import { Tabs as ShadcnTabs } from "@timeless/shadcn";
import { vm } from "@timeless/timeless";
import { onCleanup } from "@/timeless";
import { TabHeaderCore } from "@/domains/ui/tab-header";

export function TabHeader(props: { store: TabHeaderCore<any> }): any {
  const store = props.store;
  const visible_tabs = () => store.tabs.filter((tab: any) => !tab.hidden);
  const model = new vm.TabHeaderCore({ options: visible_tabs(), onChange(tab) {
    const index = store.tabs.findIndex((item: any) => item.id === tab.id);
    if (index >= 0 && index !== store.current) store.select(index);
  } });
  model.selectById(store.selectedTabId);
  onCleanup(store.onStateChange(() => { model.selectById(store.selectedTabId); }));
  return ShadcnTabs({ store: model, items: visible_tabs().map((tab: any) => ({ value: tab.id, label: tab.text })), class: "overflow-x-auto" });
}
