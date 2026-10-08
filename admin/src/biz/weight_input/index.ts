import { base, Handler } from "@/domains/base";
import { InputCore, PopoverCore } from "@/domains/ui";

interface WeightInputViewModelProps {
  defaultValue?: number;
  onChange?: (value: number, unit: "kg" | "lbs") => void;
}

export function WeightInputViewModel(props: WeightInputViewModelProps) {
  const methods = {
    handleClickNumber(num: string) {
      (() => {
        if (_text === "0") {
          _text = num;
          return;
        }
        if (_text.includes(".")) {
          const [integer, decimal] = _text.split(".");
          if (decimal.length >= 2) {
            return;
          }
          _text = _text + num;
          return;
        }
        if (_text.length >= 4) {
          return;
        }
        _text = _text + num;
      })();
      bus.emit(Events.StateChange, { ..._state });
    },
    handleClickUnit(unit: "kg" | "lbs") {
      if (_unit === unit) {
        return;
      }
      _unit = unit;
      bus.emit(Events.StateChange, { ..._state });
    },
    handleClickDot() {
      if (_text[_text.length - 1] === ".") {
        return;
      }
      if (_text.includes(".")) {
        return;
      }
      _text = _text + ".";
      bus.emit(Events.StateChange, { ..._state });
    },
    handleClickDelete() {
      _text = _text.slice(0, -1);
      if (_text.length === 0) {
        _text = "0";
      }
      bus.emit(Events.StateChange, { ..._state });
    },
  };
  const ui = {
    $input: new InputCore({
      defaultValue: "",
    }),
    $popover: new PopoverCore(),
  };
  let _text = props.defaultValue ? props.defaultValue.toString() : "0";
  let _unit = "kg";
  let _state = {
    get text() {
      return _text;
    },
    get unit() {
      return _unit;
    },
  };
  enum Events {
    StateChange,
  }
  type TheTypesOfEvents = {
    [Events.StateChange]: typeof _state;
  };
  const bus = base<TheTypesOfEvents>();

  return {
    state: _state,
    methods,
    ui,
    get value() {
      return {
        text: _text,
        unit: _unit,
      };
    },
    setValue(value: { text: string; unit: "kg" | "lbs" }) {
      _text = value.text;
      _unit = value.unit;
      bus.emit(Events.StateChange, { ..._state });
    },
    ready() {},
    onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
      return bus.on(Events.StateChange, handler);
    },
  };
}

export type WeightInputViewModel = ReturnType<typeof WeightInputViewModel>;
