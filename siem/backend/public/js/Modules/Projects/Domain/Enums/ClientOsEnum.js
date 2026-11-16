/**
 * host_os que informa el siem-client (runtime.GOOS de Go) y lo que deja su instalador en cada uno
 * (client/install/install.sh e install.ps1): el nombre del servicio / tarea y la carpeta por defecto.
 */
export const ClientOsEnum = Object.freeze({
    LINUX: "linux",
    WINDOWS: "windows",
    MACOS: "darwin",

    SERVICE_NAME: "siem-client",
    // carpeta que crea el instalador: solo si el ejecutable vive en una carpeta con este nombre se propone borrarla
    // entera (si alguien lo puso en /usr/local/bin, se borra el binario, nunca la carpeta)
    INSTALL_DIR_NAME: "siem-client",
    DEFAULT_PATH_LINUX: "/opt/siem-client/siem-client",
    DEFAULT_PATH_WINDOWS: "C:\\programs\\siem-client\\siem-client.exe",
});
