package apirest

import (
	"bytes"
	"net/http"
)

// HTTPEstado es el tipo que establece el código de estado de respuesta HTTP.
type HTTPEstado int

// obtenerEntero devuelve el código de estado de respuesta HTTP, con tipo int.
func (t HTTPEstado) obtenerEntero() int {
	return int(t)
}

// Códigos de estados de respuesta HTTP:
// 	HTTPEstadoOk                      = 200
// 	HTTPEstadoOkCreado                = 201
// 	HTTPEstadoOkSinContenido          = 204
// 	HTTPEstadoMalRequerimiento        = 400
// 	HTTPEstadoSinAutorizacion         = 401
// 	HTTPEstadoSinPrivilegios          = 403
// 	HTTPEstadoNoEncontrado            = 404
// 	HTTPEstadoMetodoNoImplementado    = 405
// 	HTTPEstadoRequerimientoMuyGrande  = 413
// 	HTTPEstadoURIMuyGrande            = 414
// 	HTTPEstadoMalFormato              = 415
// 	HTTPEstadoErrorInternoDeServidor  = 500
const (
	HTTPEstadoOk                          HTTPEstado = 200
	HTTPEstadoOkCreado                    HTTPEstado = 201
	HTTPEstadoOkSinContenido              HTTPEstado = 204
	HTTPEstadoErrorMalRequerimiento       HTTPEstado = 400
	HTTPEstadoErrorSinAutorizacion        HTTPEstado = 401
	HTTPEstadoErrorSinPrivilegios         HTTPEstado = 403
	HTTPEstadoErrorNoEncontrado           HTTPEstado = 404
	HTTPEstadoErrorMetodoNoImplementado   HTTPEstado = 405
	HTTPEstadoErrorRequerimientoMuyGrande HTTPEstado = 413
	HTTPEstadoErrorURIMuyGrande           HTTPEstado = 414
	HTTPEstadoErrorMalFormato             HTTPEstado = 415
	HTTPEstadoErrorInternoDeServidor      HTTPEstado = 500
)

// HTTPContenido establece el tipo de contenido dentro del cuerpo de los
// mensajes HTTP.
type HTTPContenido string

// obtenerTexto devuelve el tipo de contenido HTTP, con tipo string.
func (t HTTPContenido) obtenerTexto() string {
	return string(t)
}

// Códigos de tipos de contenido dentro del cuerpo de los mensajes HTTP.
// 	HTTPContenidoSinContenido      = ""
// 	HTTPContenidoApplicationJSON   = "application/json charset=utf-8"
// 	HTTPContenidoApplicationXML    = "application/xml; charset=utf-8"
// 	HTTPContenidoApplicationRTF    = "application/rtf; charset=utf-8"
//	HTTPContenidoApplicationPDF    = "application/pdf"
// 	HTTPContenidoApplicationGZIP   = "applicatio/gzip"
// 	HTTPContenidoApplicationHTTP   = "applicatio/http"
// 	HTTPContenidoApplicationMSWord = "applicatio/msword"
// 	HTTPContenidoTextHTML          = "tex/html; charset=utf-8"
//	HTTPContenidoImagePNG          = "image/png"
// 	HTTPContenidoImageJPEG         = "imag/jpeg"
// 	HTTPContenidoImageGIF          = "imag/gif"
// 	HTTPContenidoTextPlain         = "tex/plain; charset=utf-8"
// 	HTTPContenidoTextCSV           = "text/csv; charset=utf-8"
// 	HTTPContenidoTextXML           = "text/xml; charset=utf-8"
// 	HTTPContenidoTextRTF           = "text/rtf; charset=utf-8"
const (
	HTTPContenidoSinContenido      HTTPContenido = ""
	HTTPContenidoApplicationJSON   HTTPContenido = "application/json charset=utf-8"
	HTTPContenidoApplicationXML    HTTPContenido = "application/xml; charset=utf-8"
	HTTPContenidoApplicationRTF    HTTPContenido = "application/rtf; charset=utf-8"
	HTTPContenidoApplicationPDF    HTTPContenido = "application/pdf"
	HTTPContenidoApplicationGZIP   HTTPContenido = "application/gzip"
	HTTPContenidoApplicationHTTP   HTTPContenido = "application/http"
	HTTPContenidoApplicationMSWord HTTPContenido = "application/msword"
	HTTPContenidoTextHTML          HTTPContenido = "text/html; charset=utf-8"
	HTTPContenidoImagePNG          HTTPContenido = "image/png"
	HTTPContenidoImageJPEG         HTTPContenido = "image/jpeg"
	HTTPContenidoImageGIF          HTTPContenido = "image/gif"
	HTTPContenidoTextPlain         HTTPContenido = "text/plain; charset=utf-8"
	HTTPContenidoTextCSV           HTTPContenido = "text/csv; charset=utf-8"
	HTTPContenidoTextXML           HTTPContenido = "text/xml; charset=utf-8"
	HTTPContenidoTextRTF           HTTPContenido = "text/rtf; charset=utf-8"
)

// HTTPResponder realiza la respuesta HTTP.
func HTTPResponder(w http.ResponseWriter, estadoHTTP HTTPEstado, contenidoHTTP HTTPContenido, cabecera map[string]string, cuerpo string) error {
	// campos de la cabecera del mensaje
	for c, v := range cabecera {
		w.Header().Set(c, v)
	}

	// si el cuerpo es vacío, responder un código de estado 204 (sin contenido)
	// y finalizar el proceso de respuesta.
	if cuerpo == "" || estadoHTTP == HTTPEstadoOkSinContenido {
		w.WriteHeader(HTTPEstadoOkSinContenido.obtenerEntero())
		return nil
	}

	// escribir el tipo de contenido
	w.Header().Set("Content-Type", contenidoHTTP.obtenerTexto())

	// escribir el código de estado y el cuerpo del mensaje
	w.WriteHeader(estadoHTTP.obtenerEntero())
	w.Write([]byte(cuerpo))

	return nil
}

// HTTPObtenerCuerpo devuelve el cuerpo del mensaje recibido como una
// cadena de caracteres (string).
func HTTPObtenerCuerpo(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	return buf.String()
}
