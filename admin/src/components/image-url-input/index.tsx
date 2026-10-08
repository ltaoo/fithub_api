import { Show } from "@/timeless";

import { useViewModelStore } from "@/hooks";
import { Button, Input } from "@/components/ui";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore, InputCore } from "@/domains/ui";

export function ImageURLInputModel(props: {} = {}) {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $input: new InputCore({ defaultValue: "" }),
    $btn_preview: new ButtonCore({
      onClick() {
        const v = ui.$input.value;
        if (!v) {
          return;
        }
        _url = v;
        methods.refresh();
      },
    }),
  };

  let _url = "";
  let _state = {
    get value() {
      return ui.$input.value;
    },
  };
  enum Events {
    Change,
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.Change]: typeof _url;
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  return {
    shape: "custom" as const,
    type: "image_url_input",
    methods,
    ui,
    state: _state,
    get value() {
      return _state.value;
    },
    get defaultValue() {
      return "";
    },
    setValue(v: string) {
      _url = v;
      ui.$input.setValue(v);
      methods.refresh();
    },
    clear() {},
    ready() {},
    destroy() {},
    onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
      return bus.on(Events.Change, handler);
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}
type ImageURLInputModel = ReturnType<typeof ImageURLInputModel>;

export function ImageURLInputView(props: { store: ImageURLInputModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <div class="flex items-center gap-4">
      <div class="flex-1 flex items-center gap-2">
        <Input store={vm.ui.$input} />
        <Button store={vm.ui.$btn_preview}>查看</Button>
      </div>
      <Show when={state().value}>
        <img class="w-[80px] h-[80px] object-cover" src={state().value} />
      </Show>
    </div>
  );
}
