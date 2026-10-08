import { Input, Textarea, Checkbox, Select } from "@/components/ui";
import { For, Switch, createSignal, Match } from "@/timeless";

import { ObjectFieldCore, ArrayFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { TagInput } from "@/components/ui/tag-input";
import { Trash } from "@/timeless/icons";
import { ImageURLInputView } from "@/components/image-url-input";

function SingleValuesView(props: { store: SingleFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full">
      {(() => {
        if (state().input?.type === "textarea") {
          return (
            <Textarea store={store.input} />
          );
        }
        if (state().input?.shape === "tag-input") {
          return <TagInput store={store.input} />;
        }
        if (state().input?.shape === "image-upload") {
          return <ImageURLInputView store={store.input} />;
        }
        if (state().input?.shape === "input") {
          return (
            <Input store={store.input} />
          );
        }
        if (state().input?.shape === "number") {
          return (
            <Input store={store.input} />
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
          return (
            <Select store={store.input} />
          );
        }
        return null;
      })()}
    </div>
  );
}

function ArrayValuesView(props: { store: ArrayFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full space-y-3 my-2">
      <div
        class="inline-flex items-center justify-center py-2 px-4 bg-primary text-primary-foreground rounded-md hover:opacity-90 transition-colors cursor-pointer text-sm font-medium"
        onClick={() => store.append()}
      >
        添加项目
      </div>
      <div class="space-y-3">
        <For each={state().fields}>
          {(field, index) => {
            const $inner = store.mapFieldWithIndex(index());
            if (!$inner) {
              return null;
            }
            return (
              <div class="p-3 border border-border rounded-lg bg-w-bg-0 shadow-sm">
                <div class="flex justify-between items-center mb-2">
                  <div class="text-sm font-medium text-muted-foreground">
                    {field.label} {index() + 1}
                  </div>
                  <div onClick={() => store.remove(index())} class="text-red-500 hover:text-red-700 cursor-pointer p-1">
                    <Trash class="w-4 h-4" />
                  </div>
                </div>
                <Switch>
                  <Match when={$inner.field.symbol === "SingleFieldCore"}>
                    <SingleValuesView store={$inner.field} />
                  </Match>
                  <Match when={$inner.field.symbol === "ArrayFieldCore"}>
                    <ArrayValuesView store={$inner.field} />
                  </Match>
                  <Match when={$inner.field.symbol === "ObjectFieldCore"}>
                    <ObjectValuesView store={$inner.field} />
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

export function ObjectValuesView(props: { store: ObjectFieldCore<any> }) {
  const { store } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  return (
    <div class="w-full space-y-4 my-2">
      <For each={state().fields}>
        {(field) => {
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
                  <Switch>
                    <Match when={$inner.symbol === "SingleFieldCore"}>
                      {/* @ts-ignore */}
                      <SingleValuesView store={$inner} />
                    </Match>
                    <Match when={$inner.symbol === "ArrayFieldCore"}>
                      {/* @ts-ignore */}
                      <ArrayValuesView store={$inner} />
                    </Match>
                    <Match when={$inner.symbol === "ObjectFieldCore"}>
                      <div class="border border-w-fg-3 rounded-lg p-4 bg-w-bg-1">
                        {/* @ts-ignore */}
                        <ObjectValuesView store={$inner} />
                      </div>
                    </Match>
                  </Switch>
                </div>
              </div>
            </div>
          );
        }}
      </For>
    </div>
  );
}

export function EquipmentValueView(props: { store: ObjectFieldCore<any> }) {
  const { store } = props;

  return <ObjectValuesView store={store} />;
}
