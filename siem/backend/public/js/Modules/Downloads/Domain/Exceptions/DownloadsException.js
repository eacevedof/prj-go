import { ApiException } from "/js/Modules/Shared/Domain/Exceptions/ApiException.js";

export class DownloadsException extends ApiException {

    /**
     * La API respondio con status false: su codigo y su detalle viajan tal cual.
     * @param {{code: number, errors: string[]}} apiResponse
     * @throws {DownloadsException}
     */
    static fromApiResponse(apiResponse) {
        throw new DownloadsException(apiResponse.errors[0] ?? "", apiResponse.code, apiResponse.errors);
    }
}
