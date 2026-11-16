# changelog — siem-backend

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(la muestra el dashboard). Lo más nuevo arriba.

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
