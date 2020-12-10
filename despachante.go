package apirest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPEstado establece el tipo de código de estado de respuesta HTTP.
type HTTPEstado int

func (t HTTPEstado) entero() int {
	return int(t)
}

// Códigos de estados de respuesta HTTP.
const (
	HTTPEstado200OK                     HTTPEstado = 200
	HTTPEstado201Creado                 HTTPEstado = 201
	HTTPEstado204SinContenido           HTTPEstado = 204
	HTTPEstado400MalRequerimiento       HTTPEstado = 400
	HTTPEstado401SinAutorizacion        HTTPEstado = 401
	HTTPEstado403SinPrivilegios         HTTPEstado = 403
	HTTPEstado404RecursoInexistente     HTTPEstado = 404
	HTTPEstado405MetodoNoImplementado   HTTPEstado = 405
	HTTPEstado413RequerimientoMuyGrande HTTPEstado = 413
	HTTPEstado414URIMuyGrande           HTTPEstado = 414
	HTTPEstado415MalFormato             HTTPEstado = 415
	HTTPEstado500InternoDeServidor      HTTPEstado = 500
)

// HTTPContenido establece el tipo de contenido dentro del cuerpo de los
// mensajes HTTP.
type HTTPContenido string

func (t HTTPContenido) texto() string {
	return string(t)
}

// Códigos de tipos de contenido dentro del cuerpo de los mensajes HTTP.
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
	if cuerpo == nil || estadoHTTP == HTTPEstado204SinContenido {
		w.WriteHeader(HTTPEstado204SinContenido.entero())
		return nil
	}

	// escribir el tipo de contenido
	w.Header().Set("Content-Type", contenidoHTTP.texto())

	// escribir el código de estado
	w.WriteHeader(estadoHTTP.entero())

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

// HTTPRecibirJSON verifica que se reciba un objeto JSON en el cuerpo del mensaje
// con un tamaño máximo de caracteres. Se intenta convertilo al objeto
// (puntero de estructura) recibido por parámetro.
func HTTPRecibirJSON(r *http.Request, objetoPtr interface{}) error {
	cuerpo, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errDeserealizarLectura
	}

	err = r.Body.Close()
	if err != nil {
		return errDeserealizarLectura
	}

	// realizar la deserealización (casteo) al puntero de objeto recibido
	err = json.Unmarshal(cuerpo, &objetoPtr)
	if err != nil {
		return errDeserealizarConversion
	}

	return nil
}
