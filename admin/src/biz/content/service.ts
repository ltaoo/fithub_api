import dayjs from "dayjs";
import { request } from "@/biz/requests";
import { ListResponse, ListResponseWithCursor } from "@/biz/requests/types";
import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result, UnpackedResult } from "@/domains/result";
import { TheResponseOfFetchFunction } from "@/domains/request";
import { Unpacked } from "@/types";
import { parseJSONStr } from "@/utils";

export function fetchCoachContentList(body: Partial<FetchParams>) {
  return request.post<
    ListResponse<{
      id: number;
      title: string;
      video_key: string;
    }>
  >("/api/coach/content/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}
