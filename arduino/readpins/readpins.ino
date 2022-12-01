/* Poll the analog pins, print them out to TTY */

/* IMPORTANT - if using multiple Arduino boards, each must get
 *  a unique (and stable) ControllerID here. The action config
 *  will be tied to this ID. */
int controllerID = 0;

void setup(void)
{
  // initialize digital pin LED_BUILTIN as an output.
  pinMode(LED_BUILTIN, OUTPUT);
  Serial.begin(9600);
}

void loop(void)
{
  int analogVal[6];
  int pin;
  String str;

  Serial.print(str + "Controller:C" + controllerID);

  for (pin=0; pin<6; pin++) {
    analogVal[pin] = analogRead(pin);
    Serial.print(str + " A" + pin + ":" + analogVal[pin]);
  }

  Serial.println("");

  // Flash the LED, while sleeping for 0.5s
  digitalWrite(LED_BUILTIN, HIGH);  // turn the LED on (HIGH is the voltage level)
  delay(100);
  digitalWrite(LED_BUILTIN, LOW);   // turn the LED off by making the voltage LOW
  delay(400);
}
