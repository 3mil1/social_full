package errHandler

import (
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {

	var appErr *AppError
	if errors.As(err, &appErr) {
		w.WriteHeader(appErr.StatusCode)
		w.Write(appErr.Marshal())
		return
	}

	//logger.ErrorLogger.Println("Unhandled errHandler occurred: ", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(SystemError(err).Marshal())
}
