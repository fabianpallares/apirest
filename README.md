# bcrypt: Algoritmo de encriptación para Go/Golang. Es utilizado idealmente para generar, almacenar y validar contraseñas de usuarios. 

[![Go Report Card](https://goreportcard.com/badge/github.com/fabianpallares/bcrypt)](https://goreportcard.com/report/github.com/fabianpallares/bcrypt) [![GoDoc](https://godoc.org/github.com/fabianpallares/bcrypt?status.svg)](https://godoc.org/github.com/fabianpallares/bcrypt)

Siempre que se almacen contraseñas en las bases de datos, estas deberían estar encriptadas.
El algoritmo bcrypt es ideal para realizar esta tarea.

Este paquete es un envoltorio del paquete nativo de Go/Golang: "golang.org/x/crypto/bcrypt".

Simplemente expone dos funciones en español, con la idea de agilizar el desarrollo.

## Instalación:
Para instalar el paquete utilice la siguiente sentencia:
```
go get -u github.com/fabianpallares/bcrypt
```

## Generar hash:
Para generar un hash de una contraseña recibida, utilizar la siguiente función:

```GO
package main

import (
    "fmt"
    "github.com/fabianpallares/bcrypt"
)

func main() {
    hash, err := bcrypt.Encriptar("EstaEsUnaClave")
    if err != nil {
        // tratar el error
    }

    fmt.Println("El hash generado es:", hash)
}
```

## Validar contraseña recibida:
Para validar la contraseña recibida (tipo texto) contra el hash almacenado, utilizar la siguiente función:

```GO
package main

import (
    "fmt"
    "github.com/fabianpallares/bcrypt"
)

func main() {
    clave := "EstaEsUnaClave"

    hash, err := bcrypt.Encriptar(clave)
    if err != nil {
        // tratar el error
    }
    fmt.Println("El hash generado es:", hash)

    // validar la contraseña con el hash:
    if esValida := bcrypt.Validar(clave, hash); !esValida {
        fmt.Println("La clave recibida no es válida")
    }

    fmt.Println("La clave recibida es válida!")
}
```
#### Documentación:
[Documentación en godoc](https://godoc.org/github.com/fabianpallares/bcrypt)