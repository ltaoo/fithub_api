/**
 * @file 通用多选选择
 */
import { base, Handler } from "@/domains/base";
import { HttpClientCore } from "@/domains/http_client";
import { ListCore } from "@/domains/list";
import { ButtonCore, DialogCore, InputCore, PopoverCore, ScrollViewCore, SelectCore } from "@/domains/ui";

type QuizId = number;
type Quiz = {
  id: QuizId;
  content: string;
};

export function QuizSelectViewModel(props: {
  defaultValue: Quiz[];
  multiple?: boolean;
  client: HttpClientCore;
  list: ListCore<any>;
  onChange?: (list: Quiz[]) => void;
}) {
  const request = {
    quiz: {
      list: props.list,
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    select(vv: Quiz) {
      const existing = _selected.find((v) => v.id === vv.id);
      if (_multiple === false) {
        if (existing) {
          return;
        }
        _selected = [vv];
        bus.emit(Events.Change, _selected);
        methods.refresh();
        return;
      }
      if (existing) {
        _selected = _selected.filter((v) => v.id !== vv.id);
        bus.emit(Events.Change, _selected);
        methods.refresh();
        return;
      }
      const v = _list.find((v) => v.id === vv.id);
      if (!v) {
        return;
      }
      _selected.push(v);
      bus.emit(Events.Change, _selected);
      methods.refresh();
    },
    remove(vv: Quiz) {
      _selected = _selected.filter((v) => v.id !== vv.id);
      bus.emit(Events.Change, _selected);
      methods.refresh();
    },
    map_list(list: { id: QuizId }[]) {
      return _list.flatMap((a) => {
        return list.find((v) => v.id === a.id) ?? [];
      });
    },
    find(value: { id: QuizId }) {
      return _list.find((a) => a.id === value.id) ?? null;
    },
    search(value: string) {
      request.quiz.list.search({
        keyword: value,
      });
    },
    set_list(list: Quiz[]) {
      // @ts-ignore
      request.quiz.list.modifyResponse((v) => {
        return {
          ...v,
          initial: false,
          dataSource: list,
        };
      });
      methods.refresh();
    },
    clear() {
      _selected = [];
      bus.emit(Events.Change, _selected);
      methods.refresh();
    },
  };
  const ui = {
    $dropdown: new SelectCore({
      defaultValue: "",
      options: [],
    }),
    $popover: new PopoverCore(),
    $scroll: new ScrollViewCore({
      async onReachBottom() {
        await request.quiz.list.loadMore();
        ui.$scroll.finishLoadingMore();
      },
    }),
    $dialog: new DialogCore(),
    $input_search_keyword: new InputCore({ defaultValue: "" }),
    $btn_search_submit: new ButtonCore({}),
  };

  let _multiple = props.multiple ?? true;
  let _selected: Quiz[] = [];
  let _list: Quiz[] = props.list?.response.dataSource ?? [];
  let _state = {
    get value() {
      return _selected;
    },
    get selected() {
      return _selected.flatMap((item) => {
        const existing = _list.find((a) => a.id === item.id);
        if (!existing) {
          return [];
        }
        return [existing];
      });
    },
    get response() {
      return request.quiz.list.response;
    },
    get list() {
      return _list.map((v) => {
        return {
          ...v,
          selected: _state.selected
            .map((v2) => {
              return v2.id;
            })
            .includes(v.id),
        };
      });
    },
  };
  enum Events {
    Change,
    ListLoaded,
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.ListLoaded]: typeof _list;
    [Events.Change]: typeof _selected;
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  request.quiz.list.onStateChange((state) => {
    _list = state.dataSource;
    bus.emit(Events.ListLoaded, _list);
    methods.refresh();
  });
  bus.on(Events.Change, (actions) => {
    if (props.onChange) {
      props.onChange(actions);
    }
  });

  return {
    shape: "custom" as const,
    type: "multiple-select",
    state: _state,
    methods,
    request,
    ui,
    get value() {
      return _selected;
    },
    get defaultValue() {
      return props.defaultValue;
    },
    setValue(value: Quiz[]) {
      const v = _list.filter((a) => {
        return value.find((v) => v.id === a.id);
      });
      _selected = v;
      bus.emit(Events.StateChange, { ..._state });
    },
    ready() {},
    onListLoaded(handler: Handler<TheTypesOfEvents[Events.ListLoaded]>) {
      return bus.on(Events.ListLoaded, handler);
    },
    onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
      return bus.on(Events.Change, handler);
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export type QuizSelectViewModel = ReturnType<typeof QuizSelectViewModel>;
