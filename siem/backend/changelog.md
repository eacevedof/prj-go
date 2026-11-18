# changelog — siem-backend

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(la muestra el dashboard). Lo más nuevo arriba.

## 0.11.1 — 2026-11-18

1. Repaso DDD: el separador de las listas de `connections` (`,`) pasa de constante suelta del validador a
   `MonitorFieldsLimitEnum.LIST_SEPARATOR`. Sin cambio de comportamiento.

## 0.11.0 — 2026-11-18

1. **Monitor `connections`** (el siem-client pasa a ser el agente de red; lo enciende el backend, no el instalador):
   config `is_loopback_included`, `ignore_processes` e `ignore_ports` (hasta 50 de cada, como lista JSON o texto con
   comas; `MonitorFieldsValidator` los normaliza sin vacíos ni repetidos). `interval_seconds` es el muestreo, 30 s por
   defecto. Sin umbrales. Uno por agente (clave de recurso = el tipo), también con ámbito "todos los agentes". Los
   siem-client anteriores a la 0.5.0 lo ignoran con un aviso (no rompe su config).
2. **Evento `net_connection`** → source `network` (una por proceso). Gravedad en `EventSeverityResolver`:
   `is_unavailable` (sistema no soportado) warning; conexión entrante desde IP pública o puerto nuevo en escucha fuera de
   loopback, info; el resto (salientes, cierres, loopback, red local), debug. Lleva `monitor_uuid`: sale en los filtros
   de Eventos sin más.
3. **Panel**: el modal de monitores separa el intervalo (métricas y conexiones) de los umbrales (solo métricas) y trae
   los campos de conexiones; la tabla enseña su config y su intervalo.
4. CHECK de `monitor_type`, `event_type` y `source_type` editados en sus migraciones originales (sin migración nueva:
   no hay producción; la BD local se recreó).

## 0.10.0 — 2026-11-17

1. **Eventos — paginador de enlaces arriba del listado** (y repetido abajo): `← anteriores · página N · siguientes →` son
   enlaces con el hash de su página (se pueden abrir en otra pestaña), no botones; sin página anterior/siguiente el
   texto queda apagado. Pedir la página desde el de abajo sube al principio de la nueva.
2. **Eventos — velo de carga**: al cambiar de página, filtros o rango, la tabla anterior se atenúa con una barra que
   corre encima y el paginador dice "cargando…" hasta que llega la respuesta. El auto-refresco no lo enseña (misma
   consulta: parpadearía cada pocos segundos).
3. **Carpeta de instalación del agente `gosiem-client`** (antes `siem-client`): `/opt/gosiem-client`,
   `C:\programs\gosiem-client` y `/usr/local/gosiem-client` en los pasos de Descargas, en las rutas por defecto del
   comando de desinstalar (`ClientOsEnum`) y en los tests de `AgentHostDataSanitizer`. El comando de desinstalar borra
   la carpeta entera si se llama `gosiem-client` o `siem-client`: los agentes instalados antes siguen saliendo bien.

## 0.9.1 — 2026-11-16

1. **Seguridad — `host_os` y `client_path` saneados al ingestar** (`AgentHostDataSanitizer` + `AgentHostOsEnum` en el
   dominio de `agents`): el agente los informa de sí mismo y el panel monta con ellos el comando de desinstalar que el
   admin pega en un shell elevado. `host_os` fuera de `linux|windows|darwin` y una ruta que no sea absoluta y limpia
   para su SO (sin comillas, tampoco tipográficas, ni `$`, `;`, `&`…) se guardan vacías: el panel cae en la ruta por
   defecto. Test en `tests/modules/agents/domain/services`.
2. **Seguridad — front**: `ClientUninstallCommandBuilder` dobla también las comillas tipográficas U+2018-U+201B en el
   comando de PowerShell (PowerShell las toma como comilla simple: inyección por copiar y pegar).
3. Comentario de `AuthenticateDashboardService` al día: el token del panel lee y escribe (decía "solo lectura").
   Auditoría completa y reglas en el vault: `projects/siem/siem-security.md`.

## 0.9.0 — 2026-11-16

