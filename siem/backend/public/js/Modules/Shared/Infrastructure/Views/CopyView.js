import { Clipboarder } from "/js/Modules/Shared/Infrastructure/Components/Clipboarder.js";

const COPIED_LABEL = "Copiado";
const FAILED_LABEL = "No se pudo copiar";
const FEEDBACK_MS = 1500;
const CSS_COPIED = "is-copied";

/**
 * Botones "Copiar" de toda la pagina: cualquier elemento con data-copy="texto" (los pinta Htmler.getCopyButton y
 * Htmler.getCodeBlock: tokens, uuids, sha256, comandos de instalacion). Un solo listener delegado en document: vale
 * tambien para lo que se pinta despues. Confirma en el propio boton durante un momento.
 */
export class CopyView {

    static #instance = null;

    #clipboarder;

    constructor() {
        this.#clipboarder = Clipboarder.getInstance();
    }

    /**
     * @returns {CopyView}
     */
    static getInstance() {
        CopyView.#instance ??= new CopyView();
        return CopyView.#instance;
    }

    init() {
        document.addEventListener("click", (event) => {
            const buttonCopy = event.target.closest("[data-copy]");
            if (!buttonCopy) return;
            event.preventDefault();
            this.#copy(buttonCopy);
        });
    }

    async #copy(buttonCopy) {
        const isCopied = await this.#clipboarder.copy(buttonCopy.dataset.copy ?? "");
        // el texto original del boton se guarda la primera vez: dos clics seguidos no deben dejar "Copiado" fijo
        buttonCopy.dataset.copyLabel ??= buttonCopy.textContent;
        buttonCopy.textContent = isCopied ? COPIED_LABEL : FAILED_LABEL;
        buttonCopy.classList.toggle(CSS_COPIED, isCopied);
        clearTimeout(Number(buttonCopy.dataset.copyTimeoutId));
        buttonCopy.dataset.copyTimeoutId = String(setTimeout(() => {
            buttonCopy.textContent = buttonCopy.dataset.copyLabel;
            buttonCopy.classList.remove(CSS_COPIED);
        }, FEEDBACK_MS));
    }
}
