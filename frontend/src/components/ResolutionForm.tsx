import { FormEvent, useState } from "react";
import { addResolution, useResolutionsDispatch } from "../contexts/Resolutions";

function ResolutionForm() {
	const [name, setName] = useState("");
	const [userId, setUserId] = useState(0);
	const resolutionsDispatch = useResolutionsDispatch();

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const response = await fetch("http://localhost:3002/resolutions/", {
			method: "POST",
			headers: [["Content-Type", "application/json"]],
			body: JSON.stringify({ name, userId }),
		});

		const result = await response.json();

		if (response.status === 200) {
			resolutionsDispatch(addResolution(result));
		}
	};

	const onNameChange = (e: FormEvent<HTMLInputElement>) => {
		setName(e.currentTarget.value);
	};

	const onUserIdChange = (e: FormEvent<HTMLInputElement>) => {
		setUserId(parseInt(e.currentTarget.value));
	};

	return (
		<form onSubmit={onSubmit}>
			<label>
				Name
				<input type="text" onChange={onNameChange} required />
			</label>
			<label>
				UserId
				<input type="number" onChange={onUserIdChange} required />
			</label>
			<input type="submit" />
		</form>
	);
}

export default ResolutionForm;
