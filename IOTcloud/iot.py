#!/usr/bin/env python3
# coding=utf-8
import time
import aliyunsdkiotclient.AliyunIotMqttClient as AliyunIot
import configparser
from aliyunsdkcore.client import AcsClient
from aliyunsdkcore.request import CommonRequest
from aliyunsdkcore import acs_exception

import os
import threading
from threading import *
from multiprocessing import Process, Lock
import ali_speech
from ali_speech.callbacks import SpeechRecognizerCallback
from ali_speech.constant import ASRFormat
from ali_speech.constant import ASRSampleRate
import json
from pydub import AudioSegment
import logging

import pyaudio
import wave
from playsound import playsound

config = configparser.ConfigParser()
config.read('iot.cfg')

# IoT
PRODUCE_KEY = config['IOT']['productKey']
DEVICE_NAME = config['IOT']['deviceName']
DEVICE_SECRET = config['IOT']['deviceSecret']

HOST = PRODUCE_KEY + '.iot-as-mqtt.cn-shanghai.aliyuncs.com'
SUBSCRIBE_TOPIC = "/" + PRODUCE_KEY + "/" + DEVICE_NAME + "/control";

iotClient = AliyunIot.getAliyunIotMqttClient(PRODUCE_KEY,DEVICE_NAME, DEVICE_SECRET, secure_mode=3)
#iotClient.on_connect = on_connect
#iotClient.on_message = on_message
#iotClient.connect(host=HOST, port=1883, keepalive=60)

# nls
NLS_AK = config['NLS']['nlsAccessKey']
NLS_AK_SECRET = config['NLS']['nlsAccessKeySecret']
NLS_APP_KEY = config['NLS']['nlsAppKey']

# audio
audio_name = config['AUDIO']['audioName']
audio_file = config['AUDIO']['audioFile']

acsclient = AcsClient(
   NLS_AK,
   NLS_AK_SECRET,
   "cn-shanghai"
);

nlsClient = ali_speech.NlsClient()
# 设置输出日志信息的级别：DEBUG、INFO、WARNING、ERROR
nlsClient.set_log_level('INFO')


logging.basicConfig(
    level=logging.DEBUG,
    format="(%(threadName)-10s) %(message)s",
)

token, expire_time = nlsClient.create_token(NLS_AK, NLS_AK_SECRET)
print('token: %s, expire time(s): %s' % (token, expire_time))
if expire_time:
    print('token %s' % (time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(expire_time))))

messageEvent = threading.Event()
#recEvent = threading.Event()
iotMessage = None
lock = Lock()
#recStatus = 0

playMessageEvent = threading.Event()
playMessage = None

tr = None
CHUNK = 1024
class MyMusic(Thread):
    def init(self,filename):
        self.filename = filename
        self.wf = wave.open(filename, 'rb')  # (sys.argv[1], 'rb')
        self.p = pyaudio.PyAudio()
        self.stream = self.p.open(format=self.p.get_format_from_width(self.wf.getsampwidth()),
                             channels=self.wf.getnchannels(),
                             rate=self.wf.getframerate(),
                             output=True)
        self.data = self.wf.readframes(CHUNK)

        self.datas = []
        while len(self.data) > 0:
            self.data = self.wf.readframes(CHUNK)
            self.datas.append(self.data)

        self.__flag =  Event()
        #self.__flag.set()
        self.ifdo = True;
        self.restartStatus = False

    def run (self):

        while self.ifdo :
            i = 0
            for data in self.datas:
                self.__flag.wait()
                if self.restartStatus:
                    break
                    break
                #print('I am running...')
                progress = int(i/len(self.datas)*100)
                topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/property/post'
                payload_json = {
                    'id': int(time.time()),
                    'params': {
                        'progress': progress
                    },
                    'method': "thing.event.property.post"
                }
                #print('send data to iot server: ' + str(payload_json))
                iotClient.publish(topic, payload=str(payload_json))
                i = i+1
                # time.sleep(2)

                self.stream.write(data)
            self.restartStatus = False
            self.wf = wave.open(self.filename, 'rb')
            self.data = self.wf.readframes(CHUNK)
            self.datas = []
            while len(self.data) > 0:
                self.data = self.wf.readframes(CHUNK)
                self.datas.append(self.data)
        # self.data = ''
    def pause(self):
        self.__flag.clear()
        print("pause")
        topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/property/post'
        payload_json = {
            'id': int(time.time()),
            'params': {
                'pauseStatus': 0
            },
            'method': "thing.event.property.post"
        }
        print('send data to iot server: ' + str(payload_json))        
        iotClient.publish(topic, payload=str(payload_json))

    def resume(self):
        self.__flag.set()
        print("resume")
        topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/property/post'
        payload_json = {
            'id': int(time.time()),
            'params': {
                'pauseStatus': 1
            },
            'method': "thing.event.property.post"
        }
        print('send data to iot server: ' + str(payload_json))        
        iotClient.publish(topic, payload=str(payload_json))
    def stop (self):
        print('I am stopping it...')
        self.restartStatus = True
        self.__flag.clear()

        topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/property/post'
        payload_json = {
            'id': int(time.time()),
            'params': {
                'stopStatus': 0,
                'pauseStatus': 0,
                'progress': 0
            },
            'method': "thing.event.property.post"
        }
        print('send data to iot server: ' + str(payload_json))        
        iotClient.publish(topic, payload=str(payload_json))
    def restart(self):
        self.wf = wave.open(self.filename, 'rb')
        self.data = self.wf.readframes(CHUNK)
        self.datas = []
        while len(self.data) > 0:
            self.data = self.wf.readframes(CHUNK)
            self.datas.append(self.data)
        #print(self.wf.getparams())
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

        self.__flag.set()

