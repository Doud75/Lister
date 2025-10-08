<script lang="ts">
    import type { HTMLInputAttributes } from 'svelte/elements';

    type CustomProps = {
        id: string;
        label: string;
        name: string;
        type?: 'text' | 'password' | 'email' | 'number';
        togglePasswordVisibility?: boolean;
        value?: string;
    };

    let {
        id,
        label,
        name,
        type = 'text',
        placeholder = '',
        required = false,
        togglePasswordVisibility = false,
        value = $bindable(),
        class: className,
        ...rest
    } = $props<CustomProps & HTMLInputAttributes>();

    let isPasswordVisible = $state(false);

    function toggleVisibility() {
        isPasswordVisible = !isPasswordVisible;
    }
</script>

<div>
    <label for={id} class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">
        {label}
    </label>
    <div class="relative mt-2">
        <input
                {...rest}
                {id}
                {name}
                type={isPasswordVisible ? 'text' : type}
                {required}
                {placeholder}
                bind:value
                class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500 {className || ''}"
        />
        {#if type === 'password' && togglePasswordVisibility}
            <button
                    type="button"
                    onclick={toggleVisibility}
                    class="absolute inset-y-0 right-0 flex items-center pr-3"
                    aria-label={isPasswordVisible ? 'Cacher le mot de passe' : 'Afficher le mot de passe'}
            >
                {#if isPasswordVisible}
                    <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="h-5 w-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
                    >
                        <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.243 4.243L6.228 6.228"
                        />
                    </svg>
                {:else}
                    <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="h-5 w-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
                    >
                        <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.432 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
                        />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                    </svg>
                {/if}
            </button>
        {/if}
    </div>
</div>