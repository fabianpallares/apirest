package apirest

// InterceptorFunc es el tipo que deberá implementar la función que pretenda
// ser interceptora de otras funciones.
// Ejemplo de un InterceptorFunc (middleware):
// 	func unInterceptor(antes, despues string) apirest.InterceptorFunc {
// 		return func(manejadorFunc apirest.ManejadorFunc) apirest.ManejadorFunc {
// 			return func(w http.ResponseWriter, r *http.Request) {
// 				fmt.Println("antes")
// 				defer fmt.Println("después")
// 				err := manejadorFunc.ServeHTTP(w, r)
//				...
// 				return err
// 			}
// 		}
// 	}
//
type InterceptorFunc func(ManejadorFunc) ManejadorFunc

// interceptores actúa como una lista, la cual encadena las funciones
// interceptoras (middlewares) para que se procecen en el mismo orden
// en el cual fueron invocadas.
type interceptores struct {
	funciones []InterceptorFunc
}

// Agregar agrega interceptores (middlewares) a la lista.
func (lista *interceptores) Agregar(funciones ...InterceptorFunc) *interceptores {
	lista.funciones = append(lista.funciones, funciones...)
	return lista
}

// Ejecutar encadena y ejecuta todas los interceptores (middlewares).
// Ejemplo:
// 	r.Get("/uno", apirest.CrearInterceptores(f1(), f2(), f3()).Ejecutar(final).
func (lista interceptores) Ejecutar(manejadorFunc ManejadorFunc) ManejadorFunc {
	if manejadorFunc != nil {
		for i := range lista.funciones {
			manejadorFunc = lista.funciones[len(lista.funciones)-1-i](manejadorFunc)
		}
	}

	return manejadorFunc
}

// CrearInterceptores crea una lista de interceptores (middlewares).
// Los interceptores, son procesados invocando una llamada al método Ejecutar()
// de la lista.
func CrearInterceptores(lista ...InterceptorFunc) *interceptores {
	// return interceptores{append(([]InterceptorFunc)(nil), lista...)}
	return &interceptores{lista}
}
