import { FormEvent, useState } from "react";

function LoginForm() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const response = await fetch("http://localhost:3002/login/", {
			method: "POST",
			headers: [["Content-Type", "application/json"]],
			body: JSON.stringify({ username, password }),
		});

		const result = await response.json();

		if (response.status === 200) {
			console.log(result);
			localStorage.setItem("token", result.token);
		}
	};

	const onUsernameChange = (e: FormEvent<HTMLInputElement>) => {
		setUsername(e.currentTarget.value);
	};

	const onPasswordChange = (e: FormEvent<HTMLInputElement>) => {
		setPassword(e.currentTarget.value);
	};

	return (
		<form onSubmit={onSubmit}>
			<label>
				Username
				<input type="text" onChange={onUsernameChange} required />
			</label>
			<label>
				Password
				<input type="password" onChange={onPasswordChange} required />
			</label>
			<input type="submit" />
		</form>
	);
}

export default LoginForm;
