import dayjs from "dayjs";

import { request } from "@/biz/requests";
import { idsMapValue } from "@/biz/services/utils";
import { ListResponse, ListResponseWithCursor } from "@/biz/requests/types";
import { TheResponseOfFetchFunction } from "@/domains/request";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result } from "@/domains/result";
import { FetchParams } from "@/domains/list/typing";
import { parseJSONStr, seconds_to_hour_with_template, seconds_to_hour_template1 } from "@/utils";

import { WorkoutPlanStepType, WorkoutPlanSetType } from "./constants";
import {
  WorkoutPlanStepBody,
  WorkoutPlanActionPayload,
  WorkoutPlanPreviewPayload,
  WorkoutPlanStepJSON250607,
} from "./types";

export function createWorkoutPlan(body: {
  title: string;
  overview: string;
  tags: string;
  level: number;
  details: string;
  steps: WorkoutPlanStepBody[];
  estimated_duration: number;
  points: string;
  suggestions: string;
  muscle_ids: string;
  equipment_ids: string;
}) {
  return request.post<void>("/api/workout_plan/create", body);
}

export function updateWorkoutPlan(body: {
  id: number | string;
  title: string;
  overview: string;
  tags: string;
  level: number;
  steps: WorkoutPlanStepBody[];
  estimated_duration: number;
  points: string;
  suggestions: string;
  muscle_ids: string;
  equipment_ids: string;
}) {
  return request.post<void>("/api/workout_plan/update", {
    ...body,
    id: Number(body.id),
  });
}

export function fetchWorkoutPlanProfile(body: { id: number | string }) {
  return request.post<{
    id: number;
    title: string;
    overview: string;
    tags: string;
    level: number;
    estimated_duration: number;
    suggestions: string;
    steps: WorkoutPlanStepJSON250607[];
    muscle_ids: string;
    equipment_ids: string;
    creator: {
      nickname: string;
      avatar_url: string;
      is_self: boolean;
    };
    created_at: string;
  }>("/api/workout_plan/profile", { id: Number(body.id) });
}
export function fetchWorkoutPlanProfileProcess(r: TmpRequestResp<typeof fetchWorkoutPlanProfile>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  const v = r.data;
  return Result.Ok({
    id: v.id,
    title: v.title,
    overview: v.overview,
    tags: v.tags.split(",").filter(Boolean),
    level: v.level,
    steps: v.steps,
    estimated_duration: v.estimated_duration,
    estimated_duration_text: seconds_to_hour_with_template(v.estimated_duration, seconds_to_hour_template1),
    suggestions: (() => {
      const r = parseJSONStr<string[]>(v.suggestions);
      if (r.error) {
        return "";
      }
      if (r.data[0]) {
        return r.data[0];
      }
      return "";
    })(),
    creator: v.creator,
    muscle_ids: v.muscle_ids
      .split(",")
      .filter(Boolean)
      .map((v) => Number(v)),
    equipment_ids: v.equipment_ids
      .split(",")
      .filter(Boolean)
      .map((v) => Number(v)),
    created_at: dayjs(v.created_at).format("YYYY-MM-DD HH:mm"),
  });
}

export function fetchWorkoutPlanList(body: FetchParams) {
  return request.post<
    ListResponseWithCursor<{
      id: number | string;
      title: string;
      overview: string;
      tags: string;
      level: number;
      estimated_duration: number;
      details: string;
      points: string;
      suggestions: string;
      muscle_ids: string;
      equipment_ids: string;
    }>
  >("/api/workout_plan/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}
export function fetchWorkoutPlanListProcess(r: TmpRequestResp<typeof fetchWorkoutPlanList>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  return Result.Ok({
    ...r.data,
    list: r.data.list.map((plan) => {
      return {
        id: plan.id,
        title: plan.title,
        overview: plan.overview,
        tags: plan.tags.split(",").filter(Boolean),
        level: plan.level,
        estimated_duration: plan.estimated_duration,
        details: (() => {
          const r = parseJSONStr<WorkoutPlanPreviewPayload>(plan.details);
          if (r.error) {
            return null;
          }
          return r.data;
        })(),
      };
    }),
  });
}

export function createContentOfWorkoutPlan(body: { content_id: number; workout_plan_id: number; details: string }) {
  return request.post("/api/workout_plan/content/create", body);
}

export function createWorkoutPlanSet(body: {
  title: string;
  overview: string;
  icon_url?: string;
  idx: number;
  details: {
    type: number;
    id: number;
  }[];
}) {
  return request.post("/api/workout_plan_set/create", body);
}

export function updateWorkoutPlanSet(body: {
  id: number;
  title: string;
  overview: string;
  icon_url?: string;
  idx: number;
  details: string;
}) {
  return request.post("/api/workout_plan_set/update", body);
}

export function fetchWorkoutScheduleList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: string;
      title: string;
      overview: string;
      // tags: string;
    }>
  >("/api/workout_schedule/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}
