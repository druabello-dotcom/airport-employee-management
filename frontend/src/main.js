import { getNeededCheckpoints } from "./api.js";

document.getElementById("form").addEventListener("submit", event => {
	event.preventDefault();

	const file = document.getElementById("file").files[0];
	if (!file) {
		console.warn("No file selected");
		return;
	}

	console.log(getNeededCheckpoints(file));
});
