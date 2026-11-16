import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { DeleteProjectResultDto } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Baja de un proyecto: la API borra tambien sus agentes, monitores y eventos (no se puede deshacer) y devuelve los
 * agentes que tenia, para avisar de donde hay que desinstalar sus siem-client.
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
        const deletedProject = apiResponse.data ?? {};
        return DeleteProjectResultDto.fromPrimitives({
            projectUuid: deletedProject[ProjectsKeyEnum.PROJECT_UUID] ?? deleteProjectDto.getProjectUuid(),
            projectName: deletedProject[ProjectsKeyEnum.PROJECT_NAME] ?? "",
            agents: deletedProject[ProjectsKeyEnum.AGENTS] ?? [],
        });
    }
}
