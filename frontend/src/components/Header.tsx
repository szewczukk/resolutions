import { Link } from "react-router-dom";

function Header() {
	const token = localStorage.getItem("token");

	if (!token) {
		return (
			<>
				<Link to="/login">Login</Link>
				<Link to="/register" style={{ marginLeft: "12px" }}>
					Register
				</Link>
			</>
		);
	}

	return <Link to="/">Home</Link>;
}

export default Header;
