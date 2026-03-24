import { defaultConfig, getNeededCheckpoints } from "./api.js";
import { updateChart } from "./updateChart.js";

document.getElementById("uploadTrigger").addEventListener("click", () => {
	document.getElementById("file").click();
})

document.getElementById("file").addEventListener("change", (event) => {
    const file = event.target.files[0];
    const fileStatus = document.getElementById("file-name");
    if (file) {
        fileStatus.textContent = file.name;
    } else {
        fileStatus.textContent = "No file selected";
    }
});

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
