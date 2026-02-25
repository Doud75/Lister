export type Song = {
    id: number;
    title: string;
    album_name: string | null;
    song_key: string | null;
    duration_seconds: number | null;
    tempo: number | null;
    lyrics: string | null;
    links: string | null;
};

export type Interlude = {
    id: number;
    title: string;
    speaker: string | null;
    script: string | null;
    duration_seconds: number | null;
};


type SetlistItemBase = {
    id: number;
    position: number;
    title: string | null;
    duration_seconds: number | null;
    notes: string | null;
    transition_duration_seconds: number;
};

export type SetlistSongItem = SetlistItemBase & {
    item_type: 'song';
    song_id: number | null;
    tempo: number | null;
    song_key: string | null;
    links: string | null;
};

export type SetlistInterludeItem = SetlistItemBase & {
    item_type: 'interlude';
    interlude_id: number | null;
    speaker: string | null;
    script: string | null;
};

export type SetlistItem = SetlistSongItem | SetlistInterludeItem;

export type SetlistSummary = {
    id: number;
    name: string;
    color: string;
    is_archived: boolean;
    created_at: string;
};

export type SetlistDetails = SetlistSummary & {
    items: SetlistItem[];
};

export type SongPayload = {
    title: string | null;
    album_name: string | null;
    song_key: string | null;
    duration_seconds: number | null;
    tempo: number | null;
    lyrics: string | null;
    links: string | null;
};

export type BandMember = {
    id: number;
    username: string;
    role: string;
};

export type ApiError = {
    error: string;
    code?: string;
};