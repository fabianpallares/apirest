package apirest

// ErrorNuevoMalRequerimiento crea un nuevo error de tipo:
// 400 (Mal requerimiento).
func ErrorNuevoMalRequerimiento(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado400MalRequerimiento, nil, formato, args...)
}

// ErrorNuevoMalRequerimientoConOrigen crea un nuevo error con origen de tipo:
// 400 (Mal requerimiento).
func ErrorNuevoMalRequerimientoConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado400MalRequerimiento, origen, formato, args...)
}

// ErrorEsErrorMalRequerimiento devuelve el resultado de conocer si el error
// es del tipo 400 (Mal requerimiento).
func ErrorEsErrorMalRequerimiento(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado400MalRequerimiento {
		return true
	}

	return false
}

// ErrorNuevoSinAutorizacion crea un nuevo error de tipo:
// 401 (Sin autorización).
func ErrorNuevoSinAutorizacion(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado401SinAutorizacion, nil, formato, args...)
}

// ErrorNuevoSinAutorizacionConOrigen crea un nuevo error con origen de tipo:
// 401 (Sin autorización).
func ErrorNuevoSinAutorizacionConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado401SinAutorizacion, origen, formato, args...)
}

// ErrorEsErrorSinAutorizacion devuelve el resultado de conocer si el error
// es del tipo 401 (Sin autorización).
func ErrorEsErrorSinAutorizacion(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado401SinAutorizacion {
		return true
	}

	return false
}

// ErrorNuevoSinPrivilegios crea un nuevo error de tipo:
// 403 (Sin privilegios).
func ErrorNuevoSinPrivilegios(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado403SinPrivilegios, nil, formato, args...)
}

// ErrorNuevoSinPrivilegiosConOrigen crea un nuevo error con origen de tipo:
// 403 (Sin privilegios).
func ErrorNuevoSinPrivilegiosConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado403SinPrivilegios, origen, formato, args...)
}

// ErrorEsErrorSinPrivilegios devuelve el resultado de conocer si el error
// es del tipo 403 (Sin privilegios).
func ErrorEsErrorSinPrivilegios(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado403SinPrivilegios {
		return true
	}

	return false
}

// ErrorNuevoRecursoInexistente crea un nuevo error de tipo:
// 404 (Recurso inexistente).
func ErrorNuevoRecursoInexistente(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado404RecursoInexistente, nil, formato, args...)
}

// ErrorNuevoRecursoInexistenteConOrigen crea un nuevo error con origen de tipo:
// 404 (Recurso inexistente).
func ErrorNuevoRecursoInexistenteConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado404RecursoInexistente, origen, formato, args...)
}

// ErrorEsErrorRecursoInexistente devuelve el resultado de conocer si el error
// es del tipo 404 (Recurso inexistente).
func ErrorEsErrorRecursoInexistente(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado404RecursoInexistente {
		return true
	}

	return false
}

// ErrorNuevoMetodoNoImplementado crea un nuevo error de tipo:
// 405 (Método no implementado).
func ErrorNuevoMetodoNoImplementado(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado405MetodoNoImplementado, nil, formato, args...)
}

// ErrorNuevoMetodoNoImplementadoConOrigen crea un nuevo error con origen de tipo:
// 405 (Método no implementado).
func ErrorNuevoMetodoNoImplementadoConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado405MetodoNoImplementado, origen, formato, args...)
}

// ErrorEsErrorMetodoNoImplementado devuelve el resultado de conocer si el error
// es del tipo 405 (Método no implementado).
func ErrorEsErrorMetodoNoImplementado(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado405MetodoNoImplementado {
		return true
	}

	return false
}

// ErrorNuevoRequerimientoMuyGrande crea un nuevo error de tipo:
// 413 (Requerimiento muy grande).
func ErrorNuevoRequerimientoMuyGrande(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado413RequerimientoMuyGrande, nil, formato, args...)
}

// ErrorNuevoRequerimientoMuyGrandeConOrigen crea un nuevo error con origen de tipo:
// 413 (Requerimiento muy grande).
func ErrorNuevoRequerimientoMuyGrandeConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado413RequerimientoMuyGrande, origen, formato, args...)
}

// ErrorEsErrorRequerimientoMuyGrande devuelve el resultado de conocer si el error
// es del tipo 413 (Requerimiento muy grande).
func ErrorEsErrorRequerimientoMuyGrande(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado413RequerimientoMuyGrande {
		return true
	}

	return false
}

// ErrorNuevoURIMuyGrande crea un nuevo error de tipo:
// 414 (URI muy grande).
func ErrorNuevoURIMuyGrande(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado414URIMuyGrande, nil, formato, args...)
}

// ErrorNuevoURIMuyGrandeConOrigen crea un nuevo error con origen de tipo:
// 414 (URI muy grande).
func ErrorNuevoURIMuyGrandeConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado414URIMuyGrande, origen, formato, args...)
}

// ErrorEsErrorURIMuyGrande devuelve el resultado de conocer si el error
// es del tipo 414 (URI muy grande).
func ErrorEsErrorURIMuyGrande(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado414URIMuyGrande {
		return true
	}

	return false
}

// ErrorNuevoMalFormato crea un nuevo error de tipo:
// 415 (Mal formato).
func ErrorNuevoMalFormato(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado415MalFormato, nil, formato, args...)
}

// ErrorNuevoMalFormatoConOrigen crea un nuevo error con origen de tipo:
// 415 (Mal formato).
func ErrorNuevoMalFormatoConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado415MalFormato, origen, formato, args...)
}

// ErrorEsErrorMalFormato devuelve el resultado de conocer si el error
// es del tipo 415 (Mal formato).
func ErrorEsErrorMalFormato(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado415MalFormato {
		return true
	}

	return false
}

// ErrorNuevoInternoDeServidor crea un nuevo error de tipo:
// 500 (Error interno de servidor).
func ErrorNuevoInternoDeServidor(formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado500InternoDeServidor, nil, formato, args...)
}

// ErrorNuevoInternoDeServidorConOrigen crea un nuevo error con origen de tipo:
// 500 (Error interno de servidor).
func ErrorNuevoInternoDeServidorConOrigen(origen error, formato string, args ...interface{}) error {
	return errorNuevo(HTTPEstado500InternoDeServidor, origen, formato, args...)
}

// ErrorEsErrorInternoDeServidor devuelve el resultado de conocer si el error
// es del tipo 500 (Error interno de servidor).
func ErrorEsErrorInternoDeServidor(err error) bool {
	errEnvoltorio, ok := err.(*errorEnvoltorio)
	if ok && errEnvoltorio.tipo == HTTPEstado500InternoDeServidor {
		return true
	}

	return false
}
