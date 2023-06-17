import { useEffect, useState } from "react";

interface Resolution {
	userId: number;
	name: string;
	ID: number;
}

function ResolutionList() {
	const [resolutions, setResolutions] = useState<Resolution[]>([]);

	useEffect(() => {
		fetch("http://localhost:3001/").then(
			response => response.json().then(
				result => {
					console.log(result)
					setResolutions(result);
				}
			)
		)
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