const { Amqp } = require('@spectacles/brokers');

const broker = new Amqp('votes');

broker.on('error', console.error);

(async () => {

	broker.on('VOTE', (data, { ack }) => {
		ack();
		console.dir(data);
	});

	console.log('connecting');
	await broker.connect('fyko:doctordoctor@localhost//')
	console.log('connected');

	broker.subscribe('VOTE');
})();