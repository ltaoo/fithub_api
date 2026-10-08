import { $equipment_list, $muscle_list } from "@/store";
import { ViewComponentProps } from "@/store/types";
import { ImageURLInputModel } from "@/components/image-url-input";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ObjectFieldCore, SingleFieldCore, ArrayFieldCore } from "@/domains/ui/formv2";
import { InputCore } from "@/domains/ui/form/input";
import { TagInputCore } from "@/domains/ui/form/tag-input";
import { SelectCore } from "@/domains/ui";
import { Result } from "@/domains/result";
import { WorkoutActionType, WorkoutActionTypeOptions } from "@/biz/workout_action/constants";
import { fetchMuscleList, fetchMuscleListProcess, MuscleProfile } from "@/biz/muscle/services";
import { MuscleSelectViewModel } from "@/biz/muscle/muscle_select";
import { EquipmentSelectViewModel } from "@/biz/equipment/equipment_select";
import { RequestCore } from "@/domains/request";
import {
  createWorkoutAction,
  fetchWorkoutActionProfile,
  fetchWorkoutActionProfileProcess,
  getSetValueUnit,
  updateWorkoutAction,
  updateWorkoutActionIdx,
} from "@/biz/workout_action/services";
import { RefCore } from "@/domains/ui/cur";
import { ListCore } from "@/domains/list";
import { fetchEquipmentList, fetchEquipmentListProcess } from "@/biz/equipment/services";

