import jsPDF from 'jspdf';
import { formatDuration } from '$lib/utils/utils';
import type { SetlistDetails } from '$lib/types';

interface PdfOptions {
    includeNotes: boolean;
    fontSizes: {
        mainTitle: number;
        duration: number;
        itemTitle: number;
        itemNotes: number;
    };
    fileNameSuffix: string;
    lineHeightMultiplier: number;
    interludeFormat: 'speakerOnly' | 'speakerAndTitle';
}

function generatePdf(setlist: SetlistDetails, totalDurationSeconds: number, options: PdfOptions) {
    const doc = new jsPDF();

    const margin = 15;
    const lineHeight = 6;
    const pageHeight = doc.internal.pageSize.getHeight();
    let yPos = margin;

    const highlightColors = ['#6EE7B7', '#FBBF24', '#F87171', '#60A5FA', '#A78BFA', '#F472B6'];
    const speakerColors = new Map<string, string>();
    let songCounter = 1;

    const checkPageBreak = (spaceNeeded: number) => {
        if (yPos + spaceNeeded > pageHeight - margin) {
            doc.addPage();
            yPos = margin;
        }
    };

    doc.setFontSize(options.fontSizes.mainTitle);
    doc.setFont('helvetica', 'bold');
    doc.text(setlist.name, doc.internal.pageSize.getWidth() / 2, yPos, { align: 'center' });
    yPos += lineHeight * 1.5;

    doc.setFontSize(options.fontSizes.duration);
    doc.setFont('helvetica', 'normal');
    doc.text(`Durée totale : ${formatDuration(totalDurationSeconds)}`, doc.internal.pageSize.getWidth() / 2, yPos, { align: 'center' });
    yPos += lineHeight * 2;

    setlist.items.forEach((item) => {
        checkPageBreak(20 * options.lineHeightMultiplier);

        if (item.item_type === 'song') {
            doc.setFontSize(options.fontSizes.itemTitle);
            doc.setFont('helvetica', 'bold');
            doc.text(`${songCounter}. ${item.title.String}`, margin, yPos);
            songCounter++;
            yPos += lineHeight * options.lineHeightMultiplier;

            if (options.includeNotes && item.notes?.Valid && item.notes.String) {
                doc.setFontSize(options.fontSizes.itemNotes);
                doc.setFont('helvetica', 'italic');
                const notesLines = doc.splitTextToSize(item.notes.String, doc.internal.pageSize.getWidth() - margin * 2 - 5);
                checkPageBreak(notesLines.length * lineHeight * 0.8);
                doc.text(notesLines, margin + 5, yPos);
                yPos += notesLines.length * lineHeight * 0.9;
            }
        } else if (item.item_type === 'interlude') {
            const speakerName = (item.speaker?.Valid && item.speaker.String) ? item.speaker.String : null;
            const title = item.title.String || 'Interlude';
            let interludeText: string;

            if (options.interludeFormat === 'speakerAndTitle') {
                interludeText = speakerName ? `${speakerName}: ${title}` : title;
            } else {
                interludeText = speakerName || title;
            }

            const colorKey = speakerName || title;
            if (!speakerColors.has(colorKey)) {
                const color = highlightColors[speakerColors.size % highlightColors.length];
                speakerColors.set(colorKey, color);
            }
            const highlightColor = speakerColors.get(colorKey)!;

            doc.setFontSize(options.fontSizes.itemTitle);
            doc.setFont('helvetica', 'bold');
            doc.setTextColor(0, 0, 0);

            const textWidth = doc.getTextWidth(interludeText);
            doc.setFillColor(highlightColor);
            doc.rect(margin, yPos - (lineHeight * options.lineHeightMultiplier * 0.7), textWidth + 4, lineHeight * options.lineHeightMultiplier, 'F');
            doc.text(interludeText, margin + 2, yPos);
            yPos += lineHeight * options.lineHeightMultiplier;

            if (options.includeNotes && item.notes?.Valid && item.notes.String) {
                doc.setFontSize(options.fontSizes.itemNotes);
                doc.setFont('helvetica', 'italic');
                const scriptLines = doc.splitTextToSize(item.notes.String, doc.internal.pageSize.getWidth() - margin * 2 - 5);
                checkPageBreak(scriptLines.length * lineHeight * 0.8);
                doc.text(scriptLines, margin + 5, yPos);
                yPos += scriptLines.length * lineHeight * 0.9;
            }
        }
        yPos += lineHeight * 1.8 * options.lineHeightMultiplier;
    });

    const sanitizedFileName = `${setlist.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}${options.fileNameSuffix}.pdf`;
    doc.save(sanitizedFileName);
}

export function generateSetlistPdf(setlist: SetlistDetails, totalDurationSeconds: number) {
    const options: PdfOptions = {
        includeNotes: true,
        fontSizes: { mainTitle: 20, duration: 11, itemTitle: 14, itemNotes: 11 },
        fileNameSuffix: '',
        lineHeightMultiplier: 1.2,
        interludeFormat: 'speakerOnly'
    };
    generatePdf(setlist, totalDurationSeconds, options);
}

export function generateLivePdf(setlist: SetlistDetails, totalDurationSeconds: number) {
    const options: PdfOptions = {
        includeNotes: false,
        fontSizes: { mainTitle: 24, duration: 12, itemTitle: 20, itemNotes: 0 },
        fileNameSuffix: '_live',
        lineHeightMultiplier: 0.75,
        interludeFormat: 'speakerAndTitle'
    };
    generatePdf(setlist, totalDurationSeconds, options);
}