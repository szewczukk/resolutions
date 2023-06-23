import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

function UserHeading() {
	const navigate = useNavigate();
	const [username, setUsername] = useState("");

	useEffect(() => {
		const token = localStorage.getItem("token");

		fetch("http://localhost:3000/current-user/", {
			headers: [["Authorization", `Bearer ${token}`]],
		}).then((response) =>
			response
				.json()
				.then((result) => {
					setUsername(result.username);
				})
				.catch(() => {
					navigate("/login");
				})
		);
	}, []);

	return <h1>Hello, {username}!</h1>;
}

export default UserHeading;
