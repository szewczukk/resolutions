import { useEffect } from "react";
import {
	setResolutions,
	useResolutions,
	useResolutionsDispatch,
} from "../contexts/Resolutions";
import {
	completeResolution,
	deleteResolution,
} from "../contexts/Resolutions/reducer";

function ResolutionList() {
	const resolutions = useResolutions();
	const resolutionsDispatch = useResolutionsDispatch();

	const token = localStorage.getItem("token");

	useEffect(() => {
		fetch("http://localhost:3000/current-user/resolutions", {
			headers: [["Authorization", `Bearer ${token}`]],
		}).then((response) =>
			response.json().then((result) => {
				resolutionsDispatch(setResolutions(result || []));
			})
		);
	}, []);

	const onCompleteButtonClicked = (resolutionId: number) => {
		fetch(`http://localhost:3000/resolutions/${resolutionId}/complete`, {
			method: "POST",
		}).then((response) =>
			response.text().then(() => {
				resolutionsDispatch(completeResolution(resolutionId));
			})
		);
	};

	const onDeleteButtonClicked = (resolutionId: number) => {
		fetch(`http://localhost:3000/resolutions/${resolutionId}/delete`, {
			method: "POST",
		}).then((response) =>
			response.text().then(() => {
				resolutionsDispatch(deleteResolution(resolutionId));
			})
		);
	};

	return (
		<ul>
			{resolutions.map((resolution) => (
				<li
					key={resolution.id}
					style={{
						textDecoration: resolution.completed
							? "line-through"
							: "none",
					}}
				>
					{resolution.name}

					{!resolution.completed && (
						<button
							style={{ marginLeft: "12px" }}
							onClick={() =>
								onCompleteButtonClicked(resolution.id)
							}
						>
							Complete
						</button>
					)}

					<button
						style={{ marginLeft: "12px" }}
						onClick={() => onDeleteButtonClicked(resolution.id)}
					>
						Delete
					</button>
				</li>
			))}
		</ul>
	);
}

export default ResolutionList;
