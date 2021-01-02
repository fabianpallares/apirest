# apirest: Paquete para generar aplicaciones API-REST. 

El paquete apirest, tiene la funcionalidad de poder crear un aplicativo que
exponga endpoints API-REST de forma simple y sencilla.

## Instalación:
Para instalar el paquete utilice la siguiente sentencia:
```
go get -u github.com/fabianpallares/apirest
```

## Enrutador:
Una de las principales funcionalidades del paquete, es la de poder gestionar las
rutas (los endpoints) del aplicativo de una manera muy simple y completa:

```GO
package main

import (
    "fmt"
    "github.com/fabianpallares/apirest"
)

func main() {
    var r = apirest.CrearEnrutador()

    r.GET("/hola", hola)

    if err := r.IniciarPorHTTP(app.apirest.puerto); err != nil {
        apirest.Finalizar(fmt.Sprintf("No es posible iniciar el servidor: %v", err.Error()))
    }
}

func hola(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    var retorno = struct {
        Mensaje string `json:"mensaje"`
    }{Mensaje: "Hola enrutador"}

    apirest.HTTPResponder(w, apirest.HTTPEstadoOk, apirest.HTTPContenidoApplicationJSON, nil, retorno)

    return retorno, nil
}

```

Continuará... :)

#### Documentación:
[Documentación en godoc](https://godoc.org/github.com/fabianpallares/apirest)