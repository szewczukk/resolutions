import { FormEvent, useState } from "react";
import { addResolution, useResolutionsDispatch } from "../contexts/Resolutions";

function ResolutionForm() {
	const [name, setName] = useState("");
	const resolutionsDispatch = useResolutionsDispatch();

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const token = localStorage.getItem("token");

		const response = await fetch(
			"http://localhost:3000/current-user/resolutions",
			{
				method: "POST",
				headers: [
					["Content-Type", "application/json"],
					["Authorization", `Bearer ${token}`],
				],
				body: JSON.stringify({ name }),
			}
		);

		const result = await response.json();

		if (response.status === 200) {
			resolutionsDispatch(addResolution(result));
		}
	};

	const onNameChange = (e: FormEvent<HTMLInputElement>) => {
		setName(e.currentTarget.value);
	};

	return (
		<form onSubmit={onSubmit}>
			<label>
				Name
				<input type="text" onChange={onNameChange} required />
			</label>
			<input type="submit" />
		</form>
	);
}

export default ResolutionForm;
