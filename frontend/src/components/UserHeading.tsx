import { useEffect, useState } from "react";

function UserHeading() {
	const [username, setUsername] = useState("");

	useEffect(() => {
		const token = localStorage.getItem("token");

		fetch("http://localhost:3000/current-user/", {
			headers: [["Authorization", `Bearer ${token}`]],
		}).then((response) =>
			response.json().then((result) => {
				setUsername(result.username);
			})
		);
	}, []);

	return <h1>Hello, {username}!</h1>;
}

export default UserHeading;
