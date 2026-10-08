import { For, Show } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, ListView, ScrollView } from "@/components/ui";
import { WorkoutPlanPreviewCard } from "@/components/workout-plan-share-card";
import { fetchWorkoutPlanList, fetchWorkoutPlanListProcess } from "@/biz/workout_plan/services";
import { base, Handler } from "@/domains/base";
import { ButtonCore, ScrollViewCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";

function WorkoutPlanListViewModel(props: ViewComponentProps) {
  const request = {
    workout_plan: {
      list: new ListCore(
        new RequestCore(fetchWorkoutPlanList, { process: fetchWorkoutPlanListProcess, client: props.client })
      ),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $view: new ScrollViewCore({}),
    $create_btn: new ButtonCore({
      onClick: () => {
        props.history.push("root.home_layout.workout_plan_create");
      },
    }),
  };

  let _state = {
    get response() {
      return request.workout_plan.list.response;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  request.workout_plan.list.onStateChange(() => methods.refresh());

  return {
    state: _state,
    ui,
    request,
    ready() {
      request.workout_plan.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutPlanListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutPlanListViewModel, [props]);

  return (
    <>
      <ScrollView store={vm.ui.$view} class="p-4">
        <h1 class="text-2xl font-bold mb-4">计划列表</h1>
        <div>
          <Button store={vm.ui.$create_btn}>创建计划</Button>
        </div>
        <div class="py-4">
          <ListView store={vm.request.workout_plan.list}>
            <For each={state().response.dataSource}>
              {(v) => {
                return (
                  <div
                    class="py-4 border-b border-border"
                    onClick={() => {
                      props.history.push("root.home_layout.workout_plan_update", {
                        id: v.id.toString(),
                      });
                    }}
                  >
                    <div>
                      <div>{v.title}</div>
                      <div
                        onClick={(event) => {
                          event.stopPropagation();
                          props.history.push("root.home_layout.content_of_workout_plan_create", {
                            id: String(v.id),
                          });
                        }}
                      >
                        关联内容
                      </div>
                    </div>
                  </div>
                );
              }}
            </For>
          </ListView>
        </div>
      </ScrollView>
    </>
  );
}
