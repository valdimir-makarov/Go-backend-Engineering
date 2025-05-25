from kafka import KafkaConsumer, KafkaProducer
import json
import os

consumer = KafkaConsumer(
    'code-submissions',
    bootstrap_servers=['kafka:9093'],
    auto_offset_reset='earliest',
    group_id='python-executor-group'
)
producer = KafkaProducer(bootstrap_servers=['kafka:9093'])

print("Python worker started")
for message in consumer:
    submission = json.loads(message.value.decode('utf-8'))
    if submission['Language'] != 'python':
        continue
    print(f"Executing Python code: {submission['Code']}")
    response = {"status_message": "Processed", "output": "", "error": ""}
    try:
        with open('code.py', 'w') as f:
            f.write(submission['Code'])
        output = os.popen('python3 code.py').read()
        response = {"status_message": "Processed", "output": output, "error": ""}
        print(f"Result: {output}")
    except Exception as e:
        response = {"status_message": "Runtime Error", "output": "", "error": str(e)}
        print(f"Error: {e}")
    producer.send('results', json.dumps(response).encode('utf-8'), key=submission['ID'].encode('utf-8'))