// Library Include
#include <Servo.h>
#include <DHT.h>
// #include <Wire.h>
#include <Wire.h>
#include <LiquidCrystal_I2C.h>
#include <ArduinoJson.h>

// Global variable (전역변수 정의)

int piezoPin = 11;
int timeout = 0;
// DHT dht(DHTPIN, DHTTYPE);
LiquidCrystal_I2C lcd(0x27, 16, 2);
StaticJsonDocument<200> INFO;
StaticJsonDocument<200> recvDoc;
int code = 0;

void printLCD(int col, int row, char *str)
{
  for (int i = 0; i < strlen(str); i++)
  {
    lcd.setCursor(col + i, row);
    lcd.print(str[i]);
  }
}

void setup()
{
  // setup device information
  INFO["uuid"] = "DEVICE-A-UUID";
  INFO["code"] = 100;

  lcd.init();
  lcd.backlight();
  // HW I/O define
  // tone(piezoPin, 391.9954, 5000);
  // printLCD(0, 0, "Smart Farm ");
  // printLCD(0, 1, "Are You Ready? ");

  setStatus();
  Serial.begin(9600); // Serial Monitor for Debug
  while (!Serial)
    ;
  // Serial1.begin(9600); // Bluetooth Module
}

void broadcastUUID()
{
  serializeJson(INFO, Serial);
  Serial.println();
}

void recvMsg()
{
  // memset(&recvDoc, 0x00, 0);
  deserializeJson(recvDoc, Serial);
  code = recvDoc["code"];
}

StaticJsonDocument<100> sendDoc;
int preis = 0;
int is = 0;
int alarm = 0;

char *ip = "0.0.0.0";
char *t = "00.00";

void setStatus()
{
  if (preis != is)
  {
    lcd.clear();
  }

  if (is == 0)
  {
    printLCD(0, 0, "Please connect");
    printLCD(0, 1, "to internet");
  }
  else
  {
    printLCD(0, 0, ip);
    printLCD(0, 1, t);
  }

  if (alarm == 1)
  {
    tone(piezoPin, 391.9954, 2000);
  }
}
void sendStatus()
{
  sendDoc["code"] = code;
  sendDoc["is"] = is;
  sendDoc["alarm"] = alarm;
  sendDoc["ip"] = ip;
  sendDoc["t"] = t;

  serializeJson(sendDoc, Serial);
  Serial.println();
}

int timer = 0;
void loop()
{ // put your main code here, to run repeatedly:
  if (timer % 2 != 0)
  {
    recvMsg();
    if (code == 100)
    {
      broadcastUUID();
    }
    else if (code == 200)
    {
      alarm = recvDoc["alarm"];
      is = recvDoc["is"];
      ip = recvDoc["ip"];
      t = recvDoc["t"];
    }
    sendStatus();
  }
  timer++;
  setStatus();
  delay(1000);
}
