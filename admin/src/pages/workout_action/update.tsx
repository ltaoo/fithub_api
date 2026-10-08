/**
 * 健身动作编辑
 */
import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, ScrollView } from "@/components/ui";

import {
  fetchWorkoutActionProfile,
  fetchWorkoutActionProfileProcess,
  updateWorkoutAction,
} from "@/biz/workout_action/services";
import { base, Handler } from "@/domains/base";
import { ButtonCore, ScrollViewCore } from "@/domains/ui";
import { RequestCore } from "@/domains/request";

import { WorkoutActionValuesView } from "./action_form";
import { WorkoutActionEditorViewModel } from "./model";

function WorkoutActionUpdateViewModel(props: ViewComponentProps) {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { state: _state });
    },
    setState(state: Partial<typeof _state>) {
      _state = { ..._state, ...state };
      methods.refresh();
    },
  };
  const $model = WorkoutActionEditorViewModel(props);
  const ui = {
    $view: new ScrollViewCore({}),
    $values: $model.ui.$form,
    $back: new ButtonCore({
      onClick() {
        props.history.back();
      },
    }),
    $submit: new ButtonCore({
      async onClick() {
        const r = await $model.methods.update();
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["更新成功"],
        });
        // props.history.push("root.home_layout.action_list")
      },
    }),
  };

  let _loading = false;
  let _state = {
    get loading() {
      return _loading;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: {
      state: typeof _state;
    };
  };
  const bus = base<TheTypesOfEvents>();
  $model.onStateChange(() => methods.refresh());

  return {
    methods,
    ui,
    state: _state,
    async ready() {
      const id = Number(props.view.query.id);
      if (Number.isNaN(id)) {
        return;
      }
      $model.methods.fetch(id);
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutActionUpdateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutActionUpdateViewModel, [props]);

  return (
    <ScrollView store={vm.ui.$view} class="p-4">
      <div class="p-4 rounded-lg">
        <div class="flex flex-col gap-4">
          <WorkoutActionValuesView store={vm.ui.$values} />
        </div>
      </div>
      <div class="h-[68px]"></div>
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-w-bg-5 bg-w-bg-1">
        <div class="flex gap-2">
          <Button variant="subtle" store={vm.ui.$back}>
            返回
          </Button>
          <Button store={vm.ui.$submit}>提交</Button>
        </div>
      </div>
    </ScrollView>
  );
}
