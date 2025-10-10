import { describe, it, expect, vi, beforeEach } from 'vitest';
import { generateSetlistPdf } from '$lib/utils/pdfGenerator';
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

describe('generateSetlistPdf', () => {
    const mockSetlist: SetlistDetails = {
        id: 1,
        name: 'My Awesome Setlist!',
        color: '#ff0000',
        created_at: '2024-01-01T12:00:00Z',
        items: [
            {
                id: 10,
                item_type: 'song',
                title: { String: 'First Song', Valid: true },
                notes: { String: 'A little note', Valid: true }
            },
            {
                id: 11,
                item_type: 'interlude',
                title: { String: 'A quick talk', Valid: true },
                speaker: { String: 'Lead', Valid: true },
                script: { String: 'This is the original template script.', Valid: true },
                notes: { String: 'Hello PDF world', Valid: true }
            },
            {
                id: 12,
                item_type: 'song',
                title: { String: 'Second Song', Valid: true },
                notes: { String: '', Valid: false }
            }
        ]
    } as SetlistDetails;

    it('should generate the PDF with the correct filename and content', () => {
        generateSetlistPdf(mockSetlist, 300);

        expect(mockSave).toHaveBeenCalledTimes(1);

        expect(mockSave).toHaveBeenCalledWith('my_awesome_setlist_.pdf');

        const capturedTexts = mockText.mock.calls.map(call => call[0]).join('\n');

        expect(capturedTexts).toContain('My Awesome Setlist!');
        expect(capturedTexts).toContain('Durée totale : 5m 00s');

        expect(capturedTexts).toContain('1. First Song');
        expect(capturedTexts).toContain('A little note');
        expect(capturedTexts).toContain('2. Second Song');

        expect(capturedTexts).toContain('Lead');
        expect(capturedTexts).toContain('Hello PDF world');
        expect(capturedTexts).not.toContain('This is the original template script.');
        expect(capturedTexts).not.toContain('A quick talk');
    });
});