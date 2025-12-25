export type NullableString = { String: string; Valid: boolean };
export type NullableInt = { Int32: number; Valid: boolean };

export type Song = {
    id: number;
    title: string;
    album_name: NullableString;
    song_key: NullableString;
    duration_seconds: NullableInt;
    tempo: NullableInt;
    lyrics: NullableString;
    links: NullableString;
};

export type Interlude = {
    id: number;
    title: string;
    speaker: NullableString;
    script: NullableString;
    duration_seconds: NullableInt;
};


type SetlistItemBase = {
    id: number;
    position: number;
    title: NullableString;
    duration_seconds: NullableInt;
    notes: NullableString;
};

export type SetlistSongItem = SetlistItemBase & {
    item_type: 'song';
    song_id: NullableInt;
    tempo: NullableInt;
    song_key: NullableString;
    links: NullableString;
};

export type SetlistInterludeItem = SetlistItemBase & {
    item_type: 'interlude';
    interlude_id: NullableInt;
    speaker: NullableString;
    script: NullableString;
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
    color: string;
};