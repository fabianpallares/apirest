# Cambios del aplicativo

## [1.2.0] 2021-04-09
### Agregados
Se agregó la función HTTPObtenerCuerpo(r *http.Request), la cuál devuelve el cuerpo del mensaje recibido como una cadena de caracteres (string).

## [1.1.0] 2021-03-25
### Modificaciones
* Se quitaron todos los métodos json (se creó un nuevo repositorio llamado json).
* HTTPResponder ahora recibe el parámetro "cuerpo" (último parámetro) como tipo texto (antes era de tipo interface). Por lo tanto responde el texto que recibe.

## [1.0.0] 2020-12-02
### Agregados
* Primera versión del aplicativo.
