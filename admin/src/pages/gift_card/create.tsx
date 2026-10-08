/**
 * @file 礼品卡
 */
import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input, Textarea } from "@/components/ui";
import { FieldV2 } from "@/components/ui/fieldv2";
import { FieldArrV2 } from "@/components/ui/field-arrv2";
import { TagInput } from "@/components/ui/tag-input";
import { Select } from "@/components/ui/select";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ButtonCore, InputCore, SelectCore } from "@/domains/ui";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { TagInputCore } from "@/domains/ui/form/tag-input";
import { RequestCore } from "@/domains/request";
import { Result } from "@/domains/result";
import { createGiftCard, fetchGiftCardRewardList } from "@/biz/gift_card/services";
import { ListCore } from "@/domains/list";
import { GiftCardRewardSelectViewModel } from "@/biz/gift_card/reward_select";
import { For } from "@/timeless";

function GiftCardCreateViewModel(props: ViewComponentProps) {
  const request = {
    gift_card: {
      create: new RequestCore(createGiftCard, { client: props.client }),
    },
    gift_card_reward: {
      list: new ListCore(new RequestCore(fetchGiftCardRewardList, { client: props.client })),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async create() {
      const r = await ui.$values.validate();
      if (r.error) {
        return Result.Err(r.error.message);
      }
      const data = r.data;
      if (!data.num) {
        const tip = "请输入数量";
        return Result.Err(tip);
      }
      const body = {
        num: Number(data.num),
        gift_card_reward_id: Number(data.gift_card_reward),
      };
      console.log("[PAGE]quiz/create - before quiz.create.run", body);
      const r2 = await request.gift_card.create.run(body);
      if (r2.error) {
        return Result.Err(r2.error.message);
      }
      return Result.Ok(r2.data);
    },
  };
  const ui = {
    $btn_submit: new ButtonCore({
      async onClick() {
        const r = await methods.create();
        if (r.error) {
          props.app.tip({
            text: [r.error.message],
          });
          return;
        }
        props.app.tip({
          text: ["创建成功"],
        });
      },
    }),
    $select_reward: GiftCardRewardSelectViewModel({
      defaultValue: [],
      list: request.gift_card_reward.list,
    }),
    $btn_select_reward: new ButtonCore({
      onClick() {
        ui.$select_reward.request.WorkoutSchedule.list.init();
        ui.$select_reward.ui.$dialog.show();
      },
    }),
    $values: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        num: new SingleFieldCore({
          label: "数量",
          name: "num",
          input: new InputCore({ defaultValue: "" }),
        }),
        gift_card_reward: new SingleFieldCore({
          label: "奖励",
          name: "gift_card_reward",
          input: new InputCore({ defaultValue: 0 }),
        }),
      },
    }),
  };
  let _state = {
    get rewards() {
      return ui.$select_reward.state.list;
    },
  };
  enum Events {
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  ui.$select_reward.onStateChange(() => methods.refresh());
  ui.$select_reward.ui.$dialog.onOk(() => {
    const value = ui.$select_reward.value;
    if (value.length === 0) {
      return;
    }
    const v = value[0];
    ui.$select_reward.ui.$dialog.hide();
    ui.$values.fields.gift_card_reward.setValue(v.id);
  });

  return {
    methods,
    ui,
    state: _state,
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function GiftCardCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(GiftCardCreateViewModel, [props]);

  return (
    <>
      <div>
        <div>
          <div class="w-[680px] mx-auto space-y-4">
            <FieldV2 store={vm.ui.$values.fields.num}>
              <Input store={vm.ui.$values.fields.num.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.gift_card_reward}>
              <Input store={vm.ui.$values.fields.gift_card_reward.input} />
            </FieldV2>
            <div>
              <Button store={vm.ui.$btn_select_reward}>添加选项</Button>
            </div>
          </div>
          <Button store={vm.ui.$btn_submit}>创建</Button>
        </div>
      </div>
      <Dialog store={vm.ui.$select_reward.ui.$dialog}>
        <For each={state().rewards}>
          {(v) => {
            return (
              <div
                classList={{
                  "border-w-fg-3": v.selected,
                }}
                onClick={() => {
                  vm.ui.$select_reward.methods.select(v);
                }}
              >
                <div>{v.name}</div>
              </div>
            );
          }}
        </For>
      </Dialog>
    </>
  );
}
