#include <CapacitiveSensor.h>

/*
 * CapitiveSense Library Demo Sketch
 * Paul Badger 2008
 * Uses a high value resistor e.g. 10 megohm between send pin and receive pin
 * Resistor effects sensitivity, experiment with values, 50 kilohm - 50 megohm. Larger resistor values yield larger sensor values.
 * Receive pin is the sensor pin - try different amounts of foil/metal on this pin
 * Best results are obtained if sensor foil and wire is covered with an insulator such as paper or plastic sheet
 */

/* we always wait a bit between updates of the display */
unsigned long delaytime=500;
unsigned long averaged_sample_max = 20;

CapacitiveSensor   cs_4_2 = CapacitiveSensor(4,2);        // 10 megohm resistor between pins 4 & 2, pin 2 is sensor pin, add wire, foil
int ledPin = 9;

const int numReadings = 20;

int readings[numReadings];      // the readings from the analog input
int readIndex = 0;              // the index of the current reading
int total = 0;                  // the running total
int average = 0;                // the average

void setup()                    
{
   pinMode(ledPin, OUTPUT);
   cs_4_2.set_CS_AutocaL_Millis(0xFFFFFFFF);     // turn off autocalibrate on channel 1 - just as an example
   Serial.begin(9600);
   // initialize all the readings to 0:
   for (int thisReading = 0; thisReading < numReadings; thisReading++) {
      readings[thisReading] = 0;
   }
   
}

void loop()                    
{
    long start = millis();
    long total1 =  cs_4_2.capacitiveSensor(30);
    Serial.print(total1);
    Serial.print("\n");
    
    // subtract the last reading:
    total = total - readings[readIndex];
    // read from the sensor:
    readings[readIndex] = cs_4_2.capacitiveSensor(30);
    // add the reading to the total:
    total = total + readings[readIndex];
    // advance to the next position in the array:
    readIndex = readIndex + 1;
  
    // if we're at the end of the array...
    if (readIndex >= numReadings) {
      // ...wrap around to the beginning:
      readIndex = 0;
    }
  
    // calculate the average:
    average = total / numReadings;
    // send it to the computer as ASCII digits
    Serial.println(average);
    
    delay(delaytime);                             // arbitrary delay to limit data to serial port 
    
}
