import { JSX } from "@/timeless";
import { createSignal, For, Show } from "@/timeless";
import { Plus, Trash } from "@/timeless/icons";

import { useViewModelStore } from "@/hooks";

import { ArrayFieldCore, SingleFieldCore } from "@/domains/ui/formv2";

import { FieldV2 } from "./fieldv2";

export function FieldArrV2<T extends (v: number) => any>(
  props: {
    store: ArrayFieldCore<T>;
    render: (field: ReturnType<T>) => JSX.Element;
  } & JSX.HTMLAttributes<HTMLDivElement>
) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <Show when={!state().hidden}>
      <div class="header flex items-center gap-4">
        <div class="field">
          <div class="field__label">{state().label}</div>
        </div>
        <div class="flex gap-2">
          <div
            class="p-2 rounded-full bg-w-bg-5 cursor-pointer"
            onClick={() => {
              vm.append();
            }}
          >
            <Plus class="w-4 h-4 text-w-fg-0" />
          </div>
        </div>
      </div>
      <div class="mt-4">
        <For each={state().fields}>
          {(field, idx) => {
            const store = props.store.mapFieldWithIndex(idx());
            if (store === null) {
              return null;
            }
            return (
              <div>
                <div class="operations flex items-center gap-2">
                  <div
                    class="inline-block p-2 rounded-full bg-w-bg-5"
                    onClick={() => {
                      vm.remove(idx());
                    }}
                  >
                    <Trash class="w-4 h-4 text-w-fg-0" />
                  </div>
                </div>
                <div class="mt-2">{props.render(store.field)}</div>
              </div>
            );
          }}
        </For>
      </div>
    </Show>
  );
}
