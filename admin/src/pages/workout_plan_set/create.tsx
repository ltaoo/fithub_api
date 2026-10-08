import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, ListView } from "@/components/ui";

import { createWorkoutPlanSet } from "@/biz/workout_plan/services";
import { RequestCore } from "@/domains/request";
import { base, Handler } from "@/domains/base";
import { ButtonCore } from "@/domains/ui";

import { WorkoutPlanSetValueViewModel } from "./model";
import { WorkoutPlanSetValuesView } from "./values";

function WorkoutPlanSetCreateViewModel(props: ViewComponentProps) {
  const $model = WorkoutPlanSetValueViewModel(props);
  const request = {
    workout_plan: $model.request.workout_plan,
    workout_schedule: $model.request.workout_schedule,
    workout_plan_set: {
      create: new RequestCore(createWorkoutPlanSet, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    selectPlan: $model.methods.selectPlan,
    selectSchedule: $model.methods.selectSchedule,
  };
  const ui = {
    $value: $model.ui.$value,
    $btn_plan_show_dialog: $model.ui.$btn_plan_show_dialog,
    $btn_schedule_show_dialog: $model.ui.$btn_schedule_show_dialog,
    $dialog_workout_plan: $model.ui.$dialog_workout_plan,
    $dialog_workout_schedule: $model.ui.$dialog_workout_schedule,
    $btn_submit: new ButtonCore({
      async onClick() {
        const r = await $model.methods.toBody();
        if (r.error) {
          return;
        }
        const body = r.data;
        const r2 = await request.workout_plan_set.create.run(body);
        if (r2.error) {
          props.app.tip({
            text: [r2.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["创建成功"],
        });
      },
    }),
  };
  let _state = {
    get workout_plans() {
      return $model.ui.$select_workout_plan.state.list;
    },
    get selected_plan() {
      return $model.ui.$select_workout_plan.state.value;
    },
    get workout_schedules() {
      return $model.ui.$select_workout_schedule.state.list;
    },
    get selected_schedule() {
      return $model.ui.$select_workout_schedule.state.value;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  $model.onStateChange(() => methods.refresh());

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      $model.ready();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutPlanSetCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutPlanSetCreateViewModel, [props]);

  return (
    <div>
      <div>
        <div>创建训练计划合集</div>
        <div></div>
        <WorkoutPlanSetValuesView store={vm.ui.$value} />
        <Button store={vm.ui.$btn_plan_show_dialog}>添加计划</Button>
        <Button store={vm.ui.$btn_schedule_show_dialog}>添加周期</Button>
        <div>
          <Button store={vm.ui.$btn_submit}>提交</Button>
        </div>
      </div>
      <Dialog store={vm.ui.$dialog_workout_plan}>
        <div class="w-[520px]">
          <ListView store={vm.request.workout_plan.list} class="space-y-2">
            <For each={state().workout_plans}>
              {(v) => {
                return (
                  <div>
                    <div
                      classList={{
                        "p-2 border-2 rounded-lg": true,
                        "border-w-fg-1": v.selected,
                        "border-w-fg-3 ": !v.selected,
                      }}
                      onClick={() => {
                        vm.methods.selectPlan(v);
                      }}
                    >
                      {v.title}
                    </div>
                  </div>
                );
              }}
            </For>
          </ListView>
        </div>
      </Dialog>
      <Dialog store={vm.ui.$dialog_workout_schedule}>
        <div class="w-[520px]">
          <ListView store={vm.request.workout_schedule.list} class="space-y-2">
            <For each={state().workout_schedules}>
              {(data) => {
                return (
                  <div>
                    <div
                      classList={{
                        "p-2 border-2 rounded-lg": true,
                        "border-w-fg-1": data.selected,
                        "border-w-fg-3": !data.selected,
                      }}
                      onClick={() => {
                        vm.methods.selectSchedule(data);
                      }}
                    >
                      {data.title}
                    </div>
                  </div>
                );
              }}
            </For>
          </ListView>
        </div>
      </Dialog>
    </div>
  );
}
