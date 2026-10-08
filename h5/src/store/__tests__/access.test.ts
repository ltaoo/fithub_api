// @vitest-environment jsdom
import { expect, test, vi } from "vitest";
import { routes } from "../routes";
import { requires_login, login_return_path } from "../access";
import { HistoryCore } from "@/domains/history";
import { NavigatorCore } from "@/domains/navigator";
import { RouteViewCore } from "@/domains/route_view";

test("catalog is public while training, edits and personal records require login", () => {
  for (const name of ["root", "root.workout_plan_list", "root.workout_plan_profile", "root.workout_schedule_profile", "root.workout_action_list"] as const) expect(requires_login(routes[name])).toBe(false);
  for (const name of ["root.workout_day", "root.workout_day_prepare", "root.workout_plan_create", "root.workout_plan_update", "root.home_layout.mine", "root.home_layout.student_list"] as const) expect(requires_login(routes[name])).toBe(true);
  expect(login_return_path("/workout_plan_profile", { id: "7", student_id: "8" })).toBe("/workout_plan_profile?id=7&student_id=8");
});
test("internal pushes, replacements and stack resets honor the navigation guard", () => {
  const router = new NavigatorCore({ location: window.location });
  const view = new RouteViewCore({ name: "root", pathname: "/", title: "root", visible: true, parent: null });
  const history = new HistoryCore({ view, router, routes, views: { root: view } as any });
  const guard = vi.fn(() => false); history.beforeNavigate = guard;
  history.push("root.workout_day"); history.replace("root.workout_day"); history.destroyAllAndPush("root.workout_day");
  expect(guard).toHaveBeenCalledTimes(3);
  expect(history.stacks).toHaveLength(0);
});
