import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => {
    return {
        redirectTo: url.searchParams.get('redirectTo') || undefined
    };
};
