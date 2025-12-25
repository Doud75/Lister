import { describe, it, expect } from 'vitest';
import { formatDuration, formatItemDuration, calculateTotalDuration, getSongNumber } from '$lib/utils/utils';
import type { SetlistItem } from '$lib/types';

describe('formatDuration', () => {
    it('should format 0 seconds correctly', () => {
        expect(formatDuration(0)).toBe('0m 00s');
    });

    it('should format less than a minute correctly', () => {
        expect(formatDuration(45)).toBe('0m 45s');
    });

    it('should format exactly one minute correctly', () => {
        expect(formatDuration(60)).toBe('1m 00s');
    });

    it('should format minutes and seconds correctly', () => {
        expect(formatDuration(135)).toBe('2m 15s');
    });

    it('should handle null or undefined input by returning 0m 00s', () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        expect(formatDuration(null)).toBe('0m 00s');
    });
});

describe('formatItemDuration', () => {
    it('should format duration with minutes and seconds', () => {
        expect(formatItemDuration(155)).toBe('2:35');
    });

    it('should pad seconds with a leading zero', () => {
        expect(formatItemDuration(125)).toBe('2:05');
    });

    it('should return a dash for null or undefined input', () => {
        expect(formatItemDuration(null)).toBe('-');
        expect(formatItemDuration(undefined)).toBe('-');
    });
});

describe('calculateTotalDuration', () => {
    it('should return 0 for an empty array', () => {
        expect(calculateTotalDuration([])).toBe(0);
    });

    it('should calculate the total duration of all items', () => {
        const items: SetlistItem[] = [
            { id: 1, duration_seconds: { Int32: 180, Valid: true } },
            { id: 2, duration_seconds: { Int32: 240, Valid: true } },
            { id: 3, duration_seconds: { Int32: 30, Valid: true } },
        ] as SetlistItem[];

        expect(calculateTotalDuration(items)).toBe(450);
    });

    it('should handle items with no duration', () => {
        const items: SetlistItem[] = [
            { id: 1, duration_seconds: { Int32: 180, Valid: true } },
            { id: 2, duration_seconds: { Int32: 0, Valid: false } },
            { id: 3, duration_seconds: { Int32: 30, Valid: true } },
        ] as SetlistItem[];

        expect(calculateTotalDuration(items)).toBe(210);
    });
});

describe('getSongNumber', () => {
    const song1 = { id: 1, item_type: 'song' } as SetlistItem;
    const interlude1 = { id: 2, item_type: 'interlude' } as SetlistItem;
    const song2 = { id: 3, item_type: 'song' } as SetlistItem;
    const song3 = { id: 4, item_type: 'song' } as SetlistItem;
    const interlude2 = { id: 5, item_type: 'interlude' } as SetlistItem;

    const mixedList = [song1, interlude1, song2, song3, interlude2];

    it('should return 1 for the first song', () => {
        expect(getSongNumber(song1, mixedList)).toBe(1);
    });

    it('should return null for an interlude', () => {
        expect(getSongNumber(interlude1, mixedList)).toBeNull();
        expect(getSongNumber(interlude2, mixedList)).toBeNull();
    });

    it('should return correct sequential number for subsequent songs, skipping interludes', () => {
        expect(getSongNumber(song2, mixedList)).toBe(2);
        expect(getSongNumber(song3, mixedList)).toBe(3);
    });

    it('should return null if the song is not found in the list', () => {
        const unknownSong = { id: 99, item_type: 'song' } as SetlistItem;
        expect(getSongNumber(unknownSong, mixedList)).toBeNull();
    });

    it('should work correctly with an empty list', () => {
        expect(getSongNumber(song1, [])).toBeNull();
    });
});