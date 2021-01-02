package apirest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	HTTPContenidoApplicationGZIP   HTTPContenido = "applicatio/gzip"
	HTTPContenidoApplicationHTTP   HTTPContenido = "applicatio/http"
	HTTPContenidoApplicationMSWord HTTPContenido = "applicatio/msword"
	HTTPContenidoTextHTML          HTTPContenido = "tex/html; charset=utf-8"
	HTTPContenidoImagePNG          HTTPContenido = "image/png"
	HTTPContenidoImageJPEG         HTTPContenido = "imag/jpeg"
	HTTPContenidoImageGIF          HTTPContenido = "imag/gif"
	HTTPContenidoTextPlain         HTTPContenido = "tex/plain; charset=utf-8"
	HTTPContenidoTextCSV           HTTPContenido = "text/csv; charset=utf-8"
	HTTPContenidoTextXML           HTTPContenido = "text/xml; charset=utf-8"
	HTTPContenidoTextRTF           HTTPContenido = "text/rtf; charset=utf-8"
)

var (
	errDeserealizarLectura    = errors.New("No es posible deserealizar el objeto recibido")
	errDeserealizarConversion = errors.New("No es posible convertir (deserealizar) los datos al objeto recibido")
	errEncodearJSON           = errors.New("No es posible realizar la respuesta con el contenido JSON")
)

// HTTPResponder realiza la respuesta HTTP.
func HTTPResponder(w http.ResponseWriter, estadoHTTP HTTPEstado, contenidoHTTP HTTPContenido, cabecera map[string]string, cuerpo interface{}) error {
	// campos de la cabecera del mensaje
	for c, v := range cabecera {
		w.Header().Set(c, v)
	}

	// si el cuerpo es vacío, responder un código de estado 204 (sin contenido)
	// y finalizar el proceso de respuesta.
	if cuerpo == nil || estadoHTTP == HTTPEstadoOkSinContenido {
		w.WriteHeader(HTTPEstadoOkSinContenido.obtenerEntero())
		return nil
	}

	// escribir el tipo de contenido
	w.Header().Set("Content-Type", contenidoHTTP.obtenerTexto())

	// escribir el código de estado
	w.WriteHeader(estadoHTTP.obtenerEntero())

	switch contenidoHTTP {
	case HTTPContenidoApplicationJSON:
		if err := json.NewEncoder(w).Encode(cuerpo); err != nil {
			return errEncodearJSON
		}
	default:
		fmt.Fprintln(w, cuerpo)
	}

	return nil
}

// ObjetoAJSON intenta convertir el objeto recibido en un array de bytes con
// un formato JSON.
func ObjetoAJSON(objetoPtr interface{}) ([]byte, error) {
	return json.Marshal(objetoPtr)
}

// HTTPRecibirJSON verifica que se reciba un objeto JSON en el cuerpo del
// mensaje. Se intenta convertilo al objeto (puntero de estructura)
// recibido por parámetro.
func HTTPRecibirJSON(r *http.Request, objetoPtr interface{}, validarCamposDesconocidos bool) error {
	return JSONAObjeto(r.Body, objetoPtr, validarCamposDesconocidos)
}

// JSONAObjeto convierte el array de bytes recibido (que debería ser un formato
// JSON válido) e intenta completar los campos de la estructura recibida.
func JSONAObjeto(r io.Reader, objetoPtr interface{}, validarCamposDesconocidos bool) error {
	// r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	var descifrador = json.NewDecoder(r)
	if validarCamposDesconocidos {
		// no permite recibir campos desconocidos
		descifrador.DisallowUnknownFields()
	}

	var err = descifrador.Decode(objetoPtr)
	if err == nil {
		return nil
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		// verificar que se haya recibido información en el cuerpo del mensaje
		return fmt.Errorf("El formato JSON recibido es incorrecto. Contenido vacío")

	case errors.Is(err, io.ErrUnexpectedEOF):
		// verificar la lectura del cuerpo del mensaje
		return fmt.Errorf("El formato JSON recibido es incorrecto. Se ha llegado al final de la lectura de manera inesperada")

	case errors.As(err, &syntaxError):
		// verificar si el formato es correcto, si faltan dobles comillas,
		// comillas, comas, llaves, corchetes; etc.
		return fmt.Errorf("El formato JSON recibido es incorrecto. Error en la posición: %v", syntaxError.Offset)

	case errors.As(err, &unmarshalTypeError):
		// verificar si hay un error de tipo de campo, campos que contienen
		// tipos de valores erroneos
		var valor string
		switch unmarshalTypeError.Value {
		case "number":
			valor = "numérico"
		case "string":
			valor = "texto"
		case "bool":
			valor = "lógico"
		default:
			valor = unmarshalTypeError.Value
		}
		return fmt.Errorf("El formato JSON recibido es incorrecto. Error en el campo: \"%v\", tipo de valor recibido: %v, posición: %v", unmarshalTypeError.Field, valor, unmarshalTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		// verificar si se recibieron campos adicionales que no están en la
		// estructura recibida
		var campo = strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("El formato JSON recibido es incorrecto. Se ha recibido un nombre de campo inexistente: %v", campo)

	case err.Error() == "http: request body too large":
		// verificar contenido muy largo
		return fmt.Errorf("El formato JSON recibido es incorrecto. El texto recibido es demasiado grande")
	}

	// cualquier otro tipo de error
	return fmt.Errorf("El formato JSON recibido es incorrecto: %w", err)
}
