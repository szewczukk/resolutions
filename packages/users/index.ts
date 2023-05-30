import express from 'express';
import assert from 'assert';
import { hash } from 'bcrypt';
import { PrismaClient } from '@prisma/client';

const app = express();
const prisma = new PrismaClient();

app.use(express.json());

app.get('/api/v1/users', async (req, res) => {
	const users = await prisma.user.findMany({});

	res.send(users);
});

app.post('/api/v1/users', async (req, res) => {
	const { username, password } = req.body;
	assert(typeof username === 'string');
	assert(typeof password === 'string');

	const hashedPassword = await hash(password, 10);

	const user = await prisma.user.create({
		data: { username, password: hashedPassword },
	});

	res.status(201).send(user);
});

app.listen(3000, () => console.log('Listening on 3000..'));

prisma.$disconnect();
