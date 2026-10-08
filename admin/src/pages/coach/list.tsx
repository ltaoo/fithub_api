import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { useViewModel } from "@/hooks";
import { Button, ListView, ScrollView } from "@/components/ui";
import { ButtonCore, ScrollViewCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { createCoachAuthURL, fetchCoachList } from "@/biz/coach/service";

function CoachListViewModel(props: ViewComponentProps) {
  const request = {
    coach: {
      list: new ListCore(new RequestCore(fetchCoachList, { client: props.client })),
      create_auth_url: new RequestCore(createCoachAuthURL, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async create_auto_url(v: { id: number }) {
      const r = await request.coach.create_auth_url.run({ id: v.id });
      if (r.error) {
        return;
      }
      props.app.copy(r.data.url);
      props.app.tip({
        text: ["复制成功"],
      });
    },
  };
  const ui = {
    $view: new ScrollViewCore({
      async onReachBottom() {
        await request.coach.list.loadMore();
        ui.$view.finishLoadingMore();
      },
    }),
    $btn_goto_create: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.coach_create");
      },
    }),
  };
  let _state = {
    get response() {
      return request.coach.list.response;
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

  request.coach.list.onStateChange(() => methods.refresh());

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.coach.list.init();
    },
    destroy() {
      bus.destroy();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function CoachListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(CoachListViewModel, [props]);

  return (
    <div>
      <Button store={vm.ui.$btn_goto_create}>创建</Button>
      <ScrollView store={vm.ui.$view}>
        <ListView store={vm.request.coach.list} class="space-y-2">
          <For each={state().response.dataSource}>
            {(v) => {
              return (
                <div
                  class="border-2 border-w-fg-3 p-4"
                  onClick={() => {
                    props.history.push("root.home_layout.coach_content_create", {
                      id: String(v.id),
                    });
                  }}
                >
                  <div>{v.nickname}</div>
                  <div>
                    <img class="w-[68px] h-[68px]" src={v.avatar_url} />
                  </div>
                  <div
                    onClick={(event) => {
                      event.stopPropagation();
                      vm.methods.create_auto_url(v);
                    }}
                  >
                    <div>创建访问链接</div>
                  </div>
                </div>
              );
            }}
          </For>
        </ListView>
      </ScrollView>
    </div>
  );
}
