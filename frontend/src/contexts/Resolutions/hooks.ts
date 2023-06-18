import { useContext } from "react";
import {
	ResolutionContext,
	ResolutionDispatchContext,
} from "./ResolutionsProvider";
import { ActionDispatch, Store } from "./reducer";

export function useResolutions(): Store {
	return useContext(ResolutionContext);
}

export function useResolutionsDispatch(): ActionDispatch {
	return useContext(ResolutionDispatchContext);
}
