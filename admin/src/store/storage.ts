import { StorageCore } from "@/domains/storage";

export const storage = StorageCore(
  (() => {
    const DEFAULT_CACHE_VALUES = {
      user: {
        id: "",
        username: "anonymous",
        token: "",
        avatar: "",
      },
      jobs: [] as string[],
      media_manager: {
        x: 0,
        y: 0,
      },
    };
    const key = "a_global";
    const e = globalThis.localStorage.getItem(key);
    return {
      key,
      values: e ? JSON.parse(e) : {},
      defaultValues: DEFAULT_CACHE_VALUES,
      client: globalThis.localStorage,
    };
  })()
);

export const qiniu_storage = StorageCore(
  (() => {
    const QINIU_DEFAULT_CACHE_VALUES = {
      token: "",
    };
    const key = "qiniu";
    const existing = globalThis.localStorage.getItem(key);
    return {
      key,
      defaultValues: QINIU_DEFAULT_CACHE_VALUES,
      values: existing ? JSON.parse(existing) : QINIU_DEFAULT_CACHE_VALUES,
      client: globalThis.localStorage,
    };
  })()
);
