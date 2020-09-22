#!/usr/bin/env python
# coding=utf-8
import sys
import logging
import time
from proton.handlers import MessagingHandler
from proton.reactor import Container
import hashlib
import hmac
import base64

reload(sys)
sys.setdefaultencoding('utf-8')
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)
console_handler = logging.StreamHandler(sys.stdout)


def current_time_millis():
    return str(int(round(time.time() * 1000)))


def do_sign(secret, sign_content):
    m = hmac.new(secret, sign_content, digestmod=hashlib.sha1)
    return base64.b64encode(m.digest())


class AmqpClient(MessagingHandler):
    def __init__(self):
        super(AmqpClient, self).__init__()

    def on_start(self, event):
        url = "amqps://1295579658967202.iot-amqp.cn-shanghai.aliyuncs.com"
        accessKey = "LTAI4FrP4ou2YzggPqTy9Sq1"
        accessSecret = "cRrjkxUaCKIdwOrON5gjzAdCap5ltp"
        consumerGroupId = "DEFAULT_GROUP"
        clientId = "cdd9eed0-51f6-11ea-9a49-479b81bbee3f"
        signMethod = "hmacsha1"
        timestamp = current_time_millis()
        userName = clientId + "|authMode=aksign" + ",signMethod=" + signMethod \
                        + ",timestamp=" + timestamp + ",authId=" + accessKey \
                        + ",consumerGroupId=" + consumerGroupId + "|"
        signContent = "authId=" + accessKey + "&timestamp=" + timestamp
        passWord = do_sign(accessSecret.encode("utf-8"), signContent.encode("utf-8"))
        conn = event.container.connect(url, user=userName, password=passWord, heartbeat=60)
        self.receiver = event.container.create_receiver(conn)

    def on_connection_opened(self, event):
        logger.info("Connection established, remoteUrl: %s", event.connection.hostname)

    def on_connection_closed(self, event):
        logger.info("Connection closed: %s", self)

    def on_connection_error(self, event):
        logger.info("Connection error")

    def on_transport_error(self, event):
        if event.transport.condition:
            if event.transport.condition.info:
                logger.error("%s: %s: %s" % (
                    event.transport.condition.name, event.transport.condition.description,
                    event.transport.condition.info))
            else:
                logger.error("%s: %s" % (event.transport.condition.name, event.transport.condition.description))
        else:
            logging.error("Unspecified transport error")

    def on_message(self, event):
        message = event.message
        content = message.body.decode('utf-8')
        topic = message.properties.get("topic")
        message_id = message.properties.get("messageId")
        if topic == "/a1sEMfCvR11/device1/thing/event/command/post":
            print("receive message: message_id=%s, topic=%s, content=%s" % (message_id, topic, content))
        event.receiver.flow(1)


Container(AmqpClient()).run()
