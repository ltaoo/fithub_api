import { Show } from "@/timeless";

import { ViewComponentProps } from "@/store/types";
import { useViewModel } from "@/hooks";
import { Button, Dialog, Input } from "@/components/ui";
import { FieldObjV2 } from "@/components/ui/field-obj2";
import { FieldV2 } from "@/components/ui/fieldv2";
import { WorkoutActionSelectDialogView } from "@/components/workout-action-select-dialog";

import { createCoachContent, fetchCoachProfile } from "@/biz/coach/service";
import { fetchWorkoutActionList, fetchWorkoutActionListProcess } from "@/biz/workout_action/services";
import { WorkoutActionSelectDialogViewModel } from "@/biz/workout_action_select_dialog";
import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { ListCore } from "@/domains/list";
import { RequestCore } from "@/domains/request";
import { ButtonCore, InputCore, SelectCore } from "@/domains/ui";
import { ArrayFieldCore, ObjectFieldCore, SingleFieldCore } from "@/domains/ui/formv2";
import { WorkoutActionInputViewModel } from "@/biz/workout_action/workout_action_input";
import { Select } from "@/components/ui/select";
import { WorkoutActionInput } from "@/components/workout-action-input";
import { FileAboutType, ImageUploadCore } from "@/domains/ui/form/image-upload";
import { ImageUpload } from "@/components/ui/image-upload";
import { readFileAsURL, readFileAsText } from "@/utils/browser";
import { parseJSONStr } from "@/utils";

enum InfluencerContentType {
  Video = 1,
  ShortVideo = 2,
  Image = 3,
  Text = 4,
}
const InfluencerContentTypeTextMap: Record<InfluencerContentType, string> = {
  [InfluencerContentType.Video]: "视频",
  [InfluencerContentType.ShortVideo]: "短视频",
  [InfluencerContentType.Image]: "图片",
  [InfluencerContentType.Text]: "纯文本",
};

