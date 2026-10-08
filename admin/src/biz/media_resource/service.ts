import dayjs from "dayjs";

import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result } from "@/domains/result";
import { request } from "@/biz/requests";
import { ListResponse } from "@/biz/requests/types";
import { FileAboutType, FileAboutTypeTextMap } from "@/domains/ui/form/image-upload";

export function fetchQiniuToken() {
  return request.post<{ token: string }>("/api/media/qiniu_token", {});
}

export function uploadMediaResource(body: {
  type: number;
  width: number;
  height: number;
  size: number;
  duration: number;
  filename: string;
  filetype: string;
  key: string;
  hash: string;
}) {
  return request.post("/api/media/create", body);
}

export function deleteMediaResource(body: { id: number }) {
  return request.post("/api/media/delete", body);
}

export function fetchMediaResourceList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      media_type: FileAboutType;
      width: number;
      height: number;
      size: number;
      duration: number;
      filename: string;
      filetype: string;
      key: string;
      hash: string;
      url: string;
      created_at: string;
    }>
  >("/api/media/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}

export function fetchMediaResourceListProcess(r: TmpRequestResp<typeof fetchMediaResourceList>) {
  if (r.error) {
    return Result.Err(r.error);
  }
  const resp = r.data;
  return Result.Ok({
    ...resp,
    list: resp.list.map((v) => {
      return {
        ...v,
        type_text: FileAboutTypeTextMap[v.media_type],
        created_at: dayjs(v.created_at).format("YYYY-MM-DD HH:mm"),
      };
    }),
  });
}
