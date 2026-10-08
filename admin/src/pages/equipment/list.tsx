import { Card } from "@/components/ui/card";
/**
 * 器械信息管理
 */
import { For } from "@/timeless";
import { MoreHorizontal } from "@/timeless/icons";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input, ListView, ScrollView, Textarea } from "@/components/ui";
import { ImageURLInputModel, ImageURLInputView } from "@/components/image-url-input";
import { FieldObjV2 } from "@/components/ui/field-obj2";
import { FieldV2 } from "@/components/ui/fieldv2";
import { TagInput } from "@/components/ui/tag-input";
import { FieldArrV2 } from "@/components/ui/field-arrv2";

import {
  createEquipment,
  fetchEquipmentList,
  fetchEquipmentListProcess,
  updateEquipment,
} from "@/biz/equipment/services";
import { base, Handler } from "@/domains/base";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { InputCore } from "@/domains/ui/form/input";
import { RequestCore, TheResponseOfRequestCore } from "@/domains/request";
import { ScrollViewCore, SelectCore, DialogCore, ButtonCore } from "@/domains/ui";
import { TagInputCore } from "@/domains/ui/form/tag-input";
import { TmpRequestResp, UnpackedRequestPayload, RequestPayload } from "@/domains/request/utils";
import { ListCore } from "@/domains/list";
import { RefCore } from "@/domains/ui/cur";
import { TheItemTypeFromListCore } from "@/domains/list/typing";
import { Result } from "@/domains/result";
import { Unpacked } from "@/types";

import { EquipmentValueView } from "./equipment_form";

