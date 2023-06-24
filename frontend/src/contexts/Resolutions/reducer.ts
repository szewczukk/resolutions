import { Dispatch, useReducer } from "react";
import { Resolution } from "../../types";

export const setResolutions = (resolutions: Resolution[]) =>
	({
		type: "SET_RESOLUTIONS",
		payload: { resolutions },
	} as const);

export const completeResolution = (resolutionId: number) =>
	({
		type: "COMPLETE_RESOLUTION",
		payload: { resolutionId },
	} as const);

export const addResolution = (resolution: Resolution) =>
	({
		type: "ADD_RESOLUTION",
		payload: { resolution },
	} as const);

export const deleteResolution = (resolutionId: number) =>
	({
		type: "DELETE_RESOLUTION",
		payload: { resolutionId },
	} as const);

type Action =
	| ReturnType<typeof setResolutions>
	| ReturnType<typeof addResolution>
	| ReturnType<typeof completeResolution>
	| ReturnType<typeof deleteResolution>;

export type Store = Resolution[];
export type ActionDispatch = Dispatch<Action>;

const reducer = (state: Store, action: Action): Store => {
	switch (action.type) {
		case "SET_RESOLUTIONS":
			return action.payload.resolutions;

		case "ADD_RESOLUTION":
			return [...state, action.payload.resolution];

		case "COMPLETE_RESOLUTION": {
			const id = action.payload.resolutionId;

			const completedResolution = state.find((r) => r.id === id);
			if (completedResolution === undefined) {
				return state;
			}

			const newState = state.filter((resolution) => resolution.id !== id);

			newState.push({
				...completedResolution,
				completed: true,
			});

			return newState;
		}

		case "DELETE_RESOLUTION":
			return state.filter(
				(resolution) => resolution.id !== action.payload.resolutionId
			);

		default:
			return state;
	}
};

export function useResolutionsReducer(): [Store, Dispatch<Action>] {
	return useReducer(reducer, []);
}
