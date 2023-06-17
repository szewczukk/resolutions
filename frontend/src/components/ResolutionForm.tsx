import { FormEvent, useState } from "react";
import { useResolutionsDispatch } from "../contexts/ResolutionsProvider";

function ResolutionForm() {
	const [name, setName] = useState("");
	const [userId, setUserId] = useState(0);
	const resolutionsDispatch = useResolutionsDispatch();

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const response = await fetch("http://localhost:3001", {
			method: "POST",
			headers: [["Content-Type", "application/json"]],
			body: JSON.stringify({ name, userId }),
		});

		if (response.status === 200) {
			resolutionsDispatch({
				type: "ADD_RESOLUTION",
				payload: {
					resolution: { ID: Math.random() * 100, name, userId },
				},
			});
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
