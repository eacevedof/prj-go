import { DownloadsException } from "/js/Modules/Downloads/Domain/Exceptions/DownloadsException.js";
import { GetDownloadsResultDto } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsResultDto.js";
import { DownloadsReaderHttpRepository } from "/js/Modules/Downloads/Infrastructure/Repositories/DownloadsReaderHttpRepository.js";

/**
 * Los binarios del agente disponibles para descargar (los deja `make dist` del client en el backend).
 */
export class GetDownloadsService {

    #downloadsReaderHttpRepository;

    constructor() {
        this.#downloadsReaderHttpRepository = DownloadsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetDownloadsService}
     */
    static getInstance() {
        return new GetDownloadsService();
    }

    /**
     * @param {import("./GetDownloadsDto.js").GetDownloadsDto} getDownloadsDto
     * @returns {Promise<GetDownloadsResultDto>}
     * @throws {DownloadsException}
     */
    async invoke(getDownloadsDto) {
        const apiResponse = await this.#downloadsReaderHttpRepository.getDownloads(getDownloadsDto.getDashboardToken());
        if (!apiResponse.status) {
            DownloadsException.fromApiResponse(apiResponse);
        }

        const data = apiResponse.data ?? {};
        return GetDownloadsResultDto.fromPrimitives({
            endpointUrl: data.endpoint_url,
            downloads: data.downloads,
            installers: data.installers,
        });
    }
}
