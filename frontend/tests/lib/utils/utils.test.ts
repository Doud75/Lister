import { describe, it, expect } from 'vitest';
import { formatDuration, formatItemDuration, calculateTotalDuration } from '$lib/utils/utils';
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