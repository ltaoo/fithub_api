import { ViewComponentProps } from "@/store/types";

import { base, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";
import { RequestCore } from "@/domains/request";

import {
  deleteMediaResource,
  fetchMediaResourceList,
  fetchMediaResourceListProcess,
  uploadMediaResource,
} from "./service";
import { ListCore } from "@/domains/list";
import { ImageUploadCore } from "@/domains/ui/form/image-upload";
import { ButtonCore, ScrollViewCore } from "@/domains/ui";
import { QiniuOSS } from "@/biz/oss/qiniu";
import { StorageCore } from "@/domains/storage";
import { Result } from "@/domains/result";

enum MediaUploadStatus {
  Wait = 1,
  Uploading = 2,
  Success = 3,
  Failed = 4,
}

export function MediaResourceManagerModel(props: {
  app: ViewComponentProps["app"];
  storage: StorageCore<{ token: string }>;
  client: ViewComponentProps["client"];
}) {
  const request = {
    resource: {
      create: new RequestCore(uploadMediaResource, { client: props.client }),
      delete: new RequestCore(deleteMediaResource, { client: props.client }),
      list: new ListCore(
        new RequestCore(fetchMediaResourceList, { process: fetchMediaResourceListProcess, client: props.client })
      ),
    },
    qiniu: QiniuOSS({ storage: props.storage, client: props.client }),
  };
  request.qiniu.methods.set_env(props.app.prod);
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    ready() {},
    async createMediaResource(v: { hash: string; key: string }) {
      const v2 = ui.$input_file.state.value;
      if (!v2) {
        return Result.Err("找不到上传的文件");
      }
      const body = {
        hash: v.hash,
        key: v.key,
        type: v2.type,
        width: v2.width,
        height: v2.height,
        size: v2.size,
        filetype: v2.filetype,
        duration: v2.duration,
        filename: v2.filename,
      };
      const r = await request.resource.create.run(body);
      if (r.error) {
        return;
      }
      request.resource.list.refresh();
    },
    async deleteMediaResource(v: { id: number }) {
      const r = await request.resource.delete.run(v);
      if (r.error) {
        return Result.Err(r.error);
      }
      request.resource.list.deleteItem((vv) => {
        return v.id === vv.id;
      });
      return Result.Ok(null);
    },
    async handleDeleteMedia(v: { id: number }) {
      return methods.deleteMediaResource(v);
    },
    handleCopyURL(v: { url: string }) {
      props.app.copy(v.url);
      props.app.tip({
        text: ["复制成功"],
      });
    },
  };
  const ui = {
    $scroll: new ScrollViewCore({
      async onReachBottom() {
        await request.resource.list.loadMore();
        ui.$scroll.finishLoadingMore();
      },
    }),
    $input_file: ImageUploadCore({}),
    $btns_delete: new Map<number, ButtonCore>(),
    $btn_upload: new ButtonCore({
      async onClick() {
        const v = ui.$input_file.value;
        if (!v) {
          return;
        }
        _error = null;
        _status = MediaUploadStatus.Uploading;
        methods.refresh();
        ui.$btn_upload.setLoading(true);
        const r = await request.qiniu.upload(v.file, request.qiniu.filename({ scope: "media_resource" }));
        ui.$btn_upload.setLoading(false);
        if (r.error) {
          _error = r.error;
          _status = MediaUploadStatus.Failed;
          methods.refresh();
          return;
        }
        // const r = await request.resource.create
      },
    }),
  };

  let _status = MediaUploadStatus.Wait;
  let _error: BizError | null = null;
  let _state = {
    get response() {
      return request.resource.list.response;
    },
    get error() {
      return _error;
      // return request.qiniu.state
    },
    get status() {
      return _status;
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

  request.resource.list.onStateChange(() => methods.refresh());
  request.resource.create.onStateChange(() => methods.refresh());
  request.qiniu.onCompleted((v) => {
    _status = MediaUploadStatus.Success;
    methods.refresh();
    methods.createMediaResource(v);
    // request.resource.list.refresh();
  });
  request.qiniu.onError((e) => {
    _error = e;
    _status = MediaUploadStatus.Failed;
    methods.refresh();
  });

  return {
    request,
    methods,
    ui,
    state: _state,
    ready() {
      request.resource.list.init();
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export type MediaResourceManagerModel = ReturnType<typeof MediaResourceManagerModel>;
