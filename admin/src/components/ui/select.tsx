/** Select composed from Timeless primitives until the library Select renders options. */
import { ui, vm, refobj, computed, For, View, Icon } from "@timeless/timeless";
import { JSX, onCleanup } from "@/timeless";
import { SelectCore } from "@/domains/ui/select";
import { view_props } from "./timeless";

export function Select(props: { store: SelectCore<any>; position?: "popper" } & JSX.HTMLAttributes<HTMLElement>): any {
  const store = props.store;
  let syncing = false;
  let previous_options = store.options;
  const options = () => store.options.map(option => new vm.SelectItemCore({ value: option.value, label: option.label }));
  const model = new vm.SelectCore({ defaultValue: store.value, placeholder: store.placeholder,
    disabled: store.disabled, options: options(), position: "popper",
    onChange(value) { if (!syncing) store.select(value); },
  });
  onCleanup(store.onStateChange(state => {
    syncing = true;
    try {
      model.disabled = state.disabled;
      model.placeholder = state.placeholder;
      if (previous_options !== store.options) { previous_options = store.options; model.setOptions(options()); }
      model.setValue(state.value);
    } finally { syncing = false; }
  }));
  const state = refobj(model.state);
  onCleanup(model.onStateChange(next => state.as(next)));
  const forwarded = view_props(props, ["position"]);
  return ui.SelectPrimitive.Root({ store: model }, [
    ui.SelectPrimitive.Trigger({ ...forwarded, store: model, class: computed(state, s =>
      `flex h-8 w-full items-center justify-between gap-2 rounded-lg border border-input bg-background px-2.5 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring ${s.disabled ? "opacity-50 cursor-not-allowed" : ""} ${props.class || ""}`) }, [
      ui.SelectPrimitive.Value({ store: model, class: "truncate" }), Icon({ name: "chevron-down", size: 16 }),
    ]),
    ui.SelectPrimitive.Content({ store: model, attributes: { role: "listbox" },
      class: "z-[400] min-w-36 rounded-lg border border-border bg-popover text-popover-foreground shadow-md outline-none",
    }, [ui.SelectPrimitive.Viewport({ store: model, class: "max-h-72 p-1" }, [
      For({ each: computed(state, s => s.options), render(option: any) {
        const item_state = refobj(option.state);
        return ui.SelectPrimitive.Item({ select$: model, item$: option, attributes: { role: "option", "aria-selected": computed(item_state, s => s.selected) },
          class: computed(item_state, s => `flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm outline-none hover:bg-accent focus:bg-accent ${s.focused ? "bg-accent" : ""}`),
          onMounted() { return option.onStateChange((next: any) => item_state.as(next)); },
        }, [ui.SelectPrimitive.ItemIndicator({ store: option, class: "w-4" }, [Icon({ name: "check", size: 14 })]), ui.SelectPrimitive.ItemText({}, [option.label])]);
      } }),
      View({ class: computed(state, s => s.options.length ? "hidden" : "p-4 text-center text-sm text-muted-foreground") }, ["暂无选项"]),
    ])]),
  ]);
}
