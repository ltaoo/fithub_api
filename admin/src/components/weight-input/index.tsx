import { createSignal, For } from "@/timeless";

import { useViewModel, useViewModelStore } from "@/hooks";
import { Input } from "@/components/ui/input";
import * as PopoverPrimitive from "@/packages/ui/popover";
import { WeightInputViewModel } from "@/biz/weight_input";
import { base, Handler } from "@/domains/base";
import { InputCore, PopoverCore } from "@/domains/ui";

export function WeightInput(props: { store: WeightInputViewModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <>
      <div class="flex flex-col gap-4 p-4 max-w-[320px]">
        <div class="flex items-center justify-between px-2">
          <span class="text-3xl font-bold">{state().text}</span>
          <div class="flex gap-2">
            <button
              class="px-4 py-2 rounded-lg transition-colors"
              classList={{
                "bg-blue-500 text-white": state().unit === "kg",
                "bg-gray-100 hover:bg-gray-200": state().unit !== "kg",
              }}
              onClick={() => {
                vm.methods.handleClickUnit("kg");
              }}
            >
              kg
            </button>
            <button
              class="px-4 py-2 rounded-lg transition-colors"
              classList={{
                "bg-blue-500 text-white": state().unit === "lbs",
                "bg-gray-100 hover:bg-gray-200": state().unit !== "lbs",
              }}
              onClick={() => {
                vm.methods.handleClickUnit("lbs");
              }}
            >
              lbs
            </button>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-2">
          <For each={[1, 2, 3, 4, 5, 6, 7, 8, 9]}>
            {(num) => {
              return (
                <button
                  class="aspect-square flex items-center justify-center text-xl font-medium bg-white hover:bg-gray-50 active:bg-gray-100 border rounded-xl shadow-sm transition-colors"
                  onClick={() => {
                    vm.methods.handleClickNumber(num.toString());
                  }}
                >
                  {num}
                </button>
              );
            }}
          </For>
          <button
            class="aspect-square flex items-center justify-center text-xl font-medium bg-white hover:bg-gray-50 active:bg-gray-100 border rounded-xl shadow-sm transition-colors"
            onClick={() => {
              vm.methods.handleClickDot();
            }}
          >
            .
          </button>
          <button
            class="aspect-square flex items-center justify-center text-xl font-medium bg-white hover:bg-gray-50 active:bg-gray-100 border rounded-xl shadow-sm transition-colors"
            onClick={() => {
              vm.methods.handleClickNumber("0");
            }}
          >
            0
          </button>
          <button
            class="aspect-square flex items-center justify-center text-xl font-medium bg-white hover:bg-gray-50 active:bg-gray-100 border rounded-xl shadow-sm transition-colors"
            onClick={() => {
              vm.methods.handleClickDelete();
            }}
          >
            ←
          </button>
        </div>
      </div>
    </>
  );
}
