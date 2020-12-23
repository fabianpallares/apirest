package apirest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// cors almacena los nombres de los campos de la cabecera CORS.
var cors = struct {
	AccessControlAllowOrigin      string
	AccessControlAllowCredentials string
	AccessControlMaxAge           string
	AccessControlAllowMethods     string
	AccessControlAllowHeaders     string
	AccessControlExposeHeaders    string
}{"Access-Control-Allow-Origin",
	"Access-Control-Allow-Credentials",
	"Access-Control-Max-Age",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Expose-Headers",
}

// ManejadorFunc es el tipo (función) que procesa el requirimiento del recurso.
type ManejadorFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// Implementa el método ServeHTTP.
func (m ManejadorFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return m(w, r)
}

// patronDeRuta almacena un patrón de ruta único para toda la aplicación.
// 	ejemplos:
//	"/personas"
//	"/personas/{v}"
//	"/personas/{v}/datos_principales"
type patronDeRuta string

// string convierte a string el tipo patrón de ruta.
func (pr patronDeRuta) string() string {
	return string(pr)
}

// endpoint almacena un apuntador al detalle del patrón de ruta y la función
// (ManejadorFunc) a procesar.
type endpoint struct {
	detalle *patronDeRutaDetalle // apuntador al detalle del patrón de ruta al cuál pertenece el endpoint
	funcion ManejadorFunc        // función (ManejadorFunc) a procesar
}

// CORSCamposRequeridos solicita los campos CORS requeridos para poder procesar
// el endpoint solicitado.
// es el campo de cabecera: "Access-Control-Allow-Headers"
func (o *endpoint) CORSCamposRequeridos(campos ...string) *endpoint {
	for _, campo := range campos {
		var esExistente bool
		for _, campoExistente := range o.detalle.cors.camposRequeridos {
			if strings.Trim(strings.ToLower(campo), " ") == strings.Trim(strings.ToLower(campoExistente), " ") {
				esExistente = true
				break
			}
		}
		if !esExistente {
			o.detalle.cors.camposRequeridos = append(o.detalle.cors.camposRequeridos, campo)
		}
	}

	return o
}

// CORSCamposExpuestos establece los campos CORS que el endpoint expondrá para
// que el cliente recupere información. Este encabezado expone una lista blanca
// de campos que tienen permitido acceder los exploradores.
// es el campo de cabecera: "Access-Control-Expose-Headers"
func (o *endpoint) CORSCamposExpuestos(campos ...string) *endpoint {
	for _, campo := range campos {
		var esExistente bool
		for _, campoExistente := range o.detalle.cors.camposExpuestos {
			if strings.Trim(strings.ToLower(campo), " ") == strings.Trim(strings.ToLower(campoExistente), " ") {
				esExistente = true
				break
			}
		}
		if !esExistente {
			o.detalle.cors.camposExpuestos = append(o.detalle.cors.camposExpuestos, campo)
		}
	}

	return o
}

// variableDePatronDeRuta almacena los valores de una parte variable
// (la posición y el nombre) que contenga el patrón de ruta.
type variableDePatronDeRuta struct {
	posicion int
	nombre   string
}

// patronDeRutaDetalle almacena por cada patrón de ruta, todos los endpoint que
// comparten el mismo patrón de ruta, las partes variables del patrón de ruta y
// los valores CORS agrupados para que puedan ser expuestos a través de OPTIONS.
type patronDeRutaDetalle struct {
	cors struct {
		// almacena todos los métodos para una mismo patrón de ruta ("GET, POST, ...")
		// es el campo de cabecera: Access-Control-Allow-Methods
		metodosPermitidos []string

		// almacena todos los campos requeridos para un mismo patrón de ruta
		// es el campo de cabecera: Access-Control-Allow-Headers
		camposRequeridos []string

		// almacena todos los campos expuestos para un mismo patrón de ruta
		// es el campo de cabecera: Access-Control-Expose-Headers
		camposExpuestos []string
	}

	endpoints map[string]*endpoint     // cada patrón de ruta puede poseer un endpoint distinto por cada método HTTP
	variables []variableDePatronDeRuta // almacena las variables (posición y nombre) de todas las partes variables que posee el patrón de ruta
}

