import { For } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Input, ScrollView, Textarea } from "@/components/ui";
import { FieldV2 } from "@/components/ui/fieldv2";
import { FieldArrV2 } from "@/components/ui/field-arrv2";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { Result } from "@/domains/result";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { ButtonCore, InputCore, ScrollViewCore } from "@/domains/ui";
import { RequestCore } from "@/domains/request";
import { createSubscriptionPlan } from "@/biz/subscription/services";

function SubscriptionPlanCreateViewModel(props: ViewComponentProps) {
  const request = {
    subscription_plan: {
      create: new RequestCore(createSubscriptionPlan, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async submit() {
      const r = await ui.$values.validate();
      if (r.error) {
        return Result.Err(r.error.message);
      }
      const values = r.data;
      if (!values.name) {
        const tip = "请输入计划名称";
        return Result.Err(tip);
      }
      if (!values.unit_price) {
        const tip = "请输入价格";
        return Result.Err(tip);
      }
      const body = {
        name: values.name,
        details: values.details,
        unit_price: Number(values.unit_price),
        discount_policies: values.discount_policies.map((policy) => {
          return {
            name: policy.name,
            rate: Number(policy.rate),
            count_require: Number(policy.count_require),
            enabled: Number(policy.enabled),
          };
        }),
      };
      const r2 = await request.subscription_plan.create.run(body);
      if (r2.error) {
        return Result.Err(r2.error.message);
      }
      return Result.Ok(r2.data);
    },
  };
  const ui = {
    $view: new ScrollViewCore({}),
    $values: new ObjectFieldCore({
      label: "",
      name: "",
      fields: {
        name: new SingleFieldCore({
          label: "名称",
          name: "name",
          input: new InputCore({
            defaultValue: "",
          }),
        }),
        details: new SingleFieldCore({
          label: "概要",
          name: "details",
          input: new InputCore({
            defaultValue: "",
          }),
        }),
        unit_price: new SingleFieldCore({
          label: "单位价格（单位分）",
          name: "unit_price",
          input: new InputCore({
            defaultValue: 0,
          }),
        }),
        discount_policies: new ArrayFieldCore({
          label: "折扣方案",
          name: "discount_policies",
          field: () => {
            return new ObjectFieldCore({
              label: "",
              name: "",
              fields: {
                name: new SingleFieldCore({
                  label: "名称",
                  name: "name",
                  input: new InputCore({ defaultValue: "" }),
                }),
                rate: new SingleFieldCore({
                  label: "折扣比率",
                  name: "rate",
                  input: new InputCore({ defaultValue: 100 }),
                }),
                count_require: new SingleFieldCore({
                  label: "所需购买数量",
                  name: "count_require",
                  input: new InputCore({ defaultValue: 0 }),
                }),
                enabled: new SingleFieldCore({
                  label: "是否生效",
                  name: "enabled",
                  input: new InputCore({ defaultValue: 1 }),
                }),
                // enabled_at: new SingleFieldCore({
                //   label: "生效时间",
                //   name: "enabled_at",
                //   input: new InputCore({ defaultValue: 1 }),
                // }),
                // expired_at: new SingleFieldCore({
                //   label: "失效时间",
                //   name: "expired_at",
                //   input: new InputCore({ defaultValue: 1 }),
                // }),
              },
            });
          },
        }),
      },
    }),
    $btn_discount_append: new ButtonCore({
      onClick() {
        ui.$values.fields.discount_policies.append();
      },
    }),
    $btn_subscription_plan_submit: new ButtonCore({
      async onClick() {
        const r = await methods.submit();
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
  };
  let _state = {};
  enum Events {
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

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

export function SubscriptionPlanCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(SubscriptionPlanCreateViewModel, [props]);

  return (
    <ScrollView store={vm.ui.$view}>
      <div class="p-2">
        <div class="w-[680px] mx-auto space-y-4">
          <FieldV2 store={vm.ui.$values.fields.name}>
            <Input store={vm.ui.$values.fields.name.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.details}>
            <Textarea store={vm.ui.$values.fields.details.input} />
          </FieldV2>
          <FieldV2 store={vm.ui.$values.fields.unit_price}>
            <Input store={vm.ui.$values.fields.unit_price.input} />
          </FieldV2>
          <div>
            <Button store={vm.ui.$btn_discount_append}>添加折扣方案</Button>
            <div class="flex items-center gap-4 mt-2 p-2">
              <FieldArrV2
                store={vm.ui.$values.fields.discount_policies}
                render={(field) => {
                  return (
                    <div class="">
                      <FieldV2 store={field.fields.name}>
                        <Input store={field.fields.name.input} />
                      </FieldV2>
                      <FieldV2 store={field.fields.rate}>
                        <Input store={field.fields.rate.input} />
                      </FieldV2>
                      <FieldV2 store={field.fields.enabled}>
                        <Input store={field.fields.enabled.input} />
                      </FieldV2>
                      <FieldV2 store={field.fields.count_require}>
                        <Input store={field.fields.count_require.input} />
                      </FieldV2>
                    </div>
                  );
                }}
              ></FieldArrV2>
            </div>
          </div>
        </div>
      </div>
      <div>
        <Button store={vm.ui.$btn_subscription_plan_submit}>提交</Button>
      </div>
    </ScrollView>
  );
}
