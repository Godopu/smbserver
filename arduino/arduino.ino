// Library Include
// #include <DHT.h>
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

StaticJsonDocument<100> sendDoc;
int preis = 0;
int is = 0;
int al = 0;

char *ip = "0.0.0.0";
char *ti = "00.00";

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
  Serial.begin(57600); // Serial Monitor for Debug
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

void setStatus()
{
  if (preis != is)
  {
    lcd.clear();
    preis = is;
  }

  if (is == 0)
  {
    printLCD(0, 0, "Please connect");
    printLCD(0, 1, "to internet");
  }
  else
  {
    printLCD(0, 0, ip);
    printLCD(0, 1, ti);
  }

  if (al == 1)
  {
    tone(piezoPin, 391.9954, 2000);
  }
}
void sendStatus()
{
  sendDoc["code"] = code;
  sendDoc["is"] = is;
  sendDoc["al"] = al;
  sendDoc["ip"] = ip;
  sendDoc["ti"] = ti;

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
      delay(1000);
      return;
    }
    else if (code == 200)
    {
      al = recvDoc["al"];
      is = recvDoc["is"];
      ip = recvDoc["ip"];
      ti = recvDoc["ti"];
    }
    sendStatus();
  }
  timer++;
  setStatus();
  delay(1000);
}
