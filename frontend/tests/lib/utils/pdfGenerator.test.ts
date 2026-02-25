import { describe, it, expect, vi, beforeEach } from 'vitest';
import { generateSetlistPdf, generateLivePdf } from '$lib/utils/pdfGenerator';
import type { SetlistDetails } from '$lib/types';

const mockSave = vi.fn();
const mockText = vi.fn();
const mockGetTextWidth = vi.fn((text: string) => text.length * 5);

vi.mock('jspdf', () => ({
    default: vi.fn().mockImplementation(() => ({
        text: mockText,
        setFontSize: vi.fn(),
        setFont: vi.fn(),
        rect: vi.fn(),
        addPage: vi.fn(),
        save: mockSave,
        splitTextToSize: (text: string) => [text],
        getTextWidth: mockGetTextWidth,
        setFillColor: vi.fn(),
        setTextColor: vi.fn(),
        internal: {
            pageSize: {
                getWidth: () => 210,
                getHeight: () => 297,
            },
        },
    })),
}));

beforeEach(() => {
    vi.clearAllMocks();
});

describe('PDF Generation', () => {
    const mockSetlist: SetlistDetails = {
        id: 1,
        name: 'My Awesome Setlist!',
        color: '#ff0000',
        is_archived: false,
        created_at: '2024-01-01T12:00:00Z',
        items: [
            {
                id: 10,
                item_type: 'song',
                position: 0,
                title: 'First Song',
                notes: 'A little note',
                duration_seconds: null,
                tempo: null,
                song_key: null,
                links: null,
                song_id: null,
                transition_duration_seconds: 0,
            },
            {
                id: 11,
                item_type: 'interlude',
                position: 1,
                title: 'A quick talk',
                speaker: 'Lead',
                notes: 'Hello PDF world',
                duration_seconds: null,
                script: null,
                interlude_id: null,
                transition_duration_seconds: 0,
            },
            {
                id: 12,
                item_type: 'song',
                position: 2,
                title: 'Second Song',
                notes: null,
                duration_seconds: null,
                tempo: null,
                song_key: null,
                links: null,
                song_id: null,
                transition_duration_seconds: 0,
            }
        ]
    } as SetlistDetails;

    describe('Standard PDF (generateSetlistPdf)', () => {
        it('should generate the PDF with notes and speaker-only format for interludes', () => {
            generateSetlistPdf(mockSetlist, 300);

            expect(mockSave).toHaveBeenCalledWith('my_awesome_setlist_.pdf');

            const capturedTexts = mockText.mock.calls.map(call => call[0]).join('\n');

            expect(capturedTexts).toContain('My Awesome Setlist!');
            expect(capturedTexts).toContain('Durée totale : 5m 00s');
            expect(capturedTexts).toContain('1. First Song');
            expect(capturedTexts).toContain('2. Second Song');

            expect(capturedTexts).toContain('A little note');
            expect(capturedTexts).toContain('Hello PDF world');

            expect(capturedTexts).toContain('Lead');
            expect(capturedTexts).not.toContain('A quick talk');
            expect(capturedTexts).not.toContain('Lead:');
        });
    });

    describe('Live PDF (generateLivePdf)', () => {
        it('should generate a PDF with no notes and "speaker: title" format for interludes', () => {
            generateLivePdf(mockSetlist, 300);

            expect(mockSave).toHaveBeenCalledWith('my_awesome_setlist__live.pdf');

            const capturedTexts = mockText.mock.calls.map(call => call[0]).join('\n');

            expect(capturedTexts).toContain('My Awesome Setlist!');
            expect(capturedTexts).toContain('Durée totale : 5m 00s');
            expect(capturedTexts).toContain('1. First Song');
            expect(capturedTexts).toContain('2. Second Song');

            expect(capturedTexts).not.toContain('A little note');
            expect(capturedTexts).not.toContain('Hello PDF world');

            expect(capturedTexts).toContain('Lead: A quick talk');
        });
    });
});