def play_wait_for_event():
    """
    waiting for the event to be set before do anything
    :param e: Event
    :return:
    """
    logging.debug("playing")
    global playMessage
    while True:
        logging.debug("wait_for_event_timeout")
        event_is_set = e.wait()
        logging.debug("event set:%s", event_is_set)
        lock.acquire()
        logging.debug("playMessage:%s", str(playMessage))
        try:
            setjson = json.loads(iotMessage.payload.decode('utf-8'))
        except Exception as e:
            print(e)
        lock.release()
        status = setjson['params']['Status']
        #res = os.popen("ps -ef | grep arecord | grep -v grep | awk '{print $2}'")
        #ps = res.read()
        os.system("pkill -9 arecord")
        if status == 1:
            #if ps.strip():
            #os.system("pkill -9 arecord")
            recLoop = threading.Thread(
                name="recLoop",
                target=rec_wait_for_event
            )
            recLoop.start()
        elif status == 0:
            #if ps.strip():
            #os.system("pkill -9 arecord")
            process(nlsClient, NLS_APP_KEY, token)
        messageEvent.clear()

def rec_wait_for_event():
    """
    waiting for the event to be set before do anything
    :param e: Event
    :return:
    """
    logging.debug("recording")
    #event_is_set = e.wait()
    #logging.debug("event set:%s", event_is_set)
    os.system("arecord -Dhw:1,0 -f s16_le -r 16000 test.wav -vv -i")


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
            #if ps.strip():
            #os.system("pkill -9 arecord")
            recLoop = threading.Thread(
                name="recLoop",
                target=rec_wait_for_event
            )
            recLoop.start()
        elif status == 0:
            #if ps.strip():
            #os.system("pkill -9 arecord")
            process(nlsClient, NLS_APP_KEY, token)
        messageEvent.clear()

class MyCallback(SpeechRecognizerCallback):
    """
    """
    def __init__(self, name='default'):
        self._name = name
    def on_started(self, message):
        print('MyCallback.OnRecognitionStarted: %s' % message)
    def on_result_changed(self, message):
        print('MyCallback.OnRecognitionResultChanged: file: %s, task_id: %s, result: %s' % (
            self._name, message['header']['task_id'], message['payload']['result']))
    def on_completed(self, message):
        print('MyCallback.OnRecognitionCompleted: file: %s, task_id:%s, result:%s' % (
            self._name, message['header']['task_id'], message['payload']['result']))
        command = message['payload']['result'].replace('。', '')
        if command == '播放':
            tr.restart()
        elif command == '暂停':
            tr.pause()
        elif command == '恢复':
            tr.resume()
        elif command == '停止':
            tr.stop()
        print('command:   '+command)
        topic = '/sys/'+PRODUCE_KEY+'/'+DEVICE_NAME+'/thing/event/command/post'
        payload_json = {
            'id': int(time.time()),
            'params': {
                'command': command
            },
            'method': "thing.event.command.post"
        }
        print('send data to iot server: ' + str(payload_json))        
        iotClient.publish(topic, payload=str(payload_json))
    def on_task_failed(self, message):
        print('MyCallback.OnRecognitionTaskFailed: %s' % message)
    def on_channel_closed(self):
        print('MyCallback.OnRecognitionChannelClosed')
def process(client, appkey, token):
    song = AudioSegment.from_wav(audio_name).set_frame_rate(16000)
    song.export(audio_name, format='wav', bitrate='16k')
    callback = MyCallback(audio_name)
    recognizer = client.create_recognizer(callback)
    recognizer.set_appkey(appkey)
    recognizer.set_token(token)
    recognizer.set_format(ASRFormat.PCM)
    recognizer.set_sample_rate(ASRSampleRate.SAMPLE_RATE_16K)
    recognizer.set_enable_intermediate_result(False)
    recognizer.set_enable_punctuation_prediction(True)
    recognizer.set_enable_inverse_text_normalization(True)
    try:
        ret = recognizer.start()
        if ret < 0:
            return ret
        print('sending audio...')
        with open(audio_name, 'rb') as f:
            audio = f.read(3200)
            while audio:
                ret = recognizer.send(audio)
                if ret < 0:
                    break
                time.sleep(0.1)
                audio = f.read(3200)
        recognizer.stop()
    except Exception as e:
        print(e)
    finally:
        recognizer.close()
def process_multithread(client, appkey, token, number):
    thread_list = []
    for i in range(0, number):
        thread = threading.Thread(target=process, args=(client, appkey, token))
        thread_list.append(thread)
        thread.start()
    for thread in thread_list:
        thread.join()

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
    #client = AliyunIot.getAliyunIotMqttClient(PRODUCE_KEY,DEVICE_NAME, DEVICE_SECRET, secure_mode=3)
    iotClient.on_connect = on_connect
    iotClient.on_message = on_message
    iotClient.connect(host=HOST, port=1883, keepalive=60)

    tr = MyMusic()
    tr.init( audio_file )
    # tr.setDaemon(True)
    tr.start()

    messageLoop = threading.Thread(
        name="messageLoop",
        target=wait_for_event_timeout,
        args=(messageEvent, iotMessage)
    )
    messageLoop.start()

    # loop
    iotClient.loop_forever()
