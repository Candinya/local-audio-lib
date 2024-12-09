package inits

import (
	"go.uber.org/zap"
)

func Logger(debugMode bool) (l *zap.Logger, err error) {
	if debugMode {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	return l, nil
}
