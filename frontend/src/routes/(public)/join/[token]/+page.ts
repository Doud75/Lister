import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
    const { token } = params;

    try {
        const res = await fetch(`/api/invitations/${token}`);

        if (res.status === 404) {
            return { valid: false, error: 'Cette invitation est introuvable.' };
        }
        if (res.status === 410) {
            return { valid: false, error: 'Cette invitation a expiré.' };
        }
        if (!res.ok) {
            return { valid: false, error: "Cette invitation est invalide ou a expiré." };
        }

        const data = await res.json();
        return {
            valid: true as const,
            band_name: data.band_name as string,
            role: data.role as string,
            token
        };
    } catch {
        return { valid: false, error: "Impossible de charger l'invitation." };
    }
};
