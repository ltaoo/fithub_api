import { Accessor, createSignal, onMount } from "@/timeless";

export function useViewModel<
  T extends (...args: any[]) => {
    state: any;
    ready: () => void;
    onStateChange: (handler: (v: any) => any) => any;
  }
>(factory: T, props: Parameters<T>): [Accessor<ReturnType<T>["state"]>, ReturnType<T>] {
  // @ts-ignore
  const $model: ReturnType<T> = factory(...props);

  const [state, setState] = createSignal($model.state);

  $model.onStateChange((v) => {
    setState(v);
  });
  onMount(() => {
    $model.ready();
  });

  return [state, $model];
}

export function useViewModelStore<
  T extends {
    state: any;
    ready: () => void;
    onStateChange: (handler: (v: any) => any) => any;
  }
>(vm: T): [Accessor<T["state"]>, T] {
  const [state, setState] = createSignal(vm.state);

  vm.onStateChange((v) => {
    setState(v);
  });
  vm.ready();
  return [state, vm];
}
