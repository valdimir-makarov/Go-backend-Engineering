const kafka = require('kafka-node');
const client = new kafka.KafkaClient({ kafkaHost: 'kafka:9093' });
const consumer = new kafka.Consumer(
    client,
    [{ topic: 'code-submissions', partition: 0 }],
    { autoCommit: true }
);
const producer = new kafka.Producer(client);

producer.on('ready', () => console.log("Producer ready"));
producer.on('error', (err) => console.error('Producer error:', err));

console.log("Node.js worker started");
consumer.on('message', (message) => {
    const submission = JSON.parse(message.value);
    if (submission.Language !== 'nodejs') return;
    console.log(`Executing Node.js code: ${submission.Code}`);
    let response = { status_message: "Processed", output: "", error: "" };
    const fs = require('fs');
    fs.writeFileSync('code.js', submission.Code);
    const { exec } = require('child_process');
    exec('node code.js', (error, stdout, stderr) => {
        if (error) {
            response = { status_message: "Runtime Error", output: "", error: stderr };
            console.error(`Error: ${stderr}`);
        } else {
            response = { status_message: "Processed", output: stdout, error: "" };
            console.log(`Result: ${stdout}`);
        }
        producer.send([{ topic: 'results', messages: JSON.stringify(response), key: submission.ID }], (err, data) => {
            if (err) console.error('Send error:', err);
        });
    });
});

consumer.on('error', (err) => console.error('Error:', err));