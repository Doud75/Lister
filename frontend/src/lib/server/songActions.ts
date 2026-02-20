import type { RequestEvent } from '@sveltejs/kit';
import type { SongPayload } from '$lib/types';

export async function extractSongData(event: RequestEvent, existingFormData?: FormData): Promise<SongPayload> {
    const data = existingFormData ?? await event.request.formData();

    return {
        title: data.get('title') as string,
        album_name: data.get('album_name')?.toString() || null,
        song_key: data.get('song_key')?.toString() || null,
        duration_seconds: Number(data.get('duration_seconds')) || null,
        tempo: Number(data.get('tempo')) || null,
        lyrics: data.get('lyrics')?.toString() || null,
        links: data.get('links')?.toString() || null
    };
}