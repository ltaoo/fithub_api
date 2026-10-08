import { View, Text, Fragment as TTFragment, For as TTFor, Show as TTShow, Portal as TTPortal, SVG, ref } from "@timeless/timeless";
import { render as render_dom } from "@timeless/timeless-dom";
export type { JSX, JSXElement, Accessor, Component } from "solid-js";

type Scope = { mounts: (() => void)[]; cleanups: (() => void)[]; contexts: Map<symbol, unknown>; disposed: boolean };
let scope: Scope | undefined;
let tracking: Set<any> | undefined;
export function onMount(callback: () => void) { scope?.mounts.push(callback); }
export function onCleanup(callback: () => void) { scope?.cleanups.push(callback); }
function with_scope<T>(owner: Scope | undefined, callback: () => T): T {
  const previous = scope; scope = owner;
  try { return callback(); } finally { scope = previous; }
}
export function createSignal<T = undefined>(initial?: T): [() => T, (value: T | ((value: T) => T)) => T] {
  const value = ref(initial as T);
  return [() => { tracking?.add(value); return value.value; }, next => {
    const resolved = typeof next === "function" ? (next as (value: T) => T)(value.value) : next;
    if (!Object.is(value.value, resolved)) value.as(resolved);
    return resolved;
  }];
}
export function binding<T>(read: () => T) {
  const owner = scope;
  const output = ref<T>(undefined as T);
  const subscriptions = new Map<any, () => void>();
  let running = false;
  const update = () => {
    if (running || owner?.disposed) return;
    running = true;
    const previous = tracking;
    const deps = new Set<any>(); tracking = deps;
    let next: T;
    try { next = with_scope(owner, read); } finally { tracking = previous; running = false; }
    for (const [dep, unsubscribe] of subscriptions) if (!deps.has(dep)) { unsubscribe(); subscriptions.delete(dep); }
    for (const dep of deps) if (!subscriptions.has(dep)) subscriptions.set(dep, dep.subscribe({ onChange: update }));
    if (!Object.is(output.value, next!)) output.as(next!);
  };
  update();
  onCleanup(() => { subscriptions.forEach(unsubscribe => unsubscribe()); subscriptions.clear(); output.destroy(); });
  return output;
}
export function effect(callback: () => unknown) { binding(callback); }
export const createEffect = effect;
export function children(read: () => any) { return read; }
export function lazy(read: () => any) {
  let initialized = false, value: any;
  const owner = scope;
  return { __lazy: true, read() { if (!initialized) { value = with_scope(owner, read); initialized = true; } return value; } };
}
export function dynamic(read: () => any) { return { __dynamic: true, read, owner: scope }; }
export function nodes(value: any): any[] {
  if (value == null || typeof value === "boolean") return [];
  if (Array.isArray(value)) return value.flatMap(nodes);
  if (value.__lazy) return nodes(value.read());
  if (value.__dynamic) {
    return with_scope(value.owner, () => {
      const resolved = binding<any>(value.read);
      if (typeof resolved.value === "string" || typeof resolved.value === "number") {
        return [Text(resolved)];
      }
      return [TTFor({
        each: binding(() => { tracking?.add(resolved); const next = resolved.value; return next == null || typeof next === "boolean" ? [] : [next]; }),
        render: (next: any) => list_item(nodes(next)),
      })];
    });
  }
  if (typeof value === "function") return nodes(value());
  return [typeof value === "string" || typeof value === "number" ? Text(String(value)) : value];
}
export function Fragment(props: any): any { return TTFragment({}, nodes(props.children)); }
function group(children: any[]) { return children.length === 1 ? children[0] : TTFragment({}, children); }
function list_item(children: any[]) {
  const node = group(children);
  // List deletion needs a single host boundary when the row itself is reactive.
  return ["fragment", "show", "for", "match"].includes(node.t)
    ? View({ as: "tt-item", style: { display: "contents" } }, [node]) : node;
}
export function h(component: any, props: any = {}): any {
  if (typeof component === "function") {
    const owner: Scope = { mounts: [], cleanups: [], contexts: new Map(scope?.contexts), disposed: false };
    return with_scope(owner, () => {
      const result = component(props);
      // Match descriptors are consumed by Switch rather than rendered.
      if (result?.__match) return result;
      const node = group(nodes(result));
      const mounted = node.onMounted.bind(node);
      const unmounted = node.onUnmounted.bind(node);
      const destroy = node.destroy?.bind(node);
      const cleanup = () => { owner.disposed = true; owner.cleanups.splice(0).reverse().forEach(callback => callback()); };
      node.onMounted = (event: any) => { mounted(event); owner.disposed = false; owner.mounts.splice(0).forEach(callback => with_scope(owner, callback)); };
      node.onUnmounted = () => { unmounted(); cleanup(); };
      if (destroy) node.destroy = () => { destroy(); cleanup(); };
      return node;
    });
  }
  const attributes: Record<string, any> = {};
  const events: Record<string, any> = {};
  const boolean_attributes = new Set(["disabled", "checked", "selected", "multiple", "required", "readonly", "readOnly", "autofocus", "autoFocus", "hidden", "controls", "muted", "loop"]);
  for (const key of Object.keys(props)) {
    if (["children", "class", "classList", "style", "ref"].includes(key)) continue;
    if (key.startsWith("on")) events[key] = (event: any) => props[key]?.(event);
    else attributes[key === "tabIndex" ? "tabindex" : key] = binding(() => boolean_attributes.has(key) ? (props[key] ? "" : undefined) : props[key]);
  }
  const classname = binding(() => [props.class || "", ...Object.entries(props.classList || {}).filter(([, active]) => active).map(([name]) => name)].join(" "));
  // Style keys remain reactive even when a model replaces the entire style object.
  const style = binding(() => {
    const raw = props.style;
    if (typeof raw !== "string") return raw || {};
    return Object.fromEntries(raw.split(";").map(part => {
      const separator = part.indexOf(":");
      return separator < 0 ? [] : [part.slice(0, separator).trim(), part.slice(separator + 1).trim()];
    }).filter(entry => entry.length === 2));
  });
  const svg_components: Record<string, any> = { svg: SVG.SVG, path: SVG.Path, circle: SVG.Circle, rect: SVG.Rect, line: SVG.Line, polyline: SVG.Polyline, polygon: SVG.Polygon, ellipse: SVG.Ellipse, g: SVG.G, defs: SVG.Defs };
  const factory = svg_components[component] || View;
  const element = factory({ as: component, attributes, class: classname, style }, nodes(props.children));
  const properties: Record<string, any> = {};
  for (const key of ["value", "checked", "selected", "muted"]) if (key in props) properties[key] = binding(() => props[key]);
  const descriptor = Object.getOwnPropertyDescriptor(element, "$elm")!;
  Object.defineProperty(element, "$elm", {
    configurable: true, get: descriptor.get,
    set(host) {
      descriptor.set!.call(element, host);
      if (!host) return;
      const original = host.render.bind(host);
      host.render = () => {
        const dom = original();
        for (const [key, callback] of Object.entries(events)) dom.addEventListener(key.slice(2).toLowerCase(), callback);
        for (const [key, value] of Object.entries(properties)) {
          const assign = (next: any) => { if (dom[key] !== next) dom[key] = next ?? ""; };
          assign(value.value); value.subscribe({ onChange: assign });
        }
        if (props.ref) props.ref(dom);
        return dom;
      };
    },
  });
  return element;
}
export function Show(props: any): any {
  return TTShow({ when: binding(() => Boolean(props.when)), ok: () => nodes(props.children), else: () => nodes(props.fallback) });
}
export function For<T extends readonly any[]>(props: { each: T | undefined | null | false; children: (item: T[number], index: () => number) => any; fallback?: any }): any {
  const owner = scope;
  return TTFor({ each: binding(() => [...(props.each || [])]), render: (item: any, index: any) =>
    with_scope(owner, () => list_item(nodes(props.children(item, () => { tracking?.add(index); return index.value; })))) });
}
export function Match(props: any): any { return { __match: true, props }; }
export function Switch(props: any): any {
  const branches = (props.children || []).flatMap((child: any) => child.__lazy ? [child.read()] : [child]).filter((child: any) => child?.__match);
  const selected = binding(() => branches.findIndex((child: any) => child.props.when));
  return TTFor({ each: binding(() => { tracking?.add(selected); return [selected.value]; }), render: (index: number) => list_item(nodes(index < 0 ? props.fallback : branches[index].props.children)) });
}
export function Portal(props: any): any { return TTPortal({}, nodes(props.children)); }
export function createContext<T>(default_value?: T) {
  const key = Symbol();
  return { key, default_value, Provider(props: any): any { scope?.contexts.set(key, props.value); return TTFragment({}, nodes(props.children)); } };
}
export function useContext<T>(context: { key: symbol; default_value?: T }): T | undefined { return scope?.contexts.get(context.key) as T ?? context.default_value; }
export function render(callback: () => any, root: HTMLElement) {
  const node = TTFragment({}, [h(callback, {})]); render_dom(node, root);
  return () => { node.onUnmounted(); root.replaceChildren(); };
}
