import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ObjectFieldCore, SingleFieldCore, ArrayFieldCore } from "@/domains/ui/formv2";
import { InputCore } from "@/domains/ui";
import { TagInputCore } from "@/domains/ui/form/tag-input";

export function PaperValuesViewModel() {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $values: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        id: new SingleFieldCore({
          label: "id",
          name: "id",
          hidden: true,
          input: new InputCore({ defaultValue: 0 }),
        }),
        name: new SingleFieldCore({
          label: "标题",
          name: "name",
          input: new InputCore({ defaultValue: "" }),
        }),
        overview: new SingleFieldCore({
          label: "概要",
          name: "overview",
          input: new InputCore({ defaultValue: "" }),
        }),
        tags: new SingleFieldCore({
          label: "标签",
          name: "tags",
          input: new TagInputCore({}),
        }),
        duration: new SingleFieldCore({
          label: "考试时间",
          name: "duration",
          input: new InputCore({ defaultValue: 60 }),
        }),
        pass_score: new SingleFieldCore({
          label: "及格分数",
          name: "pass_score",
          input: new InputCore({ defaultValue: 60 }),
        }),
        quizzes: new ArrayFieldCore({
          label: "题目",
          name: "quizzes",
          field: () => {
            return new ObjectFieldCore({
              label: "",
              name: "",
              fields: {
                relation_id: new SingleFieldCore({
                  label: "",
                  name: "relation_id",
                  hidden: true,
                  input: new InputCore({ defaultValue: 0 }),
                }),
                id: new SingleFieldCore({
                  label: "题目id",
                  name: "id",
                  hidden: true,
                  input: new InputCore({ defaultValue: 0 }),
                }),
                content: new SingleFieldCore({
                  label: "题目",
                  name: "content",
                  input: new InputCore({ defaultValue: "" }),
                }),
                score: new SingleFieldCore({
                  label: "分数",
                  name: "score",
                  input: new InputCore({ defaultValue: 1 }),
                }),
                sort_idx: new SingleFieldCore({
                  label: "顺序",
                  name: "sort_idx",
                  input: new InputCore({ defaultValue: 1 }),
                }),
              },
            });
          },
        }),
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
