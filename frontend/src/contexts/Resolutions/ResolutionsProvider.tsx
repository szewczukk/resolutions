import { FC, PropsWithChildren, createContext } from "react";
import { ActionDispatch, Store, useResolutionsReducer } from "./reducer";

export const ResolutionContext = createContext<Store>([]);
export const ResolutionDispatchContext = createContext<ActionDispatch>(
	{} as ActionDispatch
);

const ResolutionsProvider: FC<PropsWithChildren> = ({ children }) => {
	const [store, dispatch] = useResolutionsReducer();

	return (
		<ResolutionContext.Provider value={store}>
			<ResolutionDispatchContext.Provider value={dispatch}>
				{children}
			</ResolutionDispatchContext.Provider>
		</ResolutionContext.Provider>
	);
};

export default ResolutionsProvider;
