import { useEffect } from "react";
import {
	setResolutions,
	useResolutions,
	useResolutionsDispatch,
} from "../contexts/Resolutions";

function ResolutionList() {
	const resolutions = useResolutions();
	const resolutionsDispatch = useResolutionsDispatch();

	useEffect(() => {
		const token = localStorage.getItem("token");

		fetch("http://localhost:3000/current-user/resolutions", {
			headers: [["Authorization", `Bearer ${token}`]],
		}).then((response) =>
			response.json().then((result) => {
				resolutionsDispatch(setResolutions(result));
			})
		);
	}, []);

	return (
		<ul>
			{resolutions.map((resolution) => (
				<li key={resolution.id}>{resolution.name}</li>
			))}
		</ul>
	);
}

export default ResolutionList;
