import { base, BaseDomain, Handler } from "@/domains/base";
import { BizError } from "@/domains/error";

import { debounce } from "@/utils/lodash/debounce";

export function StorageCore<T extends Record<string, any>>(props: {
  key: string;
  values: any;
  defaultValues: T;
  client: {
    setItem: (key: string, value: string) => void;
    getItem: (key: string) => void;
  };
}) {
  const methods = {
    refresh() {
      bus.emit(Events.StateChange, { ..._state });
    },
    get<K extends keyof T>(key: K, defaultValue?: T[K]) {
      const v = _values[key];
      // console.log("[DOMAIN]storage/index - get", key, v, this.values);
      if (v === undefined) {
        if (defaultValue) {
          return defaultValue[key] as T[K];
        }
        throw new Error("the default value no existing");
      }
      return v as T[K];
    },
    set: debounce(100, <K extends keyof T>(key: K, values: T[K]) => {
      // console.log("cache set", key, values);
      const next_values = {
        ..._values,
        [key]: values,
      };
      _values = next_values;
      _client.setItem(_key, JSON.stringify(_values));
      methods.refresh();
    }) as (key: keyof T, value: unknown) => void,

    merge: <K extends keyof T>(
      key: K,
      values: Partial<T[K]>,
      extra: Partial<{ reverse: boolean; limit: number }> = {}
    ) => {
      // console.log("[]merge", key, values);
      const prevValues = methods.get(key) || {};
      if (Array.isArray(prevValues)) {
        let nextValues = extra.reverse
          ? [...(values as unknown as Array<unknown>), ...prevValues]
          : [...prevValues, ...(values as unknown as Array<unknown>)];
        if (extra.limit) {
          nextValues = nextValues.slice(0, extra.limit);
        }
        methods.set(key, nextValues);
        return nextValues;
      }
      if (typeof prevValues === "object" && typeof values === "object") {
        const nextValues = {
          ...prevValues,
          ...values,
        };
        methods.set(key, nextValues);
        return nextValues;
      }
      console.warn("the params of merge must be object");
      return prevValues;
    },
    clear<K extends keyof T>(key: K) {
      const v = _values[key];
      if (v === undefined) {
        return null;
      }
      _values = {
        ..._values,
        [key]: _defaultValues[key],
      };
      _client.setItem(_key, JSON.stringify(_values));
      methods.refresh();
    },
  };
  const ui = {};

  let _key = props.key;
  let _values: T = {
    ...props.defaultValues,
    ...props.values,
  };
  let _defaultValues = props.defaultValues;
  let _client = props.client;
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

  const ins = {
    methods,
    ui,
    state: _state,
    get values() {
      return _values;
    },
    get: methods.get,
    set: methods.set,
    clear: methods.clear,
    merge: methods.merge,
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
  return ins;
}

export type StorageCore<T extends Record<string, any>> = ReturnType<typeof StorageCore<T>>;

// export class StorageCore<T extends Record<string, unknown>> extends BaseDomain<TheTypesOfEvents<T>> {
//   key: string;
//   values: T;
//   defaultValues: T;
//   client: StorageCoreProps<T>["client"];

//   get state() {
//     return {
//       values: this.values,
//     };
//   }

//   constructor(props: Partial<{ _name: string }> & StorageCoreProps<T>) {
//     super(props);

//     const { key, client, defaultValues, values } = props;
//     this.key = key;
//     this.values = values;
//     this.values = {
//       ...defaultValues,
//       ...values,
//     };
//     this.defaultValues = defaultValues;
//     this.client = client;
//   }
// }
