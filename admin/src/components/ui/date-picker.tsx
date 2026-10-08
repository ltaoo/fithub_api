import { vm } from "@timeless/timeless";
import { DatePicker as ShadcnDatePicker } from "@timeless/shadcn";
import { onCleanup } from "@/timeless";
import { DatePickerCore } from "@/domains/ui/date-picker";

export function DatePicker(props: { store: DatePickerCore }): any {
  const store = props.store;
  const model = vm.DatePickerCore({ today: store.value || new Date() });
  let syncing = false;
  if (store.value) model.setValue(store.value);
  onCleanup(model.onChange(value => {
    if (!syncing && value) { syncing = true; store.setValue(value); syncing = false; }
  }));
  onCleanup(store.$calendar.onChange(() => {
    if (!syncing && store.value) { syncing = true; model.setValue(store.value); syncing = false; }
  }));
  return ShadcnDatePicker({ store: model });
}