export function WorkoutActionEditorViewModel(props: Pick<ViewComponentProps, "client" | "app">) {
  const request = {
    action: {
      profile: new RequestCore(fetchWorkoutActionProfile, {
        process: fetchWorkoutActionProfileProcess,
        client: props.client,
      }),
      update: new RequestCore(updateWorkoutAction, { client: props.client }),
      update_idx: new RequestCore(updateWorkoutActionIdx, { client: props.client }),
      create: new RequestCore(createWorkoutAction, { client: props.client }),
    },
    muscle: {
      list: new ListCore(new RequestCore(fetchMuscleList, { process: fetchMuscleListProcess, client: props.client })),
    },
    equipment: {
      list: new ListCore(
        new RequestCore(fetchEquipmentList, { process: fetchEquipmentListProcess, client: props.client })
      ),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async toBody() {
      const r = await ui.$form.validate();
      if (r.error) {
        return Result.Err(r.error.message);
      }
      const values = r.data;
      // console.log("[PAGE]home_action_create - Saving action:", values);
      const body = {
        id: values.id,
        name: values.name,
        zh_name: values.zh_name,
        alias: values.alias.join(","),
        overview: values.overview,
        cover_url: values.cover_url,
        type: values.type ?? WorkoutActionType.RESISTANCE,
        sort_idx: Number(values.sort_idx),
        score: Number(values.score),
        level: Number(values.level),
        tags1: values.tags1.join(","),
        tags2: values.tags2.join(","),
        pattern: values.pattern.join(","),
        details: JSON.stringify({
          v: "250608",
          steps: values.steps,
        }),
        points: JSON.stringify(values.points),
        problems: JSON.stringify(values.problems),
        equipment_ids: values.equipments.map((v) => v!.id).join(","),
        muscle_ids: values.muscles.map((v) => v!.id).join(","),
        primary_muscle_ids: values.primary_muscles.map((m) => m.id).join(","),
        secondary_muscle_ids: values.secondary_muscles.map((m) => m.id).join(","),
        alternative_action_ids: "",
        advanced_action_ids: "",
        regressed_action_ids: "",
        extra_config: JSON.stringify(values.extra_config),
      };
      return Result.Ok(body);
    },
    create() {},
    async update() {
      const r = await methods.toBody();
      if (r.error) {
        return Result.Err(r.error);
      }
      const body = r.data;
      const r2 = await request.action.update.run(body);
      if (r2.error) {
        return Result.Err(r2.error);
      }
      return Result.Ok(null);
    },
    async updateIdx(v: { id: number; idx: number }) {
      const r = await request.action.update_idx.run(v);
      if (r.error) {
        return;
      }
      props.app.tip({
        text: ["更新成功"],
      });
    },
    async fetch(id: number) {
      ui.$form.clear();
      const r = await request.action.profile.run({ id });
      if (r.error) {
        return Result.Err(r.error);
      }
      const profile = r.data;
      const muscle_ids = profile.muscles.map((v) => v.id);
      const equipment_ids = profile.equipments.map((v) => v.id);
      console.log("[BIZ]workout_action/model - after profile = r.data", profile, muscle_ids, equipment_ids);
      // @ts-ignore
      const r2 = await request.muscle.list.search({ ids: muscle_ids });
      if (r2.error) {
        return;
      }
      // @ts-ignore
      const r3 = await request.equipment.list.search({ ids: equipment_ids });
      if (r3.error) {
        return;
      }
      profile.muscles = r2.data.dataSource.filter((vv) => {
        return muscle_ids.includes(vv.id);
      });
      profile.equipments = r3.data.dataSource.filter((vv) => {
        return equipment_ids.includes(vv.id);
      });
      console.log("[BIZ]workout_action/model - before ui.$form.setValue", profile);
      ui.$form.setValue(profile);
      ui.$form.refresh();
    },
  };
  const ui = {
    // $ref: new RefCore<{ id: number }>(),
    $form: new ObjectFieldCore({
      name: "",
      label: "",
      fields: {
        id: new SingleFieldCore({
          name: "id",
          label: "",
          hidden: true,
          input: new InputCore({ defaultValue: 0 }),
        }),
        zh_name: new SingleFieldCore({
          name: "zh_name",
          label: "中文名称",
          input: new InputCore({ defaultValue: "" }),
        }),
        name: new SingleFieldCore({
          name: "name",
          label: "英文名称",
          input: new InputCore({ defaultValue: "" }),
        }),
        alias: new SingleFieldCore({
          name: "alias",
          label: "别名",
          input: new TagInputCore({ defaultValue: [] }),
        }),
        overview: new SingleFieldCore({
          name: "overview",
          label: "动作概述",
          input: new InputCore({ defaultValue: "", type: "textarea" }),
        }),
        cover_url: new SingleFieldCore({
          name: "cover_url",
          label: "封面",
          input: ImageURLInputModel({}),
        }),
        sort_idx: new SingleFieldCore({
          name: "sort_idx",
          label: "排序",
          input: new InputCore({ defaultValue: 10 }),
        }),
        type: new SingleFieldCore({
          name: "type",
          label: "动作类型",
          input: new SelectCore({
            defaultValue: WorkoutActionType.RESISTANCE,
            options: WorkoutActionTypeOptions,
          }),
        }),
        level: new SingleFieldCore({
          name: "level",
          label: "动作难度",
          input: new InputCore({
            defaultValue: 1,
            type: "number",
          }),
        }),
        score: new SingleFieldCore({
          name: "score",
          label: "动作评分",
          input: new InputCore({
            defaultValue: 1,
            type: "number",
          }),
        }),
        tags1: new SingleFieldCore({
          name: "tags1",
          label: "部位标签",
          input: new TagInputCore({
            defaultValue: [],
          }),
        }),
        tags2: new SingleFieldCore({
          name: "tags2",
          label: "类型标签",
          input: new TagInputCore({
            defaultValue: [],
          }),
        }),
        pattern: new SingleFieldCore({
          name: "pattern",
          label: "动作模式标签",
          input: new TagInputCore({
            defaultValue: [],
          }),
        }),
        muscles: new SingleFieldCore({
          name: "muscles",
          label: "肌肉",
          input: MuscleSelectViewModel({
            defaultValue: [],
            list: $muscle_list,
            // onLoaded(muscles) {
            //   _muscles = muscles;
            // },
          }),
        }),
        primary_muscles: new SingleFieldCore({
          name: "primary_muscles",
          label: "主要肌肉",
          input: MuscleSelectViewModel({
            defaultValue: [],
            list: $muscle_list,
            // onLoaded(muscles) {
            //   _muscles = muscles;
            // },
          }),
        }),
        secondary_muscles: new SingleFieldCore({
          name: "secondary_muscles",
          label: "次要肌肉",
          input: MuscleSelectViewModel({
            defaultValue: [],
            list: $muscle_list,
            // onLoaded(muscles) {
            //   _muscles = muscles;
            // },
          }),
        }),
        equipments: new SingleFieldCore({
          name: "equipments",
          label: "器械",
          input: EquipmentSelectViewModel({ defaultValue: [], list: $equipment_list }),
        }),
        steps: new ArrayFieldCore({
          label: "动作步骤",
          name: "steps",
          field() {
            return new ObjectFieldCore({
              label: "",
              name: "",
              fields: {
                text: new SingleFieldCore({
                  label: "说明",
                  name: "text",
                  input: new InputCore({ defaultValue: "" }),
                }),
                imgs: new ArrayFieldCore({
                  label: "图片",
                  name: "imgs",
                  field: (index: number) => {
                    return new SingleFieldCore({
                      name: "",
                      label: "",
                      input: ImageURLInputModel({}),
                    });
                  },
                }),
                tips: new ArrayFieldCore({
                  label: "提示",
                  name: "tips",
                  field: (index: number) => {
                    return new SingleFieldCore({
                      name: "",
                      label: "",
                      input: new InputCore({ defaultValue: "", type: "textarea" }),
                    });
                  },
                }),
              },
            });
          },
        }),
        points: new ArrayFieldCore({
          name: "points",
          label: "动作要点",
          field: (index: number) => {
            return new SingleFieldCore({
              name: `point_${index}`,
              label: "",
              input: new InputCore({ defaultValue: "", type: "textarea" }),
            });
          },
        }),
        problems: new ArrayFieldCore({
          name: "problems",
          label: "常见问题",
          field: (index: number) => {
            return new ObjectFieldCore({
              name: `problem_${index}`,
              label: "",
              fields: {
                title: new SingleFieldCore({
                  name: "title",
                  label: "简要说明",
                  input: new InputCore({ defaultValue: "" }),
                }),
                reason: new SingleFieldCore({
                  name: "reason",
                  label: "详细描述",
                  input: new InputCore({ defaultValue: "", type: "textarea" }),
                }),
                solutions: new ArrayFieldCore({
                  name: "solutions",
                  label: "解决方法",
                  field: (index: number) => {
                    return new SingleFieldCore({
                      name: `solutions_${index}`,
                      label: "",
                      input: new InputCore({ defaultValue: "", type: "textarea" }),
                    });
                  },
                }),
              },
            });
          },
        }),
        extra_config: new ObjectFieldCore({
          name: "extra_config",
          label: "配置",
          fields: {
            reps_unit: new SingleFieldCore({
              name: "reps_unit",
              label: "计数单位",
              input: new SelectCore({
                defaultValue: "次",
                options: [
                  {
                    label: "次",
                    value: getSetValueUnit("次"),
                  },
                  {
                    label: "秒",
                    value: getSetValueUnit("秒"),
                  },
                  {
                    label: "米",
                    value: getSetValueUnit("米"),
                  },
                ],
              }),
            }),
          },
        }),
      },
    }),
  };

  let _muscles: MuscleProfile[] = [];
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
    get muscles() {
      return _muscles;
    },
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}
