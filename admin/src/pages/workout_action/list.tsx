/**
 * @file 健身动作列表
 */
import { For } from "@/timeless";
import { MoreHorizontal } from "@/timeless/icons";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input, ListView, ScrollView } from "@/components/ui";
import { Select } from "@/components/ui/select";

import {
  deleteWorkoutAction,
  fetchWorkoutActionList,
  fetchWorkoutActionListProcess,
} from "@/biz/workout_action/services";
import { base, Handler } from "@/domains/base";
import { ButtonCore, DialogCore, InputCore, ScrollViewCore, SelectCore } from "@/domains/ui";
import { RequestCore } from "@/domains/request";
import { ListCore } from "@/domains/list";
import { RefCore } from "@/domains/ui/cur";
import { Result } from "@/domains/result";
import { WorkoutActionType, WorkoutActionTypeOptions } from "@/biz/workout_action/constants";

import { WorkoutActionValuesView } from "./action_form";
import { WorkoutActionEditorViewModel } from "./model";

function WorkoutActionListViewModel(props: ViewComponentProps) {
  const request = {
    action: {
      list: new ListCore(
        new RequestCore(fetchWorkoutActionList, { process: fetchWorkoutActionListProcess, client: props.client }),
        {
          pageSize: 60,
        }
      ),
      delete: new RequestCore(deleteWorkoutAction, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    handleClickAction(v: { id: number }) {
      $editor = WorkoutActionEditorViewModel({
        app: props.app,
        client: props.client,
      });
      ui.$values = $editor.ui.$form;
      $editor.methods.fetch(v.id);
      ui.$dialog_update_editor.show();
      // props.history.push("root.home_layout.action_update", { id: String(v.id) });
    },
    handleClickTag(tag: { value: string }) {
      _cur_tag = tag.value;
      request.action.list.search({ tag: tag.value });
    },
    async deleteWorkoutAction(v: { id: number }) {
      const r = await request.action.delete.run({ id: v.id });
      if (r.error) {
        return Result.Err(r.error);
      }
      request.action.list.deleteItem((vv) => {
        return vv.id === v.id;
      });
      return Result.Ok(null);
    },
    async handleDeleteWorkoutAction(v: { id: number }) {
      ui.$ref_act.select(v);
      ui.$dialog_delete_confirm.show();
    },
    updateIdx(...args: Parameters<typeof $editor.methods.updateIdx>) {
      return $editor.methods.updateIdx(...args);
    },
  };

  let $editor = WorkoutActionEditorViewModel({
    app: props.app,
    client: props.client,
  });
  const ui = {
    $view: new ScrollViewCore({
      async onReachBottom() {
        await request.action.list.loadMore();
        ui.$view.finishLoadingMore();
      },
    }),
    $input_search_type: new SelectCore({
      defaultValue: WorkoutActionType.RESISTANCE,
      options: WorkoutActionTypeOptions,
      onChange(v) {
        request.action.list.search({ type: v });
      },
    }),
    $input_search_keyword: new InputCore({
      defaultValue: "",
      onEnter() {
        ui.$btn_search_submit.click();
      },
    }),
    $btn_search_reset: new ButtonCore({
      onClick() {
        ui.$input_search_keyword.reset();
        request.action.list.reset();
      },
    }),
    $btn_search_submit: new ButtonCore({
      async onClick() {
        const v = ui.$input_search_keyword.value;
        if (!v) {
          props.app.tip({
            text: ["请输入关键词"],
          });
          return;
        }
        ui.$btn_search_submit.setLoading(true);
        const r = await request.action.list.search({ keyword: v });
        ui.$btn_search_submit.setLoading(false);
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
      },
    }),
    $goto_create_btn: new ButtonCore({
      onClick() {
        props.history.push("root.home_layout.action_create");
      },
    }),
    $ref_act: new RefCore<{ id: number }>(),
    $dialog_delete_confirm: new DialogCore({
      async onOk() {
        const v = ui.$ref_act.value;
        if (!v) {
          props.app.tip({
            text: ["请选择要删除的动作"],
          });
          return;
        }
        ui.$dialog_delete_confirm.okBtn.setLoading(true);
        const r = await methods.deleteWorkoutAction(v);
        ui.$dialog_delete_confirm.okBtn.setLoading(false);
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["删除成功"],
        });
        ui.$ref_act.clear();
        ui.$dialog_delete_confirm.hide();
      },
    }),
    $dialog_create_editor: new DialogCore({}),
    $dialog_update_editor: new DialogCore({
      async onOk() {
        ui.$dialog_update_editor.okBtn.setLoading(true);
        const r = await $editor.methods.update();
        ui.$dialog_update_editor.okBtn.setLoading(false);
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["更新成功"],
        });
        ui.$dialog_update_editor.hide();
      },
    }),
    $values: $editor.ui.$form,
  };

  let _cur_tag = "";
  let _tags = ["", "胸", "肩", "腿", "背", "臀"];
  let _state = {
    get response() {
      return request.action.list.response;
    },
    get tags() {
      return _tags.map((v) => {
        return {
          value: v,
          text: v === "" ? "全部" : v,
          selected: _cur_tag === v,
        };
      });
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  request.action.list.onStateChange(() => methods.refresh());

  return {
    state: _state,
    ui,
    request,
    methods,
    ready() {
      request.action.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function WorkoutActionListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(WorkoutActionListViewModel, [props]);

  return (
    <>
      <ScrollView store={vm.ui.$view} class="p-4">
        <h1 class="text-2xl font-bold mb-4">动作列表</h1>
        <div class="flex items-center gap-2 mt-2">
          <Button store={vm.ui.$btn_search_reset}>重置</Button>
          <Button store={vm.ui.$goto_create_btn}>新增动作</Button>
        </div>
        <div class="mt-2">
          <div class="flex items-center gap-2">
            <div class="w-[180px]">
              <Select store={vm.ui.$input_search_type} />
            </div>
            <Input class="flex-1" store={vm.ui.$input_search_keyword} />
            <Button class="w-[120px]" store={vm.ui.$btn_search_submit}>
              搜索
            </Button>
          </div>
          <div class="flex gap-2 mt-2">
            {state().tags.map((tag) => {
              return (
                <div
                  class="px-4 py-2 rounded border-2 cursor-pointer"
                  classList={{
                    "border-w-fg-3 bg-w-bg-5": tag.selected,
                    "border-w-bg-3": !tag.selected,
                  }}
                  onClick={() => {
                    vm.methods.handleClickTag(tag);
                  }}
                >
                  {tag.text}
                </div>
              );
            })}
          </div>
        </div>
        <div class="py-4">
          <ListView store={vm.request.action.list}>
            <div class="actions grid grid-cols-6 gap-2 p-2">
              <For each={state().response.dataSource}>
                {(v) => (
                  <div
                    classList={{
                      "relative p-2 border-2 border-w-bg-5 rounded-md text-w-fg-1": true,
                    }}
                  >
                    <div class="text-xl">{v.zh_name}</div>
                    <input
                      class="w-[88px]"
                      value={v.idx}
                      onChange={(event) => {}}
                      onKeyDown={(event) => {
                        if (event.key === "Enter") {
                          vm.methods.updateIdx({
                            id: v.id,
                            idx: Number(event.currentTarget.value),
                          });
                        }
                      }}
                    />
                    <div class="flex items-center gap-2 mt-4">
                      <div
                        onClick={() => {
                          vm.methods.handleDeleteWorkoutAction(v);
                        }}
                      >
                        删除
                      </div>
                      <div
                        onClick={() => {
                          vm.methods.handleClickAction(v);
                        }}
                      >
                        详情
                      </div>
                    </div>
                  </div>
                )}
              </For>
            </div>
          </ListView>
        </div>
      </ScrollView>
      <Dialog store={vm.ui.$dialog_update_editor}>
        <div class="h-[80vh] overflow-y-auto">
          <div class="p-4 rounded-lg">
            <div class="flex flex-col gap-4">
              <WorkoutActionValuesView store={vm.ui.$values} />
            </div>
          </div>
        </div>
      </Dialog>
      <Dialog store={vm.ui.$dialog_delete_confirm}>
        <div>确认删除吗？</div>
      </Dialog>
    </>
  );
}
