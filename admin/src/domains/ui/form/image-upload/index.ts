import { base, Handler } from "@/domains/base";
import { ImageCore } from "@/domains/ui/image";
import { DragZoneCore } from "@/domains/ui/drag-zone";
import { BizError } from "@/domains/error";
import { readFileAsURL } from "@/utils/browser";

export enum FileAboutType {
  Unknown = 0,
  Image = 1,
  Video = 2,
  Audio = 3,
  Markdown = 4,
  Exe = 5,
  Yaml = 6,
  Json = 7,
  Zip = 8,
}
export const FileAboutTypeTextMap: Record<FileAboutType, string> = {
  [FileAboutType.Unknown]: "未知",
  [FileAboutType.Image]: "图片",
  [FileAboutType.Video]: "视频",
  [FileAboutType.Audio]: "音频",
  [FileAboutType.Exe]: "应用程序",
  [FileAboutType.Yaml]: "Yaml",
  [FileAboutType.Json]: "JSON",
  [FileAboutType.Markdown]: "markdown",
  [FileAboutType.Zip]: "压缩包",
};

type FileUploadValue = {
  type: FileAboutType;
  width: number;
  height: number;
  duration: number;
  size: number;
  filename: string;
  filetype: string;
  file: File;
};

export function ImageUploadCore(props: {
  tip?: string;
  defaultValue?: string;
  onChange?: (file: FileUploadValue) => void;
}) {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $zone: new DragZoneCore({ tip: props.tip }),
  };

  let _url = "";
  let _file_about_type = FileAboutType.Unknown;
  let _file: null | File = null;
  let _profile = {
    width: 0,
    height: 0,
    duration: 0,
  };
  const _state = {
    get value() {
      if (!_file) {
        return null;
      }
      return {
        type: _file_about_type,
        width: _profile.width,
        height: _profile.height,
        duration: _profile.duration,
        size: _file.size,
        filename: _file.name,
        filetype: _file.type,
        file: _file,
      };
    },
    get url() {
      return _url;
    },
    get type() {
      return _file_about_type;
    },
    get file() {
      return _file;
    },
  };
  enum Events {
    Change,
    StateChange,
    Error,
  }
  type TheTypesOfEvents = {
    [Events.Change]: string;
    [Events.StateChange]: typeof _state;
    [Events.Error]: BizError;
  };
  const bus = base<TheTypesOfEvents>();

  ui.$zone.onChange(async (files) => {
    const file = files[0];
    if (!file) {
      return;
    }
    const type = file.type;
    _file_about_type = (() => {
      if (type.match(/audio\//)) {
        return FileAboutType.Audio;
      }
      if (type.match(/video\//)) {
        return FileAboutType.Video;
      }
      if (type.match(/image\//)) {
        return FileAboutType.Image;
      }
      if (type.match(/x-msdownload/)) {
        return FileAboutType.Exe;
      }
      if (type.match(/x-yaml$/)) {
        return FileAboutType.Yaml;
      }
      if (type.match(/json$/)) {
        return FileAboutType.Json;
      }
      if (type.match(/zip$/)) {
        return FileAboutType.Zip;
      }
      if (type.match(/text\/markdown/)) {
        return FileAboutType.Markdown;
      }
      return FileAboutType.Unknown;
    })();
    if ([FileAboutType.Image, FileAboutType.Video].includes(_file_about_type)) {
      const r = await readFileAsURL(file);
      if (r.error) {
        bus.emit(Events.Error, r.error);
        return;
      }
      _url = r.data;
      bus.emit(Events.Change, _url);
    }
    _file = file;
    if (props.onChange && _state.value) {
      props.onChange(_state.value);
    }
    methods.refresh();
  });

  return {
    shape: "image-upload" as const,
    ui,
    state: _state,
    get value() {
      return _state.value;
    },
    get defaultValue() {
      return props.defaultValue ?? null;
    },
    setValue(url: string) {
      _url = url;
      bus.emit(Events.Change, _url);
    },
    setProfile(v: Partial<{ width: number; height: number; duration: number }>) {
      _profile = Object.assign({ ..._profile }, v);
    },
    ready() {},
    onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
      return bus.on(Events.Change, handler);
    },
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export type ImageUploadCore = ReturnType<typeof ImageUploadCore>;
