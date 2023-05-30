import express from 'express';
import assert from 'assert';
import jwt from 'jsonwebtoken';
import { hash, compare } from 'bcrypt';
import { PrismaClient } from '@prisma/client';

const app = express();
const prisma = new PrismaClient();

app.use(express.json());

app.get('/api/v1/users', async (req, res) => {
	const { username } = req.query;

	if (username === undefined) {
		const users = await prisma.user.findMany();

		res.send(users);
		return;
	}

	const user = await prisma.user.findUnique({
		where: { username: username.toString() },
	});
	res.send(user);
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

app.post('/api/v1/login', async (req, res) => {
	const { username, password } = req.body;

	assert(typeof username === 'string');
	assert(typeof password === 'string');

	const user = await prisma.user.findUnique({
		where: { username },
	});

	if (!user) {
		res.status(404).send();
		return;
	}

	const isPasswordCorrect = await compare(password, user.password);

	if (!isPasswordCorrect) {
		res.status(401).send();
		return;
	}

	const token = jwt.sign({ userId: user.id }, 'secret', { expiresIn: '14d' });

	res.status(200).send({ token });
});

app.listen(3000, () => console.log('Listening on 3000..'));

prisma.$disconnect();
