import { For, Show } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Input, ListView, ScrollView, Video } from "@/components/ui";
import { WorkoutPlanPreviewCard } from "@/components/workout-plan-share-card";
import {
  createContentOfWorkoutPlan,
  fetchWorkoutPlanList,
  fetchWorkoutPlanListProcess,
  fetchWorkoutPlanProfile,
  fetchWorkoutPlanProfileProcess,
} from "@/biz/workout_plan/services";
import { base, Handler } from "@/domains/base";
import { ButtonCore, InputCore, ScrollViewCore } from "@/domains/ui";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { Result } from "@/domains/result";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { ContentSelect } from "@/components/content-select";
import { ContentSelectViewModel } from "@/biz/content/content_select";
import { fetchCoachContentList } from "@/biz/content/service";
import { FieldObjV2 } from "@/components/ui/field-obj2";
import { FieldV2 } from "@/components/ui/fieldv2";
import { FieldArrV2 } from "@/components/ui/field-arrv2";
import { PlayerCore } from "@/domains/player";
import { minute_text_to_seconds } from "@/utils";

function ContentOfWorkoutPlanCreateViewModel(props: ViewComponentProps) {
  const request = {
    workout_plan: {
      profile: new RequestCore(fetchWorkoutPlanProfile, {
        process: fetchWorkoutPlanProfileProcess,
        client: props.client,
      }),
    },
    content: {
      list: new ListCore(new RequestCore(fetchCoachContentList, { client: props.client })),
    },
    content_of_workout_plan: {
      create: new RequestCore(createContentOfWorkoutPlan, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async ready() {
      const id = Number(props.view.query.id);
      if (Number.isNaN(id)) {
        return Result.Err("缺少 id");
      }
      const r = await request.workout_plan.profile.run({ id });
      if (r.error) {
        return Result.Err(r.error);
      }
      const { steps } = r.data;
      const existing_ids: number[] = [];
      const workout_actions: { id: number; zh_name: string }[] = [];
      for (let i = 0; i < steps.length; i += 1) {
        const { actions } = steps[i];
        for (let b = 0; b < actions.length; b += 1) {
          const act = actions[b];
          if (!existing_ids.includes(act.action_id)) {
            existing_ids.push(act.action_id);
            workout_actions.push(act.action);
          }
        }
      }
      for (let i = 0; i < workout_actions.length; i += 1) {
        const v = workout_actions[i];
        const field = ui.$form.fields.points.append();
        field.setValue({
          action_id: v.id,
          action_name: v.zh_name,
          point: 0,
        });
      }
      ui.$form.refresh();
    },
    async toBody() {
      const id = Number(props.view.query.id);
      if (Number.isNaN(id)) {
        return Result.Err("缺少 id");
      }
      const r = await ui.$form.validate();
      if (r.error) {
        return Result.Err(r.error);
      }
      if (r.data.content.length === 0) {
        return Result.Err("请选择一个内容");
      }
      return Result.Ok({
        content_id: r.data.content[0].id,
        workout_plan_id: id,
        details: JSON.stringify({
          points: r.data.points.map((v) => {
            return {
              time_text: v.point,
              time: minute_text_to_seconds(v.point),
              workout_action_id: v.action_id,
              workout_action_name: v.action_name,
            };
          }),
        }),
      });
    },
    async create() {
      const r = await methods.toBody();
      if (r.error) {
        props.app.tip({
          text: [r.error.message],
        });
        return;
      }
      const r2 = await request.content_of_workout_plan.create.run(r.data);
      if (r2.error) {
        props.app.tip({
          text: [r2.error.message],
        });
        return;
      }
      props.app.tip({
        text: ["创建成功"],
      });
    },
  };
  const ui = {
    $view: new ScrollViewCore({}),
    $form: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        content: new SingleFieldCore({
          label: "内容",
          name: "content",
          input: ContentSelectViewModel({
            defaultValue: [],
            multiple: false,
            list: request.content.list,
            onChange(list) {
              if (list.length === 0) {
                return;
              }
              const v = list[0];
              ui.$video.load(v.video_key);
            },
          }),
        }),
        points: new ArrayFieldCore({
          label: "时间点",
          name: "points",
          field() {
            return new ObjectFieldCore({
              label: "",
              name: "",
              fields: {
                action_id: new SingleFieldCore({
                  label: "",
                  name: "action_id",
                  hidden: true,
                  input: new InputCore({ defaultValue: 0, type: "number" }),
                }),
                action_name: new SingleFieldCore({
                  label: "动作",
                  name: "action_name",
                  input: new InputCore({ defaultValue: "" }),
                }),
                point: new SingleFieldCore({
                  label: "时间点",
                  name: "point",
                  input: new InputCore({ defaultValue: "" }),
                }),
              },
            });
          },
        }),
      },
    }),
    $btn_create: new ButtonCore({
      onClick() {
        methods.create();
      },
    }),
    $btn_back: new ButtonCore({
      onClick() {
        props.history.back();
      },
    }),
    $video: new PlayerCore({
      app: props.app,
    }),
  };

  let _state = {
    get profile() {
      return request.workout_plan.profile.response;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  request.workout_plan.profile.onStateChange(() => methods.refresh());

  return {
    state: _state,
    ui,
    request,
    ready() {
      methods.ready();
    },
    destroy() {
      bus.destroy();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function ContentOfWorkoutPlanCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(ContentOfWorkoutPlanCreateViewModel, [props]);

  return (
    <>
      <ScrollView store={vm.ui.$view} class="p-4">
        <h1 class="text-2xl font-bold mb-4">计划关联内容</h1>
        <Show when={state().profile}>
          <div>
            <div>{state().profile?.title}</div>
            <Video store={vm.ui.$video} />
          </div>
        </Show>
        <FieldObjV2 store={vm.ui.$form}>
          <FieldV2 store={vm.ui.$form.fields.content}>
            <ContentSelect store={vm.ui.$form.fields.content.input} />
          </FieldV2>
          <FieldArrV2
            store={vm.ui.$form.fields.points}
            render={(field) => {
              return (
                <FieldObjV2 store={field}>
                  <FieldV2 store={field.fields.action_name}>
                    <Input store={field.fields.action_name.input} />
                  </FieldV2>
                  <FieldV2 store={field.fields.point}>
                    <Input store={field.fields.point.input} />
                  </FieldV2>
                </FieldObjV2>
              );
            }}
          ></FieldArrV2>
        </FieldObjV2>
        <div class="h-[68px]"></div>
        <div class="absolute bottom-0 left-0 right-0 p-4 bg-w-bg-0 border-t border-border">
          <div class="flex gap-2">
            <Button variant="subtle" store={vm.ui.$btn_back}>
              返回
            </Button>
            <Button store={vm.ui.$btn_create}>提交</Button>
          </div>
        </div>
      </ScrollView>
    </>
  );
}
