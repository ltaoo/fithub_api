import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Input, Textarea } from "@/components/ui";
import { FieldV2 } from "@/components/ui/fieldv2";
import { FieldArrV2 } from "@/components/ui/field-arrv2";
import { TagInput } from "@/components/ui/tag-input";
import { Select } from "@/components/ui/select";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore, InputCore, SelectCore } from "@/domains/ui";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { TagInputCore } from "@/domains/ui/form/tag-input";
import { RequestCore } from "@/domains/request";
import { Result } from "@/domains/result";
import { createQuiz } from "@/biz/quiz/services";
import { QuizTypes } from "@/biz/quiz/constants";

function QuizCreateViewModel(props: ViewComponentProps) {
  const request = {
    quiz: {
      create: new RequestCore(createQuiz, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async create() {
      const r = await ui.$values.validate();
      if (r.error) {
        return Result.Err(r.error.message);
      }
      const data = r.data;
      if (!data.content) {
        const tip = "请输入题目内容";
        return Result.Err(tip);
      }
      const type = Number(data.type);
      if (!type) {
        const tip = "请选择类型";
        return Result.Err(tip);
      }
      if ([QuizTypes.Single, QuizTypes.Multiple].includes(type) && data.choices.length === 0) {
        const tip = "请输入可选项";
        return Result.Err(tip);
      }
      const body = {
        content: data.content,
        overview: data.overview,
        type,
        difficulty: Number(data.difficulty),
        tags: data.tags,
        choices: (() => {
          if ([QuizTypes.Single, QuizTypes.Multiple].includes(type)) {
            return data.choices.map((v, idx) => {
              return {
                value: idx + 1,
                text: v.text,
              };
            });
          }
          return [];
        })(),
        answer: (() => {
          if ([QuizTypes.Single, QuizTypes.Multiple].includes(type)) {
            return {
              type,
              value: data.choices
                .map((v, idx) => {
                  return {
                    idx: idx + 1,
                    correct: Number(v.correct),
                  };
                })
                .filter((v) => {
                  return v.correct;
                })
                .map((v) => {
                  return v.idx;
                }),
            };
          }
          return {
            type,
            value: [],
          };
        })(),
        analysis: data.analysis,
      };
      console.log("[PAGE]quiz/create - before quiz.create.run", body);
      const r2 = await request.quiz.create.run(body);
      if (r2.error) {
        return Result.Err(r2.error.message);
      }
      return Result.Ok(r2.data);
    },
  };
  const ui = {
    $btn_submit: new ButtonCore({
      async onClick() {
        const r = await methods.create();
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["创建成功"],
        });
      },
    }),
    $btn_add_choice: new ButtonCore({
      onClick() {
        ui.$values.fields.choices.append();
      },
    }),
    $values: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        content: new SingleFieldCore({
          label: "内容",
          name: "content",
          input: new InputCore({ defaultValue: "" }),
        }),
        overview: new SingleFieldCore({
          label: "概要",
          name: "overview",
          input: new InputCore({ defaultValue: "" }),
        }),
        type: new SingleFieldCore({
          label: "类型",
          name: "type",
          input: new SelectCore({
            defaultValue: QuizTypes.Single,
            options: [
              {
                label: "单选",
                value: QuizTypes.Single,
              },
              {
                label: "多选",
                value: QuizTypes.Multiple,
              },
              {
                label: "判断",
                value: QuizTypes.Judgment,
              },
              {
                label: "填空",
                value: QuizTypes.Fill,
              },
              {
                label: "简答",
                value: QuizTypes.Short,
              },
            ],
          }),
        }),
        difficulty: new SingleFieldCore({
          label: "难度",
          name: "difficulty",
          input: new InputCore({ defaultValue: 0 }),
        }),
        tags: new SingleFieldCore({
          label: "标签",
          name: "tags",
          input: new TagInputCore({}),
        }),
        choices: new ArrayFieldCore({
          label: "选项",
          name: "choices",
          field: () => {
            return new ObjectFieldCore({
              label: "",
              name: "",
              fields: {
                // idx: new SingleFieldCore({
                //   label: "顺序",
                //   name: "idx",
                //   input: new InputCore({ defaultValue: i }),
                // }),
                text: new SingleFieldCore({
                  label: "选项内容",
                  name: "text",
                  input: new InputCore({ defaultValue: "" }),
                }),
                correct: new SingleFieldCore({
                  label: "是否为正确答案",
                  name: "correct",
                  input: new InputCore({ defaultValue: 0 }),
                }),
              },
            });
          },
        }),
        analysis: new SingleFieldCore({
          label: "解析",
          name: "analysis",
          input: new InputCore({ defaultValue: "" }),
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

export function QuizCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(QuizCreateViewModel, [props]);

  return (
    <div>
      <div>
        <div class="w-[680px] mx-auto space-y-4">
          <FieldV2 store={vm.ui.$values.fields.content}>
            <Textarea store={vm.ui.$values.fields.content.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.overview}>
            <Textarea store={vm.ui.$values.fields.overview.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.type}>
            <Select store={vm.ui.$values.fields.type.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.difficulty}>
            <Input store={vm.ui.$values.fields.difficulty.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.tags}>
            <TagInput store={vm.ui.$values.fields.tags.input} />
          </FieldV2>
          <div>
            <Button store={vm.ui.$btn_add_choice}>添加选项</Button>
            <div class="flex items-center gap-4 mt-2 p-2">
              <FieldArrV2
                store={vm.ui.$values.fields.choices}
                render={(field) => {
                  return (
                    <div class="">
                      <FieldV2 store={field.fields.text}>
                        <Input store={field.fields.text.input} />
                      </FieldV2>
                      <FieldV2 store={field.fields.correct}>
                        <Input store={field.fields.correct.input} />
                      </FieldV2>
                    </div>
                  );
                }}
              ></FieldArrV2>
            </div>
          </div>
          <FieldV2 store={vm.ui.$values.fields.analysis}>
            <Textarea store={vm.ui.$values.fields.analysis.input} />
          </FieldV2>
        </div>
        <Button store={vm.ui.$btn_submit}>创建</Button>
      </div>
    </div>
  );
}
