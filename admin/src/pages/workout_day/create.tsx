/**
 * @file 训练日记录
 */
import { For, Show } from "@/timeless";
import { Loader } from "@/timeless/icons";

import { ViewComponentProps } from "@/store/types";
import { $workout_action_list } from "@/store";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input, ListView, ScrollView } from "@/components/ui";
import { fetchWorkoutPlanProfile, fetchWorkoutPlanProfileProcess } from "@/biz/workout_plan/services";
import { WorkoutPlanSetType } from "@/biz/workout_plan/constants";
import { CountdownViewModel } from "@/biz/countdown";
import { WorkoutActionMultipleSelectViewModel } from "@/biz/workout_action_multiple_select";
import {
  fetchWorkoutActionListByIds,
  fetchWorkoutActionListByIdsProcess,
  fetchWorkoutActionProfile,
  fetchWorkoutActionProfileProcess,
} from "@/biz/workout_action/services";
import { base, Handler } from "@/domains/base";
import { RequestCore } from "@/domains/request";
import { ButtonCore, ButtonInListCore, DialogCore, ScrollViewCore } from "@/domains/ui";
import { RefCore } from "@/domains/ui/cur";
import { WorkoutActionSelectDialogViewModel } from "@/biz/workout_action_select_dialog";
import { WorkoutActionSelectDialogView } from "@/components/workout-action-select-dialog";
import { WeightInput } from "@/components/weight-input";
import { WeightInputViewModel } from "@/biz/weight_input";
import { Countdown } from "@/components/countdown";

