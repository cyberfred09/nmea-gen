package backend

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SerialService struct {
	ctx           context.Context
	port          serial.Port
	portName      string
	baudRate      int
	isEmitting    bool
	mu            sync.Mutex
	stopChan      chan struct{}
}

func NewSerialService() *SerialService {
	return &SerialService{
		stopChan: make(chan struct{}),
	}
}

func (s *SerialService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SerialService) GetPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	return ports, nil
}

func (s *SerialService) StartEmission(portName string, baudRate int, sentences []string, frequency float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isEmitting {
		return fmt.Errorf("already emitting")
	}

	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("failed to open port %s: %w", portName, err)
	}

	s.port = port
	s.portName = portName
	s.baudRate = baudRate
	s.isEmitting = true
	s.stopChan = make(chan struct{})

	go s.emissionLoop(sentences, frequency)
	return nil
}

func (s *SerialService) StopEmission() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isEmitting {
		return nil
	}

	close(s.stopChan)
	if s.port != nil {
		s.port.Close()
	}
	s.isEmitting = false
	return nil
}

func (s *SerialService) emissionLoop(sentences []string, frequency float64) {
	interval := time.Duration(float64(time.Second) / frequency)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Mock data for NMEA
	lat := 40.4168
	lon := -3.7038
	speed := 0.0
	course := 0.0

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Update mock data slightly to simulate movement
			lat += 0.0001
			lon += 0.0001
			speed = 15.5
			course = 45.0

			for _, stype := range sentences {
				var sentence string
				switch stype {
				case "GGA":
					sentence = FormatGPGGA(lat, lon, 1, 8, 1.0, 600.0)
				case "RMC":
					sentence = FormatGPRMC(lat, lon, speed, course)
				case "VTG":
					sentence = FormatGPVTG(course, speed)
				default:
					continue
				}

				_, err := s.port.Write([]byte(sentence + "\r\n"))
				if err != nil {
					runtime.EventsEmit(s.ctx, "log", "Error writing to port: "+err.Error())
					s.StopEmission()
					return
				}
				runtime.EventsEmit(s.ctx, "nmea-sentence", sentence)
			}
		}
	}
}
