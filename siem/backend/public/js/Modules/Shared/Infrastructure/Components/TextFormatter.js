/**
 * Rellena los huecos {0}, {1}... de un texto de un *TextEnum (igual que config.Tr en el backend). Sin escapar: quien
 * lo pinte en innerHTML lo escapa con Htmler.
 */
export class TextFormatter {

    /**
     * @returns {TextFormatter}
     */
    static getInstance() {
        return new TextFormatter();
    }

    /**
     * @param {string} text con huecos {0}, {1}...
     * @param {...*} values en el orden de los huecos
     * @returns {string}
     */
    getFormatted(text, ...values) {
        return String(text ?? "").replace(/\{(\d+)\}/g, (hole, position) => String(values[Number(position)] ?? hole));
    }
}