// enrutador almacena los valores de los campos generales de CORS y todos los
// patrones de ruta de la aplicación.
type enrutador struct {
	cors struct {
		// esActivo determina que todos los endpoints utilizan CORS.
		esActivo bool

		// origenes almacena el valor por defecto del campo de la cabecera
		// "Access-Control-Allow-Origin" para todos los endpoints.
		// Este campo de la cabecera, especifica una o varias URIs que
		// pueden tener acceso al endpoint.
		// El explorador asegura esto. Para solicitudes sin credenciales, el
		// servidor debe especificar "*" como un comodín, de este modo se
		// permite a cualquier origen acceder al recurso.
		origenes []string

		// credenciales almacena el valor por defecto del campo de la cabecera
		// "Access-Control-Allow-Credentials" para todos los endpoints.
		// Indica si la respuesta puede ser expuesta cuando el campo de la
		// cabecera "credentials" es verdadera.
		// Este indica si la solicitud puede realizarse usando credenciales.
		// Las solicitudes GET no contemplan esta cabecera.
		credenciales bool

		// duracion almacena el valor por defecto del campo de la cabecera
		// "Access-Control-Max-Age" para todos los endpoints.
		// Este encabezado indica durante cuánto tiempo los resultados de la
		// solicitud pueden ser 'cacheados' por el servidor.
		// El valor se establece en segundos.
		duracion int
	}

	// mapa de patrones de rutas con su detalle
	patronesDeRutas map[patronDeRuta]*patronDeRutaDetalle
}

