import { createSignal, Show } from "@/timeless";
import { JSX } from "@/timeless";

import { useViewModelStore } from "@/hooks";

import { DragZoneCore } from "@/domains/ui/drag-zone";
import { readFileAsArrayBuffer, readFileAsURL } from "@/utils/browser";
import { ImageCore } from "@/domains/ui";
import { FileAboutType, ImageUploadCore } from "@/domains/ui/form/image-upload/index";
import { cn } from "@/utils";

import { DragZone } from "./drag-zone";
import { LazyImage } from "./image";
import { AspectRatio } from "./aspect-ratio";

export function ImageUpload(props: { store: ImageUploadCore } & JSX.HTMLAttributes<HTMLDivElement>) {
  const { store } = props;

  const [state, vm] = useViewModelStore(props.store);

  return (
    <div class={cn(props.class, "relative")}>
      <Show when={state().type === FileAboutType.Image}>
        <div class="absolute inset-0 h-full">
          <img
            class="w-full h-full object-cover"
            src={state().url}
            onLoad={(event) => {
              const elm = event.currentTarget;
              const profile = {
                width: elm.naturalWidth,
                height: elm.naturalHeight,
              };
              vm.setProfile(profile);
            }}
          />
        </div>
      </Show>
      <Show when={state().type === FileAboutType.Video}>
        <div class="absolute inset-0 flex justify-center h-full">
          <video
            src={state().url}
            class="object-cover"
            onLoad={(event) => {
              const elm = event.currentTarget;
              const profile = {
                width: elm.videoWidth,
                height: elm.videoHeight,
                duration: elm.duration,
              };
              vm.setProfile(profile);
            }}
          />
        </div>
      </Show>
      <DragZone store={store.ui.$zone}></DragZone>
    </div>
  );
}
