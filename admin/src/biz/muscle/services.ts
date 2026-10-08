import { request } from "@/biz/requests";
import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result, UnpackedResult } from "@/domains/result";

export function createMuscle(params: {
  name: string;
  zh_name: string;
  tags: string;
  overview: string;
  features: string;
}) {
  return request.post<void>("/api/muscle/create", params);
}

export function fetchMuscleList(body: Partial<FetchParams> & { ids: number[] }) {
  return request.post<{
    list: {
      id: number;
      name: string;
      zh_name: string;
      tags: string;
      overview: string;
      features: string;
    }[];
    total: number;
  }>("/api/muscle/list", {
    page_size: body.pageSize,
    page: body.page,
    ids: body.ids,
  });
}
export function fetchMuscleListProcess(r: TmpRequestResp<typeof fetchMuscleList>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  return Result.Ok({
    ...r.data,
    list: r.data.list.map((muscle) => {
      return {
        id: muscle.id,
        name: muscle.name,
        zh_name: muscle.zh_name,
        tags: muscle.tags.split(",").filter(Boolean),
        overview: muscle.overview,
        features: JSON.parse(muscle.features),
      };
    }),
  });
}
export type MuscleProfile = UnpackedResult<ReturnType<typeof fetchMuscleListProcess>>["list"][number];

export function updateMuscle(params: { id: number | string; name: string; zh_name: string }) {
  return request.post<void>("/api/muscle/update", params);
}

export function deleteMuscle(params: { id: number | string }) {
  return request.post<void>("/api/muscle/delete", params);
}

export function fetchMuscleProfileById(params: { id: number | string }) {
  return request.post<MuscleProfile>("/api/muscle/profile", params);
}
