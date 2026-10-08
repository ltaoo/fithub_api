import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, ListView } from "@/components/ui";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { fetchPaperList } from "@/biz/quiz/services";

function PaperListViewModel(props: ViewComponentProps) {
  const request = {
    paper: {
      list: new ListCore(new RequestCore(fetchPaperList, { client: props.client })),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    handleClickPaper(v: { id: number }) {
      props.history.push("root.home_layout.paper_update", {
        id: String(v.id),
      });
    },
  };
  const ui = {
    $btn_goto_quiz_list: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.quiz_list");
      },
    }),
    $btn_goto_create: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.paper_create");
      },
    }),
  };
  let _state = {
    get response() {
      return request.paper.list.response;
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

  request.paper.list.onStateChange(() => methods.refresh());

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.paper.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function PaperListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(PaperListViewModel, [props]);

  return (
    <div>
      <div></div>
      <Button store={vm.ui.$btn_goto_quiz_list}>题库</Button>
      <Button store={vm.ui.$btn_goto_create}>创建</Button>
      <div>
        <ListView store={vm.request.paper.list}>
          <For each={state().response.dataSource}>
            {(v) => {
              return (
                <div>
                  <div>{v.name}</div>
                  <div
                    onClick={() => {
                      vm.methods.handleClickPaper(v);
                    }}
                  >
                    编辑
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
