/**
 * @file 肌肉多选
 */
import { For, Show } from "@/timeless";
import { X, Check } from "@/timeless/icons";

import { useViewModelStore } from "@/hooks";
import { Dialog } from "@/components/ui";

import { MuscleSelectViewModel } from "@/biz/muscle/muscle_select";
import { cn } from "@/utils/index";

export function MuscleSelectView(props: { store: MuscleSelectViewModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <div class="muscle-select">
      <div class="flex flex-wrap items-center gap-2">
        <For each={state().value}>
          {(item) => (
            <div class="flex items-center gap-2 bg-slate-100 rounded-md p-1">
              <span class="text-sm whitespace-nowrap text-slate-400">{item.zh_name || item.name}</span>
              <button
                class="text-slate-400"
                onClick={() => {
                  vm.remove(item);
                }}
              >
                <X class="w-4 h-4" />
              </button>
            </div>
          )}
        </For>
        <button
          classList={{
            "text-slate-400 whitespace-nowrap": true,
            "cursor-not-allowed": state().disabled,
          }}
          onClick={(e) => {
            if (state().disabled) {
              return;
            }
            vm.init();
            vm.ui.$dialog.show();
          }}
        >
          选择肌肉
        </button>
      </div>
      <Dialog store={vm.ui.$dialog}>
        <div class="grid grid-cols-8 gap-2">
          <For each={state().list}>
            {(v) => {
              return (
                <div
                  classList={{
                    "relative flex cursor-default rounded-lg border-2 border-w-fg-5 select-none items-center rounded-sm p-2 text-sm outline-none cursor-pointer data-[highlighted]:bg-slate-100 focus:bg-accent focus:text-accent-foreground hover:bg-w-bg-5 data-[disabled]:pointer-events-none data-[disabled]:opacity-50":
                      true,
                    "border-w-fg-0": v.selected,
                  }}
                  onClick={() => {
                    vm.select(v);
                  }}
                >
                  <span class="absolute right-2 flex h-3.5 w-3.5 items-center justify-center">
                    <Show when={v.selected}>
                      <Check class="h-4 w-4" />
                    </Show>
                  </span>
                  <div class="py-2">{v.zh_name || v.name}</div>
                </div>
              );
            }}
          </For>
        </div>
      </Dialog>
    </div>
  );
}