1. **Eventos es la pantalla de entrada**: `#/` (y un hash desconocido) lleva al visor. Menú reordenado: 1 Eventos ·
   2 Proyectos · 3 Descargas; la lista de proyectos pasa a `#/projects` (el detalle sigue en `#/projects/<uuid>`). El
   visor pierde su cabecera editorial: arranca en "1 - consulta". Hashes en `RouteHashEnum`.
2. **Filtros múltiples Proyecto → Agente → Monitor** en el visor: tres desplegables de casillas (`<details>`), con
   resumen ("Todos", el nombre o "3 agentes"). Encadenados: los proyectos marcados acotan los agentes que se ofrecen, y
   proyectos + agentes, los monitores (uno de "todos los agentes" vale para cualquier agente de su proyecto); al
   desmarcar se quita lo que deja de ser suyo. Viven en el hash como los demás filtros (`project_uuid=a,b`…), salen en
   las pastillas con su nombre y el auto-refresco no cierra el desplegable abierto. Clic fuera lo cierra.
3. API: filtro **`monitor_uuid`** (y `not_monitor_uuid`) en `GET /api/v1/events`, `/histogram` y `/facets`;
   `GET /api/v1/events/filter-options` → `{projects, agents, monitors}` (el catálogo entero, no solo lo que tiene
   eventos como las facetas). El front lo pide aparte, sin esperar a las facetas.
4. Columna **`events.monitor_uuid`** (copia de `attributes.monitor_uuid` al insertar, índice `(monitor_uuid,
   occurred_at)`; la migración rellena los eventos que ya había): filtrar con `json_extract` sobre `attributes`
   recorría la tabla (facetas de 5 a 9 s). Excluir un monitor conserva los eventos que no son de ninguno.
5. Índice **`events(occurred_at)`**: el visor sin filtros solo acota por fecha y hacía SCAN de la tabla.
6. **Prefijo `app_` en todas las tablas** (`app_projects`, `app_agents`, `app_sources`, `app_ingest_queue`,
   `app_events`, `app_monitors`) con una migración de `RENAME` (conserva datos; SQLite reescribe las FK). `migrations`,
   la del runner, se queda como está.
7. **El panel se llama GoSiem** y se ve: logo (el favicon en línea: pulso blanco sobre grafito con el punto ember) +
   «GoSiem» en la barra y en la entrada, y en el título de la pestaña. Fuera «siem · consola».
8. **Descargas, macOS: «3.4 Si falla en Mac»** (solo con macOS elegido): el binario se compila sin Mac y sin firma de
   Apple ni prueba en un Mac real, así que se explican los fallos esperables y su salida: Gatekeeper (`xattr` o «Abrir
   igualmente»), «Killed: 9» (firma ad hoc con `codesign`), «bad CPU type» (solo Apple Silicon), sin arranque con el
   sistema (aún sin LaunchDaemon), métricas del host «no disponible» (solo Linux) y «nunca visto» (`config-sync`).
   Textos en `DownloadsTextEnum`.
9. **Repaso DDD de todo el back y el front** (reglas del vault): literales a enums (`*LimitEnum`, `*KeyEnum`, `*TextEnum`,
   `HttpEnum`, `RequestMetadataKeyEnum`…), funciones sueltas a métodos, nombres sin abreviar y por CQS (`load…OrFail`,
   fuera `process…`), orden de imports, errores en inglés en `*MessageEnum` (el del rango a medida del visor ahora sale en
   inglés) y LF en tres ficheros que eran CRLF. Sin cambios de comportamiento, rutas, claves JSON ni SQL. Lo que queda,
   en `siem-tech-debt` §5.

## 0.8.0 — 2026-11-16

1. **Borrar proyecto en un modal**: botón rojo "Borrar proyecto" en la cabecera del detalle (junto a "ver sus eventos")
   y "Borrar" (rojo) en la tabla de la portada abren el mismo `<dialog>`, que pide escribir el nombre del proyecto. Se
   quitan el apartado "4 - borrar proyecto" del detalle y la confirmación en la fila de la tabla. La cabecera del
   detalle ya no enseña el uuid del proyecto (está en la url y en la tabla).