// ServeHTTP envía la solicitud a la función cuyo patrón de ruta coincida
// con la URL de la solicitud.
func (o *enrutador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verificar la existencia de la ruta recibida
	rutaRecibida := r.RequestURI
	pos := strings.Index(rutaRecibida, "?")
	if pos > 0 {
		rutaRecibida = rutaRecibida[:pos]
	}

	detallePtr, variables, encontrado := o.buscarPatronDeRuta(rutaRecibida)
	if !encontrado {
		http.Error(w, "Recurso inexistente", http.StatusNotFound)
		return
	}

	var cabecerasCORS = make(map[string]string)
	metodoRecibido := r.Method
	if metodoRecibido == "OPTIONS" {
		if !o.cors.esActivo {
			http.Error(w, "La aplicación no implementa el método OPTIONS (No se encuentra activa la opcion CORS)", http.StatusNotFound)
			return
		}
		w.Header().Set(cors.AccessControlAllowOrigin, strings.Join(o.cors.origenes, ", "))
		w.Header().Set(cors.AccessControlAllowCredentials, strconv.FormatBool(o.cors.credenciales))
		w.Header().Set(cors.AccessControlMaxAge, strconv.Itoa(o.cors.duracion))
		w.Header().Set(cors.AccessControlAllowMethods, strings.Join(detallePtr.cors.metodosPermitidos, ", "))
		w.Header().Set(cors.AccessControlAllowHeaders, strings.Join(detallePtr.cors.camposRequeridos, ", "))
		w.Header().Set(cors.AccessControlExposeHeaders, strings.Join(detallePtr.cors.camposExpuestos, ", "))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// si no es options... verificar la existencia del método HTTP recibido
	ep, ok := detallePtr.endpoints[metodoRecibido]
	if !ok {
		http.Error(w, fmt.Sprintf("La ruta solicitada no implementa el método %v", metodoRecibido), http.StatusNotFound)
		return
	}

	// subir al contexto las cabeceras CORS y las variables de los patrones de ruta
	ctx := r.Context()

	if o.cors.esActivo {
		cabecerasCORS[cors.AccessControlAllowOrigin] = strings.Join(o.cors.origenes, ", ")
		cabecerasCORS[cors.AccessControlAllowCredentials] = strconv.FormatBool(o.cors.credenciales)
		cabecerasCORS[cors.AccessControlMaxAge] = strconv.Itoa(o.cors.duracion)
		cabecerasCORS[cors.AccessControlAllowMethods] = strings.Join(detallePtr.cors.metodosPermitidos, ", ")
		cabecerasCORS[cors.AccessControlAllowHeaders] = strings.Join(detallePtr.cors.camposRequeridos, ", ")
		cabecerasCORS[cors.AccessControlExposeHeaders] = strings.Join(detallePtr.cors.camposExpuestos, ", ")
		ctx = context.WithValue(ctx, "cors", cabecerasCORS)
	}
	if len(variables) > 0 {
		ctx = context.WithValue(ctx, "variables", variables)
	}

	ep.funcion(w, r.WithContext(ctx))
}

// CORSActivar determina que todos los recursos de la aplicación utilizarán CORS.
func (o *enrutador) CORSActivar() *enrutador {
	o.cors.esActivo = true
	return o
}

// CORSOrigenes cambia el valor de los orígenes permitidos (URIs que pueden
// tener acceso a los endpoints).
// el el campo de cabecera: "Access-Control-Allow-Origin"
// tiene como valor por defecto: "*".
func (o *enrutador) CORSOrigenes(origenes ...string) *enrutador {
	o.cors.origenes = origenes
	return o
}

// COSCredenciales cambia el valor del requerimiento de credenciales para
// consumir los recursos.
// es el campo de cabecera: "Access-Control-Allow-Credentials"
// tiene como valor por defecto: false.
func (o *enrutador) COSCredenciales(credenciales bool) *enrutador {
	o.cors.credenciales = credenciales
	return o
}

// CORSDuracion cambia el valor de la duración del tiempo que el servidor
// 'cachea' los resultados de todos los recursos.
// es el campo de cabecera: "Access-Control-Max-Age"
// tiene como valor por defecto: -1 (sin caché).
func (o *enrutador) CORSDuracion(duracion int) *enrutador {
	o.cors.duracion = duracion
	return o
}

// GET gestiona las operaciones GET de HTTP (consulta de recursos).
func (o *enrutador) GET(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("GET", ruta, funcion)
}

// POST gestiona las operaciones POST de HTTP (nuevos recursos).
func (o *enrutador) POST(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("POST", ruta, funcion)
}

// PUT gestiona las operaciones PUT de HTTP (modificacion completa de recursos).
func (o *enrutador) PUT(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("PUT", ruta, funcion)
}

// PATCH gestiona las operaciones PATCH de HTTP (modificacion parcial de recursos).
func (o *enrutador) PATCH(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("PATCH", ruta, funcion)
}

// DELETE gestiona las operaciones DELETE de HTTP (eliminación de recursos).
func (o *enrutador) DELETE(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("DELETE", ruta, funcion)
}

// IniciarPorHTTP inicia el servidor escuchando por HTTP.
func (o *enrutador) IniciarPorHTTP(puerto string) error {
	return o.iniciar("http", puerto, "", "")
}

// IniciarPorHTTPS inicia el servidor escuchando por HTTPS.
func (o *enrutador) IniciarPorHTTPS(puerto, certificadoPublico, certificadoPrivado string) error {
	return o.iniciar("https", puerto, certificadoPublico, certificadoPrivado)
}

func (o *enrutador) nuevoEndpoint(metodo string, ruta string, funcion ManejadorFunc) *endpoint {
	// convertir la ruta ingresada por el desarrollador a un patrón de ruta
	pr, variables, err := o.rutaAPatronDeRuta(ruta)
	if err != nil {
		finalizar(fmt.Sprintf("La ruta ingresada: [%v] %v, posee un error al intentar generar un patrón de ruta: %v", metodo, ruta, err))
	}

	detallePtr, ok := o.patronesDeRutas[pr]
	if !ok {
		// crear un nuevo detalle del patrón de ruta
		var detallePtr = &patronDeRutaDetalle{
			variables: variables,
		}
		detallePtr.cors.metodosPermitidos = []string{metodo}

		var epPtr = &endpoint{detalle: detallePtr, funcion: funcion} // crear un nuevo endpoint
		detallePtr.endpoints = map[string]*endpoint{metodo: epPtr}   // agregar el endpoint en el detalle del patrón de ruta
		o.patronesDeRutas[pr] = detallePtr                           // agregar el patrón de ruta en el mapa de patrones de rutas

		return epPtr
	}

	// verificar que las variables de un mismo patrón de ruta, lleven los
	// mismos nombres
	for i := 0; i < len(detallePtr.variables); i++ {
		if detallePtr.variables[i].nombre != variables[i].nombre {
			finalizar(fmt.Sprintf("Existe un patrón de ruta: %v, que contiene endpoints con distintos nombres de variables", pr))
		}
	}

	// verificar que no se pueda ingresar otro endpoint con el mismo método
	// para este patrón de ruta.
	if _, ok := o.patronesDeRutas[pr].endpoints[metodo]; ok {
		finalizar(fmt.Sprintf("La ruta ingresada: [%v] %v, ya posee un endpoint creado con el mismo método", metodo, ruta))
	}

	detallePtr.cors.metodosPermitidos = append(detallePtr.cors.metodosPermitidos, metodo) // agregar el método permitido al detalle del patrón de ruta
	var epPtr = &endpoint{detalle: detallePtr, funcion: funcion}                          // crear un nuevo endpoint
	detallePtr.endpoints[metodo] = epPtr                                                  // asignar el nuevo endpoint

	return epPtr
}

// rutaAPatronDeRuta convierte la ruta ingresada por el desarrollador de la
// aplicación a un patrón de ruta.
func (o *enrutador) rutaAPatronDeRuta(s string) (patronDeRuta, []variableDePatronDeRuta, error) {
	if s == "" {
		return "", nil, fmt.Errorf("la ruta recibida está vacía")
	}

	var partesRuta = strings.Split(s, "/")
	if partesRuta[0] == "" {
		partesRuta = partesRuta[1:]
	}
	if partesRuta[len(partesRuta)-1] == "" {
		partesRuta = partesRuta[:len(partesRuta)-1]
	}

	var partes []string
	var variables []variableDePatronDeRuta

	for pos, parte := range partesRuta {
		var p = strings.ToLower(strings.Trim(parte, " "))
		if strings.Index(p, "{") == -1 {
			partes = append(partes, p)
			continue
		}
		if len(p) == 2 {
			return "", nil, fmt.Errorf("existe al menos una parte de la ruta que no contiene un nombre de variable")
		}
		if strings.Index(p, "}") == -1 {
			return "", nil, fmt.Errorf("el nombre de variable no contiene la llave de cierre")
		}

		variables = append(variables, variableDePatronDeRuta{posicion: pos, nombre: p[1 : len(p)-1]})
		partes = append(partes, "{v}")
	}

	return patronDeRuta("/" + strings.Join(partes, "/")), variables, nil
}

// buscarPatronDeRuta busca que exista el patrón de ruta, según la ruta recibida.
func (o *enrutador) buscarPatronDeRuta(rutaRecibida string) (*patronDeRutaDetalle, map[string]string, bool) {
	var partesRutaRecibida = strings.Split(rutaRecibida, "/")
	if partesRutaRecibida[0] == "" {
		partesRutaRecibida = partesRutaRecibida[1:]
	}
	if partesRutaRecibida[len(partesRutaRecibida)-1] == "" {
		partesRutaRecibida = partesRutaRecibida[:len(partesRutaRecibida)-1]
	}

	for pr, detallePtr := range o.patronesDeRutas {
		partesPatronDeRuta := strings.Split(pr.string(), "/")[1:]

		if len(partesRutaRecibida) != len(partesPatronDeRuta) {
			continue
		}

		// Verificar cada parte de la ruta recibida con la parte del patrón de ruta actual
		encontrado := true
		for i, partePatronActual := range partesPatronDeRuta {
			switch {
			case partePatronActual == partesRutaRecibida[i]:
				continue
			case partePatronActual == "{v}" && partePatronActual != partesRutaRecibida[i]:
				continue
			default:
				encontrado = false
				break
			}
		}
		if encontrado {
			var variables = make(map[string]string, len(detallePtr.variables))
			// si se ha encontrado el patrón de ruta, crear el mapa de variables
			for _, variable := range detallePtr.variables {
				variables[variable.nombre] = partesRutaRecibida[variable.posicion]
			}

			return detallePtr, variables, true
		}
	}

	return nil, nil, false
}

// iniciar inicia el servidor escuchando por el protocolo y puerto establecido.
func (o *enrutador) iniciar(protocolo, puerto, certificadoPublico, certificadoPrivado string) error {
	if puerto != "" && string(puerto[0]) != ":" {
		puerto = ":" + puerto
	}

	if strings.ToUpper(strings.Trim(protocolo, " ")) == "HTTP" {
		return http.ListenAndServe(puerto, o)
	}

	return http.ListenAndServeTLS(puerto, certificadoPublico, certificadoPrivado, o)
}

// finalizar finaliza la ejecución del servidor.
func finalizar(mensaje string) {
	fmt.Println(mensaje)
	os.Exit(2)
}

// Finalizar finaliza la ejecución del servidor.
func Finalizar(mensaje string) {
	finalizar(mensaje)
}

// CrearEnrutador crea el enrutador para administrar las rutas de la aplicación.
func CrearEnrutador() *enrutador {
	var r = &enrutador{
		patronesDeRutas: make(map[patronDeRuta]*patronDeRutaDetalle),
	}

	r.cors.origenes, r.cors.credenciales, r.cors.duracion = []string{"*"}, false, -1
	return r
}

// ObtenerVariablesDeRuta retorna un mapa con los nombres de variables del patrón
// de ruta junto con los valores recibos de la solicitud del cliente.
func ObtenerVariablesDeRuta(r *http.Request) map[string]string {
	m, ok := r.Context().Value("variables").(map[string]string)
	if !ok {
		return map[string]string{}
	}

	return m
}

// ObtenerCORS retorna un mapa con los campos de la cabecera CORS.
func ObtenerCORS(r *http.Request) map[string]string {
	m, ok := r.Context().Value("cors").(map[string]string)
	if !ok {
		return map[string]string{}
	}

	return m
}
