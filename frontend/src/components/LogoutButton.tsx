import { useNavigate } from "react-router-dom";

function LogoutButton() {
	const navigate = useNavigate();

	const onClick = () => {
		localStorage.removeItem("token");
		navigate("/login");
	};

	return <button onClick={onClick}>Log out</button>;
}

export default LogoutButton;
