import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";
import { CreateProjectResultDto } from "/js/Modules/Projects/Application/CreateProject/CreateProjectResultDto.js";
import { ProjectsWriterHttpRepository } from "/js/Modules/Projects/Infrastructure/Repositories/ProjectsWriterHttpRepository.js";

const PROJECT_NAME_MAX_LENGTH = 100;
const PROJECT_DESCRIPTION_MAX_LENGTH = 500;

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
        if (createProjectDto.getProjectName().length > PROJECT_NAME_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_NAME_MAX_LENGTH_IS_100]);
        }
        if (createProjectDto.getProjectDescription().length > PROJECT_DESCRIPTION_MAX_LENGTH) {
            ProjectsException.badRequestCustom([ProjectsMessageEnum.PROJECT_DESCRIPTION_MAX_LENGTH_IS_500]);
        }

        const apiResponse = await this.#projectsWriterHttpRepository.createProject(
            { project_name: createProjectDto.getProjectName(), project_description: createProjectDto.getProjectDescription() },
            createProjectDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            ProjectsException.fromApiResponse(apiResponse);
        }

        const project = apiResponse.data ?? {};
        return CreateProjectResultDto.fromPrimitives({
            projectUuid: project.project_uuid,
            projectName: project.project_name,
        });
    }
}
