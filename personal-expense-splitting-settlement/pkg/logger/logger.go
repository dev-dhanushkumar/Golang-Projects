package logger

import "go.uber.org/zap"

func InitLogger() (*zap.SugaredLogger, error) {
	// Use zap.NewDevelopment() for human-readable console logs during development
	// or zap.NewProduction() for JSON output in production environments.
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	// Assign the SugaredLogger to the global variable
	sugar := logger.Sugar()
	return sugar, nil
}
