import LoginForm from "./components/LoginForm";
import ResolutionForm from "./components/ResolutionForm";
import ResolutionList from "./components/ResolutionList";
import { ResolutionsProvider } from "./contexts/Resolutions";

function App() {
	return (
		<main>
			<LoginForm />
			<ResolutionsProvider>
				<ResolutionList />
				<ResolutionForm />
			</ResolutionsProvider>
		</main>
	);
}

export default App;
