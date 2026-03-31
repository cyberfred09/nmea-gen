package main

import (
	"context"
	"nmea-gen/backend"
)

// App struct
type App struct {
	ctx           context.Context
	serialService *backend.SerialService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		serialService: backend.NewSerialService(),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.serialService.Startup(ctx)
}

// GetPorts returns a list of available serial ports
func (a *App) GetPorts() ([]string, error) {
	return a.serialService.GetPorts()
}

// StartEmission starts the NMEA emission on the specified port
func (a *App) StartEmission(portName string, baudRate int, sentences []string, frequency float64) error {
	return a.serialService.StartEmission(portName, baudRate, sentences, frequency)
}

// StopEmission stops the NMEA emission
func (a *App) StopEmission() error {
	return a.serialService.StopEmission()
}
