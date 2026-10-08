import { FetchParams } from "@/domains/list/typing";
import { request } from "@/biz/requests";
import { ListResponse } from "@/biz/requests/types";

export function fetchSubscriptionPlanList(body: FetchParams) {
  return request.post<
    ListResponse<{
      name: string;
      created_at: string;
    }>
  >("/api/subscription_plan/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}

export function createSubscriptionPlan(body: {
  name: string;
  details: string;
  unit_price: number;
  discount_policies: {
    name: string;
    rate: number;
    count_require: number;
    enabled: number;
  }[];
}) {
  return request.post("/api/subscription_plan/create", body);
}

export function updateSubscriptionPlan(body: {}) {
  return request.post("/api/subscription_plan/update", {});
}
