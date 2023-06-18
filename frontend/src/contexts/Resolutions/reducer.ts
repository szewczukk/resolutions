import { Dispatch, useReducer } from "react";
import { Resolution } from "../../types";

export const setResolutions = (resolutions: Resolution[]) =>
	({
		type: "SET_RESOLUTIONS",
		payload: { resolutions },
	} as const);

export const addResolution = (resolution: Resolution) =>
	({
		type: "ADD_RESOLUTION",
		payload: { resolution },
	} as const);

type Action =
	| ReturnType<typeof setResolutions>
	| ReturnType<typeof addResolution>;

export type Store = Resolution[];
export type ActionDispatch = Dispatch<Action>;

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

export function useResolutionsReducer(): [Store, Dispatch<Action>] {
	return useReducer(reducer, []);
}
