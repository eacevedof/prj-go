import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ApiException } from "/js/Modules/Shared/Domain/Exceptions/ApiException.js";

export class DashboardException extends ApiException {

    /**
     * @param {string[]} errors
     * @throws {DashboardException}
     */
    static badRequestCustom(errors) {
        throw new DashboardException(errors[0], HttpCodeEnum.BAD_REQUEST, errors);
    }

    /**
     * La API respondio con status false: su codigo y su detalle viajan tal cual.
     * @param {{code: number, errors: string[]}} apiResponse
     * @throws {DashboardException}
     */
    static fromApiResponse(apiResponse) {
        throw new DashboardException(apiResponse.errors[0] ?? "", apiResponse.code, apiResponse.errors);
    }
}
