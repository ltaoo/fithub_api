import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { Button } from "@/components/ui";
import { ButtonCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { fetchSubscriptionPlanList } from "@/biz/subscription/services";

function SubscriptionPlanListViewModel(props: ViewComponentProps) {
  const request = {
    subscription_plan: {
      list: new ListCore(new RequestCore(fetchSubscriptionPlanList, { client: props.client })),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $btn_goto_create_view: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.subscription_plan_create");
      },
    }),
  };
  let _state = {};
  enum Events {
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  return {
    methods,
    ui,
    state: _state,
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function SubscriptionPlanListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(SubscriptionPlanListViewModel, [props]);

  return (
    <div>
      <div>
        <Button store={vm.ui.$btn_goto_create_view}>创建</Button>
      </div>
    </div>
  );
}
