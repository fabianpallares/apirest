package apirest

import (
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	randc "crypto/rand"
	randm "math/rand"
)

type errorRastro struct {
	Paquete       string // nombre del paquete donde se originó el error
	Archivo       string // nombre del archivo donde se originó el error
	Funcion       string // nombre de la función donde se originó el error
	NroLinea      int    // número de línea donde se originó el error
	Observaciones string // observaciones adicionales del rastro actual (es opcional)
}

type errorAPIREST struct {
	estadoHTTP         HTTPEstado  // código de estado HTTP
	codigo             string      // código de error
	mensaje            string      // mensaje de error
	mensajeTecnico     string      // mensaje técnico para ser leído por el desarrollador
	valoresAdicionales []string    // lista de textos que puede ser utilizados para exponer campos con error
	uuid               string      // identificador único universal del error
	errAnterior        error       // error anterior
	rastro             errorRastro // rastro/ubicación donde se originó el error
}

// Error retorna el mensaje de error (implementa la interface error).
func (o *errorAPIREST) Error() string {
	return o.mensaje
}

// AsignarCodigo asiga el código de error.
func (o *errorAPIREST) AsignarCodigo(codigo string) *errorAPIREST {
	o.codigo = codigo
	return o
}

// AsignarMensajeTecnico asigna un mensaje técnico (un mensaje para desarrollador)
// al error actual. Además, este mensaje se incorpora a la observación del rastro
// en caso que dicha observación no posea ningún valor.
// De esta manera, al obtener los rastros del error, serán devueltos con su
// observación correspondiente.
func (o *errorAPIREST) AsignarMensajeTecnico(formato string, args ...interface{}) *errorAPIREST {
	o.mensajeTecnico = fmt.Sprintf(formato, args...)

	// asignar a la observación del rasto, el mensaje técnico.
	// si la observación ya posee un valor, no se asigna.
	if o.rastro.Observaciones == "" {
		o.rastro.Observaciones = o.mensajeTecnico
	}
	return o
}

// AsignarValoresAdicionales asiga valores (textos) a la lista adicional de valores.
// Entre otros usos, esta lista de valores puede ser utilizada para devolver los
// nombres de campos que poseen errores en una validación de entidad de negocio.
func (o *errorAPIREST) AsignarValoresAdicionales(valores ...string) *errorAPIREST {
	o.valoresAdicionales = valores
	return o
}

// AsignarUUID asiga un identificador único universal al error actual
// para que en caso de ser registrado, pueda realizarse un seguimiento.
// Para generar el UUID, se utiliza la versión 4 (random).
// 	https://es.wikipedia.org/wiki/Identificador_%C3%BAnico_universal
// 	https://tools.ietf.org/html/rfc4122
//
// 	Formato 8-4-4-4-12 (36 caracteres: 32 hexadecimales + 4 '-'):
// 		xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx donde:
// 		x es un valor hexadecimal (0, 1, 2... d, e, f).
// 		M es un valor de 1 a 5 (versión del UUID): 4 es UUID random.
// 		N es "8", "9", "a" o "b".
func (o *errorAPIREST) AsignarUUID() *errorAPIREST {
	var a [16]byte
	if _, err := io.ReadFull(randc.Reader, a[:]); err != nil {
		// si se produce un error, utilizar otro método de generación.
		randm.Seed(randm.Int63() + randm.Int63() + time.Now().UnixNano())
		for i := 0; i < 16; i++ {
			a[i] = byte(randm.Intn(256))
		}
	}

	versionUUID := byte(4)
	a[6] = (a[6] & 0x0f) | (versionUUID << 4)
	a[8] = (a[8]&(0xff>>2) | (0x02 << 6))

	b := make([]byte, 36)
	hex.Encode(b[0:8], a[0:4])
	hex.Encode(b[9:13], a[4:6])
	hex.Encode(b[14:18], a[6:8])
	hex.Encode(b[19:23], a[8:10])
	hex.Encode(b[24:], a[10:])
	for _, i := range []int{8, 13, 18, 23} {
		b[i] = '-'
	}
	o.uuid = string(b)

	return o
}

