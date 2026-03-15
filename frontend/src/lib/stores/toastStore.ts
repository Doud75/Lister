import { writable } from 'svelte/store';

type Toast = {
    id: number;
    message: string;
    type: 'error' | 'success' | 'info';
};

function createToastStore() {
    const { subscribe, update } = writable<Toast[]>([]);
    let nextId = 0;

    function add(message: string, type: Toast['type'], duration = 4000) {
        const id = nextId++;
        update((toasts) => [...toasts, { id, message, type }]);
        if (duration > 0) {
            setTimeout(() => remove(id), duration);
        }
    }

    function remove(id: number) {
        update((toasts) => toasts.filter((t) => t.id !== id));
    }

    return {
        subscribe,
        error: (msg: string) => add(msg, 'error'),
        success: (msg: string) => add(msg, 'success'),
        info: (msg: string) => add(msg, 'info'),
        remove,
    };
}

export const toastStore = createToastStore();