function EquipmentListViewModel(props: ViewComponentProps) {
  const request = {
    equipment: {
      list: new ListCore(
        new RequestCore(fetchEquipmentList, { process: fetchEquipmentListProcess, client: props.client }),
        {
          pageSize: 48,
        }
      ),
      create: new RequestCore(createEquipment, { client: props.client }),
      update: new RequestCore(updateEquipment, { client: props.client }),
    },
  };
  type TheEquipment = TheItemTypeFromListCore<typeof request.equipment.list>;
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async toBody() {
      const r = await ui.$form.validate();
      if (r.error) {
        return Result.Err(r.error);
      }
      const values = r.data;
      console.log("[PAGE]muscle/list - toBody", values);
      const { name, zh_name, alias, overview, tags, sort_idx, medias } = values;
      const body = {
        name,
        zh_name,
        alias,
        overview,
        tags: tags.join(","),
        sort_idx,
        medias: JSON.stringify(medias),
      };
      return Result.Ok(body);
    },
    async create() {
      const r = await methods.toBody();
      if (r.error) {
        props.app.tip({
          text: [r.error.message],
        });
        return;
      }
      const body = r.data;
      const r2 = await request.equipment.create.run(body);
      if (r2.error) {
        props.app.tip({
          text: [r2.error.message],
        });
        return;
      }
      props.app.tip({
        text: ["创建成功"],
      });
      request.equipment.list.refresh();
    },
    async update() {
      const vv = ui.$ref.value;
      if (!vv) {
        return;
      }
      const r = await methods.toBody();
      if (r.error) {
        props.app.tip({
          text: [r.error.message],
        });
        return;
      }
      const body = r.data;
      // console.log(body);
      const r2 = await request.equipment.update.run({
        ...body,
        id: vv.id,
      });
      if (r2.error) {
        props.app.tip({
          text: [r2.error.message],
        });
        return;
      }
      props.app.tip({
        text: ["编辑成功"],
      });
      ui.$dialog_update.hide();
    },
    showDialogEdit(v: TheEquipment) {
      ui.$ref.select(v);
      console.log(v.medias);
      ui.$form.setValue(v);
      ui.$dialog_update.show();
    },
  };
  const ui = {
    $view: new ScrollViewCore({}),
    $ref: new RefCore<TheEquipment>(),
    $btn_show_create_dialog: new ButtonCore({
      onClick() {
        ui.$dialog_create.show();
      },
    }),
    $dialog_create: new DialogCore({
      title: "创建",
      async onOk() {
        methods.create();
      },
    }),
    $dialog_update: new DialogCore({
      title: "编辑",
      async onOk() {
        methods.update();
      },
    }),
    $form: new ObjectFieldCore({
      name: "",
      label: "",
      fields: {
        name: new SingleFieldCore({
          name: "name",
          label: "英文名称",
          input: new InputCore({ defaultValue: "" }),
        }),
        zh_name: new SingleFieldCore({
          name: "zh_name",
          label: "中文名称",
          input: new InputCore({ defaultValue: "" }),
        }),
        alias: new SingleFieldCore({
          name: "alias",
          label: "别名",
          input: new InputCore({ defaultValue: "" }),
        }),
        overview: new SingleFieldCore({
          name: "overview",
          label: "概述",
          input: new InputCore({ defaultValue: "", type: "textarea" }),
        }),
        tags: new SingleFieldCore({
          name: "tags",
          label: "标签",
          input: new TagInputCore({}),
        }),
        sort_idx: new SingleFieldCore({
          name: "sort_idx",
          label: "排序",
          input: new InputCore({ defaultValue: 0, type: "number" }),
        }),
        medias: new ObjectFieldCore({
          name: "medias",
          label: "",
          fields: {
            pics: new ArrayFieldCore({
              name: "pics",
              label: "图片",
              field() {
                return new SingleFieldCore({
                  name: "",
                  label: "",
                  input: ImageURLInputModel(),
                });
              },
            }),
          },
        }),
      },
    }),
  };
  let _loading = false;
  let _state = {
    get loading() {
      return _loading;
    },
    get list() {
      return request.equipment.list.response.dataSource;
    },
  };

  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  request.equipment.list.onStateChange(() => bus.emit(Events.StateChange, { ..._state }));
  ui.$dialog_create.onCancel(() => {
    ui.$form.clear();
  });
  ui.$dialog_update.onCancel(() => {
    ui.$form.clear();
  });

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.equipment.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function EquipmentListView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(EquipmentListViewModel, [props]);

  return (
    <>
      <ScrollView store={vm.ui.$view} class="p-4">
        <h1 class="text-2xl font-bold mb-4">器械列表</h1>
        <div>
          <Button store={vm.ui.$btn_show_create_dialog}>创建器械</Button>
        </div>
        <ListView store={vm.request.equipment.list} class="equipments grid grid-cols-2 md:grid-cols-4 xl:grid-cols-6 gap-3 py-4">
          <For each={state().list}>
            {(v) => {
              return (
                <Card
                  class="relative p-3 flex justify-between text-foreground"
                  onClick={() => {}}
                >
                  <div class="absolute inset-0 p-2">
                    <div class="overflow-hidden">
                      <div class="">{v.zh_name}</div>
                    </div>
                    <div>{v.sort_idx}</div>
                    <div class="flex gap-2">
                      <div
                        class="px-2 py-1 text-sm border-2 border-w-fg-3 rounded-full text-center"
                        onClick={() => {
                          vm.methods.showDialogEdit(v);
                        }}
                      >
                        编辑
                      </div>
                    </div>
                  </div>
                  <div class="w-full" style="padding-bottom: 100%"></div>
                </Card>
              );
            }}
          </For>
        </ListView>
      </ScrollView>
      <Dialog store={vm.ui.$dialog_create}>
        <div class="w-[720px]">
          <FieldObjV2 store={vm.ui.$form}>
            <FieldV2 store={vm.ui.$form.fields.zh_name}>
              <Input store={vm.ui.$form.fields.zh_name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.name}>
              <Input store={vm.ui.$form.fields.name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.alias}>
              <Input store={vm.ui.$form.fields.alias.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.overview}>
              <Textarea store={vm.ui.$form.fields.overview.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.tags}>
              <TagInput store={vm.ui.$form.fields.tags.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.sort_idx}>
              <Input store={vm.ui.$form.fields.sort_idx.input} />
            </FieldV2>
            <FieldArrV2
              store={vm.ui.$form.fields.medias.fields.pics}
              render={(field) => {
                return (
                  <FieldV2 store={field}>
                    <ImageURLInputView store={field.input} />
                  </FieldV2>
                );
              }}
            ></FieldArrV2>
          </FieldObjV2>
        </div>
      </Dialog>
      <Dialog store={vm.ui.$dialog_update}>
        <div class="w-[720px]">
          <FieldObjV2 store={vm.ui.$form}>
            <FieldV2 store={vm.ui.$form.fields.zh_name}>
              <Input store={vm.ui.$form.fields.zh_name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.name}>
              <Input store={vm.ui.$form.fields.name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.alias}>
              <Input store={vm.ui.$form.fields.alias.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.overview}>
              <Textarea store={vm.ui.$form.fields.overview.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.tags}>
              <TagInput store={vm.ui.$form.fields.tags.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.sort_idx}>
              <Input store={vm.ui.$form.fields.sort_idx.input} />
            </FieldV2>
            <FieldObjV2 store={vm.ui.$form.fields.medias}>
              <FieldArrV2
                store={vm.ui.$form.fields.medias.fields.pics}
                render={(field) => {
                  return (
                    <FieldV2 store={field}>
                      <ImageURLInputView store={field.input} />
                    </FieldV2>
                  );
                }}
              ></FieldArrV2>
            </FieldObjV2>
          </FieldObjV2>
        </div>
      </Dialog>
    </>
  );
}