// AsignarObservacionAlRastro asiga una observación al rastro del error actual.
func (o *errorAPIREST) AsignarObservacionAlRastro(formato string, args ...interface{}) *errorAPIREST {
	o.rastro.Observaciones = fmt.Sprintf(formato, args...)
	return o
}

// ObtenerHTTPEstado devuelve el código de estado HTTP.
func (o *errorAPIREST) ObtenerHTTPEstado() HTTPEstado {
	return o.estadoHTTP
}

// ObtenerCodigo devuelve el código del error.
func (o *errorAPIREST) ObtenerCodigo() string {
	return o.codigo
}

// ObtenerMensajeTecnico devuelve el mensaje técnico del error (un mensaje para desarrollador).
func (o *errorAPIREST) ObtenerMensajeTecnico() string {
	return o.mensajeTecnico
}

// ObtenerUUID devuelve el identificador único universal del error.
func (o *errorAPIREST) ObtenerUUID() string {
	return o.uuid
}

// ObtenerValoresAdicionales devuelve los valores adicionales del error.
func (o *errorAPIREST) ObtenerValoresAdicionales() []string {
	return o.valoresAdicionales
}

// ObtenerRastro devuelve el ratro, sólo del error actual.
func (o *errorAPIREST) ObtenerRastro() *errorRastro {
	return &o.rastro
}

// asignarRastro establece la ubicación donde se ha generado el error.
func (o *errorAPIREST) asignarRastro(profundidad int) {
	_, paquete, _, ok := runtime.Caller(0)
	if ok {
		o.rastro.Paquete = paquete
	}

	pc, archivo, linea, ok := runtime.Caller(profundidad + 1)
	if ok {
		o.rastro.Archivo = archivo[strings.LastIndex(archivo, "/")+1 : len(archivo)-3]
		o.rastro.Funcion = runtime.FuncForPC(pc).Name()
		o.rastro.NroLinea = linea
	}
}

func errorNuevo(estadoHTTP HTTPEstado, formato string, args ...interface{}) *errorAPIREST {
	err := &errorAPIREST{
		estadoHTTP: estadoHTTP,
		mensaje:    fmt.Sprintf(formato, args...),
	}
	err.asignarRastro(2)

	return err
}

func errorBuscarTipo(err error, estadoHTTP HTTPEstado) (*errorAPIREST, bool) {
	for {
		errAPIREST, ok := err.(*errorAPIREST)
		if !ok {
			return nil, false
		}
		if errAPIREST.estadoHTTP == estadoHTTP {
			return errAPIREST, true
		}
		err = errAPIREST.errAnterior
	}
}

// -----------------------------------------------------------------------------
// Obtener el primer error que sea de tipo errorAPIREST.

// ErrorEsAPIREST busca y verifica en la cadena de errores internos y devuelve
// el error cuando encuentra un tipo de error determinado (estadoHTTP con valor).
func ErrorEsAPIREST(err error) (*errorAPIREST, bool) {
	for {
		errAPIREST, ok := err.(*errorAPIREST)
		if !ok {
			return nil, false
		}
		if errAPIREST.estadoHTTP.obtenerEntero() != 0 {
			return errAPIREST, true
		}
		err = errAPIREST.errAnterior
	}
}

// -----------------------------------------------------------------------------
// Agregar un rastro al error recibido.

// ErrorNuevoRastro devuelve un nuevo error donde se agrega el ratro.
func ErrorNuevoRastro(err error) *errorAPIREST {
	errAPIREST := &errorAPIREST{errAnterior: err}
	errAPIREST.asignarRastro(1)

	return errAPIREST
}

// -----------------------------------------------------------------------------
// Devolver una lista de los rastros por los que pasó el error actual.

// ErrorObtenerRastros devuelve la lista de rastros/ubicaciones por los que
// pasó el error actual. Devuelve el trazado del error actual.
func ErrorObtenerRastros(err error) []errorRastro {
	var rastros []errorRastro
	for {
		errAPIREST, ok := err.(*errorAPIREST)
		if !ok {
			break
		}

		rastros = append(rastros, errAPIREST.rastro)
		err = errAPIREST.errAnterior
	}

	return rastros
}

// -----------------------------------------------------------------------------
// Error mal requerimiento

