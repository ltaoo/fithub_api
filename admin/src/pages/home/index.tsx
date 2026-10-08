import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
/**
 * @file 首页
 */
import { createSignal, For, Show } from "@/timeless";
import { Send, FileSearch, RefreshCcw, AlertTriangle, Loader, Bird, BarChart } from "@/timeless/icons";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { ScrollView } from "@/components/ui";
import { base, Handler } from "@/domains/base";
import { ScrollViewCore } from "@/domains/ui";

function HomeIndexPageViewModel(props: ViewComponentProps) {
  const $view = new ScrollViewCore();
  const _state = {};
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  return {
    state: _state,
    ui: {
      $view: $view,
    },
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export const HomeIndexPage = (props: ViewComponentProps) => {
  const { app, history, client, view } = props;

  const [state, $model] = useViewModel(HomeIndexPageViewModel, [props]);

  return (
    <>
      <ScrollView store={$model.ui.$view} class="h-screen p-6 md:p-8">
        <div class="flex items-center space-x-4">
          <h1 class="text-2xl">
            <div>数据统计</div>
          </h1>
          <div class="flex items-center space-x-2"></div>
        </div>
        <div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-4">
          <Card><CardHeader><CardTitle class="text-base">健身总天数</CardTitle></CardHeader><CardContent><p class="text-sm text-muted-foreground">暂无统计数据</p></CardContent></Card>
          <Card><CardHeader><CardTitle class="text-base">管理入口</CardTitle></CardHeader><CardContent><p class="text-sm text-muted-foreground">从左侧导航管理动作、训练计划和学员。</p></CardContent></Card>
        </div>
      </ScrollView>
    </>
  );
};
