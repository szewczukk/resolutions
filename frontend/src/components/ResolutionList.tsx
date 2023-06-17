import { useEffect } from "react";
import { Resolution } from "../types";
import {
	useResolutions,
	useResolutionsDispatch,
} from "../contexts/ResolutionsProvider";

function ResolutionList() {
	const resolutions = useResolutions();
	const resolutionsDispatch = useResolutionsDispatch();

	useEffect(() => {
		fetch("http://localhost:3001/").then((response) =>
			response.json().then((result) => {
				console.log(result);
				resolutionsDispatch({
					type: "SET_RESOLUTIONS",
					payload: { resolutions: result as Resolution[] },
				});
			})
		);
	}, []);

	return (
		<ul>
			{resolutions.map((resolution) => (
				<li key={resolution.ID}>{resolution.name}</li>
			))}
		</ul>
	);
}

export default ResolutionList;
