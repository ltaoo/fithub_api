import { Input, Textarea, Checkbox, Select } from "@/components/ui";
import { For, Switch, Match, createSignal } from "@/timeless";
import { Trash, ArrowUp, ArrowDown } from "@/timeless/icons";

import { TagInput } from "@/components/ui/tag-input";
import { EquipmentSelectView } from "@/components/equipment-select";
import { MuscleSelectView } from "@/components/muscle-select";

import { ObjectFieldCore, ArrayFieldCore, SingleFieldCore } from "@/domains/ui/formv2";

function WorkoutPlanSetInputView(props: { store: SingleFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full">
      {(() => {
        if (state().input?.type === "textarea") {
          return (
            <Textarea
              class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              store={store.input}
            />
          );
        }
        if (state().input?.type === "equipment_select") {
          return <EquipmentSelectView store={store.input} />;
        }
        if (state().input?.type === "muscle_select") {
          return <MuscleSelectView store={store.input} />;
        }
        if (state().input?.shape === "tag-input") {
          return <TagInput store={store.input} />;
        }
        if (state().input?.shape === "input") {
          return (
            <Input
              class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              store={store.input}
            />
          );
        }
        if (state().input?.shape === "number") {
          return (
            <Input
              class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              store={store.input}
            />
          );
        }
        if (state().input?.shape === "checkbox") {
          return (
            <div class="flex items-center">
              <Checkbox store={store.input} aria-label={state().label} />
              <span class="ml-2 text-foreground">{state().label}</span>
            </div>
          );
        }
        if (state().input?.shape === "select") {
          return <Select store={store.input} />;
        }
        return null;
      })()}
    </div>
  );
}

function WorkoutPlanSetArrView(props: { store: ArrayFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full space-y-3 my-2">
      <div
        class="inline-flex items-center justify-center py-2 px-4 bg-primary text-primary-foreground rounded-md hover:opacity-90 transition-colors cursor-pointer text-sm font-medium"
        onClick={() => store.append()}
      >
        添加
      </div>
      <div class="space-y-3">
        <For each={state().fields}>
          {(field, index) => {
            const $inner = store.mapFieldWithIndex(index());
            if (!$inner) {
              return null;
            }
            return (
              <div class="p-3 border border-border rounded-lg shadow-sm">
                <div class="flex justify-between items-center mb-2">
                  <div class="text-sm font-medium text-muted-foreground">
                    {field.label} {index() + 1}
                  </div>
                  <div class="flex items-center">
                    <div
                      onClick={() => {
                        store.insertBefore(field.id);
                      }}
                    >
                      在前面插入
                    </div>
                    <div
                      onClick={() => {
                        store.insertAfter(field.id);
                      }}
                    >
                      在后面插入
                    </div>
                    <div
                      class="text-red-500 hover:text-red-700 cursor-pointer p-1"
                      onClick={() => {
                        store.remove(field.id);
                      }}
                    >
                      <Trash class="w-4 h-4" />
                    </div>
                    <div
                      class="text-blue-500 hover:text-blue-700 cursor-pointer p-1"
                      onClick={() => {
                        store.upIdx(field.id);
                      }}
                    >
                      <ArrowUp class="w-4 h-4" />
                    </div>
                    <div
                      class="text-blue-500 hover:text-blue-700 cursor-pointer p-1"
                      onClick={() => {
                        store.downIdx(field.id);
                      }}
                    >
                      <ArrowDown class="w-4 h-4" />
                    </div>
                  </div>
                </div>
                <Switch>
                  <Match when={$inner.field.symbol === "SingleFieldCore"}>
                    <WorkoutPlanSetInputView store={$inner.field} />
                  </Match>
                  <Match when={$inner.field.symbol === "ArrayFieldCore"}>
                    <WorkoutPlanSetArrView store={$inner.field} />
                  </Match>
                  <Match when={$inner.field.symbol === "ObjectFieldCore"}>
                    <WorkoutPlanSetObjView store={$inner.field} />
                  </Match>
                </Switch>
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );
}

export function WorkoutPlanSetObjView(props: { store: ObjectFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full space-y-4 my-2">
      <For each={state().fields}>
        {(field) => {
          if (field.hidden) {
            return null;
          }
          const $inner = store.mapFieldWithName(field.name);
          if (!$inner) {
            return null;
          }
          return (
            <div class="w-full">
              <div class="flex mb-2">
                <label class="block w-16 pt-2 mr-4 text-sm font-medium text-foreground mb-1 whitespace-nowrap">
                  {field.label}
                </label>
                <div class="flex-1 w-0">
                  {$inner.symbol === "ArrayFieldCore" ? (
                    <WorkoutPlanSetArrView store={$inner} />
                  ) : $inner.symbol === "SingleFieldCore" ? (
                    <WorkoutPlanSetInputView store={$inner} />
                  ) : $inner.symbol === "ObjectFieldCore" ? (
                    <div class="border border-border rounded-lg p-4 bg-muted">
                      <WorkoutPlanSetObjView store={$inner} />
                    </div>
                  ) : null}
                </div>
              </div>
            </div>
          );
        }}
      </For>
    </div>
  );
}

export function WorkoutPlanSetValuesView(props: { store: ObjectFieldCore<any> }) {
  const { store } = props;

  return (
    <div class="w-[780px] mx-auto">
      <WorkoutPlanSetObjView store={store} />
    </div>
  );
}
