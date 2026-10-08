import { Card as ShadcnCard, CardHeader as ShadcnCardHeader, CardContent as ShadcnCardContent, CardTitle as ShadcnCardTitle, Badge as ShadcnBadge } from "@timeless/shadcn";
import { JSX, nodes } from "@/timeless";
import { view_props } from "./timeless";
export function Card(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnCard(view_props(props), nodes(props.children)); }
export function CardHeader(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnCardHeader(view_props(props), nodes(props.children)); }
export function CardContent(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnCardContent(view_props(props), nodes(props.children)); }
export function CardTitle(props: JSX.HTMLAttributes<HTMLDivElement>): any { return ShadcnCardTitle(view_props(props), nodes(props.children)); }
export function Badge(props: JSX.HTMLAttributes<HTMLSpanElement>): any { return ShadcnBadge(view_props(props), nodes(props.children)); }
