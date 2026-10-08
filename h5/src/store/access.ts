import { RouteConfig } from "@/domains/route_view/utils";

export function requires_login(route: Pick<RouteConfig<string>, "options"> | undefined) {
  return !!route?.options?.require?.includes("login");
}

// Only accept known internal routes; never redirect a login to an external URL.
export function login_return_path(pathname: string, query: Record<string, string>) {
  const search = new URLSearchParams(query).toString();
  return pathname + (search ? "?" + search : "");
}
