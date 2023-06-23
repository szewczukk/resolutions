import { createBrowserRouter, Link, RouterProvider } from "react-router-dom";
import LoginForm from "./components/LoginForm";
import ResolutionForm from "./components/ResolutionForm";
import ResolutionList from "./components/ResolutionList";
import { ResolutionsProvider } from "./contexts/Resolutions";
import LogoutButton from "./components/LogoutButton";
import UserHeading from "./components/UserHeading";
import RegisterForm from "./components/RegisterForm";
import UserScoreTable from "./components/UserScoreTable";
import Header from "./components/Header";

const router = createBrowserRouter([
	{
		path: "/",
		element: (
			<>
				<Link to="/table" style={{ marginRight: "12px" }}>
					User scores table
				</Link>
				<LogoutButton />
				<ResolutionsProvider>
					<UserHeading />
					<ResolutionList />
					<ResolutionForm />
				</ResolutionsProvider>
			</>
		),
	},
	{
		path: "/login",
		element: (
			<>
				<Link to="/table">User scores table</Link>
				<LoginForm />
			</>
		),
	},
	{
		path: "/register",
		element: (
			<>
				<Link to="/table">User scores table</Link>
				<RegisterForm />
			</>
		),
	},
	{
		path: "/table",
		element: (
			<>
				<Header />
				<UserScoreTable />
			</>
		),
	},
]);

function App() {
	return (
		<main>
			<RouterProvider router={router} />
		</main>
	);
}

export default App;
