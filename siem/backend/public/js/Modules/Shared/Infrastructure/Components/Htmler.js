const HTML_ENTITIES = Object.freeze({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    "\"": "&quot;",
    "'": "&#39;",
});
const COPY_LABEL = "Copiar";

/**
 * Todo texto que viene de la API o del usuario se escapa antes de entrar en innerHTML (XSS). Tambien monta las piezas
 * de html que se repiten en todas las pantallas: pastilla de estado, boton Copiar y bloque de codigo copiable (los
 * clics de copiar los atiende CopyView).
 */
export class Htmler {

    /**
     * @returns {Htmler}
     */
    static getInstance() {
        return new Htmler();
    }

    /**
     * @param {*} text
     * @returns {string}
     */
    getEscaped(text) {
        return String(text ?? "").replace(/[&<>"']/g, (char) => HTML_ENTITIES[char]);
    }

    /**
     * Pastilla de estado (ok, online, warning, critical...): la clase badge-<valor> le da color y forma.
     * @param {string} value
     * @returns {string}
     */
    getBadge(value) {
        const escapedValue = this.getEscaped(value);
        return `<span class="badge badge-${escapedValue}">${escapedValue}</span>`;
    }

    /**
     * @param {string} text lo que se copia
     * @param {string} label
     * @returns {string}
     */
    getCopyButton(text, label = COPY_LABEL) {
        return `<button type="button" class="button button-small button-copy" data-copy="${this.getEscaped(text)}">${this.getEscaped(label)}</button>`;
    }

    /**
     * Un valor tecnico (uuid, sha256, token) en monoespaciada con su boton Copiar al lado.
     * @param {string} value
     * @returns {string}
     */
    getCopyableValue(value) {
        return `<span class="copyable"><code>${this.getEscaped(value)}</code>${this.getCopyButton(value)}</span>`;
    }

    /**
     * Bloque de comandos o de configuracion con su boton Copiar (copia el bloque entero, tal cual).
     * @param {string} code
     * @param {string} title que es (bash, PowerShell, .env...)
     * @returns {string}
     */
    getCodeBlock(code, title = "") {
        return `<div class="code-block">
            <div class="code-head"><span>${this.getEscaped(title)}</span>${this.getCopyButton(code)}</div>
            <pre><code>${this.getEscaped(code)}</code></pre>
        </div>`;
    }
}
