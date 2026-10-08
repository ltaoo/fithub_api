import dayjs from "dayjs";
import { request } from "@/biz/requests";
import { ListResponse, ListResponseWithCursor } from "@/biz/requests/types";
import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result, UnpackedResult } from "@/domains/result";
import { TheResponseOfFetchFunction } from "@/domains/request";
import { Unpacked } from "@/types";
import { parseJSONStr } from "@/utils";

export function fetchCoachList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      nickname: string;
      avatar_url: string;
    }>
  >("/api/coach/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}

export function createCoach(body: { nickname: string; avatar_url: string; bio?: string }) {
  return request.post("/api/coach/create", body);
}

export function fetchCoachProfile(body: { id: number }) {
  return request.post<{
    id: number;
    nickname: string;
    avatar_url: string;
  }>("/api/admin/coach/profile", { id: body.id });
}

export function createCoachContent(body: {
  coach_id: number;
  content_type: number;
  title: string;
  description: string;
  content_url: string;
  video_key: string;
  // workout_action_id: number;
  // start_point: number;
}) {
  return request.post("/api/coach/content/create", body);
}

export function createCoachAuthURL(body: { id: number }) {
  return request.post<{ url: string }>("/api/admin/coach/auth_url", body);
}
