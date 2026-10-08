/**
 * @file 后台/首页布局
 */
import { For, JSX, createSignal, onMount } from "@/timeless";
import { Portal } from "@/timeless";
import {
  Film,
  Users,
  FolderInput,
  Home,
  Bot,
  Flame,
  LogOut,
  Settings,
  Tv,
  File,
  CircuitBoard,
  Subtitles,
  AlarmClock,
  Folder,
  Sparkles,
  HeartCrack,
  HardDrive,
} from "@/timeless/icons";

import { ViewComponent, ViewComponentProps } from "@/store/types";
import { $media_resource } from "@/store";
import { PageKeys } from "@/store/routes";
import { Show } from "@/packages/ui/show";
import { Button, Checkbox, Dialog, DropdownMenu, Input, KeepAliveRouteView, Textarea } from "@/components/ui";
import { MediaResourceManageView } from "@/components/media-resource-view";
import { DragContainer } from "@/components/drag-container";

import {
  ButtonCore,
  CheckboxCore,
  DialogCore,
  DropdownMenuCore,
  InputCore,
  MenuCore,
  MenuItemCore,
} from "@/domains/ui";
import { RequestCore } from "@/domains/request";
import { cn, sleep } from "@/utils/index";
import { __VERSION__ } from "@/constants/index";

function Page(props: ViewComponentProps) {
  const { app, history, client, storage, pages, view } = props;
}

export const HomeLayout: ViewComponent = (props) => {
  const { app, history, client, storage, pages, view } = props;

  const logoutBtn = new ButtonCore({
    async onClick() {
      logoutBtn.setLoading(true);
      app.$user.logout();
      await sleep(2000);
      logoutBtn.setLoading(false);
    },
  });

  const userDropdown = new DropdownMenuCore({
    align: "start",
    side: "bottom",
    items: [
      new MenuItemCore({
        icon: <Settings class="w-4 h-4" />,
        label: "设置",
        onClick() {},
      }),
      new MenuItemCore({
        icon: <LogOut class="w-4 h-4" />,
        label: "退出登录",
        async onClick() {
          logoutBtn.setLoading(true);
          app.$user.logout();
          await sleep(2000);
          userDropdown.hide();
          logoutBtn.setLoading(false);
        },
      }),
    ],
  });
  const $canRegisterCheckbox = new CheckboxCore();
  const $noNeedCode = new CheckboxCore();

  const [curSubView, setCurSubView] = createSignal(view.curView);
  const [subViews, setSubViews] = createSignal(view.subViews);

  view.onSubViewsChange((nextSubViews) => {
    setSubViews(nextSubViews);
  });
  view.onCurViewChange((nextCurView) => {
    setCurSubView(nextCurView);
  });

  const [menus, setMenus] = createSignal<
    { text: string; icon: JSX.Element; badge?: boolean; url?: PageKeys; onClick?: () => void }[]
  >([
    {
      text: "首页",
      icon: <Home class="w-6 h-6" />,
      url: "root.home_layout.index",
    },
    {
      text: "肌肉列表",
      icon: <HeartCrack class="w-6 h-6" />,
      url: "root.home_layout.muscle",
    },
    {
      text: "器械列表",
      icon: <HardDrive class="w-6 h-6" />,
      url: "root.home_layout.equipment",
    },
    {
      text: "动作列表",
      icon: <HardDrive class="w-6 h-6" />,
      url: "root.home_layout.action_list",
    },
    {
      text: "订阅计划",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.subscription_plan_list",
    },
    {
      text: "用户管理",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.coach_list",
    },
    {
      text: "题库",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.paper_list",
    },
    {
      text: "礼品卡",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.gift_card_list",
    },
    {
      text: "训练计划列表",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.workout_plan_list",
    },
    {
      text: "训练计划集合列表",
      icon: <Tv class="w-6 h-6" />,
      url: "root.home_layout.plan_set_list",
    },
    {
      text: "学员列表",
      icon: <Users class="w-6 h-6" />,
      url: "root.home_layout.student_list",
    },
  ]);
  const [curRouteName, setCurRouteName] = createSignal(history.$router.name);

  // onMount(() => {
  // console.log("[PAGE]home/layout onMount", history.$router.href);
  // });
  history.onRouteChange(({ name }) => {
    setCurRouteName(name);
  });

  return (
    <>
      <div class="flex w-full h-full bg-w-bg-0">
        <div class="w-56 shrink-0 border-r border-border bg-muted/30 p-3">
          <div class="flex flex-col justify-between h-full w-full">
            <div class="flex-1 space-y-1 w-full overflow-y-auto">
              <div class="px-3 py-4 mb-2"><div class="text-lg font-semibold tracking-tight">Fithub</div><div class="text-xs text-muted-foreground">健身管理后台</div></div>
              <For each={menus()}>
                {(menu) => {
                  const { icon, text, url, badge, onClick } = menu;
                  return (
                    <Menu
                      app={app}
                      icon={icon}
                      history={history}
                      highlight={(() => {
                        return curRouteName() === url;
                      })()}
                      url={url}
                      badge={badge}
                      onClick={onClick}
                    >
                      {text}
                    </Menu>
                  );
                }}
              </For>
            </div>
            <div class="flex justify-center">
              {/* <Button class="" store={logoutBtn} variant="subtle" icon={<LogOut class="w-4 h-4" />}>
                退出登录
              </Button>
              <Button class="" store={settingsBtn} variant="subtle" icon={<Settings class="w-4 h-4" />}>
                设置
              </Button> */}
            </div>
          </div>
        </div>
        <div class="flex-1 min-w-0">
          <div class="relative w-full h-full">
            <For each={subViews()}>
              {(subView, i) => {
                const routeName = subView.name;
                const PageContent = pages[routeName as Exclude<PageKeys, "root">];
                return (
                  <KeepAliveRouteView
                    class={cn(
                      "absolute inset-0",
                      "data-[state=open]:animate-in data-[state=open]:fade-in",
                      "data-[state=closed]:animate-out data-[state=closed]:fade-out"
                    )}
                    store={subView}
                    index={i()}
                  >
                    <PageContent
                      app={app}
                      client={client}
                      storage={storage}
                      pages={pages}
                      history={history}
                      view={subView}
                    />
                  </KeepAliveRouteView>
                );
              }}
            </For>
          </div>
        </div>
      </div>
      <div class="absolute z-50 right-8 top-6">
        <DropdownMenu store={userDropdown}>
          <div class="w-12 h-12 rounded-full bg-slate-300"></div>
        </DropdownMenu>
      </div>
      <Portal>
        <DragContainer title="图片" storage={props.storage}>
          <MediaResourceManageView store={$media_resource}></MediaResourceManageView>
        </DragContainer>
      </Portal>
    </>
  );
};

function Menu(
  props: Pick<ViewComponentProps, "app" | "history"> & {
    highlight?: boolean;
    url?: PageKeys;
    icon: JSX.Element;
    badge?: boolean;
  } & JSX.HTMLAttributes<HTMLDivElement>
) {
  const button = new ButtonCore({ onClick() {
    if (props.url) props.history.push(props.url);
    else if (typeof props.onClick === "function") (props.onClick as () => void)();
  } });
  return <Button store={button} variant="ghost" icon={props.icon}
    class={cn("w-full justify-start px-3 h-10 text-sm", props.highlight ? "bg-accent text-accent-foreground" : "text-muted-foreground")}
    aria-current={props.highlight ? "page" : undefined}>
    {props.children}
    <Show when={props.badge}><span class="ml-auto h-2 w-2 rounded-full bg-destructive" /></Show>
  </Button>;
}
