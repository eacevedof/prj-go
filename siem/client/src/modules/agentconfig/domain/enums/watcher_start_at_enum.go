package enums

// WatcherStartAtEnum desde donde se leen los ficheros que ya existen la PRIMERA vez que se vigila una carpeta:
// end = solo lo que se escriba a partir de ahora (no reenvia el historico); beginning = todo el fichero.
// Los ficheros que aparecen despues siempre se leen desde el principio.
var WatcherStartAtEnum = struct {
    END       string
    BEGINNING string
    VALUES    []string
}{
    END:       "end",
    BEGINNING: "beginning",
    VALUES:    []string{"end", "beginning"},
}
