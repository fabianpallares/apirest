package apirest

import (
	"fmt"
	"runtime"
	"strings"
)

// errorEnvoltorio contiene una descripción de un error junto con información
// sobre donde se creó el error.
// Puede estar incrustado en tipos de envoltorio personalizados para agregar
// información adicional que este paquete de errores puede comprender.
type errorEnvoltorio struct {
	tipo    HTTPEstado // tipo de error, utilizado para generar la respuesta HTTP
	origen  error      // el error de origen es utilizado para almacenar un error generado por otro paquete
	mensaje string     // mensaje de error
	paquete string     // nombre del paquete donde se generó el error
	archivo string     // nombre del archivo donde se generó el error
	linea   int        // número de línea donde se generó el error
}

// Error implementa la interface error.
func (e *errorEnvoltorio) Error() string {
	return e.mensaje
}

// ObtenerUbicacion devuelve la ubicación donde se ha generado el error.
func (e *errorEnvoltorio) ObtenerUbicacion() (string, string, int) {
	return e.paquete, e.archivo, e.linea
}

// establecerUbicacion establece la ubicación donde se ha generado el error.
func (e *errorEnvoltorio) establecerUbicacion(profundidad int) {
	_, paquete, _, ok := runtime.Caller(0)
	if ok {
		e.paquete = paquete
	}

	_, archivo, linea, ok := runtime.Caller(profundidad + 1)
	if ok {
		pos := strings.LastIndex(archivo, "/")
		e.archivo = archivo[pos+1 : len(archivo)-3]
		e.linea = linea
	}
}

func errorNuevo(tipo HTTPEstado, origen error, formato string, args ...interface{}) error {
	ne := &errorEnvoltorio{
		tipo:    tipo,
		origen:  origen,
		mensaje: fmt.Sprintf(formato, args...),
	}
	ne.establecerUbicacion(2)

	return ne
}

// ErrorObtenerOrigen devuelve el error de origen.
func ErrorObtenerOrigen(err error) error {
	if errEnvoltorio, ok := err.(*errorEnvoltorio); ok {
		return errEnvoltorio.origen
	}

	return nil
}