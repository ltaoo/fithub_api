/**
 * @file 礼品卡奖励
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
import { QuizTypes } from "@/biz/quiz/constants";
import { createGiftCardReward } from "@/biz/gift_card/services";
import { SubscriptionPlanSelectViewModel } from "@/biz/subscription/subscription_select";
import { ListCore } from "@/domains/list";
import { fetchSubscriptionPlanList } from "@/biz/subscription/services";
import { For } from "@/timeless";

function GiftCardRewardCreateViewModel(props: ViewComponentProps) {
  const request = {
    gift_card_reward: {
      create: new RequestCore(createGiftCardReward, { client: props.client }),
    },
    subscription_plan: {
      list: new ListCore(new RequestCore(fetchSubscriptionPlanList, { client: props.client })),
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
      if (!data.name) {
        const tip = "请输入题目内容";
        return Result.Err(tip);
      }
      const body = {
        name: data.name,
        overview: data.overview,
        details: {
          day_count: Number(data.details.day_count),
          subscription_plan_id: Number(data.details.subscription),
        },
      };
      console.log("[PAGE]quiz/create - before quiz.create.run", body);
      const r2 = await request.gift_card_reward.create.run(body);
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
    $btn_select_subscription_plan: new ButtonCore({
      onClick() {
        ui.$select_subscription_plan.request.WorkoutSchedule.list.init();
        ui.$select_subscription_plan.ui.$dialog.show();
      },
    }),
    $select_subscription_plan: SubscriptionPlanSelectViewModel({
      defaultValue: [],
      multiple: false,
      list: request.subscription_plan.list,
    }),
    $values: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        name: new SingleFieldCore({
          label: "名称",
          name: "name",
          input: new InputCore({ defaultValue: "" }),
        }),
        overview: new SingleFieldCore({
          label: "概要",
          name: "overview",
          input: new InputCore({ defaultValue: "" }),
        }),
        details: new ObjectFieldCore({
          label: "",
          name: "details",
          fields: {
            subscription: new SingleFieldCore({
              label: "订阅",
              name: "subscription",
              input: new InputCore({ defaultValue: 0 }),
            }),
            day_count: new SingleFieldCore({
              label: "天数",
              name: "day_count",
              input: new InputCore({
                defaultValue: 30,
              }),
            }),
          },
        }),
      },
    }),
  };
  let _state = {
    get subscription_plan() {
      return ui.$select_subscription_plan.state.list;
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

  ui.$select_subscription_plan.onStateChange(() => methods.refresh());
  ui.$select_subscription_plan.ui.$dialog.onOk(() => {
    const value = ui.$select_subscription_plan.value;
    if (value.length === 0) {
      return;
    }
    const v = value[0];
    ui.$select_subscription_plan.ui.$dialog.hide();
    ui.$values.fields.details.fields.subscription.setValue(v.id);
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

export function GiftCardRewardCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(GiftCardRewardCreateViewModel, [props]);

  return (
    <>
      <div>
        <div>
          <div class="w-[680px] mx-auto space-y-4">
            <FieldV2 store={vm.ui.$values.fields.name}>
              <Textarea store={vm.ui.$values.fields.name.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$values.fields.overview}>
              <Textarea store={vm.ui.$values.fields.overview.input} />
            </FieldV2>
            <div>
              <Button store={vm.ui.$btn_select_subscription_plan}>选择订阅</Button>
              <FieldV2 store={vm.ui.$values.fields.details.fields.subscription}>
                <Input store={vm.ui.$values.fields.details.fields.subscription.input} />
              </FieldV2>
              <FieldV2 store={vm.ui.$values.fields.details.fields.day_count}>
                <Input store={vm.ui.$values.fields.details.fields.day_count.input} />
              </FieldV2>
            </div>
          </div>
          <Button store={vm.ui.$btn_submit}>创建</Button>
        </div>
      </div>
      <Dialog store={vm.ui.$select_subscription_plan.ui.$dialog}>
        <For each={state().subscription_plan}>
          {(plan) => {
            return (
              <div
                classList={{
                  "border-w-fg-3": plan.selected,
                }}
                onClick={() => {
                  vm.ui.$select_subscription_plan.methods.select(plan);
                }}
              >
                <div>{plan.name}</div>
              </div>
            );
          }}
        </For>
      </Dialog>
    </>
  );
}
