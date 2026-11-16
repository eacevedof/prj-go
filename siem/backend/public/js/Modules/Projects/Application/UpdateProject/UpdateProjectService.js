import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";
import { ProjectsLimitEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsLimitEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { UpdateProjectResultDto } from "/js/Modules/Projects/Application/UpdateProject/UpdateProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Cambia nombre y descripcion de un proyecto. Aqui solo lo evidente: el resto (nombre duplicado, caracteres
 * permitidos...) lo aplica la API y vuelve como 400 / 404 / 409.
 */
export class UpdateProjectService {

    #projectsWriterHttpRepository;

    constructor() {
        this.#projectsWriterHttpRepository = ProjectsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {UpdateProjectService}
     */
    static getInstance() {
        return new UpdateProjectService();
    }

    /**
     * @param {import("./UpdateProjectDto.js").UpdateProjectDto} updateProjectDto
     * @returns {Promise<UpdateProjectResultDto>}
     * @throws {ProjectsException}
     */
    async invoke(updateProjectDto) {
        if (!updateProjectDto.getProjectUuid()) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_UUID_IS_REQUIRED]);
        }
        if (!updateProjectDto.getProjectName()) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_IS_REQUIRED]);
        }
        if (updateProjectDto.getProjectName().length > ProjectsLimitEnum.PROJECT_NAME_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_MAX_LENGTH_IS_100]);
        }
        if (updateProjectDto.getProjectDescription().length > ProjectsLimitEnum.PROJECT_DESCRIPTION_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_DESCRIPTION_MAX_LENGTH_IS_500]);
        }

        const apiResponse = await this.#projectsWriterHttpRepository.updateProjectByProjectUuid(
            updateProjectDto.getProjectUuid(),
            {
                [ProjectsKeyEnum.PROJECT_NAME]: updateProjectDto.getProjectName(),
                [ProjectsKeyEnum.PROJECT_DESCRIPTION]: updateProjectDto.getProjectDescription(),
            },
            updateProjectDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            ProjectsException.fromApiResponse(apiResponse);
        }

        const project = apiResponse.data ?? {};
        return UpdateProjectResultDto.fromPrimitives({
            projectUuid: project[ProjectsKeyEnum.PROJECT_UUID],
            projectName: project[ProjectsKeyEnum.PROJECT_NAME],
        });
    }
}
