import { For, Show } from "@/timeless";

import { useViewModelStore } from "@/hooks";
import { Dialog } from "@/components/ui";

import { ContentSelectViewModel } from "@/biz/content/content_select";
import { Check } from "@/timeless/icons";

export function ContentSelect(props: { store: ContentSelectViewModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <>
      <div>
        <Show
          when={state().value.length}
          fallback={
            <div
              onClick={() => {
                vm.init();
                vm.ui.$dialog.show();
              }}
            >
              点击选择
            </div>
          }
        >
          <For each={state().value}>
            {(v) => {
              return <div>{v.title}</div>;
            }}
          </For>
          <div
            onClick={() => {
              vm.init();
              vm.ui.$dialog.show();
            }}
          >
            切换
          </div>
        </Show>
      </div>
      <Dialog store={vm.ui.$dialog}>
        <For each={state().list}>
          {(v) => {
            return (
              <div
                class="relative p-4"
                onClick={() => {
                  vm.select(v);
                }}
              >
                <div>{v.title}</div>
                <Show when={v.selected}>
                  <div class="absolute">
                    <Check class="w-4 h-4 text-w-fg-0" />
                  </div>
                </Show>
              </div>
            );
          }}
        </For>
      </Dialog>
    </>
  );
}
