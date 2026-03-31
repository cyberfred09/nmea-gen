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
	return serial.GetPortsList()
}

func (s *SerialService) StartEmission(portName string, baudRate int, sentences []string, frequency float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isEmitting {
		return fmt.Errorf("already emitting")
	}

	mode := &serial.Mode{BaudRate: baudRate}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}

	s.port = port
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

	// Initial Mock State
	lat, lon := 40.4168, -3.7038
	speed, course := 12.5, 225.0
	depth := 25.0
	windSpeed, windAngle := 15.0, 45.0
	waterTemp := 18.5
	xte := 0.02

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Simulate movement and variations
			lat += 0.0001
			lon += 0.0001
			depth += (float64(time.Now().UnixNano()%10) - 5) * 0.1
			if depth < 2.0 { depth = 2.0 }
			windAngle += (float64(time.Now().UnixNano()%10) - 5) * 0.5
			xte += (float64(time.Now().UnixNano()%10) - 5) * 0.001

			for _, stype := range sentences {
				var toSend []string

				switch stype {
				case "GGA": toSend = append(toSend, FormatGPGGA(lat, lon))
				case "RMC": toSend = append(toSend, FormatGPRMC(lat, lon, speed, course))
				case "GLL": toSend = append(toSend, FormatGPGLL(lat, lon))
				case "GSA": toSend = append(toSend, FormatGPGSA())
				case "GSV": toSend = append(toSend, FormatGPGSV()...)
				case "MWV": toSend = append(toSend, FormatIIMWV(windAngle, windSpeed))
				case "DBT": toSend = append(toSend, FormatIIDBT(depth))
				case "DPT": toSend = append(toSend, FormatIIDPT(depth))
				case "VHW": toSend = append(toSend, FormatIIVHW(course, speed))
				case "HDM": toSend = append(toSend, FormatIIHDM(course))
				case "HDT": toSend = append(toSend, FormatIIHDT(course))
				case "MTW": toSend = append(toSend, FormatIIMTW(waterTemp))
				case "APB": toSend = append(toSend, FormatGPAPB(xte))
				case "BWC": toSend = append(toSend, FormatGPBWC(lat, lon))
				case "BOD": toSend = append(toSend, FormatGPBOD())
				case "XTE": toSend = append(toSend, FormatGPXTE(xte))
				case "AIS": toSend = append(toSend, FormatAIVDM())
				}

				for _, sentence := range toSend {
					s.port.Write([]byte(sentence + "\r\n"))
					runtime.EventsEmit(s.ctx, "nmea-sentence", sentence)
				}
			}
		}
	}
}
