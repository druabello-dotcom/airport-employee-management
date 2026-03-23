import { getNeededCheckpoints } from "./api.js";
import { updateChart } from "./updateChart.js";

document.getElementById("uploadTrigger").addEventListener("click", () => {
	document.getElementById("file").click();
})

document.getElementById("form").addEventListener("submit", async event => {
	event.preventDefault();

	const file = document.getElementById("file").files[0];
	if (!file) {
		console.warn("No file selected");
		return;
	}

	let results = await getNeededCheckpoints(file);
	updateChart(results);
});
