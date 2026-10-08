import { JSX } from "@/timeless";
import { Skeleton as ShadcnSkeleton } from "@timeless/shadcn";

import { cn } from "@/utils";

function Skeleton(props: {} & JSX.HTMLAttributes<HTMLDivElement>): any {
  return ShadcnSkeleton({ class: cn("w-full h-full", props.class) });
}

export { Skeleton };
