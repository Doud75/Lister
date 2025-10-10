import jsPDF from 'jspdf';
import { formatDuration } from '$lib/utils/utils';
import type { SetlistDetails } from '$lib/types';

export function generateSetlistPdf(setlist: SetlistDetails, totalDurationSeconds: number) {
    const doc = new jsPDF();

    const margin = 15;
    const lineHeight = 6;
    const pageHeight = doc.internal.pageSize.getHeight();
    let yPos = margin;

    const highlightColors = ['#FDFFB6', '#CAFFBF', '#9BF6FF', '#A0C4FF', '#BDB2FF', '#FFC6FF'];
    const speakerColors = new Map<string, string>();
    let songCounter = 1;

    const checkPageBreak = (spaceNeeded: number) => {
        if (yPos + spaceNeeded > pageHeight - margin) {
            doc.addPage();
            yPos = margin;
        }
    };

    doc.setFontSize(20);
    doc.setFont('helvetica', 'bold');
    doc.text(setlist.name, doc.internal.pageSize.getWidth() / 2, yPos, { align: 'center' });
    yPos += lineHeight * 1.5;

    doc.setFontSize(11);
    doc.setFont('helvetica', 'normal');
    doc.text(`Durée totale : ${formatDuration(totalDurationSeconds)}`, doc.internal.pageSize.getWidth() / 2, yPos, { align: 'center' });
    yPos += lineHeight * 2;

    setlist.items.forEach((item) => {
        checkPageBreak(20);

        if (item.item_type === 'song') {
            doc.setFontSize(16);
            doc.setFont('helvetica', 'bold');
            doc.text(`${songCounter}. ${item.title.String}`, margin, yPos);
            songCounter++;
            yPos += lineHeight * 1.2;

            if (item.notes?.Valid && item.notes.String) {
                doc.setFontSize(11);
                doc.setFont('helvetica', 'italic');
                const notesLines = doc.splitTextToSize(item.notes.String, doc.internal.pageSize.getWidth() - margin * 2 - 5);
                checkPageBreak(notesLines.length * lineHeight * 0.8);
                doc.text(notesLines, margin + 5, yPos);
                yPos += notesLines.length * lineHeight * 0.9;
            }

        } else if (item.item_type === 'interlude') {
            const speakerName = (item.speaker?.Valid && item.speaker.String) ? item.speaker.String : "Interlude";

            if (!speakerColors.has(speakerName)) {
                const color = highlightColors[speakerColors.size % highlightColors.length];
                speakerColors.set(speakerName, color);
            }
            const highlightColor = speakerColors.get(speakerName)!;

            doc.setFontSize(14);
            doc.setFont('helvetica', 'normal');

            const textWidth = doc.getTextWidth(speakerName);
            doc.setFillColor(highlightColor);
            doc.rect(margin, yPos - 4.5, textWidth + 4, lineHeight + 1, 'F');

            doc.text(speakerName, margin + 2, yPos);
            yPos += lineHeight * 1.2;

            if (item.notes?.Valid && item.notes.String) {
                doc.setFontSize(11);
                doc.setFont('helvetica', 'italic');
                const scriptLines = doc.splitTextToSize(item.notes.String, doc.internal.pageSize.getWidth() - margin * 2 - 5);
                checkPageBreak(scriptLines.length * lineHeight * 0.8);
                doc.text(scriptLines, margin + 5, yPos);
                yPos += scriptLines.length * lineHeight * 0.9;
            }
        }
        yPos += lineHeight * 1.8;
    });

    const sanitizedFileName = `${setlist.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}.pdf`;
    doc.save(sanitizedFileName);
}