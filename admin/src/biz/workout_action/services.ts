import { request } from "@/biz/requests";
import { idsMapValue } from "@/biz/services/utils";
import { ListResponseWithCursor } from "@/biz/requests/types";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result, UnpackedResult } from "@/domains/result";
import { FetchParams } from "@/domains/list/typing";
import { parseJSONStr } from "@/utils";

import { WorkoutActionSteps, WorkoutActionProblems, WorkoutActionStepsJSON250608 } from "./types";
import { WorkoutActionType } from "./constants";

export type SetValueUnit = "公斤" | "磅" | "秒" | "分" | "次" | "千米" | "米" | "千卡";
export function getSetValueUnit(v: SetValueUnit): SetValueUnit {
  return v;
}

type PartialWorkoutAction = {
  id: number;
  name: string;
  zh_name: string;
  alias: string;
  idx: number;
  type: string;
  overview: string;
  muscle_ids: string;
  equipment_ids: string;
};

/**
 * 获取健身动作列表
 * @returns
 */
export function fetchWorkoutActionList(body: FetchParams & { type?: string; keyword?: string; tag?: string }) {
  return request.post<ListResponseWithCursor<PartialWorkoutAction>>("/api/workout_action/list", {
    page_size: body.pageSize,
    page: body.page,
    // next_marker: body.next_marker,
    type: body.type === "all" ? undefined : body.type,
    keyword: body.keyword,
    tag: body.tag,
  });
}
export function fetchWorkoutActionListProcess(r: TmpRequestResp<typeof fetchWorkoutActionList>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  return Result.Ok({
    ...r.data,
    list: r.data.list.map((action) => {
      return {
        id: action.id,
        name: action.name,
        zh_name: action.zh_name,
        type: action.type,
        overview: action.overview,
        idx: action.idx,
        // tags1: action.tags1.split(",").filter(Boolean),
        // tags2: action.tags2.split(",").filter(Boolean),
        // level: action.level,
        equipments: idsMapValue(action.equipment_ids),
        muscles: idsMapValue(action.muscle_ids),
      };
    }),
  });
}
export type WorkoutActionProfile = UnpackedResult<ReturnType<typeof fetchWorkoutActionListProcess>>["list"][number];

export function fetchWorkoutActionListByIds(body: { ids: number[] }) {
  return request.post<{
    list: PartialWorkoutAction[];
  }>("/api/workout_action/list_by_ids", {
    ids: body.ids,
  });
}
export function fetchWorkoutActionListByIdsProcess(r: TmpRequestResp<typeof fetchWorkoutActionListByIds>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  return Result.Ok({
    list: r.data.list.map((v) => {
      return {
        id: v.id,
        name: v.name,
        zh_name: v.zh_name,
        type: v.type,
        overview: v.overview,
        idx: v.idx,
        equipments: idsMapValue(v.equipment_ids),
        muscles: idsMapValue(v.muscle_ids),
      };
    }),
  });
}

export function deleteWorkoutAction(v: { id: number }) {
  return request.post("/api/workout_action/delete", { id: v.id });
}

type WorkoutActionBody = {
  name: string;
  zh_name: string;
  alias: string;
  overview: string;
  cover_url: string;
  type: string;
  level: number;
  tags1: string;
  tags2: string;
  details: string;
  points: string;
  problems: string;
  equipment_ids: string;
  muscle_ids: string;
  primary_muscle_ids: string;
  secondary_muscle_ids: string;
  alternative_action_ids: string;
  advanced_action_ids: string;
  regressed_action_ids: string;
  pattern: string;
  sort_idx: number;
  score: number;
  extra_config: string;
};
/**
 * 创建健身动作
 * @param body
 * @returns
 */
export function createWorkoutAction(body: WorkoutActionBody) {
  return request.post<void>("/api/workout_action/create", {
    ...body,
  });
}

/**
 * 获取健身动作详情
 * @param body
 * @returns
 */
export function fetchWorkoutActionProfile(body: { id: number | string }) {
  return request.post<{
    id: number;
    name: string;
    zh_name: string;
    alias: string;
    overview: string;
    type: string;
    sort_idx: number;
    level: number;
    score: number;
    tags1: string;
    tags2: string;
    pattern: string;
    details: string;
    steps: string;
    points: string;
    problems: string;
    extra_config: string;
    equipment_ids: string;
    muscle_ids: string;
    primary_muscle_ids: string;
    secondary_muscle_ids: string;
    alternative_action_ids: string;
    advanced_action_ids: string;
    regressed_action_ids: string;
  }>("/api/workout_action/profile", {
    id: Number(body.id),
  });
}
export function fetchWorkoutActionProfileProcess(r: TmpRequestResp<typeof fetchWorkoutActionProfile>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  const data = r.data;
  return Result.Ok({
    id: data.id,
    name: data.name,
    zh_name: data.zh_name,
    alias: data.alias.split(",").filter(Boolean),
    overview: data.overview,
    sort_idx: data.sort_idx,
    type: data.type,
    level: data.level,
    score: data.score,
    tags1: data.tags1.split(",").filter(Boolean),
    tags2: data.tags2.split(",").filter(Boolean),
    pattern: data.pattern?.split(",").filter(Boolean),
    steps: (() => {
      if (data.steps) {
        const r = parseJSONStr<WorkoutActionStepsJSON250608>(data.steps);
        if (r.error) {
          return [];
        }
        return r.data.steps;
      }
      const r = parseJSONStr<WorkoutActionSteps & WorkoutActionStepsJSON250608>(data.details);
      console.log("[SERVICE]fetchWorkoutActionProfileProcess", r);
      if (r.error) {
        return [];
      }
      if (!r.data.v) {
        const data = r.data as WorkoutActionSteps;
        return [
          {
            text: data.start_position ?? data.startPosition,
            imgs: [],
            tips: ["起始姿势"],
          },
          ...data.steps.map((text) => {
            return {
              text,
              imgs: [],
              tips: [],
            };
          }),
        ];
      }
      const d = r.data as WorkoutActionStepsJSON250608;
      return d.steps;
    })(),
    extra_config: (() => {
      return {
        reps_unit: getSetValueUnit("次"),
      };
    })(),
    points: (() => {
      const r = parseJSONStr<string[]>(data.points);
      if (r.error) {
        return [];
      }
      return r.data;
    })(),
    problems: (() => {
      const r = parseJSONStr<WorkoutActionProblems[]>(data.problems);
      if (r.error) {
        return [];
      }
      return r.data.map((problem) => {
        return {
          title: problem.title,
          reason: problem.reason,
          solutions: problem.solutions ?? [],
        };
      });
    })(),
    equipments: idsMapValue(data.equipment_ids),
    muscles: idsMapValue(data.muscle_ids),
    primary_muscles: idsMapValue(data.primary_muscle_ids),
    secondary_muscles: idsMapValue(data.secondary_muscle_ids),
    alternative_actions: idsMapValue(data.alternative_action_ids),
    advanced_actions: idsMapValue(data.advanced_action_ids),
    regressed_actions: idsMapValue(data.regressed_action_ids),
  });
}

/**
 * 更新健身动作
 * @param body
 * @returns
 */
export function updateWorkoutAction(body: { id: number } & WorkoutActionBody) {
  return request.post<void>("/api/workout_action/update", {
    ...body,
  });
}

export function updateWorkoutActionIdx(body: { id: number; idx: number }) {
  return request.post<void>("/api/workout_action/update_idx", {
    ...body,
  });
}
