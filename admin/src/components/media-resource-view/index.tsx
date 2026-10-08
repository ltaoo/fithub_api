import { For, Show } from "@/timeless";

import { useViewModel, useViewModelStore } from "@/hooks";
import { Button, ListView, ScrollView } from "@/components/ui";
import { ImageUpload } from "@/components/ui/image-upload";

import { MediaResourceManagerModel } from "@/biz/media_resource/media_resource";
import { FileAboutType } from "@/domains/ui/form/image-upload";

export function MediaResourceManageView(props: { store: MediaResourceManagerModel }) {
  const [state, vm] = useViewModelStore(props.store);

  return (
    <div>
      <div>
        <ImageUpload class="h-[240px]" store={vm.ui.$input_file}></ImageUpload>
        <Button class="mt-2 w-full" store={vm.ui.$btn_upload}>
          上传
        </Button>
        <Show when={state().error}>
          <div class="text-center text-red-500">{state().error?.message}</div>
        </Show>
      </div>
      <div>
        <ScrollView store={vm.ui.$scroll}>
          <ListView store={vm.request.resource.list}>
            <For each={state().response.dataSource}>
              {(v) => {
                return (
                  <div class="w-[320px]">
                    <div>{v.type_text}</div>
                    <Show when={v.media_type === FileAboutType.Image}>
                      <img class="w-[128px] h-[128px] object-cover" src={v.url} />
                    </Show>
                    <div class="flex items-center justify-between">
                      <div class="truncate">
                        <div class="truncate">{v.filename}</div>
                        <div>{v.created_at}</div>
                      </div>
                      <div class="flex items-center gap-2">
                        <div
                          onClick={() => {
                            vm.methods.handleCopyURL(v);
                          }}
                        >
                          复制链接
                        </div>
                        <div
                          onClick={() => {
                            vm.methods.handleDeleteMedia(v);
                          }}
                        >
                          删除
                        </div>
                      </div>
                    </div>
                  </div>
                );
              }}
            </For>
          </ListView>
        </ScrollView>
      </div>
    </div>
  );
}
