export function formatDuration(seconds: number): string {
    if (!seconds || seconds === 0) {
        return '0m 00s';
    }
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}m ${remainingSeconds.toString().padStart(2, '0')}s`;
}