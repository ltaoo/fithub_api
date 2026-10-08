import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, ListView } from "@/components/ui";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { fetchQuizList } from "@/biz/quiz/services";
import { fetchGiftCardList } from "@/biz/gift_card/services";

function GiftCardListViewModel(props: ViewComponentProps) {
  const request = {
    gift_card: {
      list: new ListCore(new RequestCore(fetchGiftCardList, { client: props.client })),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $btn_goto_create: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.gift_card_create");
      },
    }),
    $btn_goto_create_reward: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.gift_card_reward_create");
      },
    }),
  };
  let _state = {
    get response() {
      return request.gift_card.list.response;
    },
  };
  enum Events {
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  request.gift_card.list.onStateChange(() => methods.refresh());

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.gift_card.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function GiftCardListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(GiftCardListViewModel, [props]);

  return (
    <div>
      <div></div>
      <Button store={vm.ui.$btn_goto_create}>创建</Button>
      <Button store={vm.ui.$btn_goto_create_reward}>创建奖励</Button>
      <div>
        <ListView store={vm.request.gift_card.list}>
          <For each={state().response.dataSource}>
            {(v) => {
              return (
                <div>
                  <div>
                    <div>{v.code}</div>
                  </div>
                </div>
              );
            }}
          </For>
        </ListView>
      </div>
    </div>
  );
}
