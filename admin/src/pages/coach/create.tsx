import { createCoach } from "@/biz/coach/service";
import { Button, Input, Textarea } from "@/components/ui";
import { FieldObjV2 } from "@/components/ui/field-obj2";
import { FieldV2 } from "@/components/ui/fieldv2";
import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { RequestCore } from "@/domains/request";
import { ButtonCore, InputCore } from "@/domains/ui";
import { ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { useViewModel } from "@/hooks";
import { ViewComponentProps } from "@/store/types";

function InfluencerCreateViewModel(props: ViewComponentProps) {
  const request = {
    coach: {
      create: new RequestCore(createCoach, { client: props.client }),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $form: new ObjectFieldCore({
      name: "",
      label: "",
      fields: {
        nickname: new SingleFieldCore({
          name: "nickname",
          label: "昵称",
          input: new InputCore({ defaultValue: "" }),
        }),
        avatar_url: new SingleFieldCore({
          name: "avatar_url",
          label: "头像",
          input: new InputCore({ defaultValue: "" }),
        }),
        bio: new SingleFieldCore({
          name: "bio",
          label: "简介",
          input: new InputCore({ defaultValue: "" }),
        }),
      },
    }),
    $btn_create: new ButtonCore({
      async onClick() {
        const r = await ui.$form.validate();
        if (r.error) {
          return;
        }
        const body = {
          nickname: r.data.nickname,
          avatar_url: r.data.avatar_url,
          bio: r.data.bio,
        };
        const r2 = await request.coach.create.run(body);
        if (r2.error) {
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
    destroy() {
      bus.destroy();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function InfluencerCreateView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(InfluencerCreateViewModel, [props]);

  return (
    <div>
      <div>创建</div>
      <FieldObjV2 store={vm.ui.$form}>
        <FieldV2 store={vm.ui.$form.fields.nickname}>
          <Input store={vm.ui.$form.fields.nickname.input} />
        </FieldV2>
        <FieldV2 store={vm.ui.$form.fields.avatar_url}>
          <Input store={vm.ui.$form.fields.avatar_url.input} />
        </FieldV2>
        <FieldV2 store={vm.ui.$form.fields.bio}>
          <Textarea store={vm.ui.$form.fields.bio.input} />
        </FieldV2>
      </FieldObjV2>
      <Button store={vm.ui.$btn_create}>创建</Button>
    </div>
  );
}
