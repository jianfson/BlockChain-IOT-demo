#!/usr/bin/env python3
# coding=utf-8

import time
import aliyunsdkiotclient.AliyunIotMqttClient as AliyunIot
import configparser

import os
import threading
from threading import *
from multiprocessing import Process, Lock
import json
import logging

config = configparser.ConfigParser()
config.read('iot.cfg')

# IoT
PRODUCE_KEY = config['register']['productKey']
DEVICE_NAME = config['register']['deviceName']
DEVICE_SECRET = config['register']['deviceSecret']

HOST = PRODUCE_KEY + '.iot-as-mqtt.cn-shanghai.aliyuncs.com'
SUBSCRIBE_TOPIC = "/" + PRODUCE_KEY + "/" + DEVICE_NAME + "/control"

iotClient = AliyunIot.getAliyunIotMqttClient(PRODUCE_KEY,DEVICE_NAME, DEVICE_SECRET, secure_mode=3)

logging.basicConfig(
    level=logging.DEBUG,
    format="(%(threadName)-10s) %(message)s",
)

messageEvent = threading.Event()
iotMessage = None
lock = Lock()

def wait_for_event_timeout(e, t):
    """
    wait t seconds and then timeout
    :param e: Event
    :param t: (seconds)
    :return:
    """
    global iotMessage
    while True:
        logging.debug("wait_for_event_timeout")
        event_is_set = e.wait()
        logging.debug("event set:%s", event_is_set)
        lock.acquire()
        logging.debug("iotMessage:%s", str(iotMessage))
        try:
            setjson = json.loads(iotMessage.payload.decode('utf-8'))
        except Exception as e:
            print(e)
        lock.release()
        status = setjson['params']['flag']
        #res = os.popen("ps -ef | grep arecord | grep -v grep | awk '{print $2}'")
        #ps = res.read()
        os.system("pkill -9 arecord")
        if status == 1:
            topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/property/post'
            info = 'filename = '+str(self.filename)+'\n' +'channels = '+str(self.wf.getnchannels()) +'\n'+'framerate = '+str(self.wf.getframerate())
            payload_json = {
                'id': int(time.time()),
                'params': {
                    'information': info,
                    'stopStatus': 1,
                    'pauseStatus': 1
                },
                'method': "thing.event.property.post"
            }
            print('send data to iot server: ' + str(payload_json))        
            iotClient.publish(topic, payload=str(payload_json))
        messageEvent.clear()

def on_connect(client, userdata, flags, rc):
    print('subscribe '+SUBSCRIBE_TOPIC)
    client.subscribe(topic=SUBSCRIBE_TOPIC)


def on_message(client, userdata, msg):
    global iotMessage
    #print('receive message topic :'+ msg.topic)
    if msg.topic[-5:] == 'reply':
        return
    print('receive message topic :'+ msg.topic)
    try:
        lock.acquire()
        iotMessage = msg
        lock.release()
    except Exception as e:
        print(e)
    messageEvent.set()

if __name__ == '__main__':
    iotClient.on_connect = on_connect
    iotClient.on_message = on_message
    iotClient.connect(host=HOST, port=1883, keepalive=60)

    messageLoop = threading.Thread(
        name="messageLoop",
        target=wait_for_event_timeout,
        args=(messageEvent, iotMessage)
    )
    messageLoop.start()

    # loop
    iotClient.loop_forever()
