import { FormEvent, useState } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";

function LoginForm() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const navigate = useNavigate();

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const response = await fetch("http://localhost:3000/register/", {
			method: "POST",
			headers: [["Content-Type", "application/json"]],
			body: JSON.stringify({ username, password }),
		});

		const result = await response.json();

		if (response.status === 200) {
			localStorage.setItem("token", result.token);
			navigate("/login");
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
			<input type="submit" value="Register" />
			<Link to="/login" style={{ marginLeft: "12px" }}>
				Login
			</Link>
		</form>
	);
}

export default LoginForm;
