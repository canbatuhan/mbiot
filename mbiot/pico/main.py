import rp2
import time
import network
from umqtt.simple import MQTTClient
from mbiot.sensor import SoilMoistureSensor
from mbiot.configs import *

rp2.country("TR")

class EndDevice:
    def __init__(self) -> None:
        self.__wlan = network.WLAN(network.STA_IF)
        self.__mqtt_client = MQTTClient(
            client_id=CLIENT_ID,
            server=MQTT_PUBLIC)
        self.__sensor = SoilMoistureSensor(SENSOR_ID, ANALOG_PIN)
        
    def __timeout_ms(self, duration) -> None:
        time.sleep(duration/1000)
        
    def __connect_to_wifi(self) -> None:
        poll_time = 2 # in seconds
        while not self.__wlan.isconnected():
            self.__wlan.active(True)
            self.__wlan.connect(WIFI_AP, WIFI_PW)
            self.__timeout_ms(poll_time*1000)
            if poll_time < 4096:
                poll_time *= 2
        
    def __connect_to_mqtt(self) -> None:
        self.__mqtt_client.connect()
        self.__timeout_ms(2000) # sleep 2 seconds
            
    def __start(self) -> None:
        try:
            self.__connect_to_wifi()
            self.__connect_to_mqtt()
        except OSError as err:
            self.__timeout_ms(15000) # sleep 15 seconds
            self.__start() # Restart
        
    def run(self) -> None:
        self.__start()
        while True:
            while self.__wlan.isconnected():
                package = self.__sensor.generate_package()
                self.__mqtt_client.publish(MQTT_TOPIC, str(package))
                self.__timeout_ms(1000) # sleep 1 seconds
            self.__start() # Restart
            
    
if __name__ == "__main__":
    end_device = EndDevice()
    end_device.run()
    
