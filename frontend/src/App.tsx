import ResolutionForm from "./components/ResolutionForm";
import ResolutionList from "./components/ResolutionList";
import ResolutionsProvider from "./contexts/ResolutionsProvider";

function App() {
	return (
		<>
			<ResolutionsProvider>
				<ResolutionList />
				<ResolutionForm />
			</ResolutionsProvider>
		</>
	);
}

export default App;
