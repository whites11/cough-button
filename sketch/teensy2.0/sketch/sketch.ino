const int ledPin = 1;
const int btnPin = 3;

void setup() {
  pinMode(ledPin, OUTPUT);
  pinMode(btnPin, INPUT_PULLUP);

  // turn led off.
  digitalWrite(ledPin, LOW);

  // initialize serial port connection.
  Serial.begin(57600);
}

void loop() {
  if (!Serial.dtr()) {
    digitalWrite(ledPin, HIGH);
    delay(150);
    digitalWrite(ledPin, LOW);
    delay(150);
    return;
  }

  if (Serial.available()) {
    String s = Serial.readStringUntil('\n');
  
    if (s == "ping") {
        Serial.println("pong");
    } else if (s == "muted") {
        digitalWrite(ledPin, LOW);
    } else if (s == "unmuted") {
        digitalWrite(ledPin, HIGH);
    } else if (s != "") {
        Serial.println("Unknown command: " + s);
    }
  }
    
  int pressed = digitalRead(btnPin);

  if (pressed == LOW) {
    Serial.println("toggle");
    delay(100);
  }

  delay(100);
}
