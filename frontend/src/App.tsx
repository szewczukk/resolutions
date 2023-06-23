import { createBrowserRouter, RouterProvider } from "react-router-dom";
import LoginForm from "./components/LoginForm";
import ResolutionForm from "./components/ResolutionForm";
import ResolutionList from "./components/ResolutionList";
import { ResolutionsProvider } from "./contexts/Resolutions";
import LogoutButton from "./components/LogoutButton";
import UserHeading from "./components/UserHeading";
import RegisterForm from "./components/RegisterForm";

const router = createBrowserRouter([
	{
		path: "/",
		element: (
			<>
				<ResolutionsProvider>
					<UserHeading />
					<ResolutionList />
					<ResolutionForm />
					<LogoutButton />
				</ResolutionsProvider>
			</>
		),
	},
	{
		path: "/login",
		element: (
			<>
				<LoginForm />
			</>
		),
	},
	{
		path: "/register",
		element: (
			<>
				<RegisterForm />
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
