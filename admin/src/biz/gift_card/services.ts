import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result } from "@/domains/result";
import { request } from "@/biz/requests";
import { ListResponse } from "@/biz/requests/types";

export function createGiftCardReward(body: { name: string; overview: string; details: Record<string, any> }) {
  return request.post("/api/gift_card/create_reward", {
    name: body.name,
    overview: body.overview,
    details: JSON.stringify(body.details),
  });
}

export function createGiftCard(body: { num: number; gift_card_reward_id: number }) {
  return request.post("/api/gift_card/create", {
    num: body.num,
    gift_card_reward_id: body.gift_card_reward_id,
  });
}

export function fetchGiftCardList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      code: string;
      gift_card_reward: { name: string };
    }>
  >("/api/gift_card/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}

export function fetchGiftCardRewardList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      name: string;
    }>
  >("/api/gift_card/reward_list", {
    page_size: body.pageSize,
    page: body.page,
  });
}
