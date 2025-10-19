import { writable } from 'svelte/store';
import type { Component } from 'svelte';

type ModalState = {
    component: Component<any> | null;
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
                component,
                props: { ...props, close },
                isOpen: true,
            });
        },
        close,
    };
}

export const modalStore = createModalStore();
