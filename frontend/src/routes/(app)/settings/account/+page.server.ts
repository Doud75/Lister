import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch }) => {
        const data = await request.formData();
        const currentPassword = data.get('current_password');
        const newPassword = data.get('new_password');
        const confirmPassword = data.get('confirm_password');

        if (!currentPassword || !newPassword || !confirmPassword) {
            return fail(400, { error: 'Tous les champs sont requis.' });
        }

        if (newPassword !== confirmPassword) {
            return fail(400, { error: 'Les nouveaux mots de passe ne correspondent pas.' });
        }

        const response = await fetch('/api/user/password', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                current_password: currentPassword,
                new_password: newPassword
            })
        });

        if (!response.ok) {
            if (response.status === 401) {
                return fail(401, { error: 'Votre ancien mot de passe est incorrect.' });
            }
            const result = await response.json();
            return fail(response.status, { error: result.error || 'Une erreur est survenue.' });
        }

        throw redirect(303, '/');
    }
};