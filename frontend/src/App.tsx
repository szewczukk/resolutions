import { createBrowserRouter, RouterProvider } from "react-router-dom";
import LoginForm from "./components/LoginForm";
import ResolutionForm from "./components/ResolutionForm";
import ResolutionList from "./components/ResolutionList";
import { ResolutionsProvider } from "./contexts/Resolutions";
import LogoutButton from "./components/LogoutButton";

const router = createBrowserRouter([
	{
		path: "/",
		element: (
			<>
				<ResolutionsProvider>
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
]);

function App() {
	return (
		<main>
			<RouterProvider router={router} />
		</main>
	);
}

export default App;
