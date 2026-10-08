import { ViewComponentProps } from "@/store/types";

import { useViewModel } from "@/hooks";
import { Button, ScrollView } from "@/components/ui";

import { base, Handler } from "@/domains/base";
import { ButtonCore, ScrollViewCore } from "@/domains/ui";

function WorkoutPlanSetListViewModel(props: ViewComponentProps) {
  const ui = {
    $view: new ScrollViewCore(),
    $btn_goto_create: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.plan_set_create");
      },
    }),
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

export function WorkoutPlanSetListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutPlanSetListViewModel, [props]);

  return (
    <ScrollView store={vm.ui.$view}>
      <div class="p-4">
        <div></div>
        <Button store={vm.ui.$btn_goto_create}>创建集合</Button>
      </div>
    </ScrollView>
  );
}
