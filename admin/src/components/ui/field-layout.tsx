import { FieldGroup as ShadcnFieldGroup, FieldSet as ShadcnFieldSet, FieldLegend as ShadcnFieldLegend, FieldDescription as ShadcnFieldDescription } from "@timeless/shadcn";
import { JSX, nodes } from "@/timeless";
import { view_props } from "./timeless";
export function FieldGroup(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnFieldGroup(view_props(props), nodes(props.children)); }
export function FieldSet(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnFieldSet(view_props(props), nodes(props.children)); }
export function FieldLegend(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnFieldLegend(view_props(props), nodes(props.children)); }
export function FieldDescription(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnFieldDescription(view_props(props), nodes(props.children)); }
