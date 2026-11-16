import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";
import { ProjectsLimitEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsLimitEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { CreateProjectResultDto } from "/js/Modules/Projects/Application/CreateProject/CreateProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Alta de un proyecto desde el panel. Aqui solo lo evidente (nombre obligatorio y su longitud): el resto de reglas
 * (nombre duplicado...) las aplica la API y vuelven como 400 / 409.
 */
export class CreateProjectService {

    #projectsWriterHttpRepository;

    constructor() {
        this.#projectsWriterHttpRepository = ProjectsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {CreateProjectService}
     */
    static getInstance() {
        return new CreateProjectService();
    }

    /**
     * @param {import("./CreateProjectDto.js").CreateProjectDto} createProjectDto
     * @returns {Promise<CreateProjectResultDto>}
     * @throws {ProjectsException}
     */
    async invoke(createProjectDto) {
        if (!createProjectDto.getProjectName()) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_IS_REQUIRED]);
        }
        if (createProjectDto.getProjectName().length > ProjectsLimitEnum.PROJECT_NAME_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_MAX_LENGTH_IS_100]);
        }
        if (createProjectDto.getProjectDescription().length > ProjectsLimitEnum.PROJECT_DESCRIPTION_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_DESCRIPTION_MAX_LENGTH_IS_500]);
        }

        const apiResponse = await this.#projectsWriterHttpRepository.createProject(
            {
                [ProjectsKeyEnum.PROJECT_NAME]: createProjectDto.getProjectName(),
                [ProjectsKeyEnum.PROJECT_DESCRIPTION]: createProjectDto.getProjectDescription(),
            },
            createProjectDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            ProjectsException.fromApiResponse(apiResponse);
        }

        const project = apiResponse.data ?? {};
        return CreateProjectResultDto.fromPrimitives({
            projectUuid: project[ProjectsKeyEnum.PROJECT_UUID],
            projectName: project[ProjectsKeyEnum.PROJECT_NAME],
        });
    }
}
