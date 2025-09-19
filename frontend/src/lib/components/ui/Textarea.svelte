<script lang="ts">
    let {
        id,
        label,
        name,
        placeholder = '',
        required = false,
        rows = 4,
        value = $bindable()
    } = $props<{
        id: string;
        label: string;
        name: string;
        placeholder?: string;
        required?: boolean;
        rows?: number;
        value?: string;
    }>();
</script>

<div>
    <label for={id} class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200"
    >{label}</label
    >
    <div class="mt-2">
		<textarea
                {id}
                {name}
                {required}
                {placeholder}
                {rows}
                bind:value
                class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
        ></textarea>
    </div>
</div>