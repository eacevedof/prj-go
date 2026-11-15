/**
 * Copia texto al portapapeles: la API moderna (pide contexto seguro: https o localhost) y, si no esta o falla, el
 * metodo antiguo con un textarea fuera de pantalla. No captura mas alla de ese reintento: devuelve si se copio.
 */
export class Clipboarder {

    /**
     * @returns {Clipboarder}
     */
    static getInstance() {
        return new Clipboarder();
    }

    /**
     * @param {string} text
     * @returns {Promise<boolean>}
     */
    async copy(text) {
        if (navigator.clipboard?.writeText) {
            try {
                await navigator.clipboard.writeText(text);
                return true;
            }
            catch (error) {
                // sin permiso o sin contexto seguro: se intenta el metodo antiguo
            }
        }
        return this.#copyWithTextarea(text);
    }

    #copyWithTextarea(text) {
        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.setAttribute("readonly", "");
        textarea.style.position = "fixed";
        textarea.style.left = "-9999px";
        document.body.append(textarea);
        textarea.select();
        let isCopied = false;
        try {
            isCopied = document.execCommand("copy");
        }
        catch (error) {
            isCopied = false;
        }
        textarea.remove();
        return isCopied;
    }
}
