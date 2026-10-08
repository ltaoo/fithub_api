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

function QuizListViewModel(props: ViewComponentProps) {
  const request = {
    quiz: {
      list: new ListCore(new RequestCore(fetchQuizList, { client: props.client })),
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
        props.history.push("root.home_layout.quiz_create");
      },
    }),
  };
  let _state = {
    get response() {
      return request.quiz.list.response;
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

  request.quiz.list.onStateChange(() => methods.refresh());

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.quiz.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function QuizListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(QuizListViewModel, [props]);

  return (
    <div>
      <div></div>
      <Button store={vm.ui.$btn_goto_create}>创建</Button>
      <div>
        <ListView store={vm.request.quiz.list}>
          <For each={state().response.dataSource}>
            {(v) => {
              return (
                <div>
                  <div>
                    <div>{v.content}</div>
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
