import { writable } from 'svelte/store';
import type { Component } from 'svelte';

type ModalState = {
    component: Component<Record<string, unknown>> | null;
    props: Record<string, unknown>;
    isOpen: boolean;
};

function createModalStore() {
    const { subscribe, set } = writable<ModalState>({
        component: null,
        props: {},
        isOpen: false,
    });

    const close = () => set({ component: null, props: {}, isOpen: false });

    return {
        subscribe,
        open: <
            P extends Record<string, unknown> & { close?: () => void }
        >(
            component: Component<P>,
            props: Omit<P, 'close'>
        ) => {
            set({
                component: component as unknown as Component<Record<string, unknown>>,
                props: { ...props, close },
                isOpen: true,
            });
        },
        close,
    };
}

export const modalStore = createModalStore();
