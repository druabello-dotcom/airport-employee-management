import { getNeededCheckpoints } from "./api.js";

document.getElementById("form").addEventListener("submit", async event => {
	event.preventDefault();

	const file = document.getElementById("file").files[0];
	if (!file) {
		console.warn("No file selected");
		return;
	}

	let results = await getNeededCheckpoints(file);
	console.log(results);
});
