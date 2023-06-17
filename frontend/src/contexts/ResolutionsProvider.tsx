import {
	Dispatch,
	FC,
	PropsWithChildren,
	createContext,
	useContext,
	useReducer,
} from "react";
import { Resolution } from "../types";

const ResolutionContext = createContext<Store>([]);
const ResolutionDispatchContext = createContext<Dispatch<Action>>(
	{} as Dispatch<Action>
);

const ResolutionsProvider: FC<PropsWithChildren> = ({ children }) => {
	const [store, dispatch] = useReducer(reducer, initialValue);

	return (
		<ResolutionContext.Provider value={store}>
			<ResolutionDispatchContext.Provider value={dispatch}>
				{children}
			</ResolutionDispatchContext.Provider>
		</ResolutionContext.Provider>
	);
};

export default ResolutionsProvider;

export function useResolutions(): Store {
	return useContext(ResolutionContext);
}

export function useResolutionsDispatch(): Dispatch<Action> {
	return useContext(ResolutionDispatchContext);
}

type Action =
	| {
			type: "SET_RESOLUTIONS";
			payload: { resolutions: Resolution[] };
	  }
	| {
			type: "ADD_RESOLUTION";
			payload: { resolution: Resolution };
	  };

type Store = Resolution[];
const initialValue: Store = [];

const reducer = (state: Store, action: Action): Store => {
	switch (action.type) {
		case "SET_RESOLUTIONS":
			return action.payload.resolutions;

		case "ADD_RESOLUTION":
			return [...state, action.payload.resolution];

		default:
			return state;
	}
};
