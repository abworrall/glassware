package controller

import(
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.bug.st/serial"  // https://pkg.go.dev/go.bug.st/serial

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/event"
	"github.com/abworrall/glassware/pkg/sensor"
)

// A SerialController assumes a serial TTY, that has line-based output
// like this:
//
//  Controller:C0 A0:681 A1:587 A2:509 A3:440 A4:369 A5:326
//  Controller:C0 A0:682 A1:587 A2:508 A3:438 A4:367 A5:323
//  Controller:C0 A0:678 A1:584 A2:505 A3:436 A4:366 A5:324
//
// A new line should show up every ~second, with space-separated
// key-val pairs, where keys are unique and vals are integers. Each
// key-val pair is a sensor that the controller is reporting on. Note
// that sensors might be reporting noise; the controller has no way to
// know.
//
// The arduino/readpins.ino sketch provides this output, with the
// six anolog pins reporting values in the range [0,1000].
type SerialController struct {
	Config config.Config
	
	PortName string                     // e.g. "/dev/ttyUSB0"
	ControllerID string                 // e.g. "C0" - must be unique & stable over reboots
	Sensors map[string]sensor.Sensor
}

func NewSerialController(portName string) Controller {
	return &SerialController{
		PortName: portName,
		Sensors: map[string]sensor.Sensor{},
	}
}

func (sc *SerialController)String() string {
	return fmt.Sprintf("SerialController[%s[%s] %d sensors]", sc.PortName, sc.ControllerID, len(sc.Sensors))
}

func (sc *SerialController)Start(c config.Config, eventsOut chan<- event.Event) {
	sc.Config = c
	log.Printf("(controller at %s starting)", sc.PortName)

	port, err := serial.Open(sc.PortName, &serial.Mode{BaudRate: 9600})
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		sc.ProcessReadings(scanner.Text(), eventsOut)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (sc *SerialController)ProcessReadings(str string, eventsOut chan<- event.Event) {
	if sc.Config.Verbosity > 0 {
		log.Printf("SerialController(%s), read(%s)= %q", sc.ControllerID, sc.PortName, str)
	}

	for _, keyval := range strings.Fields(strings.TrimSpace(str)) {
		bits := strings.Split(keyval, ":")

		if bits[0] == "Controller" {
			if sc.ControllerID == "" {
				sc.ControllerID = bits[1]
			}
			continue
		}

		// abw FIXME hack
		if bits[0] != "A0" { continue }

		// Autocreate DropRestoreSensors for all keys
		if _,exists := sc.Sensors[bits[0]]; !exists {
			sc.Sensors[bits[0]] = sensor.NewDropRestoreSensor(sc.ControllerID + "/" + bits[0])
		}

		if i, err := strconv.Atoi(bits[1]); err == nil {
			sc.Sensors[bits[0]].ProcessNewReading(i, eventsOut)
		}
	}
}

func ListSerialControllers() []string {
	ret := []string{}

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		ret = append(ret, port)
	}

	return ret
}