export function WorkoutDayUpdateViewModel(props: ViewComponentProps) {
  const request = {
    workout_plan: {
      profile: new RequestCore(fetchWorkoutPlanProfile, {
        process: fetchWorkoutPlanProfileProcess,
        client: props.client,
      }),
    },
    workout_action: {
      list_by_id: new RequestCore(fetchWorkoutActionListByIds, {
        process: fetchWorkoutActionListByIdsProcess,
        client: props.client,
      }),
      profile: new RequestCore(fetchWorkoutActionProfile, {
        process: fetchWorkoutActionProfileProcess,
        client: props.client,
      }),
    },
    workout_day: {},
  };
  const ui = {
    $view: new ScrollViewCore(),
    //     $workout_action_dialog: new DialogCore({}),
    //     $workout_action_dialog_btn: new ButtonCore({
    //       onClick() {
    //         ui.$action_select_dialog.show();
    //       },
    //     }),
    $workout_plan_dialog_btn: new ButtonCore(),
    $action_select_dialog: new DialogCore({}),
    $action_select_view: new ScrollViewCore(),
    $action_select: WorkoutActionMultipleSelectViewModel({
      list: $workout_action_list,
      defaultValue: [],
      client: props.client,
    }),
    $start_btn: new ButtonCore(),
    $cur_step_ref: new RefCore<{ id?: number | string; idx: number }>(),
    $workout_action_change_btn: new ButtonInListCore<{ id?: number | string; idx: number }>({
      onClick(v) {
        ui.$action_select.request.action.list.init();
        ui.$cur_step_ref.select(v);
        ui.$workout_action_dialog.ui.$dialog.show();
      },
    }),
    $workout_action_dialog: WorkoutActionSelectDialogViewModel({
      defaultValue: [],
      client: props.client,
      list: $workout_action_list,
      onOk(acts) {
        const v = ui.$cur_step_ref.value;
        console.log("[PAGE]workout_day/create", acts, v);
        if (acts.length === 0) {
          props.app.tip({
            text: ["请选择动作"],
          });
          return;
        }
        if (v) {
          const step = _steps[v.idx];
          if (acts.length === 1) {
            _steps = [
              ..._steps.slice(0, v.idx),
              {
                id: step.id,
                idx: step.idx,
                sets: step.sets.map((set) => {
                  return {
                    actions: acts.map((act) => {
                      return {
                        id: Number(act.id),
                        zh_name: act.zh_name,
                        reps: 0,
                      };
                    }),
                  };
                }),
              },
              ..._steps.slice(v.idx + 1),
            ];
          }
          ui.$cur_step_ref.clear();
          ui.$workout_action_dialog.methods.clear();
          ui.$workout_action_dialog.ui.$dialog.hide();
          bus.emit(Events.StateChange, { ..._state });
          return;
        }
        _steps = [
          ..._steps,
          {
            idx: _steps.length,
            sets: [
              {
                actions: acts.map((act) => {
                  return {
                    id: Number(act.id),
                    zh_name: act.zh_name,
                    reps: 0,
                  };
                }),
              },
            ],
          },
        ];
        ui.$workout_action_dialog.methods.clear();
        ui.$workout_action_dialog.ui.$dialog.hide();
        bus.emit(Events.StateChange, { ..._state });
      },
    }),
    $weight_input_dialog: new DialogCore({
      onOk() {
        const v = ui.$weight_input.value;
        console.log("[PAGE]workout_day/create", v);
        ui.$weight_input_dialog.hide();
      },
    }),
    $weight_input: WeightInputViewModel({}),
    $countdown: CountdownViewModel({}),
    $workout_action_profile_dialog: new DialogCore({ footer: false }),
  };
  const methods = {
    async showWorkoutActionProfile(id: number | string) {
      ui.$workout_action_profile_dialog.show();
      const r = await request.workout_action.profile.run({ id });
      if (r.error) {
        props.app.tip({
          text: ["获取动作详情失败", r.error.message],
        });
        return;
      }
    },
  };
  let _steps: {
    id?: number | string;
    idx: number;
    sets: {
      actions: {
        id: number;
        zh_name: number | string;
        reps: number;
      }[];
    }[];
  }[] = [];
  let _state = {
    get steps() {
      return _steps;
    },
    get actions() {
      return ui.$action_select.state.actions;
    },
    get selectedActions() {
      return ui.$action_select.state.value;
    },
    get loading() {
      return request.workout_action.profile.loading;
    },
    get curWorkoutAction() {
      return request.workout_action.profile.response;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  ui.$action_select.onStateChange((v) => {
    bus.emit(Events.StateChange, { ..._state });
  });
  request.workout_action.profile.onStateChange((v) => {
    bus.emit(Events.StateChange, { ..._state });
  });
  return {
    state: _state,
    ui,
    request,
    methods,
    async ready() {
      //       ui.$countdown.start();
      // const workout_plan_id = props.view.query.workout_plan_id;
      // if (!workout_plan_id) {
      //   return;
      // }
      // const r = await request.workout_plan.profile.run({ id: Number(workout_plan_id) });
      // if (r.error) {
      //   props.app.tip({
      //     text: ["获取计划内容失败"],
      //   });
      //   return;
      // }
      // const { steps } = r.data;
      // //       const actions: Record<string, boolean> = [];
      // _steps = steps.map((step, idx) => {
      //   return {
      //     id: step.id,
      //     idx,
      //     sets: (() => {
      //       if (step.set_type === WorkoutPlanSetType.Normal) {
      //         //       actions[step.action.id] = true;
      //         const r1 = new Array(step.set_count).fill(0).map((_, index) => {
      //           return {
      //             actions: [
      //               {
      //                 id: Number(step.action.id),
      //                 zh_name: step.action.zh_name,
      //                 reps: step.reps,
      //               },
      //             ],
      //           };
      //         });
      //         return r1;
      //       }
      //       if (step.set_type === WorkoutPlanSetType.Combo) {
      //         const r2 = new Array(step.set_count).fill(0).map((_, index) => {
      //           return {
      //             actions: step.actions.map((act) => {
      //               //     actions[act.action.id] = true;
      //               return {
      //                 id: Number(act.action.id),
      //                 zh_name: act.action.zh_name,
      //                 reps: act.reps,
      //               };
      //             }),
      //           };
      //         });
      //         return r2;
      //       }
      //       if (step.set_type === WorkoutPlanSetType.Free) {
      //         const r3 = step.sets3.map((set) => {
      //           return {
      //             actions: set.actions.map((act) => {
      //               //     actions[act.action.id] = true;
      //               return {
      //                 id: Number(act.action.id),
      //                 zh_name: act.action.zh_name,
      //                 reps: act.reps,
      //               };
      //             }),
      //           };
      //         });
      //         return r3;
      //       }
      //       return [];
      //     })(),
      //   };
      // });
      // bus.emit(Events.StateChange, { ..._state });
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutDayUpdateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutDayUpdateViewModel, [props]);

  return (
    <>
      <ScrollView store={vm.ui.$view} class="p-4">
        <div class="h-[68px]">
          <h1 class="text-2xl font-bold mb-4">训练日</h1>
        </div>
        <div class="bg-w-bg-0 p-4 rounded-lg">
          <div class="flex gap-2">
            <Button variant="subtle" store={vm.ui.$workout_action_dialog.ui.$show_btn}>
              添加训练动作
            </Button>
            <Button variant="subtle" store={vm.ui.$workout_plan_dialog_btn}>
              选择训练计划
            </Button>
            <Button variant="subtle" store={vm.ui.$start_btn}>
              开始计时
            </Button>
          </div>

          <div class="mt-4 space-y-4">
            <For each={state().steps}>
              {(step, idx) => {
                return (
                  <div>
                    <Button store={vm.ui.$workout_action_change_btn.bind({ id: step.id, idx: idx() })}>修改动作</Button>
                    <div class="mt-4 space-y-4">
                      <For each={step.sets}>
                        {(set) => {
                          return (
                            <div class="flex items-center gap-2">
                              <div class="flex items-center gap-2 p-4 bg-gray-100 rounded-md">
                                <div class="space-y-2">
                                  <For each={set.actions}>
                                    {(action) => {
                                      return (
                                        <div class="gap-2">
                                          <div
                                            onClick={() => {
                                              vm.methods.showWorkoutActionProfile(action.id);
                                            }}
                                          >
                                            {action.zh_name}
                                          </div>
                                          <div class="flex items-center gap-2">
                                            <input
                                              class="border border-input rounded-md p-2"
                                              placeholder={String(action.reps)}
                                            />
                                            <div
                                              onClick={() => {
                                                vm.ui.$weight_input_dialog.show();
                                              }}
                                            >
                                              <input class="border border-input rounded-md p-2" placeholder="" />
                                            </div>
                                          </div>
                                        </div>
                                      );
                                    }}
                                  </For>
                                </div>
                                <div>
                                  <input
                                    type="checkbox"
                                    onChange={(e) => {
                                      console.log(e.target.checked);
                                      if (e.target.checked) {
                                        vm.ui.$countdown.start();
                                      }
                                    }}
                                  />
                                  完成
                                </div>
                              </div>
                            </div>
                          );
                        }}
                      </For>
                    </div>
                  </div>
                );
              }}
            </For>
          </div>
        </div>
      </ScrollView>
      <WorkoutActionSelectDialogView store={vm.ui.$workout_action_dialog} />
      <Dialog store={vm.ui.$weight_input_dialog}>
        <div class="w-[340px]">
          <WeightInput store={vm.ui.$weight_input} />
        </div>
      </Dialog>
      <Dialog store={vm.ui.$workout_action_profile_dialog}>
        <div class="w-[520px]">
          <Show
            when={!state().loading}
            fallback={
              <div class="flex items-center justify-center">
                <Loader class="w-8 h-8 animate-spin" />
              </div>
            }
          >
            <div>
              <div>{state().curWorkoutAction?.zh_name}</div>
              <div>{state().curWorkoutAction?.overview}</div>
            </div>
          </Show>
        </div>
      </Dialog>
      <div class="absolute left-1/2 bottom-10 -translate-x-1/2">
        <div class="flex items-center gap-2 bg-w-bg-0 p-4 rounded-lg">
          <Countdown store={vm.ui.$countdown} />
        </div>
      </div>
    </>
  );
}
