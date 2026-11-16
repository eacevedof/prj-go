import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";
import { UpdateProjectResultDto } from "/js/Modules/Projects/Application/UpdateProject/UpdateProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";

const PROJECT_NAME_MAX_LENGTH = 100;
const PROJECT_DESCRIPTION_MAX_LENGTH = 500;

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
        if (updateProjectDto.getProjectName().length > PROJECT_NAME_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_MAX_LENGTH_IS_100]);
        }
        if (updateProjectDto.getProjectDescription().length > PROJECT_DESCRIPTION_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_DESCRIPTION_MAX_LENGTH_IS_500]);
        }

        const apiResponse = await this.#projectsWriterHttpRepository.updateProjectByProjectUuid(
            updateProjectDto.getProjectUuid(),
            { project_name: updateProjectDto.getProjectName(), project_description: updateProjectDto.getProjectDescription() },
            updateProjectDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            ProjectsException.fromApiResponse(apiResponse);
        }

        const project = apiResponse.data ?? {};
        return UpdateProjectResultDto.fromPrimitives({
            projectUuid: project.project_uuid,
            projectName: project.project_name,
        });
    }
}