function InfluencerContentUploadViewModel(props: ViewComponentProps) {
  const request = {
    coach: {
      profile: new RequestCore(fetchCoachProfile, { client: props.client }),
    },
    content: {
      create: new RequestCore(createCoachContent, { client: props.client }),
    },
    workout_action: {
      list: new ListCore(
        new RequestCore(fetchWorkoutActionList, { process: fetchWorkoutActionListProcess, client: props.client })
      ),
    },
  };
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    async create() {
      const r = await ui.$form.validate();
      if (r.error) {
        return;
      }
      const values = r.data;
      // if (values.workout_action_id.length === 0) {
      //   props.app.tip({
      //     text: ["请选择关联动作"],
      //   });
      //   return;
      // }
      // if (values.workout_action_id.length !== 1) {
      //   props.app.tip({
      //     text: ["只能关联一个动作"],
      //   });
      //   return;
      // }
      const body = {
        coach_id: Number(props.view.query.id),
        content_type: Number(values.content_type),
        title: values.title,
        description: values.description,
        content_url: values.content_url,
        video_key: values.video_key,
        // workout_action_id: Number(values.workout_action_id[0].id),
        // start_point: values.video_start_point,
      };
      const r2 = await request.content.create.run(body);
      if (r2.error) {
        return;
      }
      props.app.tip({
        text: ["创建成功"],
      });
    },
  };

  const $workout_action_select = WorkoutActionSelectDialogViewModel({
    defaultValue: [],
    list: request.workout_action.list,
    client: props.client,
  });
  const ui = {
    $form: new ObjectFieldCore({
      name: "",
      label: "",
      fields: {
        title: new SingleFieldCore({
          name: "title",
          label: "标题",
          input: new InputCore({ defaultValue: "" }),
        }),
        description: new SingleFieldCore({
          name: "description",
          label: "描述",
          input: new InputCore({ defaultValue: "" }),
        }),
        content_type: new SingleFieldCore({
          name: "content_type",
          label: "类型",
          input: new SelectCore({
            defaultValue: InfluencerContentType.Video,
            options: [
              InfluencerContentType.Video,
              InfluencerContentType.ShortVideo,
              InfluencerContentType.Image,
              InfluencerContentType.Text,
            ].map((t) => {
              return {
                value: t,
                label: InfluencerContentTypeTextMap[t],
              };
            }),
          }),
        }),
        content_url: new SingleFieldCore({
          name: "content_url",
          label: "来源链接",
          input: new InputCore({ defaultValue: "" }),
        }),
        // video_key: new ArrayFieldCore({
        //   name: "video_key",
        //   label: "",
        //   field() {
        //     return new SingleFieldCore({
        //       name: "",
        //       label: "",
        //       input: new InputCore({ defaultValue: "" }),
        //     });
        //   },
        // }),
        video_key: new SingleFieldCore({
          name: "video_key",
          label: "视频地址",
          input: new InputCore({ defaultValue: "" }),
        }),
        image_keys: new ArrayFieldCore({
          name: "image_keys",
          label: "",
          field() {
            return new SingleFieldCore({
              name: "",
              label: "",
              input: new InputCore({ defaultValue: "" }),
            });
          },
        }),
        workout_action_id: new SingleFieldCore({
          label: "关联动作",
          name: "workout_action_id",
          input: WorkoutActionInputViewModel({
            $select: $workout_action_select,
          }),
        }),
        video_start_point: new SingleFieldCore({
          label: "视频空降点",
          name: "video_start_point",
          input: new InputCore({
            defaultValue: 0,
            type: "number",
          }),
        }),
      },
    }),
    $workout_action_select,
    $json: ImageUploadCore({
      defaultValue: "",
      async onChange(v) {
        if (v.type === FileAboutType.Json) {
          const file = v.file;
          const r = await readFileAsText(file);
          if (r.error) {
            props.app.tip({
              text: [r.error.message],
            });
            return;
          }
          const r2 = parseJSONStr<{
            url: string;
            author: string;
            filename: string;
            video_key: string;
            thumbnail: string;
          }>(r.data);
          if (r2.error) {
            props.app.tip({
              text: [r2.error.message],
            });
            return;
          }
          const data = r2.data;
          const prefix = "//static.fithub.top/";
          ui.$form.setValue({
            title: data.filename,
            content_url: data.url,
            video_key: prefix + data.video_key,
            // cover_image_url: prefix + data.thumbnail,
          });
        }
      },
    }),
    $btn_create: new ButtonCore({
      onClick() {
        methods.create();
      },
    }),
  };
  let _state = {
    get profile() {
      return request.coach.profile.response;
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

  request.coach.profile.onStateChange(() => methods.refresh());

  return {
    methods,
    ui,
    state: _state,
    ready() {
      const id = Number(props.view.query.id);
      if (Number.isNaN(id)) {
        return;
      }
      request.coach.profile.run({ id });
    },
    destroy() {
      bus.destroy();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export function InfluencerContentUploadView(props: ViewComponentProps) {
  const [state, vm] = useViewModel(InfluencerContentUploadViewModel, [props]);

  return (
    <>
      <div>
        <div class="p-4">
          <Show when={state().profile}>
            <div class="flex items-center gap-2">
              <div>
                <img class="w-[60px] h-[60px]" src={state().profile?.avatar_url} />
              </div>
              <div>{state().profile?.nickname}</div>
            </div>
          </Show>
        </div>
        {/* The BEST Deadlift Tutorial | Step by Step */}
        {/* https://www.youtube.com/shorts/McCDaAsSeRc */}
        <div class="relative">
          <ImageUpload class="w-[320px] h-[120px]" store={vm.ui.$json} />
        </div>
        <div class="shadow-xl p-4">
          <FieldObjV2 store={vm.ui.$form}>
            <FieldV2 store={vm.ui.$form.fields.content_type}>
              <Select store={vm.ui.$form.fields.content_type.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.title}>
              <Input store={vm.ui.$form.fields.title.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.description}>
              <Input store={vm.ui.$form.fields.description.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.content_url}>
              <Input store={vm.ui.$form.fields.content_url.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.video_key}>
              <Input store={vm.ui.$form.fields.video_key.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.content_url}>
              <WorkoutActionInput store={vm.ui.$form.fields.workout_action_id.input} />
            </FieldV2>
            <FieldV2 store={vm.ui.$form.fields.video_start_point}>
              <Input store={vm.ui.$form.fields.video_start_point.input} />
            </FieldV2>
          </FieldObjV2>
        </div>
        <Button store={vm.ui.$btn_create}>创建</Button>
      </div>
      <Dialog store={vm.ui.$workout_action_select.ui.$dialog}>
        <WorkoutActionSelectDialogView store={vm.ui.$workout_action_select} />
      </Dialog>
    </>
  );
}
