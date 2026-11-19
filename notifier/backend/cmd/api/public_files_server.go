package main

import (
    "bytes"
    "html"
    "io/fs"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "sync"
    "time"

    "notifier-backend/src/core/boot"
    "notifier-backend/src/modules/shared/domain/enums"
)

const (
    indexUrlPath = "/index.html"
    jsFolder     = "js"
    // cada cuanto se vuelve a montar la portada con la lista de modulos (un cambio en public/js se ve en segundos)
    indexHtmlTtl = 15 * time.Second
    // cada cuanto se comprueba (un stat) si un .js o .css servido desde memoria ha cambiado en disco
    frontFileTtl = 10 * time.Second
)

// PublicFilesServer lo que en task-php hace nginx: si la url es un fichero de public/ se sirve tal cual; si no, pasa
// al front controller (try_files $uri @rewriteapp). Los ficheros ocultos (.env, .git...) nunca: 404.
type PublicFilesServer struct {
    indexMain http.Handler

    indexHtmlMutex   sync.Mutex
    indexHtml        []byte
    indexHtmlBuiltAt time.Time

    frontFilesMutex sync.Mutex
    frontFiles      map[string]*frontFile
}

func PublicFilesServerGetInstance(indexMain http.Handler) *PublicFilesServer {
    return &PublicFilesServer{indexMain: indexMain}
}

func (s *PublicFilesServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
    if strings.Contains(request.URL.Path, "/.") {
        http.NotFound(responseWriter, request)
        return
    }

    urlPath := request.URL.Path
    // la portada es el front: public/index.html
    if urlPath == "/" {
        urlPath = indexUrlPath
    }
    pathFile := s.getPathPublicFile(urlPath)
    if pathFile == "" {
        s.indexMain.ServeHTTP(responseWriter, request)
        return
    }

    // el front son ~130 modulos sin build ni hash en el nombre: sin esto el navegador decide por su cuenta cuanto
    // los da por buenos y, tras actualizar, mezcla modulos nuevos con viejos. no-cache = guardalos, pero pregunta
    // siempre si han cambiado (304 si no): nunca se sirve una version antigua.
    responseWriter.Header().Set(enums.HttpEnum.HEADER_CACHE_CONTROL, "no-cache")
    if urlPath == indexUrlPath {
        if indexHtml := s.getIndexHtmlWithModulePreloads(pathFile); indexHtml != nil {
            responseWriter.Header().Set(enums.HttpEnum.HEADER_CONTENT_TYPE, enums.HttpEnum.CONTENT_TYPE_HTML)
            _, _ = responseWriter.Write(indexHtml)
            return
        }
    }
    if cachedFile := s.getCachedFrontFile(urlPath, pathFile); cachedFile != nil {
        // ServeContent pone el Content-Type por la extension y responde 304 a If-Modified-Since
        http.ServeContent(responseWriter, request, filepath.Base(pathFile), cachedFile.modifiedAt, bytes.NewReader(cachedFile.content))
        return
    }
    http.ServeFile(responseWriter, request, pathFile)
}

// frontFile un estatico del front en memoria.
type frontFile struct {
    content    []byte
    modifiedAt time.Time
    checkedAt  time.Time
}

// getCachedFrontFile los .js y .css del front desde memoria (son ~130 ficheros pequeños, medio MB): cargar el panel
// eran 130 lecturas de disco, lentas en una carpeta montada (Docker Desktop). Cada frontFileTtl se mira si el fichero
// ha cambiado (un stat) y solo entonces se relee: un cambio en public/ se ve en segundos, sin reiniciar. nil para lo
// que no se cachea (los binarios de downloads/, imagenes...) o si la lectura falla: lo sirve http.ServeFile.
func (s *PublicFilesServer) getCachedFrontFile(urlPath string, pathFile string) *frontFile {
    if !strings.HasSuffix(urlPath, ".js") && !strings.HasSuffix(urlPath, ".css") {
        return nil
    }
    s.frontFilesMutex.Lock()
    cachedFile := s.frontFiles[pathFile]
    s.frontFilesMutex.Unlock()
    if cachedFile != nil && time.Since(cachedFile.checkedAt) < frontFileTtl {
        return cachedFile
    }

    // el disco se toca sin el cerrojo: las comprobaciones de los 130 modulos van en paralelo, no en fila. Cada
    // frontFile es inmutable (se sustituye entero): dos peticiones a la vez leen lo mismo y gana la ultima.
    fileInfo, err := os.Stat(pathFile)
    if err != nil {
        return nil
    }
    freshFile := &frontFile{modifiedAt: fileInfo.ModTime(), checkedAt: time.Now()}
    if cachedFile != nil && fileInfo.ModTime().Equal(cachedFile.modifiedAt) {
        freshFile.content = cachedFile.content
    } else if freshFile.content, err = os.ReadFile(pathFile); err != nil {
        return nil
    }

    s.frontFilesMutex.Lock()
    if s.frontFiles == nil {
        s.frontFiles = map[string]*frontFile{}
    }
    s.frontFiles[pathFile] = freshFile
    s.frontFilesMutex.Unlock()
    return freshFile
}

