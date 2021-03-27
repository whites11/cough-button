const int ledPin = 1;
const int btnPin = 3;

void setup() {
  pinMode(ledPin, OUTPUT);
  pinMode(btnPin, INPUT_PULLUP);

  // turn led off.
  digitalWrite(ledPin, LOW);

  // initialize serial port connection.
  Serial.begin(38400);
}

void loop() {
  int pressed = digitalRead(btnPin);

  if (pressed == LOW) {
    Serial.print("toggle");
    delay(500);
  }

  if (Serial.available() > 0) {
    // read the incoming byte:
    byte incomingByte = Serial.read();

    if (incomingByte == 0x00) {
      // Mic is muted.
      digitalWrite(ledPin, LOW);
    } else {
      // Mic is not muted.
      digitalWrite(ledPin, HIGH);
    }
  }

  delay(100);
}
