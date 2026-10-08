/**
 * @file 训练计划 集合 创建
 */
import { ViewComponentProps } from "@/store/types";

import { base, Handler } from "@/domains/base";
import { Result } from "@/domains/result";
import { ButtonCore, DialogCore, InputCore, SelectCore } from "@/domains/ui";
import { ImageUploadCore } from "@/domains/ui/form/image-upload";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { RefCore } from "@/domains/ui/cur";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { TheItemTypeFromListCore } from "@/domains/list/typing";
import {
  fetchWorkoutPlanList,
  fetchWorkoutPlanListProcess,
  fetchWorkoutScheduleList,
} from "@/biz/workout_plan/services";
import { WorkoutScheduleSelectViewModel } from "@/biz/workout_plan/workout_schedule_select";
import { WorkoutPlanSelectViewModel } from "@/biz/workout_plan/workout_plan_select";
import { ImageURLInputModel } from "@/components/image-url-input";
import { TagInputCore } from "@/domains/ui/form/tag-input";

export function WorkoutPlanSetValueViewModel(props: ViewComponentProps) {
  const request = {
    workout_plan: {
      list: new ListCore(
        new RequestCore(fetchWorkoutPlanList, { process: fetchWorkoutPlanListProcess, client: props.client })
      ),
    },
    workout_schedule: {
      list: new ListCore(new RequestCore(fetchWorkoutScheduleList, { client: props.client })),
    },
  };
  type WorkoutPlan = TheItemTypeFromListCore<typeof request.workout_plan.list>;
  type WorkoutSchedule = TheItemTypeFromListCore<typeof request.workout_schedule.list>;
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    selectPlan(...v: Parameters<typeof ui.$select_workout_plan.select>) {
      ui.$select_workout_plan.select(...v);
    },
    selectSchedule(...v: Parameters<typeof ui.$select_workout_schedule.select>) {
      ui.$select_workout_schedule.select(...v);
    },
    async toBody() {
      const r = await ui.$value.validate();
      if (r.error) {
        return Result.Err(r.error.message);
      }
      const { title, overview, icon_url, idx, content } = r.data;
      const details = [];
      for (let i = 0; i < content.length; i += 1) {
        const v = content[i];
        details.push({
          type: v.type ?? 1,
          id: Number(v.id),
          title: v.title,
          overview: v.overview,
        });
      }
      const body = {
        title,
        overview,
        idx: Number(idx),
        icon_url,
        details,
      };
      return Result.Ok(body);
    },
  };
  const ui = {
    $btn_plan_show_dialog: new ButtonCore({
      onClick() {
        ui.$dialog_workout_plan.show();
      },
    }),
    $btn_schedule_show_dialog: new ButtonCore({
      onClick() {
        ui.$dialog_workout_schedule.show();
      },
    }),
    $dialog_workout_plan: new DialogCore({
      onOk() {
        const selected_value = ui.$select_workout_plan.value;
        if (!selected_value) {
          props.app.tip({
            text: ["请先选择计划"],
          });
          return;
        }
        for (let i = 0; i < selected_value.length; i += 1) {
          const vv = selected_value[i];
          const field = ui.$value.fields.content.append();
          console.log("[PAGE]workout_plan_set/model - after content.append()", i, vv);
          field.setValue({
            id: vv.id,
            title: vv.title,
            type: 1,
          });
        }
        ui.$dialog_workout_plan.hide();
        methods.refresh();
      },
    }),
    $dialog_workout_schedule: new DialogCore({
      onOk() {
        const selected_value = ui.$select_workout_schedule.value;
        if (!selected_value) {
          alert(1);
          props.app.tip({
            text: ["请先选择计划"],
          });
          return;
        }
        for (let i = 0; i < selected_value.length; i += 1) {
          const vv = selected_value[i];
          const field = ui.$value.fields.content.append();
          field.setValue({
            id: vv.id,
            title: vv.title,
            type: 2,
          });
        }
        ui.$dialog_workout_schedule.hide();
        methods.refresh();
      },
    }),
    $ref_plan: new RefCore<WorkoutPlan>(),
    $ref_schedule: new RefCore<WorkoutSchedule>(),
    $ref_selected_plan: new RefCore<WorkoutPlan>(),
    $ref_selected_schedule: new RefCore<WorkoutSchedule>(),
    $select_workout_plan: WorkoutPlanSelectViewModel({
      defaultValue: [],
      list: request.workout_plan.list,
    }),
    $select_workout_schedule: WorkoutScheduleSelectViewModel({
      defaultValue: [],
      list: request.workout_schedule.list,
    }),
    $value: new ObjectFieldCore({
      label: "值",
      name: "value",
      fields: {
        title: new SingleFieldCore({
          label: "标题",
          name: "title",
          input: new InputCore({ defaultValue: "" }),
        }),
        overview: new SingleFieldCore({
          label: "描述",
          name: "overview",
          input: new InputCore({ defaultValue: "" }),
        }),
        icon_url: new SingleFieldCore({
          label: "图标",
          name: "icon_url",
          input: ImageURLInputModel({}),
        }),
        idx: new SingleFieldCore({
          label: "索引",
          name: "idx",
          input: new InputCore({ defaultValue: 0 }),
        }),
        content: new ArrayFieldCore({
          label: "计划",
          name: "content",
          field() {
            return new ObjectFieldCore({
              label: "计划或周期",
              name: "",
              fields: {
                type: new SingleFieldCore({
                  label: "类型",
                  name: "type",
                  input: new SelectCore({
                    defaultValue: 1,
                    options: [
                      {
                        value: 1,
                        label: "训练计划",
                      },
                      {
                        value: 2,
                        label: "周期计划",
                      },
                    ],
                  }),
                }),
                title: new SingleFieldCore({
                  label: "标题",
                  name: "title",
                  input: new InputCore({ defaultValue: "" }),
                }),
                overview: new SingleFieldCore({
                  label: "概要",
                  name: "overview",
                  input: new InputCore({ defaultValue: "" }),
                }),
                id: new SingleFieldCore({
                  label: "",
                  name: "id",
                  hidden: true,
                  input: new InputCore({ defaultValue: 0 }),
                }),
                tags: new SingleFieldCore({
                  label: "标签",
                  name: "tags",
                  input: new TagInputCore({}),
                }),
              },
            });
          },
        }),
      },
    }),
  };
  let _state = {
    get selected_plan() {
      return ui.$ref_selected_plan.value;
    },
    get selected_schedule() {
      return ui.$ref_selected_schedule.value;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  ui.$select_workout_plan.onStateChange(() => methods.refresh());
  ui.$select_workout_schedule.onStateChange(() => methods.refresh());

  return {
    methods,
    request,
    ui,
    state: _state,
    ready() {
      ui.$select_workout_plan.init();
      ui.$select_workout_schedule.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}