// getIndexHtmlWithModulePreloads index.html con un <link rel="modulepreload"> por cada modulo de public/js antes de
// </head>. Sin ellos el navegador descubre los import nivel a nivel (pide un modulo, lo lee, pide los suyos...) y
// cargar el panel son segundos; con ellos los pide todos en paralelo desde el principio. Hace de "build" sin build:
// la lista sale de la carpeta, no hay que mantenerla. nil si algo falla (se sirve el fichero tal cual).
func (s *PublicFilesServer) getIndexHtmlWithModulePreloads(pathIndexFile string) []byte {
    s.indexHtmlMutex.Lock()
    defer s.indexHtmlMutex.Unlock()
    if s.indexHtml != nil && time.Since(s.indexHtmlBuiltAt) < indexHtmlTtl {
        return s.indexHtml
    }

    rawIndexHtml, err := os.ReadFile(pathIndexFile)
    if err != nil {
        return nil
    }
    modulePreloads := strings.Builder{}
    for _, moduleUrlPath := range s.getModuleUrlPaths() {
        modulePreloads.WriteString(`    <link rel="modulepreload" href="` + html.EscapeString(moduleUrlPath) + `">` + "\n")
    }
    indexHtml := strings.Replace(string(rawIndexHtml), "</head>", modulePreloads.String()+"</head>", 1)

    s.indexHtml = []byte(indexHtml)
    s.indexHtmlBuiltAt = time.Now()
    return s.indexHtml
}

// getModuleUrlPaths /js/main.js, /js/Modules/..., por orden alfabetico (misma portada mientras no cambie la carpeta).
func (s *PublicFilesServer) getModuleUrlPaths() []string {
    moduleUrlPaths := []string{}
    pathJs := filepath.Join(boot.BootPathsEnum.PATH_PUBLIC, jsFolder)
    _ = filepath.WalkDir(pathJs, func(pathFile string, dirEntry fs.DirEntry, err error) error {
        if err != nil || dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".js") {
            return nil
        }
        relativePath, err := filepath.Rel(boot.BootPathsEnum.PATH_PUBLIC, pathFile)
        if err == nil {
            moduleUrlPaths = append(moduleUrlPaths, "/"+filepath.ToSlash(relativePath))
        }
        return nil
    })
    sort.Strings(moduleUrlPaths)
    return moduleUrlPaths
}

// getPathPublicFile la ruta del fichero dentro de public/, o "" si no existe o no es un fichero.
func (s *PublicFilesServer) getPathPublicFile(urlPath string) string {
    relativePath := filepath.FromSlash(filepath.Clean("/" + urlPath))
    pathFile := filepath.Join(boot.BootPathsEnum.PATH_PUBLIC, relativePath)
    if !strings.HasPrefix(pathFile, boot.BootPathsEnum.PATH_PUBLIC+string(filepath.Separator)) {
        return ""
    }
    // un .js o .css que ya esta en memoria y recien comprobado no toca el disco
    if s.isFrontFileFresh(pathFile) {
        return pathFile
    }
    fileInfo, err := os.Stat(pathFile)
    if err != nil || fileInfo.IsDir() {
        return ""
    }
    return pathFile
}

func (s *PublicFilesServer) isFrontFileFresh(pathFile string) bool {
    s.frontFilesMutex.Lock()
    defer s.frontFilesMutex.Unlock()
    cachedFile := s.frontFiles[pathFile]
    return cachedFile != nil && time.Since(cachedFile.checkedAt) < frontFileTtl
}