// ErrorNuevoMalRequerimiento crea un error de tipo:
// 400 (Mal requerimiento).
func ErrorNuevoMalRequerimiento(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorMalRequerimiento, formato, args...)
}

// ErrorEsMalRequerimiento verifica que el error sea del tipo:
// 400 (Mal requerimiento).
func ErrorEsMalRequerimiento(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorMalRequerimiento)
}

// -----------------------------------------------------------------------------
// Error sin autorización

// ErrorNuevoSinAutorizacion crea un error de tipo:
// 401 (Sin autorización).
func ErrorNuevoSinAutorizacion(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorSinAutorizacion, formato, args...)
}

// ErrorEsSinAutorizacion verifica que el error sea del tipo:
// 401 (Sin autorización).
func ErrorEsSinAutorizacion(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorSinAutorizacion)
}

// -----------------------------------------------------------------------------
// Error sin privilegios

// ErrorNuevoSinPrivilegios crea un error de tipo:
// 403 (Sin privilegios).
func ErrorNuevoSinPrivilegios(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorSinPrivilegios, formato, args...)
}

// ErrorEsSinPrivilegios verifica que el error sea del tipo:
// 403 (Sin privilegios).
func ErrorEsSinPrivilegios(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorSinPrivilegios)
}

// -----------------------------------------------------------------------------
// Error no encontrado

// ErrorNuevoNoEncontrado crea un error de tipo:
// 404 (No encontrado).
func ErrorNuevoNoEncontrado(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorNoEncontrado, formato, args...)
}

// ErrorEsNoEncontrado verifica que el error sea del tipo:
// 404 (No encontrado).
func ErrorEsNoEncontrado(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorNoEncontrado)
}

// -----------------------------------------------------------------------------
// Error método no implementado

// ErrorNuevoMetodoNoImplementado crea un error de tipo:
// 405 (Método no implementado).
func ErrorNuevoMetodoNoImplementado(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorMetodoNoImplementado, formato, args...)
}

// ErrorEsMetodoNoImplementado verifica que el error sea del tipo:
// 405 (Método no implementado).
func ErrorEsMetodoNoImplementado(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorMetodoNoImplementado)
}

// -----------------------------------------------------------------------------
// Error requerimiento muy grande

// ErrorNuevoRequerimientoMuyGrande crea un error de tipo:
// 413 (Requerimiento muy grande).
func ErrorNuevoRequerimientoMuyGrande(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorRequerimientoMuyGrande, formato, args...)
}

// ErrorEsRequerimientoMuyGrande verifica que el error sea del tipo:
// 413 (Requerimiento muy grande).
func ErrorEsRequerimientoMuyGrande(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorRequerimientoMuyGrande)
}

// -----------------------------------------------------------------------------
// Error uri muy grande

// ErrorNuevoURIMuyGrande crea un error de tipo:
// 414 (URI muy grande).
func ErrorNuevoURIMuyGrande(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorURIMuyGrande, formato, args...)
}

// ErrorEsURIMuyGrande verifica que el error sea del tipo:
// 414 (URI muy grande).
func ErrorEsURIMuyGrande(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorURIMuyGrande)
}

// -----------------------------------------------------------------------------
// Error mal formato

// ErrorNuevoMalFormato crea un error de tipo:
// 415 (mal formato).
func ErrorNuevoMalFormato(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorMalFormato, formato, args...)
}

// ErrorEsMalFormato verifica que el error sea del tipo:
// 415 (mal formato).
func ErrorEsMalFormato(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorMalFormato)
}

// -----------------------------------------------------------------------------
// Error interno del servidor

// ErrorNuevoInternoDeServidor crea un error de tipo:
// 500 (error interno del servidor).
func ErrorNuevoInternoDeServidor(formato string, args ...interface{}) *errorAPIREST {
	return errorNuevo(HTTPEstadoErrorInternoDeServidor, formato, args...)
}

// ErrorEsInternoDeServidor verifica que el error sea del tipo:
// 500 (error interno del servidor).
func ErrorEsInternoDeServidor(err error) (*errorAPIREST, bool) {
	return errorBuscarTipo(err, HTTPEstadoErrorInternoDeServidor)
}