2. **Qué queda instalado**: al borrar se vuelve a la portada con un aviso que se queda hasta "Entendido": por cada
   agente, su servidor (host · SO · ip), la ruta de su ejecutable (Copiar) y el comando para desinstalarlo (bash:
   servicio systemd + carpeta; PowerShell: tarea programada + carpeta). Si el binario no está en una carpeta
   `siem-client` (p. ej. `/usr/local/bin`) se borra solo el binario, nunca la carpeta. Un agente que nunca conectó
   solo dice dónde estaría.
3. API: `DELETE /api/v1/projects/:project_uuid` devuelve `{project_uuid, project_name, agents: [{agent_name, host_name,
   host_os, host_ip, client_path}]}` (leídos antes de la cascada).
4. Columna `agents.client_path` (migración `20261116100000-add-client-path-to-agents.sql`, sin borrar la BD): ruta del
   ejecutable que el siem-client ≥ 0.4.0 manda en cada envío (`agent.client_path`); sale en el detalle del proyecto.
5. Portada sin cabecera editorial ("Lo que vigilan tus agentes, de un vistazo" + entradilla): arranca en "1 - servidor"
   bajo las migas, como el detalle de un proyecto.
6. Front: `ProjectDeleteDialogView` y `ProjectDeletedView` sustituyen a `ProjectDangerView`; un solo
   `DeleteProjectController` (fuera `DeleteListedProjectController` y `DELETE_LISTED_PROJECT_REQUESTED`); servicio de
   dominio `ClientUninstallCommandBuilder`. Botón `button-danger-fill`.

## 0.7.0 — 2026-11-16

1. **Crear monitor desde el agente**: botón "Crear monitor" en las acciones de cada fila de agentes; abre un modal
   (`<dialog>` nativo, patrón task-js) con el agente fijado. En "monitores", "Crear monitor para todos los agentes"
   abre el mismo modal para los generales (`agent_uuid` vacío), y "Editar" también edita en el modal. Se quita el
   formulario fijo del apartado 4 y el select de ámbito: el ámbito lo decide desde dónde se abre y no se cambia al
   editar (como el tipo). "Borrar proyecto" pasa a ser el apartado 4.
2. **Monitores repetidos**: además del nombre (único por proyecto), no puede haber dos monitores que vigilen el mismo
   recurso en un agente que aplicaría los dos. Recurso: cpu, memoria y gpu, uno por agente · disco, la ruta · red, la
   interfaz · carpeta, ruta + glob · log, ruta + glob + regex. Rutas sin separador final y, las de Windows, sin
   distinguir mayúsculas. Uno de "todos" choca con cualquiera del proyecto; uno de un agente, con los suyos y con los
   de todos. El modal avisa mientras se escribe (con el monitor que choca) y no deja enviar; la API responde 409
   (`POST` y `PUT`). Servicio de dominio `MonitorResourceKeyer` (Go, con tests) y su gemelo `MonitorDuplicateFinder`
   (front).

## 0.6.0 — 2026-11-16

1. **Proyectos en tabla CRUD** (patrón task-js) en lugar de tarjetas: columnas ID, UUID (con Copiar), proyecto (nombre
   con su semáforo y la descripción debajo; si es larga, recortada y entera en el tooltip), warnings, errores y
   críticos de 24 h, y acciones (Editar, Borrar). Filtros en el front: texto (id, uuid, nombre, descripción; con
   debounce) y salud; contador "(2 de 4)". Filtros y confirmación pendiente sobreviven al auto-refresco.
