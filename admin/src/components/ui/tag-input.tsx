import { createSignal, For } from "@/timeless";
import { X } from "@/timeless/icons";

import { TagInputCore } from "@/domains/ui/form/tag-input";

export function TagInput(props: { store: TagInputCore }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);
  store.onStateChange((v) => setState(v));

  return (
    <div class="flex flex-col gap-2 border border-w-fg-2 rounded-lg bg-w-bg-1">
      <div class="flex flex-wrap gap-2 p-2 min-h-10 items-center">
        <For each={state().value}>
          {(tag, index) => {
            return (
              <div class="flex items-center gap-1 bg-w-bg-5 px-2 py-1 rounded-md">
                <span>{tag}</span>
                <button
                  type="button"
                  class="text-w-fg-1 hover:text-w-fg-0 cursor-pointer"
                  onClick={() => {
                    store.removeTag(tag);
                  }}
                >
                  <X size={14} />
                </button>
              </div>
            );
          }}
        </For>
        <input
          type="text"
          value={state().inputValue}
          class="flex-grow outline-none min-w-20 bg-transparent"
          placeholder={"输入标签并按回车添加..."}
          onInput={(e) => {
            store.input(e.target.value);
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              store.addTag();
            }
          }}
        />
      </div>
    </div>
  );
}
