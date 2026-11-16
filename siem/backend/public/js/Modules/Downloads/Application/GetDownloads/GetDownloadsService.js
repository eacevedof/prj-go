import { DownloadsKeyEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsKeyEnum.js";
import { GetDownloadsResultDto } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsResultDto.js";
import { DownloadsReaderHttpRepository } from "/js/Modules/Downloads/Infrastructure/Repositories/DownloadsReaderHttpRepository.js";
import { DownloadsException } from "/js/Modules/Downloads/Domain/Exceptions/DownloadsException.js";

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

        const downloadsCatalog = apiResponse.data ?? {};
        return GetDownloadsResultDto.fromPrimitives({
            endpointUrl: downloadsCatalog[DownloadsKeyEnum.ENDPOINT_URL],
            downloads: downloadsCatalog[DownloadsKeyEnum.DOWNLOADS],
            installers: downloadsCatalog[DownloadsKeyEnum.INSTALLERS],
        });
    }
}
