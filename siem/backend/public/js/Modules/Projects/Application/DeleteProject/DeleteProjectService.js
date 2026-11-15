import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";
import { DeleteProjectResultDto } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";

/**
 * Baja de un proyecto: la API borra tambien sus agentes, monitores y eventos. No se puede deshacer.
 */
export class DeleteProjectService {

    #projectsWriterHttpRepository;

    constructor() {
        this.#projectsWriterHttpRepository = ProjectsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {DeleteProjectService}
     */
    static getInstance() {
        return new DeleteProjectService();
    }

    /**
     * @param {import("./DeleteProjectDto.js").DeleteProjectDto} deleteProjectDto
     * @returns {Promise<DeleteProjectResultDto>}
     * @throws {ProjectsException}
     */
    async invoke(deleteProjectDto) {
        if (!deleteProjectDto.getProjectUuid()) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_UUID_IS_REQUIRED]);
        }

        const apiResponse = await this.#projectsWriterHttpRepository.deleteProjectByProjectUuid(
            deleteProjectDto.getProjectUuid(), deleteProjectDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            ProjectsException.fromApiResponse(apiResponse);
        }
        return DeleteProjectResultDto.fromPrimitives({ projectUuid: deleteProjectDto.getProjectUuid() });
    }
}
