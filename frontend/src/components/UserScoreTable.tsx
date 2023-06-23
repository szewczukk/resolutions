import { useEffect, useState } from "react";

type UserScore = {
	userId: number;
	username: string;
	score: number;
};

function UserScoreTable() {
	const [userScores, setUserScores] = useState<UserScore[]>([]);

	useEffect(() => {
		fetch("http://localhost:3000/users").then((response) =>
			response.json().then((result) => setUserScores(result))
		);
	}, []);

	return (
		<table style={{ display: "block" }}>
			<thead>
				<td>Username</td>
				<td>Score</td>
			</thead>
			<tbody>
				{userScores.map((score) => (
					<tr key={score.userId}>
						<td>{score.username}</td>
						<td>{score.score}</td>
					</tr>
				))}
			</tbody>
		</table>
	);
}

export default UserScoreTable;
