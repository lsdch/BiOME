import { Comment, computed, getCurrentInstance, Text, type Slot, type VNode } from 'vue';

export function hasSlotContent(slot: Slot | undefined | null, props: any = {}) {
  return !isSlotEmpty(slot, props);
}

export function isSlotEmpty(slot: Slot | undefined | null, props: any = {}) {
  return isVNodeEmpty(slot?.(props));
}

export function isVNodeEmpty(vnode: VNode | VNode[] | undefined | null) {
  return (
    !vnode ||
    asArray(vnode).every(
      (vnode) =>
        vnode.type === Comment ||
        (vnode.type === Text && !vnode.children?.length),
    )
  );
}

export function asArray<T>(arg: T | T[] | null) {
  return Array.isArray(arg) ? arg : arg !== null ? [arg] : [];
}

/**
 * Check if the current component has a listener with the given name.
 * @param listenerName  The name of the listener to check. E.g. 'onClick', 'onUpdate:modelValue', etc.
 */
export function hasEventListener(listenerName: string) {
  return computed(() => !!getCurrentInstance()?.vnode.props?.[listenerName]);
}