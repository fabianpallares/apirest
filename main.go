package apirest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CORS contiene los nombres de campos de cabecera HTTP para que cada respuesta
// de un endpoint del servidor, pueda cumplir con el "intercambio de recursos
// de orígenes cruzados".
//
// AccessControlAllowCredentials:
//	almacena el valor por defecto del campo de
// 	la cabecera "Access-Control-Allow-Credentials" para todos los recursos.
// 	Indica si la respuesta puede ser expuesta cuando el campo de la
// 	cabecera "credentials" es verdadera.
// 	Este indica si la solicitud puede realizarse usando credenciales.
// 	Las solicitudes GET no contemplan esta cabecera.
//
// AccessControlMaxAge:
// 	almacena el valor por defecto del campo de la cabecera
// 	"Access-Control-Max-Age" para todos los recursos.
// 	Este encabezado indica durante cuánto tiempo los resultados de la
// 	solicitud pueden ser 'cacheados' por el servidor.
// 	El valor se establece en segundos.
var CORS = struct {
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

// campoCabeceraHTTP almacena el nombre y el valor de un campo de una cabecera
// de un endpoint.
type campoCabeceraHTTP struct {
	nombre string
	valor  interface{}
}

// endpoint almacena las cabeceras a exponer y el nombre de función a procesar.
type endpoint struct {
	// QUITAR cors     map[string]interface{} // campos de cabecera CORS del endpoint
	// QUITAR cabecera map[string]interface{} // campos de cabecera del endpoint
	funcion ManejadorFunc // función a procesar
}

// // CampoCabeceraHTTP almacena un nuevo campo ala cabecera HTTP (campos y valores)
// // que el endpoint expone al momento de responder.
// func (o *endpoint) CampoCabeceraHTTP(nombre string, valor interface{}) *endpoint {
// 	o.cabecera[nombre] = valor

// 	return o
// }

// // CabeceraHTTP almacena la cabecera HTTP (campos y valores) que el endpoint
// // expone al momento de responder.
// func (o *endpoint) CabeceraHTTP(cabecera ...campoCabeceraHTTP) *endpoint {
// 	for _, campo := range cabecera {
// 		o.cabecera[campo.nombre] = campo.valor
// 	}

// 	return o
// }

// variableDePatronDeRuta almacena los valores de una parte variable
// (la posición y el nombre) que contenga el patrón de ruta.
type variableDePatronDeRuta struct {
	posicion int
	nombre   string
}

type patronDeRutaDetalle struct {
	// QUITAR metodosPermitos []string                 // métodos HTTP permitidos (agrupados) para un mismo patrón de ruta
	variables []variableDePatronDeRuta // almacena las variables (posición y nombre) de todas las partes variables que posee el patrón de ruta
	endpoints map[string]endpoint      // cada patrón de ruta puede poseer un endpoint distinto por cada método HTTP
}

// enrutador almacena todos los patrones de ruta, junto con todos los endpoints
// de la aplicación.
type enrutador struct {
	patronesDeRutas map[string]patronDeRutaDetalle // mapa de patrones de rutas con su detalle
	// errores []error // almacena los errores de generación de patrones de rutas
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

	detalle, variables, encontrado := o.buscarPatronDeRuta(rutaRecibida)
	if !encontrado {
		http.Error(w, "Recurso inexistente", http.StatusNotFound)
		return
	}

	// Verificar la existencia del método HTTP recibido
	metodoRecibido := r.Method
	// if rt.cors.activo && metodoRecibido == "OPTIONS" {
	// 	w.Header().Set(accessControlAllowOrigin, rt.cors.origenes)
	// 	w.Header().Set(accessControlAllowCredentials, strconv.FormatBool(rt.cors.credenciales))
	// 	w.Header().Set(accessControlMaxAge, strconv.Itoa(rt.cors.duracion))
	// 	w.Header().Set(accessControlAllowMethods, strings.Join(patronPtr.cors.metodos, ", "))
	// 	w.Header().Set(accessControlAllowHeaders, strings.Join(patronPtr.cors.camposRequeridos, ", "))
	// 	w.Header().Set(accessControlExposeHeaders, strings.Join(patronPtr.cors.camposExpuestos, ", "))

	// 	w.WriteHeader(http.StatusNoContent)
	// 	return
	// }

	// Buscar la ruta por el método.
	ep, ok := detalle.endpoints[metodoRecibido]
	if !ok {
		http.Error(w, fmt.Sprintf("La ruta solicitada no implementa el método %v", metodoRecibido), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	if len(variables) > 0 {
		ctx = context.WithValue(ctx, "variables", variables)
	}

	ep.funcion(w, r.WithContext(ctx))
}

// OPTIONS gestiona las operaciones OPTIONS de HTTP (utilizado principalmente para solicitudes CORS).
func (o *enrutador) OPTIONS(ruta string, funcion ManejadorFunc) *endpoint {
	return o.nuevoEndpoint("OPTIONS", ruta, funcion)
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
func (o *enrutador) EscucharHTTPS(puerto, certificadoPublico, certificadoPrivado string) error {
	return o.iniciar("https", puerto, certificadoPublico, certificadoPrivado)
}

// NuevoCampoCabeceraHTTP devuelve un nuevo campo de cabecera HTTP para ser
// asignado a un endpoint
func (o *enrutador) NuevoCampoCabeceraHTTP(nombre string, valor interface{}) campoCabeceraHTTP {
	return campoCabeceraHTTP{nombre, valor}
}

func (o *enrutador) nuevoEndpoint(metodo string, ruta string, funcion ManejadorFunc) *endpoint {
	// convertir la ruta ingresada por el desarrollador a un patrón de ruta
	patron, variables, err := o.rutaAPatronDeRuta(ruta)
	if err != nil {
		o.salir(fmt.Sprintf("La ruta ingresada: [%v] %v, posee un error al intentar generar un patrón de ruta: %v", metodo, ruta, err))
		// o.errores = append(o.errores, fmt.Errorf("La ruta ingresada: [%v] %v, posee un error al intentar generar un patrón de ruta: %w", metodo, ri, err))
	}

	if _, ok := o.patronesDeRutas[patron]; !ok {
		// crear un nuevo detalle del patrón de ruta
		var detalle = patronDeRutaDetalle{
			// QUITAR metodosPermitos: []string{metodo},
			variables: variables,
		}
		// var ep = endpoint{cabecera: make(map[string]interface{}), funcion: funcion} // crear un nuevo endpoint
		var ep = endpoint{funcion: funcion}                 // crear un nuevo endpoint
		detalle.endpoints = map[string]endpoint{metodo: ep} // agregar el endpoint en el detalle del patrón de ruta
		o.patronesDeRutas[patron] = detalle                 // agregar el patrón de ruta en el mapa de patrones de rutas

		return &ep
	}

	// verificar que no se pueda ingresar otro endpoint con el mismo método
	// para este patrón de ruta.
	if _, ok := o.patronesDeRutas[patron].endpoints[metodo]; ok {
		o.salir(fmt.Sprintf("La ruta ingresada: [%v] %v, ya posee un endpoint creado con el mismo método", metodo, ruta))
	}

	var detalle = o.patronesDeRutas[patron]
	// QUITAR  detalle.metodosPermitos = append(detalle.metodosPermitos, metodo) // agregar el método permitido al detalle del patrón de ruta
	o.patronesDeRutas[patron] = detalle // actualizar el detalle del patrón de ruta

	// var ep = endpoint{cabecera: make(map[string]interface{}), funcion: funcion} // crear un nuevo endpoint
	var ep = endpoint{funcion: funcion}              // crear un nuevo endpoint
	o.patronesDeRutas[patron].endpoints[metodo] = ep // asignar el nuevo endpoint

	return &ep
}

// rutaAPatronDeRuta convierte la ruta ingresada por el desarrollador del
// aplicativo a un patrón de ruta.
func (o *enrutador) rutaAPatronDeRuta(s string) (string, []variableDePatronDeRuta, error) {
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

	return "/" + strings.Join(partes, "/"), variables, nil
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

	for patron, detalle := range o.patronesDeRutas {
		partesPatronDeRuta := strings.Split(patron, "/")

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
			var variables = make(map[string]string, len(detalle.variables))
			// si se ha encontrado el patrón de ruta, crear el mapa de variables
			for _, variable := range detalle.variables {
				variables[variable.nombre] = partesPatronDeRuta[variable.posicion]
			}

			return &detalle, variables, true
		}
	}

	return nil, nil, false
}

func (o *enrutador) mostrarPatronesDeRutas() {
	for rp, rpd := range o.patronesDeRutas {
		fmt.Printf("[%v]\n", rp)
		// QUITAR fmt.Printf("\t%v\n", rpd.metodosPermitos)
		fmt.Printf("\t%v\n", rpd.variables)
		for metodo, ep := range rpd.endpoints {
			fmt.Printf("\t\t[%v]\n", metodo)
			// for c, v := range ep.cabecera {
			// 	fmt.Printf("\t\t\t [%v] = %v\n", c, v)
			// }
			fmt.Printf("\t\tFunción: %v\n", ep.funcion)
		}
	}
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

// salir finaliza la ejecución del servidor.
func (o *enrutador) salir(mensaje string) {
	fmt.Println(mensaje)
	os.Exit(2)
}

// CrearEnrutador crea el enrutador para administrar las rutas del aplicativo.
func CrearEnrutador() *enrutador {
	return &enrutador{
		patronesDeRutas: make(map[string]patronDeRutaDetalle),
	}
}
