This simple sketch just polls the analog pins, and prints our their
values twice per second. You'll want to use the arduino IDE to upload
it into an actual Arduino board.

Analog pins (A0-A5) are assumed to be hooked up to Force Sensitive
Resistors (FSRs) with a 10Kohm pull-down resistor, so that voltage
goes up with weight. Each FSR pin should read ~0 with no force, and up
to ~1000 at max force.

Note that if a pin has nothing attached to it, it will generate some
random signal; and when there is a big voltage change on one pin, the
other unwired pins will follow suit. There isn't a simple way to tell
when a pin has a real analog signal, vs. when it is unwired.

Any similar analog circuit should work fine here, as long as it
generates a similar range of voltages at the analog pin.

Digital pins (D2-D13) are ignored, but could be used for load cells at
some point.
