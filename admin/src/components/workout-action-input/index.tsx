import { For } from "@/timeless";

import { useViewModelStore } from "@/hooks";

import { WorkoutActionInputViewModel } from "@/biz/workout_action/workout_action_input";
import { Button } from "@/components/ui";

export function WorkoutActionInput(props: { store: WorkoutActionInputViewModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <div>
      <Button store={vm.ui.$btn_show_dialog}>选择</Button>
      <div>
        <For each={state().value}>
          {(v) => {
            return (
              <div>
                <div>{v.zh_name}</div>
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );
}
