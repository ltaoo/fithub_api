import { ViewComponentProps } from "@/store/types";
import { base, Handler } from "@/domains/base";

import { WorkoutPlanSetValueViewModel } from "./model";
import { WorkoutPlanSetValuesView } from "./values";
import { useViewModel } from "@/hooks";

function WorkoutPlanSetUpdateViewModel(props: ViewComponentProps) {
  const $value_model = WorkoutPlanSetValueViewModel(props);
  const ui = {
    $value: $value_model.ui.$value,
  };
  let _state = {};
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  return {
    ui,
    state: _state,
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutPlanSetUpdateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutPlanSetUpdateViewModel, [props]);

  return (
    <div>
      <div>
        <div>创建训练计划合集</div>
        <div></div>
        <WorkoutPlanSetValuesView store={vm.ui.$value} />
      </div>
    </div>
  );
}
