import pika
import json
import smtp.lib
from email.mime.Text import MIMEText

connection = pika.BlockingConnection(pika.ConnectionParameters( host = 'localhost' ))

channel = connection.channel()
channel.queue_declare(queue='email')

print ' [*] Waiting for messages. To exit press CTRL+C'

def callback(channel, method, properties, body):
    print ' [x] Received %r ' % body
    parsed = json.loads(body)
    msg = MIMEText()
    msg['From'] = 'Me'
    msg['To'] = parsed['email']
    msg['Subject'] = parsed['message']
    s = smtplib.SMTP('localhost')
    s.sendmail('Me', parsed['email'], msg.as_string())
    s.quit()

channel.basic_consume(callback, queue='email',no_ack=True)

channel.start_consuming()