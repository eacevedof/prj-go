const BYTE_UNITS = Object.freeze(["B", "KB", "MB", "GB"]);
const BYTES_PER_UNIT = 1024;

/**
 * Numeros para pintar: separador de miles del navegador y tamaños en B/KB/MB/GB.
 */
export class NumberFormatter {

    /**
     * @returns {NumberFormatter}
     */
    static getInstance() {
        return new NumberFormatter();
    }

    /**
     * @param {*} value
     * @returns {string}
     */
    getNumber(value) {
        return (Number(value) || 0).toLocaleString();
    }

    /**
     * @param {*} bytes
     * @returns {string} "512 B", "12.4 MB"
     */
    getBytes(bytes) {
        let value = Number(bytes) || 0;
        let unitIndex = 0;
        while (value >= BYTES_PER_UNIT && unitIndex < BYTE_UNITS.length - 1) {
            value /= BYTES_PER_UNIT;
            unitIndex++;
        }
        return `${unitIndex === 0 ? value : value.toFixed(1)} ${BYTE_UNITS[unitIndex]}`;
    }
}
