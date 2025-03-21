package utils

import "go.uber.org/zap"

func logger(msg string) {

	sugar, err := zap.NewProduction()
	sugartwo := sugar.Sugar()
	if err != nil {

		panic("error in the error Logger")

	}

	sugartwo.Info("The error occured in the File:->", msg)
}
