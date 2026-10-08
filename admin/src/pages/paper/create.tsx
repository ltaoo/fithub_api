import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input, ListView, Textarea } from "@/components/ui";
import { FieldV2 } from "@/components/ui/fieldv2";
import { TagInput } from "@/components/ui/tag-input";
import { FieldArrV2 } from "@/components/ui/field-arrv2";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore, InputCore } from "@/domains/ui";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { TagInputCore } from "@/domains/ui/form/tag-input";
import { RequestCore } from "@/domains/request";
import { createPaper, fetchQuizList } from "@/biz/quiz/services";
import { Result } from "@/domains/result";
import { QuizSelectViewModel } from "@/biz/quiz_select";
import { ListCore } from "@/domains/list";
import { PaperValuesViewModel } from "./paper";

function PaperCreateViewModel(props: ViewComponentProps) {
  const request = {
    quiz: {
      list: new ListCore(new RequestCore(fetchQuizList, { client: props.client })),
    },
    paper: {
      create: new RequestCore(createPaper, { client: props.client }),
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
      if (!data.name) {
        const tip = "请输入标题";
        return Result.Err(tip);
      }
      if (data.quizzes.length === 0) {
        const tip = "请添加题目";
        return Result.Err(tip);
      }
      const body = {
        name: data.name,
        overview: data.overview,
        tags: data.tags,
        pass_score: Number(data.pass_score),
        duration: Number(data.duration),
        quiz_list: (() => {
          return data.quizzes.map((quiz) => {
            return {
              relation_id: Number(quiz.relation_id),
              id: Number(quiz.id),
              score: Number(quiz.score),
              sort_idx: Number(quiz.sort_idx),
            };
          });
        })(),
      };
      const r2 = await request.paper.create.run(body);
      if (r2.error) {
        return Result.Err(r2.error.message);
      }
      return Result.Ok(r2.data);
    },
  };
  const _$value = PaperValuesViewModel();
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
    $btn_add_quiz: new ButtonCore({
      onClick() {
        ui.$select_quiz.request.quiz.list.init();
        ui.$select_quiz.ui.$dialog.show();
      },
    }),
    $select_quiz: QuizSelectViewModel({
      defaultValue: [],
      list: request.quiz.list,
      client: props.client,
    }),
    $values: _$value.ui.$values,
  };
  let _state = {
    get quizzes() {
      return ui.$select_quiz.state.list;
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

  ui.$select_quiz.onStateChange(() => methods.refresh());
  ui.$select_quiz.ui.$dialog.onOk(() => {
    for (let i = 0; i < ui.$select_quiz.value.length; i += 1) {
      const v = ui.$select_quiz.value[i];
      const field = ui.$values.fields.quizzes.append();
      field.setValue(v);
    }
    ui.$select_quiz.ui.$dialog.hide();
  });

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

export function PaperCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(PaperCreateViewModel, [props]);

  return (
    <>
      <div>
        <div></div>
        <div>
          <div class="w-[680px] mx-auto space-y-4">
            <FieldV2 store={vm.ui.$values.fields.name}>
              <Input store={vm.ui.$values.fields.name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.overview}>
              <Textarea store={vm.ui.$values.fields.overview.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.pass_score}>
              <Input store={vm.ui.$values.fields.pass_score.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.duration}>
              <Input store={vm.ui.$values.fields.duration.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.tags}>
              <TagInput store={vm.ui.$values.fields.tags.input} />
            </FieldV2>
            <div>
              <Button store={vm.ui.$btn_add_quiz}>添加题目</Button>
              <div class="space-y-2 mt-2 p-2">
                <FieldArrV2
                  store={vm.ui.$values.fields.quizzes}
                  render={(field) => {
                    return (
                      <div class="">
                        <FieldV2 store={field.fields.content}>
                          <Input store={field.fields.content.input} />
                        </FieldV2>
                        <FieldV2 store={field.fields.score}>
                          <Input store={field.fields.score.input} />
                        </FieldV2>
                        <FieldV2 store={field.fields.sort_idx}>
                          <Input store={field.fields.sort_idx.input} />
                        </FieldV2>
                      </div>
                    );
                  }}
                ></FieldArrV2>
              </div>
            </div>
          </div>
          <Button store={vm.ui.$btn_submit}>创建</Button>
        </div>
      </div>
      <Dialog store={vm.ui.$select_quiz.ui.$dialog}>
        <div class="w-[520px]">
          <ListView store={vm.ui.$select_quiz.request.quiz.list} class="space-y-2">
            <For each={state().quizzes}>
              {(v) => {
                return (
                  <div
                    classList={{
                      "p-4 border-2 border-w-bg-5 rounded-lg": true,
                      "border-w-fg-2 bg-w-bg-5": v.selected,
                    }}
                    onClick={() => {
                      vm.ui.$select_quiz.methods.select(v);
                    }}
                  >
                    <div>{v.content}</div>
                  </div>
                );
              }}
            </For>
          </ListView>
        </div>
      </Dialog>
    </>
  );
}