2. Formulario **plegable** "Nuevo proyecto", el mismo para editar ("Editar proyecto #3", uuid oculto, "Cancelar
   edición"): foco en el nombre al desplegar y Ctrl+Enter guarda. Descripción en textarea (hasta 500).
3. Borrar desde la tabla con confirmación en la fila ("Sí, borrar"); se queda en la portada. El borrado desde el
   detalle (escribiendo el nombre) sigue igual.
4. API: `PUT /api/v1/projects/:project_uuid` `{project_name, project_description}` (404 si no existe, 409 si el nombre
   es de otro proyecto; cambiar solo mayúsculas vale). `GET /api/v1/dashboard` añade `project_id` a cada proyecto.
5. Front: la tabla es del módulo Projects (escucha `DASHBOARD_RENDERED`); Dashboard solo pinta el servidor. Casos de
   uso nuevos `UpdateProject` y `FilterProjects`; borrar desde la lista es su propio controller
   (`DeleteListedProjectController`), sin flag. Componente `TextFormatter` para los `{0}` de los textos.
6. **Agentes en tabla CRUD**, igual que proyectos, en el detalle de un proyecto: ID, UUID (Copiar), agente (nombre con
   su estado; host y último latido debajo), warnings, errores y críticos de 24 h, acciones (Ver detalle, Editar,
   Borrar). "Ver detalle" despliega bajo la fila lo que antes era la tarjeta: latido, host, cliente y la tabla de sus
   fuentes. Filtros en el front (texto en id, uuid, nombre, host, ip y versión; estado). Filas abiertas, filtros y
   confirmación sobreviven al auto-refresco; al cambiar de proyecto se vacían.
7. Formulario plegable "Nuevo agente" que también edita ("Editar agente #4"): nombre, latido y, en edición, **Activo**
   (desactivado, su token recibe 401 hasta reactivarlo; no se borra). El token del alta se sigue enseñando una vez.
8. API: `PUT /api/v1/agents/:agent_uuid` `{agent_name, heartbeat_interval_seconds, is_enabled}` (404, 409 si el nombre
   es de otro agente del proyecto). `GET /api/v1/projects/:uuid` añade `agent_id` e `is_enabled` a cada agente.
9. **Migas de pan** en todas las pantallas (Proyectos · Proyectos / <nombre> · Eventos · Descargas), en lugar de los
   eyebrow y del "← Proyectos"; las pinta `NavigationView` y el detalle las completa con el nombre.
10. Fix: el front validaba el latido entre 5 y 3600 s y el nombre desde 1 carácter, pero la API exige 10-86400 s y 3-100
    caracteres (un latido de 5 pasaba el front y la API lo rechazaba). Ahora los límites viven en `AgentFieldsEnum`.
11. Detalle de proyecto renumerado: 2 agentes (con su formulario), 3 monitores, 4 nuevo monitor, 5 borrar proyecto.

## 0.5.1 — 2026-11-16

1. Ejemplos neutros: el placeholder de "nuevo proyecto" y la ayuda de `project-create` dicen `Mi app` en vez de
   un producto de Lazarus.

## 0.5.0 — 2026-11-16

1. **Rediseño UI/UX del panel** con el lenguaje de Ventriloc (refero.design): casi monocromo (grafito sobre blanco,
   superficies ash/fog/ivory) con un único acento naranja Ember y latón como secundario; sin sombras ni bordes de caja.
2. Barra lateral oscura sustituida por una **cabecera**: marca a la izquierda, menú en pastilla al centro, refresco y
   "Salir" a la derecha; contenido centrado a 1200 px.
3. Cada pantalla abre con una **cabecera editorial** (titular grande + entradilla); en Eventos es compacta, en una línea,
   para que consulta, histograma y tabla queden en el primer pliegue.
4. Botones rectos (grafito relleno o contorno), tarjetas de datos de radio 20 (cifras del servidor, proyectos,
   histograma), paneles con radio asimétrico, pastillas de estado y filtros redondeadas, tablas sin caja con filetes.
5. Gravedades en escala cálida (gris → latón → ember → rojo tostado) que mantienen la forma de cada pastilla; tema oscuro
   cálido por tokens. Solo `index.html` y `siem.css`: ids y clases que usa el JS, sin cambios.
6. Visor de eventos con **un solo scroll** (el de la página): la tabla ya no tiene altura máxima con scroll propio y
   su cabecera se pega bajo la barra superior (`--topbar-h`); en estrecho vuelve el desplazamiento horizontal. Fix: el
   recorte a 2 líneas del mensaje iba en el `<td>` (dejaba de ser celda, asomaba media línea y descuadraba la fila):
   ahora va en un `div` interior.
7. **Favicon** de observabilidad (pulso blanco sobre grafito con el punto ember): `public/favicon.svg` y `public/favicon.ico`
   (16/32/48 px, generado del SVG); antes `GET /favicon.ico` daba 404.
8. Descargas: instalación asistida (`install.sh`) en el paso 3, trabajo previo que no había subido versión.

## 0.4.0 — 2026-11-15

1. **Panel con escritura** (Bearer del dashboard): `POST /api/v1/projects`, `DELETE /api/v1/projects/:uuid`,
   `POST /api/v1/projects/:uuid/agents` (devuelve el `api_token` una vez), `DELETE /api/v1/agents/:uuid`.
2. **Monitores** (recursos a vigilar por proyecto, módulo `monitors`, tabla `monitors`): `cpu`, `memory`, `disk`, `gpu`,
   `network`, `folder`, `log_pattern`, con umbrales warning/critical. CRUD desde el panel y `GET /api/v1/agent-config`
   (Bearer del agente) para que el agente los descargue.
3. Evento `host_metric` (fuente `metric`): la gravedad sale de los umbrales del monitor, que pone el backend al
   procesar el lote (por `attributes.monitor_uuid`); lo que mande el agente en esas claves se descarta.
4. **Visor tipo Kibana**: filtros con rango `from`/`to`, varios valores por campo e inversos (`not_<campo>`),
   `source_type`; `GET /api/v1/events/histogram` (barras por tramo y gravedad) y `GET /api/v1/events/facets` (top de
   valores por campo, con el parámetro por el que se filtra).
5. **Descargas**: `GET /api/v1/downloads` lista los binarios del agente de `public/downloads/` (tamaño, sha256) y la
   url de ingesta; se sirven en `/downloads/<fichero>` sin token.
6. Front: rediseño, altas desde el panel, visor con histograma, facetas y pastillas de filtro, pantalla de descargas
   con instrucciones de instalación copiables.
7. Migraciones originales cambiadas otra vez (sin producción): hay que borrar la BD local.

## 0.3.0 — 2026-11-15

1. Eventos de Docker: tipos `docker_state`, `docker_stats`, `docker_event`, `docker_log`; fuentes `docker`, `docker_stats`,
   `docker_log`; columna `events.attributes` (JSON, máx. 4000 bytes) que la API devuelve como objeto. Migraciones
   originales cambiadas (sin producción): hay que borrar la BD local.
2. Gravedad de Docker en `EventSeverityResolver`: OOM o `dead` critical; parado con exit ≠ 0 o `unhealthy` error; parado
   con exit 0, reiniciando o reiniciado warning; recursos debug salvo memoria ≥ 85 % (warning) / ≥ 95 % (error) o
   cpu ≥ 90 % (warning); daemon: `oom` critical, `die` con exit ≠ 0 error, `kill` warning.
3. Front: los tipos nuevos en el filtro del visor de eventos.

## 0.2.0 — 2026-11-15

1. Front rehecho con el patrón de `task-js` (maestro de los fronts): `public/js/Modules/{Shared,Session,Dashboard,Events}`
   con casos de uso (`Dto` + `ResultDto` + `Service`), enums, excepciones, controllers por evento (`Dispatcher`/
   `Subscriber`), repositorios `*ReaderHttpRepository` y vistas con referencias al DOM; `index.html` con secciones fijas.
2. Fechas siempre en UTC, formato `YYYY-MM-DD HH:MM:SS` (24 h): ya no se pasan a la hora local del navegador; los
   "hace X" van al lado. La portada añade la hora del servidor y la de arranque.
3. Un refresco fallido (backend caído) deja lo último pintado con el error encima y sigue reintentando cada 10 s.

## 0.1.0 — 2026-11-15

1. Primera versión, esqueleto de `task-go`: módulos `projects`, `agents`, `ingest`, `events`, `dashboard` (+ `devops`,
   `mailer`, `shared`).
2. `POST /api/v1/ingest-batches` (Bearer del agente): latido + lote crudo a `ingest_queue`, 202; un reenvío con el mismo
   `batch_uuid` no duplica.
3. Worker en la API (cada `APP_QUEUE_WORKER_INTERVAL_SECONDS`): cola → `events` + `sources`, una transacción por lote,
   3 intentos y `failed`.
4. Dashboard en `/` (token `APP_DASHBOARD_TOKEN`): servidor, semáforo por proyecto, agentes y fuentes, visor de eventos.
5. Consola: `project-create`, `agent-create`, `ingest-queue-process`, `events-purge`, `run-migrations`.
6. SQLite en WAL; log de SQL solo con `APP_DEBUG=1`.
