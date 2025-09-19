import type {SetlistItem} from "$lib/types";

export function formatDuration(seconds: number): string {
    if (!seconds || seconds === 0) {
        return '0m 00s';
    }
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}m ${remainingSeconds.toString().padStart(2, '0')}s`;
}

export function formatItemDuration(seconds: number | null | undefined): string {
    if (seconds === null || seconds === undefined) return '-';
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

export function calculateTotalDuration(items: SetlistItem[]): number {
    if (!items) return 0;
    return items.reduce((total, item) => total + (item.duration_seconds?.Int32 ?? 0), 0);
}