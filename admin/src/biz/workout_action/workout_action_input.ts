import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { WorkoutActionSelectDialogViewModel } from "@/biz/workout_action_select_dialog";
import { ButtonCore } from "@/domains/ui";

export function WorkoutActionInputViewModel(props: { $select: WorkoutActionSelectDialogViewModel }) {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $btn_show_dialog: new ButtonCore({
      onClick() {
        props.$select.init();
        props.$select.ui.$dialog.show();
      },
    }),
  };

  let _value: { id: number | string; zh_name: string }[] = [];
  let _state = {
    get value() {
      return _value;
    },
  };
  enum Events {
    Change,
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.Change]: typeof _value;
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  //   props.$select.onChange((v) => {
  //     _value = v;
  //     methods.refresh();
  //   });
  props.$select.ui.$dialog.onOk(() => {
    _value = props.$select.value;
    methods.refresh();
    props.$select.ui.$dialog.hide();
  });

  return {
    shape: "custom" as const,
    type: "workout_action" as const,
    methods,
    ui,
    state: _state,
    get value() {
      return _state.value;
    },
    get defaultValue() {
      return [];
    },
    setValue(v: { id: number | string; zh_name: string }[]) {
      _value = v;
      methods.refresh();
    },
    ready() {},
    destroy() {
      bus.destroy();
    },
    onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
      return bus.on(Events.Change, handler);
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export type WorkoutActionInputViewModel = ReturnType<typeof WorkoutActionInputViewModel>;
